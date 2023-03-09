package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/thewalpa/wellensittich/commands"
)

func main() {
	// read token file
	token, err := os.ReadFile("token")
	if err != nil {
		fmt.Println("Token file could not be read. Stopping here!")
		return
	}

	// create discord session with token
	var s *discordgo.Session
	s, err = discordgo.New("Bot " + string(token))
	if err != nil {
		fmt.Printf("Discord session could not be created: %v\n", err)
		return
	}

	// bot needs to know about message content
	s.Identify.Intents |= discordgo.IntentMessageContent
	s.Identify.Intents |= discordgo.IntentsGuildVoiceStates

	// open and defer closing the session
	err = s.Open()
	if err != nil {
		fmt.Printf("Opening session produced error: %v\n", err)
		return
	}
	defer s.Close()

	// register commands
	//registeredCommands := make([]*discordgo.ApplicationCommand, len(commands.Commands))
	commandMap := make(map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate), len(commands.Commands))
	for _, v := range commands.Commands {
		_, err := s.ApplicationCommandCreate(s.State.User.ID, "", v.Command)
		if err != nil {
			fmt.Printf("Cannot create %v command: %v\n", v.Command.Name, err)
		}
		commandMap[v.Command.Name] = v.Handler
		//registeredCommands[i] = cmd
	}

	// one handler for all the commands
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandMap[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})

	// end program when keyboard interrupt
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}
