package interfaces

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
)

// docker-hub image with gpu: onerahmet/openai-whisper-asr-webservice:latest-gpu

type WhisperAsrWebserviceProvider struct {
	host string
}

func NewWhisperAsrWebserviceProvider(host string) *WhisperAsrWebserviceProvider {
	return &WhisperAsrWebserviceProvider{
		host: host,
	}
}

func (wp *WhisperAsrWebserviceProvider) SpeechToText(content []byte) (string, error) {
	url := fmt.Sprintf("%v/asr", wp.host)
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)
	part, err := writer.CreateFormFile("audio_file", "out.ogg")
	if err != nil {
		return "", err
	}
	_, err = io.Copy(part, bytes.NewReader(content))
	if err != nil {
		panic(err)
	}
	err = writer.Close()
	if err != nil {
		panic(err)
	}

	response, err := http.Post(url+"?language=de", writer.FormDataContentType(), &requestBody)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	return strings.ReplaceAll(string(body), "\n", " "), err
}
