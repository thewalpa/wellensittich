package commands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/thewalpa/wellensittich/discordws"
	voice "github.com/thewalpa/wellensittich/interfaces/voice"
	ytDlp "github.com/thewalpa/wellensittich/interfaces/yt-dlp"
	"github.com/thewalpa/wellensittich/util"
)

func playCommandHandler(s *discordws.WellensittichSession, i *discordgo.InteractionCreate) {
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

	// if it never happens can be removed
	if i.Type != discordgo.InteractionApplicationCommand {
		// this should never happen
		fmt.Println("playCommandHandler: Interaction was not an InteractionApplicationCommand")
		err := ic.DefaulInteractionAnswer("Mystical error happened, don't do that again.")
		if err != nil {
			fmt.Println("playCommandHandler:", err)
		}
		return
	}

	iacd := i.ApplicationCommandData()

	// check input
	ytLink := ""
	for _, o := range iacd.Options {
		if o.Name == "link" {
			ytLink = o.StringValue()
		}
	}
	if !wasPlayRequested(ytLink) {
		err := ic.DefaulInteractionAnswer("No play requested.")
		if err != nil {
			fmt.Println("playCommandHandler:", err)
		}
		return
	}

	err = ic.DeferAnswer()
	if err != nil {
		fmt.Println("playCommandHandler:", err)
		return
	}

	videoInfo, err := ytDlp.GetVideoInfo(ytLink)
	if err != nil {
		fmt.Println("playCommandHandler:", err)
		err := ic.UpdateAnswer("Could not find the requested play.")
		if err != nil {
			fmt.Println("playCommandHandler:", err)
		}
		return
	}

	ytp := voice.NewStreamPlayer(videoInfo.Url)
	vc.VoiceSender.EnqueuePlay(util.NewPlay(videoInfo.Title, ytp, videoInfo.Duration))

	// success
	err = ic.UpdateAnswer("Successfully enqueued the requested play: " + videoInfo.Title)
	if err != nil {
		fmt.Println("playCommandHandler:", err)
	}
}

func wasPlayRequested(link string) bool {
	return link != ""
}
