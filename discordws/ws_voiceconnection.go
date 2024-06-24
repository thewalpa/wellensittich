package discordws

import (
	"fmt"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
)

type WellensittichVoiceConnection struct {
	wsSession *WellensittichSession
	GuildID   string
	*discordgo.VoiceConnection
	// Handles voice receiving for the VoiceConnection
	VoiceReceiver *VoiceReceiver
	// Handles voice sending for the VoiceConnection
	VoiceSender  *VoiceSender
	connected    bool
	reconnecting bool
	mu           sync.Mutex
}

func NewWellensittichVoiceConnection(gID string, wss *WellensittichSession, vc *discordgo.VoiceConnection, wspq *PlayQueueModel) *WellensittichVoiceConnection {
	return &WellensittichVoiceConnection{
		VoiceConnection: vc,
		VoiceSender:     NewVoiceSender(gID, wspq),
		VoiceReceiver:   NewVoiceReceiver(),
		wsSession:       wss,
		GuildID:         gID,
		connected:       true,
		reconnecting:    false,
	}
}

func (wsvc *WellensittichVoiceConnection) IsReconnecting() bool {
	wsvc.mu.Lock()
	defer wsvc.mu.Unlock()
	return wsvc.reconnecting
}

// wrapper around discordgo.VoiceConnection.Disconnect()
func (wsvc *WellensittichVoiceConnection) Disconnect() error {
	wsvc.mu.Lock()
	if !wsvc.connected && !wsvc.reconnecting {
		fmt.Println("already disconnected")
		return nil
	} else {
		wsvc.connected = false
		wsvc.reconnecting = false
	}
	wsvc.mu.Unlock()
	wsvc.VoiceSender.Stop()
	wsvc.VoiceReceiver.Stop()
	wsvc.wsSession.mu.Lock()
	delete(wsvc.wsSession.WsVoiceConnections, wsvc.GuildID)
	wsvc.wsSession.mu.Unlock()
	err := wsvc.VoiceConnection.Disconnect()
	return err
}

func (wsvc *WellensittichVoiceConnection) VoiceStateJoin() {
	wsvc.mu.Lock()
	if wsvc.reconnecting {
		wsvc.connected = true
		wsvc.reconnecting = false
	}
	wsvc.mu.Unlock()
}

func (wsvc *WellensittichVoiceConnection) VoiceStateLeave() {
	wsvc.mu.Lock()
	if !wsvc.connected {
		return
	}
	if wsvc.reconnecting {
		return
	}
	wsvc.connected = false
	wsvc.reconnecting = true
	wsvc.mu.Unlock()

	for i := 0; i < 5; i++ {
		time.Sleep(time.Second)

		wsvc.mu.Lock()
		connected := wsvc.connected
		wsvc.mu.Unlock()

		if connected {
			return
		}
	}

	// we actually left, in case of manual disconnect this should be called immediately - how to find out that this happened???

	wsvc.wsSession.mu.Lock()
	delete(wsvc.wsSession.WsVoiceConnections, wsvc.GuildID)
	wsvc.wsSession.mu.Unlock()

	wsvc.mu.Lock()
	wsvc.reconnecting = false
	wsvc.mu.Unlock()

	wsvc.wsSession.EventGuildFeatureFunction("botVoiceLeave", wsvc.GuildID, "")
	wsvc.VoiceSender.Stop()
	wsvc.VoiceReceiver.Stop()
}

func (wsvc *WellensittichVoiceConnection) VoiceStateMove() {
	// stop listening
	wsvc.wsSession.EventGuildFeatureFunction("botVoiceMove", wsvc.GuildID, "")
	wsvc.VoiceReceiver.Stop()
}
