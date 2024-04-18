package commands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/thewalpa/wellensittich/discordws"
	"github.com/thewalpa/wellensittich/util"
)

func transcribeCommandHandler(s *discordws.WellensittichSession, i *discordgo.InteractionCreate) {
	ic := util.InteractionContext{Session: s.Session, Interaction: i}

	g, err := s.State.Guild(i.GuildID)
	if err != nil {
		fmt.Printf("Could not find guild. WHY: %v\n", err)
		return
	}

	voiceID, _ := util.VoiceChannel(g, i.Member.User.ID)

	vc, ok := s.WsVoiceConnections[g.ID]
	if voiceID != "" && ok && vc.ChannelID != voiceID {
		err = ic.DefaulInteractionAnswer("The bot is in a different voice channel.")
		if err != nil {
			fmt.Printf("Error responding to interaction: %v\n", err)
		}
		return
	}

	answer := "Succesfully deactivated transcribe feature."
	isEnabled, err := s.ToggleGuildFeature("transcribe", i.GuildID, i.ChannelID)
	if err != nil {
		answer = "Error: not implemented."
	}

	if isEnabled {
		answer = "Successfully activated transcribe feature."
	}

	err = ic.DefaulInteractionAnswer(answer)
	if err != nil {
		fmt.Printf("Error responding to interaction: %v\n", err)
	}
}
