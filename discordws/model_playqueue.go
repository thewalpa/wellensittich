package discordws

import (
	"fmt"
	"sync"

	"github.com/thewalpa/wellensittich/util"
)

type PlayQueueModel struct {
	GuildID     string
	View        *PlayQueueView
	mu          sync.Mutex
	queue       []*Play
	currentPlay *Play
}

func NewPlayQueueModel(guildID string, wss *WellensittichSession) *PlayQueueModel {
	return &PlayQueueModel{
		GuildID: guildID,
		View:    NewPlayQueue(wss),
	}
}

func (mpq *PlayQueueModel) UpdateMessage(ic *util.InteractionContext) {
	mess, err := ic.GetResponse()
	if err != nil {
		fmt.Println("PlayQueueModel/UpdateMessage:", err)
		return
	}
	mpq.View.MessageID = mess.ID
	mpq.View.ChannelID = mess.ChannelID
	mpq.View.ic = ic
	mpq.updateView()

}

func (mpq *PlayQueueModel) updateView() {
	if mpq.View == nil {
		return
	}
	mpq.View.Update(mpq)
}

func (mpq *PlayQueueModel) GetQueueInfo(limit int) ([]util.PlayInfo, int) {
	mpq.mu.Lock()
	defer mpq.mu.Unlock()
	info := []util.PlayInfo{}
	for i, p := range mpq.queue {
		if i+1 == limit {
			break
		}
		info = append(info, p.PlayInfo)
	}
	size := len(mpq.queue)
	if mpq.currentPlay != nil {
		info = append([]util.PlayInfo{mpq.currentPlay.PlayInfo}, info...)
		size++
	}
	return info, size
}

func (mpq *PlayQueueModel) Enqueue(p *Play) bool {
	defer mpq.updateView()
	mpq.mu.Lock()
	defer mpq.mu.Unlock()
	mpq.queue = append(mpq.queue, p)
	return true
}

func (mpq *PlayQueueModel) DiscardTo(i int) bool {
	mpq.mu.Lock()
	defer mpq.mu.Unlock()
	if len(mpq.queue) <= i {
		return false
	}
	defer func() { go mpq.updateView() }()
	mpq.queue = mpq.queue[i:]
	return true
}

func (mpq *PlayQueueModel) Dequeue() (*Play, bool) {
	mpq.mu.Lock()
	defer mpq.mu.Unlock()
	if len(mpq.queue) == 0 {
		return nil, false
	}
	defer func() { go mpq.updateView() }()
	play := mpq.queue[0]
	mpq.queue = mpq.queue[1:]
	return play, true
}

func (mpq *PlayQueueModel) Unqueue(i int) bool {
	mpq.mu.Lock()
	defer mpq.mu.Unlock()
	if i < 1 || len(mpq.queue) < i {
		return false
	}
	defer func() { go mpq.updateView() }()
	mpq.queue = append(mpq.queue[:i], mpq.queue[i+1:]...)
	return true
}

func (mpq *PlayQueueModel) Reset() {
	defer mpq.updateView()
	mpq.mu.Lock()
	defer mpq.mu.Unlock()
	mpq.queue = []*Play{}
}

type Play struct {
	Sound util.SoundPlayer
	util.PlayInfo
}

func NewPlay(name string, sound util.SoundPlayer, length uint32) *Play {
	return &Play{
		Sound: sound,
		PlayInfo: util.PlayInfo{
			Name:   name,
			Length: length,
		},
	}
}
