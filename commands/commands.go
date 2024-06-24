package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/thewalpa/wellensittich/discordws"
)

type WellensittichCommand struct {
	Command *discordgo.ApplicationCommand
	Handler func(s *discordws.WellensittichSession, i *discordgo.InteractionCreate)
}

var Commands = []*WellensittichCommand{
	testCommand(),
	tsLadyCommand(),
	joinVoice(),
	leaveVoice(),
	listenCommand(),
	transcribeCommand(),
	playCommand(),
	skipCommand(),
	stopCommand(),
	pauseCommand(),
	resumeCommand(),
	queueCommand(),
	removeCommand(),
	pinMusicCommand(),
}

func pinMusicCommand() *WellensittichCommand {
	return &WellensittichCommand{
		Command: &discordgo.ApplicationCommand{
			Name:        "pin-music",
			Description: "WIP: Pins the music view to this channel.",
		},
		Handler: pinMusicCommandHandler,
	}
}

func removeCommand() *WellensittichCommand {
	minValue := 1.0
	return &WellensittichCommand{
		Command: &discordgo.ApplicationCommand{
			Name:        "remove",
			Description: "WIP: Remove a play with index.",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					MinValue:    &minValue,
					Required:    true,
					Name:        "index",
					Description: "Index of play to remove.",
				},
			},
		},
		Handler: removeCommandHandler,
	}
}

func queueCommand() *WellensittichCommand {
	return &WellensittichCommand{
		Command: &discordgo.ApplicationCommand{
			Name:        "queue",
			Description: "WIP: Shows information about the queue.",
		},
		Handler: queueCommandHandler,
	}
}

func pauseCommand() *WellensittichCommand {
	return &WellensittichCommand{
		Command: &discordgo.ApplicationCommand{
			Name:        "pause",
			Description: "WIP: Pauses the current play.",
		},
		Handler: pauseCommandHandler,
	}
}

func resumeCommand() *WellensittichCommand {
	return &WellensittichCommand{
		Command: &discordgo.ApplicationCommand{
			Name:        "resume",
			Description: "WIP: Resumes the current play.",
		},
		Handler: resumeCommandHandler,
	}
}

func skipCommand() *WellensittichCommand {
	return &WellensittichCommand{
		Command: &discordgo.ApplicationCommand{
			Name:        "skip",
			Description: "WIP: Skips the current play.",
		},
		Handler: skipCommandHandler,
	}
}

func stopCommand() *WellensittichCommand {
	return &WellensittichCommand{
		Command: &discordgo.ApplicationCommand{
			Name:        "stop",
			Description: "WIP: Stop playing (deletes the queue).",
		},
		Handler: stopCommandHandler,
	}
}

func playCommand() *WellensittichCommand {
	return &WellensittichCommand{
		Command: &discordgo.ApplicationCommand{
			Name:        "play",
			Description: "WIP: Enqueue a play.",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "link",
					Description: "Specifies a YouTube (or other platform) video link.",
					Type:        discordgo.ApplicationCommandOptionString,
				},
				{
					Name:        "query",
					Description: "Specifies a YouTube search query.",
					Type:        discordgo.ApplicationCommandOptionString,
				},
			},
		},
		Handler: playCommandHandler,
	}
}

func listenCommand() *WellensittichCommand {
	return &WellensittichCommand{
		Command: &discordgo.ApplicationCommand{
			Name:        "listen",
			Description: "Will enable features to listen to your voice channel.",
		},
		Handler: listenVoiceCommandHandler,
	}
}

func testCommand() *WellensittichCommand {
	return &WellensittichCommand{
		Command: &discordgo.ApplicationCommand{
			Name:        "test",
			Description: "Just a test command",
		},
		Handler: testCommandHandler,
	}
}

func joinVoice() *WellensittichCommand {
	return &WellensittichCommand{
		Command: &discordgo.ApplicationCommand{
			Name:        "join-voice",
			Description: "Join your voice channel.",
		},
		Handler: joinVoiceCommandHandler,
	}
}

func leaveVoice() *WellensittichCommand {
	return &WellensittichCommand{
		Command: &discordgo.ApplicationCommand{
			Name:        "leave-voice",
			Description: "Leave your voice channel.",
		},
		Handler: leaveVoiceCommandHandler,
	}
}

func tsLadyCommand() *WellensittichCommand {
	return &WellensittichCommand{
		Command: &discordgo.ApplicationCommand{
			Name:        "ts-lady",
			Description: "Let the TeamSpeak lady join your voice channel.",
		},
		Handler: tsLadyCommandHandler,
	}
}

func transcribeCommand() *WellensittichCommand {
	return &WellensittichCommand{
		Command: &discordgo.ApplicationCommand{
			Name:        "transcribe",
			Description: "Will toggle the transcribe feature.",
		},
		Handler: transcribeCommandHandler,
	}
}
