package interfaces

import (
	"io"

	"github.com/thewalpa/dca"
)

type StreamPlayer struct {
	dlLink string
}

func NewStreamPlayer(dlLink string) *StreamPlayer {
	return &StreamPlayer{
		dlLink: dlLink,
	}
}

func (ytp *StreamPlayer) Play(sendChan chan []byte, done chan struct{}, stop chan struct{}, pause chan bool, resume chan bool) error {
	encodingSession, err := dca.EncodeFile(ytp.dlLink, dca.StdEncodeOptions)
	if err != nil {
		return err
	}
	defer encodingSession.Cleanup()

	for {
		select {
		case <-stop:
			return nil
		case <-done:
			return nil
		case <-pause:
			<-resume
		default:
			opus, err := encodingSession.OpusFrame()
			if err != nil {

				if err == io.EOF {
					return nil
				}
				return err
			}
			sendChan <- opus
		}
	}
}
