package commands

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/thewalpa/wellensittich/discordws"
	"github.com/thewalpa/wellensittich/util"
)

func queueCommandHandler(s *discordws.WellensittichSession, i *discordgo.InteractionCreate) {
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

	queueInfo, queueLen := vc.VoiceSender.GetQueueInfo(11)
	if len(queueInfo) == 0 {
		err = ic.DefaulInteractionAnswer("Nothing in the queue.")
		if err != nil {
			fmt.Printf("Error responding to interaction: %v\n", err)
		}
		return
	}
	// get buttIDs from helper function
	buttIDs := util.QueueButtonsCustomIDs()
	// create button labels 1 to n
	buttLabels := make([]string, len(queueInfo)-1)
	for i := range len(buttLabels) {
		buttLabels[i] = strconv.Itoa(i + 1)
	}
	// create interaction answer message
	sb := strings.Builder{}
	sb.WriteString("The current queue:\n")
	for i, play := range queueInfo {
		sb.WriteString(fmt.Sprintf("%d: %s - %d seconds\n", i, play.Name, play.Length))
	}
	if queueLen > len(queueInfo) {
		sb.WriteString(fmt.Sprintf("and %d more...", queueLen-len(queueInfo)))
	}
	// send interaction answer
	err = ic.ButtonInteractionAnswer(sb.String(), buttLabels, buttIDs[:len(buttLabels)])
	if err != nil {
		fmt.Println("queueCommandHandler:", err)
	}
}
