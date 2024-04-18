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
	// buffer, err := assets.LoadSoundFile("assets/ts-sounds-dca/greeting.dca")
	// if err != nil {
	// 	fmt.Printf("Error loading the sound file: %v\n", err)
	// 	return
	// }

	g, err := s.State.Guild(i.GuildID)
	if err != nil {
		fmt.Printf("Could not find guild. WHY: %v\n", err)
		return
	}

	voiceID, ok := util.VoiceChannel(g, i.Member.User.ID)
	if !ok {
		err = ic.DefaulInteractionAnswer("You are not in a voice chat.")
		if err != nil {
			fmt.Printf("Error responding to interaction: %v\n", err)
		}
		return
	}

	vc, ok := s.WsVoiceConnections[g.ID]
	if ok && vc.ChannelID != voiceID {
		err = ic.DefaulInteractionAnswer("The bot is already in a different voice channel.")
		if err != nil {
			fmt.Printf("Error responding to interaction: %v\n", err)
		}
		return
	}

	if ok && vc.ChannelID == voiceID {
		err = ic.DefaulInteractionAnswer("The bot is already in your voice channel.")
		if err != nil {
			fmt.Printf("Error responding to interaction: %v\n", err)
		}
		return
	}

	if vc == nil {
		vc, err = s.ChannelVoiceJoin(i.GuildID, voiceID, false, true)
		if err != nil {
			fmt.Printf("Error joining voice channel: %v\n", err)
			return
		}
	}

	soundPlayer := interfaces.NewDCAFilePlayer("assets/ts-sounds-dca/greeting.dca")
	err = vc.VoiceSender.EnqueuePlay(discordws.NewPlay("ts-sounds/greeting", soundPlayer, 0))
	if err != nil {
		answer = err.Error()
	}

	err = ic.DefaulInteractionAnswer(answer)
	if err != nil {
		fmt.Printf("Error responding to interaction: %v\n", err)
	}
}
