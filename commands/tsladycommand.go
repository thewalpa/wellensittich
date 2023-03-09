package commands

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/thewalpa/wellensittich/assets"
	"github.com/thewalpa/wellensittich/util"
)

func tsLadyCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	ic := util.InteractionContext{Session: s, Interaction: i}

	buffer, err := assets.LoadSoundFile("assets/ts-sounds-dca/greeting.dca")
	if err != nil {
		fmt.Printf("Error loading the sound file: %v\n", err)
		return
	}

	g, err := s.State.Guild(i.GuildID)
	if err != nil {
		fmt.Printf("Could not find guild. WHY EVER: %v\n", err)
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

	vc, err := s.ChannelVoiceJoin(i.GuildID, voiceID, false, true)
	if err != nil {
		fmt.Printf("Error joining voice channel: %v\n", err)
		return
	}

	err = ic.DefaulInteractionAnswer("Welcome to TeamSpeak!")
	if err != nil {
		fmt.Printf("Error responding to interaction: %v\n", err)
	}

	time.Sleep(time.Millisecond * 100)
	err = util.PlaySound(vc, buffer)
	if err != nil {
		fmt.Printf("Error playing sound: %v\n", err)
		return
	}
	time.Sleep(time.Millisecond * 100)

	err = vc.Disconnect()
	if err != nil {
		fmt.Printf("Error disconnecting: %v\n", err)
		return
	}
}
