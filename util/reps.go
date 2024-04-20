package util

import "fmt"

type PlayInfo struct {
	Name   string //Name to display for play in queue
	Length uint32 //Length of play in seconds
}

// String representation of a PlayInfo
func (p PlayInfo) String() string {
	return fmt.Sprintf("%s - %s", p.Name, FormatSeconds(p.Length))
}
