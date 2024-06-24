package util

type SoundPlayer interface {
	// needs to return on closing the first channel, optionally on the second
	Play(sendChan chan []byte, done chan struct{}, stop chan struct{}, pause chan bool, resume chan bool) error
}
