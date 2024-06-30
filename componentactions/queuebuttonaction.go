package componentactions

import (
	"fmt"
	"strconv"

	"github.com/bwmarrin/discordgo"
	"github.com/thewalpa/wellensittich/discordws"
	"github.com/thewalpa/wellensittich/util"
)

func (cID customID) queueButtonActionHandler(s *discordws.WellensittichSession, i *discordgo.InteractionCreate) {
	ic := util.InteractionContext{Session: s.Session, Interaction: i}

	g, err := s.State.Guild(i.GuildID)
	if err != nil {
		fmt.Printf("queueButtonAction: Could not find guild. WHY: %v\n", err)
		return
	}

	voiceID, ok := util.VoiceChannel(g, i.Member.User.ID)
	if !ok {
		err = ic.DefaulInteractionAnswer("You are not in a voice channel.")
		if err != nil {
			fmt.Printf("queueButtonAction: Error responding to interaction: %v\n", err)
		}
		return
	}

	vc, ok := s.WsVoiceConnections[g.ID]
	if vc == nil || ok && vc.ChannelID != voiceID {
		err = ic.DefaulInteractionAnswer("You are not in the same voice channel as the bot.")
		if err != nil {
			fmt.Printf("queueButtonAction: Error responding to interaction: %v\n", err)
		}
		return
	}

	playStr := cID[len(cID)-1]
	playNum, err := strconv.Atoi(string(playStr))
	if err != nil {
		err = ic.DefaulInteractionAnswer("Unexpected error.")
		if err != nil {
			fmt.Printf("queueButtonAction: Error responding to interaction: %v\n", err)
		}
		return
	}

	err = vc.VoiceSender.SkipTo(playNum)
	message := fmt.Sprintf("Successfully skipped to the %d. play.", playNum+1)
	// we can't skip there, so display current queue and show new buttons
	if err != nil {
		message = "The queue is not long enough anymore. Try again!"
	}

	err = ic.DefaulInteractionAnswer(message)
	if err != nil {
		fmt.Println("queueButtonAction:", err)
	}
}
