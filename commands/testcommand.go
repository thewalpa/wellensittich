package commands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/thewalpa/wellensittich/discordws"
)

func testCommandHandler(s *discordws.WellensittichSession, i *discordgo.InteractionCreate) {

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			Title:    "test",
			CustomID: "555",
			Content:  "Test succesful!",
		},
	})
	fmt.Println(err)
}
