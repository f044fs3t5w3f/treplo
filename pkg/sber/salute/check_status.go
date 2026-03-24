package salute

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type responseApiCheckStatus struct {
	Status int `json:"status"`
	Result struct {
		ID             string `json:"id"`
		CreatedAt      string `json:"created_at"`
		UpdatedAt      string `json:"updated_at"`
		Status         string `json:"status"`
		ResponseFileID string `json:"response_file_id"`
	} `json:"result"`
}

func (s *SpeachService) CheckStatus(ctx context.Context, saluteTaskId string) (string, string, error) {
	token, err := s.tokenStorage.GetToken()
	if err != nil {
		return "", "", fmt.Errorf("tokenStorage.GetToken: %w", err)
	}
	s.wg.Add(1)
	defer s.wg.Done()

	url := "https://smartspeech.sber.ru/rest/v1/task:get?id=" + saluteTaskId

	client := &http.Client{}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)

	if err != nil {
		return "", "", fmt.Errorf("http.NewRequestWithContext: %w", err)
	}
	req.Header.Add("Accept", "application/octet-stream")
	req.Header.Add("Authorization", "Bearer "+token)

	res, err := client.Do(req)
	if err != nil {
		return "", "", fmt.Errorf("client.Do: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return "", "", fmt.Errorf("io.ReadAll: %w", err)
		}
		return "", "", fmt.Errorf("salute speach request failded: %s, %s", res.Status, string(body))
	}

	var response responseApiCheckStatus
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return "", "", fmt.Errorf("decoder.Decode: %w", err)
	}
	return response.Result.Status, response.Result.ResponseFileID, nil
}
