package discordws

import (
	"fmt"
	"time"
)

type VoiceSender struct {
	playQueue   *PlayQueue
	currentPlay *Play
	done        chan struct{}
	stop        chan struct{}
	paused      chan bool // Channel to manage pause state
	resume      chan bool // Channel to manage resume state
}

func NewVoiceSender() *VoiceSender {
	return &VoiceSender{
		playQueue: &PlayQueue{},
		done:      make(chan struct{}),
		stop:      make(chan struct{}),
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
	vs.playQueue.reset()
	close(vs.stop)
	vs.stop = make(chan struct{})
}

// skips the current play
func (vs *VoiceSender) SkipPlay() {
	close(vs.stop)
	vs.stop = make(chan struct{})
}

// remove the play at index i
func (vs *VoiceSender) RemovePlay(i int) error {
	if i == 0 {
		vs.SkipPlay()
		return nil
	}
	if !vs.playQueue.unqueue(i) {
		return fmt.Errorf("queue not long enough")
	}
	return nil
}

// skips play and continues with nth play in queue
func (vs *VoiceSender) SkipTo(n int) error {
	if !vs.playQueue.discardTo(n) {
		return fmt.Errorf("queue not long enough")
	}
	vs.SkipPlay()
	return nil
}

// pauses playing
func (vs *VoiceSender) PausePlaying() {
	select {
	case vs.paused <- true:
	default:
	}
}

// resumes playing
func (vs *VoiceSender) ResumePlaying() {
	select {
	case vs.resume <- true:
	default:
	}
}

func (vs *VoiceSender) GetQueueInfo(limit int) ([]PlayInfo, int) {
	vs.playQueue.mu.Lock()
	defer vs.playQueue.mu.Unlock()
	if vs.currentPlay == nil {
		return []PlayInfo{}, 0
	}
	info := []PlayInfo{vs.currentPlay.PlayInfo}
	for i, p := range vs.playQueue.queue {
		if i+1 == limit {
			break
		}
		info = append(info, p.PlayInfo)
	}
	return info, len(vs.playQueue.queue)
}

// This will run as long as the VoiceConnection is running
func (vs *VoiceSender) Start(wsvc *WellensittichVoiceConnection) {
	fmt.Printf("Starting VoiceSender for %v\n", wsvc.GuildID)
	for {
		select {
		case <-vs.done:
			return
		default:
			if rp, ok := vs.playQueue.dequeue(); ok {
				vs.currentPlay = rp
				err := rp.Sound.Play(wsvc.OpusSend, vs.done, vs.stop, vs.paused, vs.resume)
				if err != nil {
					fmt.Println(err)
				}
				vs.currentPlay = nil
			} else {
				time.Sleep(time.Millisecond * 100) // Sleep briefly to prevent busy looping
			}
		}
	}
}

func (vs *VoiceSender) EnqueuePlay(p *Play) error {
	fmt.Printf("Enqueued a Play\n")
	if !vs.playQueue.enqueue(p) {
		return fmt.Errorf("could not add to play queue")
	}
	return nil
}
