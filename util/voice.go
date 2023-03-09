package util

import (
	"github.com/bwmarrin/discordgo"
)

func VoiceChannel(g *discordgo.Guild, authorID string) (string, bool) {
	for _, vs := range g.VoiceStates {
		if vs.UserID == authorID {
			return vs.ChannelID, true
		}
	}
	return "", false
}

func PlaySound(vc *discordgo.VoiceConnection, buffer [][]byte) error {
	vc.Speaking(true)

	for _, buff := range buffer {
		vc.OpusSend <- buff
	}

	vc.Speaking(false)

	return nil
}
