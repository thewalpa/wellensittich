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

func (pq *PlayQueueView) Update(pqm *PlayQueueModel) {
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
	queueInfo, queueLen := pqm.GetQueueInfo(11)
	if len(queueInfo) != 0 {
		// get buttIDs from helper function
		//buttIDs := util.QueueButtonsCustomIDs()
		// create button labels 1 to n
		buttLabels := make([]string, len(queueInfo)-1)
		for i := range len(buttLabels) {
			buttLabels[i] = strconv.Itoa(i + 1)
		}
		sb := strings.Builder{}
		sb.WriteString("The current queue:\n")
		for i, play := range queueInfo {
			sb.WriteString(fmt.Sprintf("%d: %s\n", i, play))
		}
		if queueLen > len(queueInfo) {
			sb.WriteString(fmt.Sprintf("and %d more...", queueLen-len(queueInfo)))
		}
		queueEmbed := pq.queueEmbed((sb.String()))
		messageEdit.SetEmbeds([]*discordgo.MessageEmbed{&queueEmbed})
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
