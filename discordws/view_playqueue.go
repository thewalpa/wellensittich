package discordws

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

type PlayQueueView struct {
	ChannelID string
	MessageID string
	session   *WellensittichSession
	i         int
}

func (pq *PlayQueueView) Update() {
	if pq.session == nil || pq.ChannelID == "" || pq.MessageID == "" {
		return
	}
	message, err := pq.session.ChannelMessage(pq.ChannelID, pq.MessageID)
	if err != nil {
		fmt.Println(err)
		return
	}
	message.Content = fmt.Sprintf("this was edited %d", pq.i)
	pq.i++
	messageEdit := discordgo.NewMessageEdit(pq.ChannelID, pq.MessageID)
	messageEdit.SetContent(message.Content)
	messageEdit.SetEmbeds(message.Embeds)
	_, err = pq.session.ChannelMessageEditComplex(messageEdit)
	if err != nil {
		fmt.Println(err)
	}
}

func (pq *PlayQueueView) Embed(content string) discordgo.MessageEmbed {
	return discordgo.MessageEmbed{
		Type:        discordgo.EmbedTypeRich,
		Title:       "Wellensittich Play Queue",
		Description: "test",
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "test",
				Value: content,
			},
		},
	}
}
