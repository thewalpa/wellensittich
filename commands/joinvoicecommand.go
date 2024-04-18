package commands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/thewalpa/wellensittich/discordws"
	"github.com/thewalpa/wellensittich/util"
)

func joinVoiceCommandHandler(s *discordws.WellensittichSession, i *discordgo.InteractionCreate) {
	ic := util.InteractionContext{Session: s.Session, Interaction: i}

	g, err := s.State.Guild(i.GuildID)
	if err != nil {
		fmt.Printf("Could not find guild. WHY: %v\n", err)
		return
	}

	voiceID, ok := util.VoiceChannel(g, i.Member.User.ID)
	if !ok {
		err = ic.DefaulInteractionAnswer("You are not in a voice channel.")
		if err != nil {
			fmt.Printf("Error responding to interaction: %v\n", err)
		}
		return
	}

	vc, ok := s.WsVoiceConnections[g.ID]
	reconnecting := false
	if ok {
		reconnecting = vc.IsReconnecting()
	}

	if ok && !reconnecting && vc.ChannelID != voiceID {
		err = ic.DefaulInteractionAnswer("The bot is already in a different voice channel.")
		if err != nil {
			fmt.Printf("Error responding to interaction: %v\n", err)
		}
		return
	}

	if ok && !reconnecting && vc.ChannelID == voiceID {
		err = ic.DefaulInteractionAnswer("The bot is already in your voice channel.")
		if err != nil {
			fmt.Printf("Error responding to interaction: %v\n", err)
		}
		return
	}

	if reconnecting || vc == nil {
		_, err = s.ChannelVoiceJoin(i.GuildID, voiceID, false, false)
		if err != nil {
			fmt.Printf("Error joining voice channel: %v\n", err)
			err = ic.DefaulInteractionAnswer("Could not join voice channel in time.")
			if err != nil {
				fmt.Printf("Error responding to interaction: %v\n", err)
			}
			return
		}
	}

	err = ic.DefaulInteractionAnswer("Successfull!")
	if err != nil {
		fmt.Printf("Error responding to interaction: %v\n", err)
	}
}
