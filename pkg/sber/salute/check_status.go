package salute

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (s *SpeachService) CheckStatus(saluteTaskId string) (string, error) {
	token, err := s.tokenStorage.GetToken()
	if err != nil {
		return "", fmt.Errorf("tokenStorage.GetToken: %w", err)
	}

	url := "https://smartspeech.sber.ru/rest/v1/task:get?id=" + saluteTaskId

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		return "", fmt.Errorf("http.NewRequest: %w", err)
	}
	req.Header.Add("Accept", "application/octet-stream")
	req.Header.Add("Authorization", "Bearer "+token)

	res, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("client.Do: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return "", fmt.Errorf("io.ReadAll: %w", err)
		}
		return "", fmt.Errorf("salute speach request failded: %s, %s", res.Status, string(body))
	}
	var response responseApiCreateRecognizeTask
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return "", fmt.Errorf("decoder.Decode: %w", err)
	}
	return response.Result.Status, nil
}
