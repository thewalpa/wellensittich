package discordws

import (
	"fmt"
	"sync"
)

type SoundPlayer interface {
	// needs to return on closing the first channel, optionally on the second
	Play(sendChan chan []byte, done chan struct{}, stop chan struct{}, pause chan bool, resume chan bool) error
}

type PlayInfo struct {
	Name   string //Name to display for play in queue
	Length uint32 //Length of play in seconds
}

type Play struct {
	Sound SoundPlayer
	PlayInfo
}

func NewPlay(name string, sound SoundPlayer, length uint32) *Play {
	return &Play{
		Sound: sound,
		PlayInfo: PlayInfo{
			Name:   name,
			Length: length,
		},
	}
}

type PlayQueue struct {
	mu    sync.Mutex
	queue []*Play
}

func (pq *PlayQueue) enqueue(p *Play) bool {
	pq.mu.Lock()
	defer pq.mu.Unlock()
	pq.queue = append(pq.queue, p)
	return true
}

func (pq *PlayQueue) discardTo(i int) bool {
	pq.mu.Lock()
	defer pq.mu.Unlock()
	if len(pq.queue) <= i {
		return false
	}
	pq.queue = pq.queue[i:]
	return true
}

func (pq *PlayQueue) dequeue() (*Play, bool) {
	pq.mu.Lock()
	defer pq.mu.Unlock()
	if len(pq.queue) == 0 {
		return nil, false
	}
	play := pq.queue[0]
	pq.queue = pq.queue[1:]
	return play, true
}

func (pq *PlayQueue) unqueue(i int) bool {
	pq.mu.Lock()
	defer pq.mu.Unlock()
	if i < 1 || len(pq.queue) < i {
		return false
	}
	pq.queue = append(pq.queue[:i], pq.queue[i+1:]...)
	return true
}

func (pq *PlayQueue) reset() {
	pq.mu.Lock()
	defer pq.mu.Unlock()
	pq.queue = []*Play{}
}

type VoiceSender struct {
	playQueue     *PlayQueue
	currentPlay   *Play
	done          chan struct{}
	stop          chan struct{}
	newPlay       chan struct{}
	paused        chan bool // Channel to manage pause state
	resume        chan bool // Channel to manage resume state
	currentPlayMu sync.Mutex
}

func NewVoiceSender() *VoiceSender {
	return &VoiceSender{
		playQueue: &PlayQueue{},
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
	vs.playQueue.reset()
	close(vs.stop)
	vs.stop = make(chan struct{})
}

// skips the current play
func (vs *VoiceSender) SkipPlay() {
	vs.currentPlay = nil //prevent race condition
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
			} else {
				vs.currentPlay = nil
				<-vs.newPlay // Wait for new item if the queue is empty
			}
		}
	}
}

func (vs *VoiceSender) EnqueuePlay(p *Play) error {
	if !vs.playQueue.enqueue(p) {
		return fmt.Errorf("could not add to play queue")
	}
	select {
	case vs.newPlay <- struct{}{}:
	default:
	}
	fmt.Printf("Enqueued a Play\n")
	return nil
}
