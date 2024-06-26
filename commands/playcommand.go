package commands

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/thewalpa/wellensittich/discordws"
	voice "github.com/thewalpa/wellensittich/interfaces/voice"
	ytDlp "github.com/thewalpa/wellensittich/interfaces/yt-dlp"
	"github.com/thewalpa/wellensittich/util"
)

const (
	LINK            = "link"
	YT_PLAYLIST     = "yt-playlist"
	YT_MUSIC_SEARCH = "yt-music-search"
	YT_SEARCH       = "yt-search"
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

	link := ""
	query := ""
	typed := ""
	queryType := ""
	// get input
	for _, o := range iacd.Options {
		if o.Name == "query" {
			query = o.StringValue()
		}
		if o.Name == "type" {
			typed = o.StringValue()
		}
	}

	// return if empty query
	if !wasPlayRequested(query) {
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

	// get query type
	if typed != "" {
		queryType = typed
	} else {
		isLink := strings.HasPrefix(query, "https://")
		isPlaylist := isLink && strings.Contains(query, "list=") && strings.Contains(query, "youtube.com")
		if isPlaylist {
			queryType = YT_PLAYLIST
		} else if isLink {
			queryType = LINK
		} else {
			queryType = YT_MUSIC_SEARCH
		}
	}

	queryAnswer := ""
	// handle different query types
	switch queryType {
	case YT_PLAYLIST:
		link = query
		err = ic.UpdateAnswer("Starting to enqueue plays for the requested playlist.")
		if err != nil {
			fmt.Println("playCommandHandler:", err)
		}
		videoInfoChan := make(chan ytDlp.VideoInfo)
		errorChan := make(chan error)
		go ytDlp.StreamPlaylistVideoInfos(link, videoInfoChan, errorChan)
		i := 0
		for {
			select {
			case videoInfo, ok := <-videoInfoChan:
				if !ok {
					videoInfoChan = nil
				} else {
					ytp := voice.NewStreamPlayer(videoInfo.Url)
					vc.VoiceSender.EnqueuePlay(discordws.NewPlay(videoInfo.Title, ytp, videoInfo.Duration))
					i++
				}
			case err, ok := <-errorChan:
				if !ok {
					errorChan = nil
				} else {
					fmt.Println("Error:", err)
				}
			}
			if videoInfoChan == nil && errorChan == nil {
				break
			}
		}
		// success
		err = ic.UpdateAnswer(fmt.Sprintf("Successfully enqueued %v plays for the requested playlist, check with /queue or /pin-music.", i))
		if err != nil {
			fmt.Println("playCommandHandler:", err)
		}
		// return early
		return
	case LINK:
		link = query
	case YT_SEARCH:
		searchresult, err := s.YoutubeMusicProvider.SearchPlay(query, "generic")
		if err != nil {
			fmt.Println("playCommandHandler:", err)
			err = ic.UpdateAnswer("Could not get results for search query.")
			if err != nil {
				fmt.Println("playCommandHandler:", err)
			}
			return
		}
		link = searchresult.URL
		queryAnswer = fmt.Sprintf("For your search query: %s", query)
	case YT_MUSIC_SEARCH:
		searchresult, err := s.YoutubeMusicProvider.SearchPlay(query, "music")
		if err != nil {
			fmt.Println("playCommandHandler:", err)
			err = ic.UpdateAnswer("Could not get results for search query.")
			if err != nil {
				fmt.Println("playCommandHandler:", err)
			}
			return
		}
		link = searchresult.URL
		queryAnswer = fmt.Sprintf("For your search query: %s", query)
	}

	// if no playlist do this
	videoInfo, err := ytDlp.GetVideoInfo(link)
	if err != nil {
		fmt.Println("playCommandHandler:", err)
		err := ic.UpdateAnswer("Could not find the requested play.")
		if err != nil {
			fmt.Println("playCommandHandler:", err)
		}
		return
	}

	ytp := voice.NewStreamPlayer(videoInfo.Url)
	vc.VoiceSender.EnqueuePlay(discordws.NewPlay(videoInfo.Title, ytp, videoInfo.Duration))

	// success
	err = ic.UpdateAnswer(queryAnswer + "\nSuccessfully enqueued the requested play: " + videoInfo.Title)
	if err != nil {
		fmt.Println("playCommandHandler:", err)
	}
}

func wasPlayRequested(query string) bool {
	return query != ""
}
