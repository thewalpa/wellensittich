package discordws

import (
	"sync"

	"github.com/bwmarrin/discordgo"
	"github.com/thewalpa/wellensittich/config"
	gif "github.com/thewalpa/wellensittich/interfaces/gif"
	musicsearch "github.com/thewalpa/wellensittich/interfaces/musicsearch"
	speechtotext "github.com/thewalpa/wellensittich/interfaces/speechtotext"
)

type WSGuildFeature interface {
	EventHandlers(string, string, string)
	//DeactivationHandlers(string, string, string)
}

const (
	TRANSCRIBE_FEATURE_NAME = "transcribe"
	TSLADY_FEATURE_NAME     = "tsLady"
	PLAYQUEUE_FEATURE_NAME  = "playQueue"
)

type guildFeatureConfig struct {
	enabled map[string]bool
}

type WellensittichSession struct {
	mu sync.Mutex
	*discordgo.Session

	// Use instead of discordgo.Session.VoiceConnections because we need VoiceReceiver and VoiceSender
	WsVoiceConnections map[string]*WellensittichVoiceConnection

	WsPlayQueues map[string]*PlayQueue

	// commandMap
	commandMap         map[string]func(wss *WellensittichSession, i *discordgo.InteractionCreate)
	componentActionMap map[string]func(wss *WellensittichSession, i *discordgo.InteractionCreate)

	guildFeatureConfigs map[string]*guildFeatureConfig
	guildFeatures       map[string]WSGuildFeature

	// interfaces
	SpeechToTextProvider speechtotext.SpeechToTextProvider
	GifProvider          gif.GifProvider
	YoutubeMusicProvider musicsearch.MusicSearchProvider
}

func NewWellensittichSession(s *discordgo.Session, wsc config.WellensittichConfig) *WellensittichSession {
	wss := &WellensittichSession{
		Session:              s,
		WsVoiceConnections:   make(map[string]*WellensittichVoiceConnection),
		WsPlayQueues:         make(map[string]*PlayQueue),
		guildFeatureConfigs:  make(map[string]*guildFeatureConfig),
		guildFeatures:        make(map[string]WSGuildFeature),
		SpeechToTextProvider: speechtotext.NewWhisperAsrWebserviceProvider(wsc.WhisperHost),
		GifProvider:          gif.NewTenorProvider(wsc.TenorKey),
		YoutubeMusicProvider: musicsearch.NewYoutubeMusicSearch(wsc.TenorKey),
	}
	// register features
	wss.guildFeatures[TRANSCRIBE_FEATURE_NAME] = NewTranscribeFeature(wss)

	// Debug func to print the state of the session every 5 seconds
	// go func() {
	// 	for {
	// 		time.Sleep(5 * time.Second)
	// 		fmt.Printf("WSVCs: %v\n", wss.WsVoiceConnections)
	// 		fmt.Printf("DSVCs: %v\n", wss.Session.VoiceConnections)
	// 	}
	// }()
	return wss
}

// Init the session handlers, needs information from main
func (wss *WellensittichSession) InitSession(
	commandMap map[string]func(wss *WellensittichSession, i *discordgo.InteractionCreate),
	componentActionMap map[string]func(wss *WellensittichSession, i *discordgo.InteractionCreate)) {

	wss.commandMap = commandMap
	wss.componentActionMap = componentActionMap
	wss.AddHandler(wss.interactionCreateEventHandler)
	wss.AddHandler(wss.voiceStateUpdateEventHandler)
}

func (wss *WellensittichSession) ChannelVoiceJoin(gID, cID string, mute, deaf bool) (*WellensittichVoiceConnection, error) {
	vc, err := wss.Session.ChannelVoiceJoin(gID, cID, mute, deaf)
	if err != nil {
		return nil, err
	}
	wss.mu.Lock()
	defer wss.mu.Unlock()
	var wspq *PlayQueue
	if wspqVal, ok := wss.WsPlayQueues[gID]; !ok {
		wspq = NewPlayQueueModel(gID, wss)
		wss.WsPlayQueues[gID] = wspq
	} else {
		wspq = wspqVal
	}
	if wsvc, ok := wss.WsVoiceConnections[gID]; !ok {
		wsvc = NewWellensittichVoiceConnection(gID, wss, vc, wspq)
		wss.WsVoiceConnections[gID] = wsvc
		// log info
		//wsvc.LogLevel = 2
		go wsvc.VoiceSender.Start(wsvc)
		return wsvc, nil
	} else {
		return wsvc, nil
	}
}

func (wss *WellensittichSession) GetPlayQueue(gID string) *PlayQueue {
	wss.mu.Lock()
	defer wss.mu.Unlock()
	var wspq *PlayQueue
	if wspqVal, ok := wss.WsPlayQueues[gID]; !ok {
		wspq = NewPlayQueueModel(gID, wss)
		wss.WsPlayQueues[gID] = wspq
	} else {
		wspq = wspqVal
	}

	return wspq
}

func (wss *WellensittichSession) ToggleGuildFeature(featureName, gID, channelID string) (enabled bool, err error) {
	wss.mu.Lock()
	defer wss.mu.Unlock()

	gfc, ok := wss.guildFeatureConfigs[gID]
	if !ok {
		gfc = &guildFeatureConfig{
			enabled: make(map[string]bool),
		}
		wss.guildFeatureConfigs[gID] = gfc
	}

	enabled = !gfc.enabled[featureName]
	gfc.enabled[featureName] = enabled
	if enabled {
		wss.guildFeatures[featureName].EventHandlers("toggleOn", gID, channelID)
	} else {
		wss.guildFeatures[featureName].EventHandlers("toggleOff", gID, channelID)
	}
	return
}

func (wss *WellensittichSession) IsFeatureEnabledForGuild(featureName, guildID string) bool {
	wss.mu.Lock()
	defer wss.mu.Unlock()

	gfc, exists := wss.guildFeatureConfigs[guildID]
	if !exists {
		return false
	}

	return gfc.enabled[featureName]
}

func (wss *WellensittichSession) EventGuildFeatureFunction(source, gID, channelID string) {
	for k, v := range wss.guildFeatures {
		if wss.IsFeatureEnabledForGuild(k, gID) {
			v.EventHandlers(source, gID, channelID)
		}
	}
}
