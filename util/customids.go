package util

import "strconv"

const QUEUE_BACKWARDS_CID = "queue_backwards"
const QUEUE_FORWARDS_CID = "queue_forwards"
const QUEUE_PAUSE_CID = "queue_pause"
const QUEUE_RESUME_CID = "queue_resume"
const QUEUE_SKIP_CID = "queue_skip"
const QUEUE_STOP_CID = "queue_stop"
const QUEUE_SKIPTO_CID = "queue_skipto"
const QUEUE_SKIPTO_MODAL_CID = "queue_skipto_modal"
const QUEUE_SKIPTO_MODAL_TEXT_CID = "queue_skipto_modal_text"

// code to make sure component custom IDs are always the same

func QueueButtonsCustomIDs() []string {
	n := 10
	customIDs := make([]string, n)
	for i := range n {
		customIDs[i] = "queue_skip_button_" + strconv.Itoa(i)
	}
	return customIDs
}
