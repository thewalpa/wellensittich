package discordws

import (
	"fmt"
	"sync"
)

type VoiceSender struct {
	GuildID   string
	playQueue *PlayQueue
	//currentPlay     *Play
	done            chan struct{}
	stop            chan struct{}
	newPlay         chan struct{}
	paused          chan bool // Channel to manage pause state
	resume          chan bool // Channel to manage resume state
	immediatePlayMu sync.Mutex
}

func NewVoiceSender(guildID string, playQueueModel *PlayQueue) *VoiceSender {
	return &VoiceSender{
		playQueue: playQueueModel,
		GuildID:   guildID,
		done:      make(chan struct{}),
		stop:      make(chan struct{}),
		newPlay:   make(chan struct{}, 1),
		paused:    make(chan bool, 1), // Buffer to prevent blocking
		resume:    make(chan bool, 1),
	}
}

// clean up
func (vs *VoiceSender) Stop() {
	close(vs.stop)
	close(vs.done)
}

// stops playing by deleting queue
func (vs *VoiceSender) StopPlaying() {
	vs.playQueue.Reset()
	close(vs.stop)
	vs.stop = make(chan struct{})
}

// skips the current play
func (vs *VoiceSender) SkipPlay() {
	vs.playQueue.SetCurrentPlay(nil) //prevent race condition
	close(vs.stop)
	vs.stop = make(chan struct{})
}

// remove the play at index i
func (vs *VoiceSender) RemovePlay(i int) error {
	if i == 0 {
		vs.SkipPlay()
		return nil
	}
	if !vs.playQueue.Unqueue(i) {
		return fmt.Errorf("queue not long enough")
	}
	return nil
}

// skips play and continues with nth play in queue
func (vs *VoiceSender) SkipTo(n int, useStartIdx bool) error {
	if !vs.playQueue.DiscardTo(n, useStartIdx) {
		return fmt.Errorf("queue not long enough")
	}
	vs.SkipPlay()
	return nil
}

// pauses playing
func (vs *VoiceSender) PausePlaying() {
	select {
	case vs.paused <- true:
		vs.playQueue.SetPaused(true)
	default:
	}
}

// resumes playing
func (vs *VoiceSender) ResumePlaying() {
	select {
	case vs.resume <- true:
		vs.playQueue.SetPaused(false)
	default:
	}
}

// This will run as long as the VoiceConnection is running
func (vs *VoiceSender) Start(wsvc *WellensittichVoiceConnection) {
	fmt.Printf("Starting VoiceSender for %v\n", wsvc.GuildID)
	for {
		select {
		case <-vs.done:
			return
		default:
			if rp, ok := vs.playQueue.Dequeue(); ok {
				vs.playQueue.SetCurrentPlay(rp)
				err := rp.Sound.Play(wsvc.OpusSend, vs.done, vs.stop, vs.paused, vs.resume)
				if err != nil {
					fmt.Println(err)
				}
			} else {
				vs.playQueue.SetCurrentPlay(nil)
				<-vs.newPlay // Wait for new item if the queue is empty
			}
		}
	}
}

func (vs *VoiceSender) EnqueuePlay(p *Play) error {
	if !vs.playQueue.Enqueue(p) {
		return fmt.Errorf("could not add to play queue")
	}
	select {
	case vs.newPlay <- struct{}{}:
	default:
	}
	fmt.Printf("Enqueued a Play\n")
	return nil
}

func (vs *VoiceSender) PlayImmediately(p *Play, wsvc *WellensittichVoiceConnection) error {
	vs.immediatePlayMu.Lock()
	defer vs.immediatePlayMu.Unlock()
	vs.PausePlaying()
	err := p.Sound.Play(wsvc.OpusSend, vs.done, vs.stop, vs.paused, vs.resume)
	if err != nil {
		return err
	}
	vs.ResumePlaying()
	return nil
}
