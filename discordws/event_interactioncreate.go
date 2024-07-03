package discordws

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func (wss *WellensittichSession) interactionCreateEventHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		if h, ok := wss.commandMap[i.ApplicationCommandData().Name]; ok {
			h(wss, i)
		} else {
			fmt.Println("The command does not exist.", i.ApplicationCommandData().Name)
		}
	case discordgo.InteractionMessageComponent:
		if h, ok := wss.componentActionMap[i.MessageComponentData().CustomID]; ok {
			h(wss, i)
		} else {
			fmt.Println("The component customID does not exist.", i.MessageComponentData().CustomID)
		}
	case discordgo.InteractionModalSubmit:
		if h, ok := wss.modalSubmitMap[i.ModalSubmitData().CustomID]; ok {
			h(wss, i)
		} else {
			fmt.Println("The modalsubmit customID does not exist.", i.ModalSubmitData().CustomID)
		}
	}
}
