package util

import "fmt"

func FormatSeconds(seconds uint32) string {
	hours := seconds / 3600          // calculate the total hours
	minutes := (seconds % 3600) / 60 // calculate the remaining minutes
	seconds = seconds % 60           // calculate the remaining seconds
	if hours > 0 {
		return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
	} else {
		return fmt.Sprintf("%02d:%02d", minutes, seconds)
	}
}
