package util

import (
	"github.com/bwmarrin/discordgo"
	"github.com/pion/rtp"
)

func VoiceChannel(g *discordgo.Guild, authorID string) (string, bool) {
	for _, vs := range g.VoiceStates {
		if vs.UserID == authorID {
			return vs.ChannelID, true
		}
	}
	return "", false
}

func CreatePionRTPPacket(p *discordgo.Packet) *rtp.Packet {
	return &rtp.Packet{
		Header: rtp.Header{
			Version:        2,
			PayloadType:    0x78,
			SequenceNumber: p.Sequence,
			Timestamp:      p.Timestamp,
			SSRC:           p.SSRC,
		},
		Payload: p.Opus,
	}
}
