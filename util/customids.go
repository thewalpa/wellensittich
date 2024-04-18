package util

import "strconv"

// code to make sure component custom IDs are always the same

func QueueButtonsCustomIDs() []string {
	n := 10
	customIDs := make([]string, n)
	for i := range n {
		customIDs[i] = "queue_skip_button_" + strconv.Itoa(i)
	}
	return customIDs
}
