package discordws

import "sync"

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
