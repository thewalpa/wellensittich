package commands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/thewalpa/wellensittich/discordws"
	interfaces "github.com/thewalpa/wellensittich/interfaces/voice"
	"github.com/thewalpa/wellensittich/util"
)

func tsLadyCommandHandler(s *discordws.WellensittichSession, i *discordgo.InteractionCreate) {
	ic := util.InteractionContext{Session: s.Session, Interaction: i}
	answer := "Welcome to TeamSpeak!"
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
	if vc == nil || ok && vc.ChannelID != voiceID {
		err = ic.DefaulInteractionAnswer("You are not in the same voice channel as the bot.")
		if err != nil {
			fmt.Printf("Error responding to interaction: %v\n", err)
		}
		return
	}

	soundPlayer := interfaces.NewDCAFilePlayer("assets/ts-sounds-dca/greeting.dca")
	err = vc.VoiceSender.PlayImmediately(util.NewPlay("ts-sounds/greeting", soundPlayer, 0), vc)
	if err != nil {
		answer = err.Error()
	}

	err = ic.DefaulInteractionAnswer(answer)
	if err != nil {
		fmt.Printf("Error responding to interaction: %v\n", err)
	}
}
