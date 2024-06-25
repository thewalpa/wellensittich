package util

import (
	"github.com/bwmarrin/discordgo"
)

type InteractionContext struct {
	Session     *discordgo.Session
	Interaction *discordgo.InteractionCreate
}

func (ic *InteractionContext) GetResponse() (*discordgo.Message, error) {
	return ic.Session.InteractionResponse(ic.Interaction.Interaction)
}

func (ic *InteractionContext) GeneralInteractionAnswer(data *discordgo.InteractionResponseData) error {
	return ic.Session.InteractionRespond(ic.Interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: data,
	})
}

func (ic *InteractionContext) DefaulInteractionAnswer(message string) error {
	return ic.Session.InteractionRespond(ic.Interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: message,
		},
	})
}

// send this if your real answer to the interaction can potentially time out (3s)
func (ic *InteractionContext) DeferAnswer() error {
	return ic.Session.InteractionRespond(ic.Interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})
}

// if you deferred your answer before use this to update your answer
func (ic *InteractionContext) UpdateAnswer(message string) error {
	_, err := ic.Session.InteractionResponseEdit(ic.Interaction.Interaction, &discordgo.WebhookEdit{
		Content: &message,
	})
	return err
}

func (ic *InteractionContext) UpdateAnswerComplex(messageEdit *discordgo.WebhookEdit) error {
	_, err := ic.Session.InteractionResponseEdit(ic.Interaction.Interaction, messageEdit)
	return err
}

// use this function if you need some simple buttons in your answer
// there can be only 5 buttons in a row, this function takes care of this and adds new rows if needed
func (ic *InteractionContext) ButtonInteractionAnswer(message string, buttonLabels, buttonsIDs []string) error {
	actionRowComps := []discordgo.MessageComponent{}
	buttonComps := []discordgo.MessageComponent{}
	for i, bStr := range buttonsIDs {
		buttonComps = append(buttonComps, discordgo.Button{
			Label:    buttonLabels[i],
			Style:    discordgo.SecondaryButton,
			CustomID: bStr,
		})

		// After every 5 buttons, or on the last iteration, add the buttons to an action row.
		if (i+1)%5 == 0 || i == len(buttonsIDs)-1 {
			actionRow := discordgo.ActionsRow{
				Components: buttonComps,
			}
			actionRowComps = append(actionRowComps, actionRow)
			// Reset buttonComps for the next group of buttons.
			buttonComps = []discordgo.MessageComponent{}
		}
	}
	return ic.Session.InteractionRespond(ic.Interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content:    message,
			Components: actionRowComps,
		},
	})
}
