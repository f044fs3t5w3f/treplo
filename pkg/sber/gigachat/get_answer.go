package gigachat

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type Message struct {
	Role        string   `json:"role"`
	Content     string   `json:"content"`
	Attachments []string `json:"attachments"`
}

type responseApicompletions struct {
	Choices []struct {
		Message Message `json:"message"`
	} `json:"choices"`
}

func (g *GigaChatService) GetAnswer(ctx context.Context, messages []Message) (string, error) {
	url := "https://gigachat.devices.sberbank.ru/api/v1/chat/completions"
	client := &http.Client{}
	token, err := g.tokenStorage.GetToken()
	if err != nil {
		return "", fmt.Errorf("tokenStorage.GetToken: %w", err)
	}

	requestData := struct {
		Model             string `json:"model"`
		Messages          []Message
		Stream            bool `json:"stream"`
		RepetitionPenalty int  `json:"repetition_penalty"`
	}{
		Model: "GigaChat",
		// Model:             "GigaChat-Pro",
		RepetitionPenalty: 1,
		Messages:          messages,
	}

	buf := &bytes.Buffer{}
	encoder := json.NewEncoder(buf)
	err = encoder.Encode(requestData)
	if err != nil {
		return "", fmt.Errorf("encoder.Encode: %w ", err)
	}

	request, err := http.NewRequest(http.MethodPost, url, buf)
	if err != nil {
		return "", fmt.Errorf("http.NewRequest: %w ", err)
	}
	// uuid := generateRqUID()
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Add("Accept", "application/json")
	request.Header.Add("Authorization", "Bearer "+token)
	// request.Header.Add("RqUID", uuid)
	response, err := client.Do(request)

	// TODO: maybe make retriable
	if err != nil {
		return "", fmt.Errorf("client.Do: %w ", err)
	}
	defer response.Body.Close()
	if response.StatusCode > 299 {
		errorContent, _ := io.ReadAll(response.Body)
		return "", fmt.Errorf("error while performing request: %s", errorContent)
	}
	var responseStruct responseApicompletions

	decoder := json.NewDecoder(response.Body)
	err = decoder.Decode(&responseStruct)
	if err != nil {
		return "", fmt.Errorf("error while unmarshalling response: %w", err)
	}
	if len(responseStruct.Choices) == 0 {
		return "", errors.New("empty response from service")
	}
	return responseStruct.Choices[0].Message.Content, nil
}
