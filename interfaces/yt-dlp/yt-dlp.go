package interfaces

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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

func StreamPlaylistVideoInfos(url string, videoInfoChan chan VideoInfo, errorChan chan error) {
	defer close(videoInfoChan)
	defer close(errorChan)
	cmd := exec.Command("yt-dlp", "--skip-download", "--print-json", "-S", "+size", "--format", "wa*", url)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		errorChan <- err
		return
	}
	// start command
	if err := cmd.Start(); err != nil {
		return
	}

	done := make(chan bool)
	go func() {
		defer close(done)
		decoder := json.NewDecoder(stdout)
		for {
			resp := VideoInfo{}
			if err := decoder.Decode(&resp); err != nil {
				if err == io.EOF {
					break
				}
				errorChan <- err
				return
			}

			dlLink := resp.Url
			if dlLink == "" {
				errorChan <- fmt.Errorf("no link found for format")
				return
			}

			videoInfoChan <- resp
		}
	}()

	// wait for decoding and sending to channel to finish
	<-done
	// wait for command to finish
	if err := cmd.Wait(); err != nil {
		fmt.Println(err)
	}
}

func GetPlaylistVideoInfos(url string) ([]VideoInfo, error) {
	cmd := exec.Command("yt-dlp", "--skip-download", "--print-json", "-S", "+size", "--format", "wa*", url)
	out := bytes.Buffer{}
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return nil, err
	}

	videoInfos := []VideoInfo{}
	decoder := json.NewDecoder(&out)

	for decoder.More() {
		var resp VideoInfo
		if err := decoder.Decode(&resp); err != nil {
			return nil, err
		}

		dlLink := resp.Url
		if dlLink == "" {
			return nil, fmt.Errorf("no link for format found for one of the videos")
		}

		videoInfos = append(videoInfos, resp)
	}

	return videoInfos, nil
}
