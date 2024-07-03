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
	componentActions = append(componentActions, queueControlButtonActions()...)
	return componentActions
}

type WellensittichModalSubmit struct {
	CustomID string
	Handler  func(s *discordws.WellensittichSession, i *discordgo.InteractionCreate)
}

func ModalSubmitActions() []*WellensittichModalSubmit {
	return []*WellensittichModalSubmit{
		{
			CustomID: util.QUEUE_SKIPTO_MODAL_CID,
			Handler:  skipToModalHandler,
		},
	}
}

func queueControlButtonActions() []*WellensittichComponentAction {
	return []*WellensittichComponentAction{
		{
			CustomID: util.QUEUE_PAUSE_CID,
			Handler:  pauseButtonHandler,
		},
		{
			CustomID: util.QUEUE_RESUME_CID,
			Handler:  resumeButtonHandler,
		},
		{
			CustomID: util.QUEUE_SKIP_CID,
			Handler:  skipButtonHandler,
		},
		{
			CustomID: util.QUEUE_SKIPTO_CID,
			Handler:  skipToButtonHandler,
		},
		{
			CustomID: util.QUEUE_STOP_CID,
			Handler:  stopButtonHandler,
		},
	}
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
