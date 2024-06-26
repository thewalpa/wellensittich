package interfaces

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type YoutubeSearchList struct {
	Items []struct {
		ID struct {
			VideoID string `json:"videoId"`
		} `json:"id"`
	} `json:"items"`
}

type YoutubeMusicSearch struct {
	api_key string
}

func NewYoutubeMusicSearch(api_key string) *YoutubeMusicSearch {
	return &YoutubeMusicSearch{
		api_key: api_key,
	}
}

func (yts *YoutubeMusicSearch) SearchPlay(query string, options ...string) (MusicSearchResult, error) {
	// default is 10 for music
	videoCategoryId := "&videoCategoryId=10"
	for _, o := range options {
		if o == "generic" {
			videoCategoryId = ""
		}
	}
	// youtube search API with some args
	url := "https://youtube.googleapis.com/youtube/v3/search" + fmt.Sprintf("?part=snippet%s&type=video&q=%v&key=%v", videoCategoryId, url.QueryEscape(query), yts.api_key)
	response, err := http.Get(url)
	if err != nil {
		return MusicSearchResult{}, err
	}
	if response.StatusCode > 400 {
		return MusicSearchResult{}, fmt.Errorf(response.Status)
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return MusicSearchResult{}, err
	}
	SearchResults := YoutubeSearchList{}
	err = json.Unmarshal(body, &SearchResults)
	//fmt.Println(string(body))
	if err != nil {
		return MusicSearchResult{}, err
	}
	if len(SearchResults.Items) == 0 {
		return MusicSearchResult{}, nil
	}
	return MusicSearchResult{
		URL: "https://www.youtube.com/watch?v=" + SearchResults.Items[0].ID.VideoID,
	}, nil
}
