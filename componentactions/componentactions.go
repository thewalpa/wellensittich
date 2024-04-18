package componentactions

import (
	"github.com/bwmarrin/discordgo"
	"github.com/thewalpa/wellensittich/discordws"
	"github.com/thewalpa/wellensittich/util"
)

type WellensittichComponentAction struct {
	CustomID string
	Handler  func(s *discordws.WellensittichSession, i *discordgo.InteractionCreate)
}

func ComponentActions() []*WellensittichComponentAction {
	componentActions := []*WellensittichComponentAction{}
	componentActions = append(componentActions, queueSkipButtonActions()...)
	return componentActions
}

func queueSkipButtonActions() []*WellensittichComponentAction {
	queueSkipButtonActions := []*WellensittichComponentAction{}
	for _, id := range util.QueueButtonsCustomIDs() {
		cID := customID(id)
		queueSkipButtonActions = append(queueSkipButtonActions, &WellensittichComponentAction{
			CustomID: id,
			Handler:  cID.queueButtonActionHandler,
		})
	}
	return queueSkipButtonActions
}
