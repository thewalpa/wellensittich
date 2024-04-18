package interfaces

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
)

type VideoInfo struct {
	Title    string `json:"title"`
	Url      string `json:"url"`
	Duration uint32 `json:"duration"`
}

func GetVideoInfo(url string) (VideoInfo, error) {
	cmd := exec.Command("yt-dlp", "--skip-download", "--print-json", "--flat-playlist", "-S", "+size", "--format", "wa*", url)
	out := bytes.Buffer{}
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return VideoInfo{}, err
	}

	resp := VideoInfo{}
	err = json.Unmarshal(out.Bytes(), &resp)
	if err != nil {
		return VideoInfo{}, err
	}

	dlLink := resp.Url
	if dlLink == "" {
		return VideoInfo{}, fmt.Errorf("no link for format found")
	}

	return resp, nil
}
