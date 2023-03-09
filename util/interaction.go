package util

import "github.com/bwmarrin/discordgo"

type InteractionContext struct {
	Session     *discordgo.Session
	Interaction *discordgo.InteractionCreate
}

func (ic *InteractionContext) DefaulInteractionAnswer(message string) error {
	return ic.Session.InteractionRespond(ic.Interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: message,
		},
	})
}
