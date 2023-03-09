package commands

import "github.com/bwmarrin/discordgo"

type WellensittichCommand struct {
	Command *discordgo.ApplicationCommand
	Handler func(s *discordgo.Session, i *discordgo.InteractionCreate)
}

var Commands = []*WellensittichCommand{
	testCommand(),
	tsLadyCommand(),
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

func tsLadyCommand() *WellensittichCommand {
	return &WellensittichCommand{
		Command: &discordgo.ApplicationCommand{
			Name:        "ts-lady",
			Description: "Let the TeamSpeak lady join your voice channel.",
		},
		Handler: tsLadyCommandHandler,
	}
}
