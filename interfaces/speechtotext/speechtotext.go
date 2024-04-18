package interfaces

type SpeechToTextProvider interface {
	SpeechToText([]byte) (string, error)
}
