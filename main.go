package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/thewalpa/wellensittich/commands"
	"github.com/thewalpa/wellensittich/componentactions"
	"github.com/thewalpa/wellensittich/config"
	"github.com/thewalpa/wellensittich/discordws"
)

func main() {
	// read config file
	configFile, err := os.ReadFile("config.json")
	config := config.WellensittichConfig{}
	if err != nil {
		fmt.Println("Config file (config.json) could not be read. Stopping!")
		return
	}
	err = json.Unmarshal(configFile, &config)
	if err != nil {
		fmt.Println(err)
		return
	}

	// create discord session with token
	var s *discordgo.Session
	s, err = discordgo.New("Bot " + string(config.Token))
	if err != nil {
		fmt.Printf("Discord session could not be created: %v\n", err)
		return
	}

	wss := discordws.NewWellensittichSession(s, config)

	// bot needs to know about message content and voice states
	wss.Identify.Intents |= discordgo.IntentMessageContent
	wss.Identify.Intents |= discordgo.IntentsGuildVoiceStates

	// open and defer closing the session
	err = wss.Open()
	if err != nil {
		fmt.Printf("Opening session produced error: %v\n", err)
		return
	}
	defer wss.Close()

	doRegister := true

	// Collect commands and register them at discord
	commandMap := make(map[string]func(wss *discordws.WellensittichSession, i *discordgo.InteractionCreate), len(commands.Commands))
	for _, v := range commands.Commands {
		if doRegister {
			go wss.ApplicationCommandCreate(wss.State.User.ID, config.DevServer, v.Command)
		}
		commandMap[v.Command.Name] = v.Handler
	}
	// Collect component actions
	componentActions := componentactions.ComponentActions()
	componentActionMap := make(map[string]func(wss *discordws.WellensittichSession, i *discordgo.InteractionCreate), len(componentActions))
	for _, v := range componentActions {
		componentActionMap[v.CustomID] = v.Handler
	}

	// Init WellenSittichSession with the collected commands
	wss.InitSession(commandMap, componentActionMap)

	// end program when keyboard interrupt
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}
