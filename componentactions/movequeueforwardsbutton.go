package componentactions

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/thewalpa/wellensittich/discordws"
	"github.com/thewalpa/wellensittich/util"
)

func (cID customID) moveQueueViewForwardsButtonHandler(s *discordws.WellensittichSession, i *discordgo.InteractionCreate) {
	ic := util.InteractionContext{Session: s.Session, Interaction: i}

	g, err := s.State.Guild(i.GuildID)
	if err != nil {
		fmt.Printf("queueButtonAction: Could not find guild. WHY: %v\n", err)
		return
	}

	voiceID, ok := util.VoiceChannel(g, i.Member.User.ID)
	if !ok {
		err = ic.DefaulInteractionAnswer("You are not in a voice channel.")
		if err != nil {
			fmt.Printf("queueButtonAction: Error responding to interaction: %v\n", err)
		}
		return
	}

	vc, ok := s.WsVoiceConnections[g.ID]
	if vc == nil || ok && vc.ChannelID != voiceID {
		err = ic.DefaulInteractionAnswer("You are not in the same voice channel as the bot.")
		if err != nil {
			fmt.Printf("queueButtonAction: Error responding to interaction: %v\n", err)
		}
		return
	}

	wspq := s.GetPlayQueue(g.ID)
	mok := wspq.MoveQueueViewForwards()
	if !mok {
		err = ic.DefaulInteractionAnswer("Not possible to move forwards.")
		if err != nil {
			fmt.Printf("queueButtonAction: Error responding to interaction: %v\n", err)
		}
		return
	}

	err = ic.NoAnswer()
	if err != nil {
		fmt.Printf("queueButtonAction: Error responding to interaction: %v\n", err)
	}
}
