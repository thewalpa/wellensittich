package commands

import (
	"fmt"
	"strconv"
	"strings"

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
	queueInfo, queueLen := wspq.GetQueueInfo(11)
	if len(queueInfo) == 0 {
		err = ic.DefaulInteractionAnswer("Nothing in the queue.")
		if err != nil {
			fmt.Printf("Error responding to interaction: %v\n", err)
		}
		return
	}
	// get buttIDs from helper function
	//buttIDs := util.QueueButtonsCustomIDs()
	// create button labels 1 to n
	buttLabels := make([]string, len(queueInfo)-1)
	for i := range len(buttLabels) {
		buttLabels[i] = strconv.Itoa(i + 1)
	}
	// create interaction answer message
	sb := strings.Builder{}
	sb.WriteString("The current queue:\n")
	for i, play := range queueInfo {
		sb.WriteString(fmt.Sprintf("%d: %s\n", i, play))
	}
	if queueLen > len(queueInfo) {
		sb.WriteString(fmt.Sprintf("and %d more...", queueLen-len(queueInfo)))
	}

	queue_embed := wspq.View.Embed((sb.String()))
	err = ic.Session.InteractionRespond(ic.Interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				&queue_embed,
			},
		},
	})
	if err != nil {
		fmt.Println("queueCommandHandler:", err)
		return
	}

	mess, err := ic.GetResponse()
	if err != nil {
		fmt.Println("queueCommandHandler:", err)
		return
	}
	wspq.UpdateMessage(mess.ID, mess.ChannelID)

	// send interaction answer
	// err = ic.ButtonInteractionAnswer(sb.String(), buttLabels, buttIDs[:len(buttLabels)])
	// if err != nil {
	// 	fmt.Println("queueCommandHandler:", err)
	// }
}
