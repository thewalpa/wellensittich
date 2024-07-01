package discordws

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/thewalpa/wellensittich/util"
)

type PlayQueueView struct {
	ChannelID string
	MessageID string
	session   *WellensittichSession
	ic        *util.InteractionContext
}

func NewPlayQueue(wss *WellensittichSession) *PlayQueueView {
	return &PlayQueueView{
		session: wss,
	}
}

func (pq *PlayQueueView) Update(pqm *PlayQueue) {
	if pq.session == nil || pq.ChannelID == "" || pq.MessageID == "" {
		return
	}
	message, err := pq.session.ChannelMessage(pq.ChannelID, pq.MessageID)
	if err != nil {
		fmt.Println(err)
		return
	}
	messageEdit := discordgo.NewMessageEdit(pq.ChannelID, pq.MessageID)
	message.Content = "Hier könnte Ihre Werbung stehen!"
	messageEdit.SetContent(message.Content)
	queueInfo, queueLen := pqm.GetQueueInfo(11, true)
	if len(queueInfo) != 0 {
		startIdx := pqm.GetStartIdx()
		// get buttIDs from helper function
		buttIDs := util.QueueButtonsCustomIDs()
		// create button labels 1 to n
		buttLabels := make([]string, len(queueInfo)-1)
		for i := range len(buttLabels) {
			buttLabels[i] = strconv.Itoa(i + 1 + startIdx)
		}
		sb := strings.Builder{}
		sb.WriteString("The current queue:\n")
		for i, play := range queueInfo {
			idx := i
			if idx != 0 {
				idx += startIdx
			}
			sb.WriteString(fmt.Sprintf("%d: %s\n", idx, play))
		}
		if queueLen > len(queueInfo) {
			sb.WriteString(fmt.Sprintf("and %d more...", queueLen-len(queueInfo)))
		}
		queueEmbed := pq.queueEmbed((sb.String()))
		messageEdit.SetEmbeds([]*discordgo.MessageEmbed{&queueEmbed})
		messageEdit.Components = util.ActionRowComps(buttLabels, buttIDs[:len(buttLabels)])
		audioControlButtons := discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{
					CustomID: util.QUEUE_BACKWARDS_CID,
					Style:    discordgo.SecondaryButton,
					Emoji:    &discordgo.ComponentEmoji{Name: "◀️"},
				},
				discordgo.Button{
					CustomID: util.QUEUE_FORWARDS_CID,
					Style:    discordgo.SecondaryButton,
					Emoji:    &discordgo.ComponentEmoji{Name: "▶️"},
				},
			},
		}
		*messageEdit.Components = append(*messageEdit.Components, audioControlButtons)
	} else {
		messageEdit.SetEmbeds([]*discordgo.MessageEmbed{})
	}
	err = pq.ic.UpdateMessageComplex(messageEdit)
	if err != nil {
		fmt.Println(err)
	}
}

func (pq *PlayQueueView) queueEmbed(content string) discordgo.MessageEmbed {
	return discordgo.MessageEmbed{
		Type:  discordgo.EmbedTypeRich,
		Title: "Play Queue",
		Fields: []*discordgo.MessageEmbedField{
			{
				Value: content,
			},
		},
	}
}
