package util

import "strconv"

const QUEUE_BACKWARDS_CID = "queue_backwards"
const QUEUE_FORWARDS_CID = "queue_forwards"
const QUEUE_PAUSE_CID = "queue_pause"
const QUEUE_RESUME_CID = "queue_resume"

// code to make sure component custom IDs are always the same

func QueueButtonsCustomIDs() []string {
	n := 10
	customIDs := make([]string, n)
	for i := range n {
		customIDs[i] = "queue_skip_button_" + strconv.Itoa(i)
	}
	return customIDs
}
