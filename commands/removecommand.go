package commands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/thewalpa/wellensittich/discordws"
	"github.com/thewalpa/wellensittich/util"
)

func removeCommandHandler(s *discordws.WellensittichSession, i *discordgo.InteractionCreate) {
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
	if vc == nil || ok && vc.ChannelID != voiceID {
		err = ic.DefaulInteractionAnswer("You are not in the same voice channel as the bot.")
		if err != nil {
			fmt.Printf("Error responding to interaction: %v\n", err)
		}
		return
	}

	iacd := i.ApplicationCommandData()

	// check input
	index := -1
	for _, o := range iacd.Options {
		if o.Name == "index" {
			index = int(o.IntValue())
		}
	}
	if index == -1 {
		err := ic.DefaulInteractionAnswer("Index was not provided.")
		if err != nil {
			fmt.Println("removeCommandHandler:", err)
		}
		return
	}

	err = vc.VoiceSender.RemovePlay(index)
	if err != nil {
		err = ic.DefaulInteractionAnswer("Unexpected error.")
		if err != nil {
			fmt.Printf("removeCommandHandler: Error responding to interaction: %v\n", err)
		}
		return
	}

	// success
	err = ic.DefaulInteractionAnswer("Successfully removed the chosen play from queue.")
	if err != nil {
		fmt.Println("playCommandHandler:", err)
	}
}
