package componentactions

import (
	"github.com/bwmarrin/discordgo"
	"github.com/thewalpa/wellensittich/discordws"
	"github.com/thewalpa/wellensittich/util"
)

type customID string

type WellensittichComponentAction struct {
	CustomID string
	Handler  func(s *discordws.WellensittichSession, i *discordgo.InteractionCreate)
}

func ComponentActions() []*WellensittichComponentAction {
	componentActions := []*WellensittichComponentAction{}
	componentActions = append(componentActions, queueSkipButtonActions()...)
	componentActions = append(componentActions, moveQueueButtonActions()...)
	return componentActions
}

func moveQueueButtonActions() []*WellensittichComponentAction {
	return []*WellensittichComponentAction{
		{
			CustomID: util.QUEUE_BACKWARDS_CID,
			Handler:  customID(util.QUEUE_BACKWARDS_CID).moveQueueViewBackwardsButtonHandler,
		},
		{
			CustomID: util.QUEUE_FORWARDS_CID,
			Handler:  customID(util.QUEUE_FORWARDS_CID).moveQueueViewForwardsButtonHandler,
		},
	}
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
