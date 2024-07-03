package componentactions

import (
	"fmt"
	"strconv"

	"github.com/bwmarrin/discordgo"
	"github.com/thewalpa/wellensittich/discordws"
	"github.com/thewalpa/wellensittich/util"
)

func skipToModalHandler(s *discordws.WellensittichSession, i *discordgo.InteractionCreate) {
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

	data := i.ModalSubmitData()
	// what, loool
	text := data.Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value

	number, err := strconv.Atoi(text)
	if err != nil || number < 1 {
		fmt.Println("skipToModalHandler: failed to parse number", err)
		ic.DefaulInteractionAnswer("Please enter a valid number.")
	}

	err = vc.VoiceSender.SkipTo(number-1, false)
	if err != nil {
		fmt.Println("skipToModalHandler: failed to skip", err)
	}

	ic.NoAnswer()
}
