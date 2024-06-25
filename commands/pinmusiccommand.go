package commands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/thewalpa/wellensittich/discordws"
	"github.com/thewalpa/wellensittich/util"
)

func pinMusicCommandHandler(s *discordws.WellensittichSession, i *discordgo.InteractionCreate) {
	ic := util.InteractionContext{Session: s.Session, Interaction: i}

	g, err := s.State.Guild(i.GuildID)
	if err != nil {
		fmt.Printf("Could not find guild. WHY: %v\n", err)
		return
	}

	wspq := s.GetPlayQueue(g.ID)
	err = ic.DeferAnswer()
	if err != nil {
		fmt.Println("queueCommandHandler:", err)
		return
	}

	wspq.UpdateMessage(&ic)
	// send interaction answer
	// err = ic.ButtonInteractionAnswer(sb.String(), buttLabels, buttIDs[:len(buttLabels)])
	// if err != nil {
	// 	fmt.Println("queueCommandHandler:", err)
	// }
}
