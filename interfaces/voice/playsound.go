package interfaces

import (
	"fmt"

	"github.com/thewalpa/wellensittich/assets"
)

type DCAFilePlayer struct {
	path string
}

func NewDCAFilePlayer(path string) *DCAFilePlayer {
	return &DCAFilePlayer{
		path: path,
	}
}

func (p *DCAFilePlayer) Play(sendChan chan []byte, done chan struct{}, stop chan struct{}, pause chan bool, resume chan bool) error {
	buff, err := assets.LoadSoundFile(p.path)
	if err != nil {
		return err
	}
	fmt.Println("Loaded File!")
	for _, buff := range buff {
		select {
		case <-done:
			return nil
		default:
			sendChan <- buff
		}
	}
	return nil
}
