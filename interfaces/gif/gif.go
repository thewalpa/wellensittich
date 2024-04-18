package interfaces

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type GifProvider interface {
	SearchGif(string) (string, error)
}

type TenorGif struct {
	ItemURL string `json:"itemurl"`
	URL     string `json:"url"`
}

type TenorSearchResponse struct {
	Results []TenorGif `json:"results"`
}

type TenorProvider struct {
	api_key string
}

func NewTenorProvider(key string) *TenorProvider {
	return &TenorProvider{
		api_key: key,
	}
}

func (tp *TenorProvider) SearchGif(query string) (string, error) {
	url := "https://tenor.googleapis.com/v2/search" + fmt.Sprintf("?q=%v&key=%v&limit=1", url.QueryEscape(query), tp.api_key)
	response, err := http.Get(url)
	if err != nil {
		return "", err
	}
	if response.StatusCode > 400 {
		return "", fmt.Errorf(response.Status)
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	SearchResults := TenorSearchResponse{}
	err = json.Unmarshal(body, &SearchResults)
	if err != nil {
		fmt.Println(string(body))
		return "", err
	}
	if len(SearchResults.Results) == 0 {
		return "", nil
	}
	return SearchResults.Results[0].ItemURL, nil
}
