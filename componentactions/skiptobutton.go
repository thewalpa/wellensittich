package componentactions

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/thewalpa/wellensittich/discordws"
	"github.com/thewalpa/wellensittich/util"
)

func skipToButtonHandler(s *discordws.WellensittichSession, i *discordgo.InteractionCreate) {
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

	err = s.InteractionRespond(
		i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseModal,
			Data: &discordgo.InteractionResponseData{
				CustomID: util.QUEUE_SKIPTO_MODAL_CID,
				Title:    "Skip to:",
				Components: []discordgo.MessageComponent{
					discordgo.ActionsRow{
						Components: []discordgo.MessageComponent{
							discordgo.TextInput{
								CustomID:    util.QUEUE_SKIPTO_MODAL_TEXT_CID,
								Label:       "Index",
								Style:       discordgo.TextInputShort,
								Placeholder: "1-999",
								Required:    true,
								MaxLength:   3,
								MinLength:   1,
							},
						},
					},
				},
			},
		},
	)
	if err != nil {
		fmt.Println("playCommandHandler:", err)
	}
}
