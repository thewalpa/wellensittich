package discordws

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"os"
	"sync"

	"github.com/bwmarrin/discordgo"
	"github.com/pion/webrtc/v3/pkg/media/oggwriter"
	"github.com/thewalpa/wellensittich/util"
)

type bufferEvent struct {
	SSRC   uint32
	Buffer *bytes.Buffer
}

type TranscribeFeature struct {
	guildContexts         map[string]chan struct{}
	guildResponseChannels map[string]string
	session               *WellensittichSession
	eventHandlers         map[string]func(string, string)
	//deactivationHandlers  map[string]func(string, string)
}

func NewTranscribeFeature(session *WellensittichSession) *TranscribeFeature {
	tf := &TranscribeFeature{
		guildContexts:         make(map[string]chan struct{}),
		guildResponseChannels: make(map[string]string),
		session:               session,
		eventHandlers:         make(map[string]func(string, string)),
	}
	tf.init()
	return tf
}

func (tf *TranscribeFeature) EventHandlers(source, gID, channelID string) {
	if handler, ok := tf.eventHandlers[source]; ok {
		handler(gID, channelID)
	}
}

func (tf *TranscribeFeature) init() {
	// Activate work when toggling and listeningCommand
	tf.eventHandlers["toggleOn"] = func(guildID, channelID string) {
		tf.guildResponseChannels[guildID] = channelID
		tf.activate(guildID, channelID)
	}
	tf.eventHandlers["listenCommandOn"] = tf.activate
	// Deactivate work when toggling and listeningCommand
	tf.eventHandlers["toggleOff"] = tf.deactivate
	tf.eventHandlers["listenCommandOff"] = tf.deactivate

	// VoiceStateUpdate
	tf.eventHandlers["botVoiceLeave"] = tf.deactivate
	tf.eventHandlers["botVoiceMove"] = tf.deactivate
}

func (tf *TranscribeFeature) activate(guildID, channelID string) {
	// if a context exists, was already activated
	if done, ok := tf.guildContexts[guildID]; ok && done != nil {
		return
	}
	// check for voice connection, we can also only activate if the listening command was activated
	wsvc, exists := tf.session.WsVoiceConnections[guildID]
	if exists && wsvc.VoiceReceiver.IsListenPassive() {
		done := make(chan struct{})
		tf.guildContexts[guildID] = done
		go tf.transcribeLoop(done, wsvc)
	}
}

func (tf *TranscribeFeature) deactivate(guildID, channelID string) {
	// if there is a context, we deactivate, good faith that we do not do that without a reason
	if done, ok := tf.guildContexts[guildID]; ok && done != nil {
		close(done)
		delete(tf.guildContexts, guildID)
		delete(tf.guildResponseChannels, guildID)
	}
}

func (tf *TranscribeFeature) transcribeLoop(done chan struct{}, wsvc *WellensittichVoiceConnection) {
	fmt.Println("Starting transcribe loop.")
	defer fmt.Println("Ending transcribe loop.")

	oggReceiver := NewSyncedOGGFileSoundReceiver()

	wsvc.VoiceReceiver.Subscribe(TRANSCRIBE_FEATURE_NAME, oggReceiver)
	defer wsvc.VoiceReceiver.Unsubscribe(TRANSCRIBE_FEATURE_NAME)

	for {
		select {
		case bufferEvent := <-oggReceiver.EventChan:
			if bufferEvent.Buffer == nil {
				break
			}
			buffer := bufferEvent.Buffer
			content := buffer.Bytes()

			// debugging
			randomFileName := fmt.Sprintf("tmp/%d.ogg", rand.Int())
			go os.WriteFile(randomFileName, content, 0664)

			content, err := io.ReadAll(buffer)
			if err != nil {
				break
			}
			text, err := tf.session.SpeechToTextProvider.SpeechToText(content)
			if err != nil {
				fmt.Println("Error SpeechToTextProvider", err)
				break
			}
			fmt.Println(text, ":", randomFileName)
			if text == "" {
				break
			}
			gifLink, err := tf.session.GifProvider.SearchGif(text)
			if err != nil {
				fmt.Println("Error GifProvider", err)
				break
			}
			if gifLink == "" {
				break
			}
			_, err = tf.session.ChannelMessageSend(tf.guildResponseChannels[wsvc.GuildID], gifLink)
			if err != nil {
				fmt.Println("Error sending message to channel", tf.guildResponseChannels[wsvc.GuildID], err)
			}
		case <-done:
			return
		}
	}
}

// func (tf *TranscribeFeature) voiceStateUpdateEventHandler(s *discordgo.Session, vsu *discordgo.VoiceStateUpdate) {
// 	if vsu.UserID != s.State.User.ID || vsu.BeforeUpdate == nil {
// 		return
// 	}

// 	// if bot leaves or moves the feature has to be deactivated
// 	if vsu.ChannelID == "" && vsu.BeforeUpdate.ChannelID != "" {
// 		tf.deactivate(vsu.GuildID, "")
// 	} else if vsu.ChannelID != "" && vsu.BeforeUpdate.ChannelID != "" && vsu.ChannelID != vsu.BeforeUpdate.ChannelID {
// 		tf.deactivate(vsu.GuildID, "")
// 		return
// 	}
// }

type SyncedOGGFileSoundReceiver struct {
	OGGBuffers            map[uint32]*bytes.Buffer
	bufferStartTimeStamps map[uint32]uint32
	oggFiles              map[uint32]*oggwriter.OggWriter
	EventChan             chan bufferEvent
	mu                    sync.Mutex
}

func NewSyncedOGGFileSoundReceiver() *SyncedOGGFileSoundReceiver {
	return &SyncedOGGFileSoundReceiver{
		OGGBuffers:            make(map[uint32]*bytes.Buffer),
		bufferStartTimeStamps: make(map[uint32]uint32),
		oggFiles:              make(map[uint32]*oggwriter.OggWriter),
		EventChan:             make(chan bufferEvent, 1),
	}
}

func (osr *SyncedOGGFileSoundReceiver) Before() {
	// go func() {
	// 	for {
	// 		time.Sleep(4 * time.Second)
	// 		fmt.Println(osr.oggFiles)
	// 	}
	// }()
}

func (osr *SyncedOGGFileSoundReceiver) After() {
	osr.mu.Lock()
	defer osr.mu.Unlock()
	for _, file := range osr.oggFiles {
		err := file.Close()
		if err != nil {
			fmt.Println(err)
		}
	}
}

func (osr *SyncedOGGFileSoundReceiver) HandleOpusPacket(p discordgo.Packet, done chan struct{}) {
	osr.mu.Lock()
	defer osr.mu.Unlock()
	file, ok := osr.oggFiles[p.SSRC]
	if !ok {
		newBuffer := &bytes.Buffer{}
		newFile, err := oggwriter.NewWith(newBuffer, 48000, 2)
		if err != nil {
			fmt.Println(err)
			return
		}
		osr.OGGBuffers[p.SSRC] = newBuffer
		osr.oggFiles[p.SSRC] = newFile
		osr.bufferStartTimeStamps[p.SSRC] = p.Timestamp
		file = newFile
	}

	rtp := util.CreatePionRTPPacket(&p)
	err := file.WriteRTP(rtp)
	if err != nil {
		fmt.Println(err)
	}

	// in seconds
	dur := float32(p.Timestamp-osr.bufferStartTimeStamps[p.SSRC]) / 80000

	// some kind of an end signal for voice
	if len(p.Opus) == 3 && p.Opus[0] == 248 && p.Opus[1] == 255 && p.Opus[2] == 254 {
		if dur > 0.5 {
			bufferEvent := bufferEvent{
				SSRC:   p.SSRC,
				Buffer: osr.OGGBuffers[p.SSRC],
			}
			osr.EventChan <- bufferEvent
		}
		delete(osr.oggFiles, p.SSRC)
	}
}
