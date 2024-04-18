package discordws

import (
	"github.com/bwmarrin/discordgo"
)

func (wss *WellensittichSession) voiceStateUpdateEventHandler(s *discordgo.Session, vsu *discordgo.VoiceStateUpdate) {
	if vsu.UserID != s.State.User.ID {
		return
	}

	// check WsVoiceConnection stuff

	wsvc, ok := wss.WsVoiceConnections[vsu.GuildID]
	if !ok {
		return
	}

	if vsu.BeforeUpdate == nil {
		vsu.BeforeUpdate = &discordgo.VoiceState{}
	}
	if vsu.ChannelID != "" && vsu.BeforeUpdate.ChannelID == "" {
		// join
		wsvc.VoiceStateJoin()
	} else if vsu.ChannelID == "" && vsu.BeforeUpdate.ChannelID != "" {
		// leave
		wsvc.VoiceStateLeave()
	} else if vsu.ChannelID != "" && vsu.BeforeUpdate.ChannelID != "" && vsu.ChannelID != vsu.BeforeUpdate.ChannelID {
		// move
		wsvc.VoiceStateMove()
	}
}
