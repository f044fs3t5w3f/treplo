package salute

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type responseApiUploadFile struct {
	Status int `json:"status"`
	Result struct {
		RequestFileID string `json:"request_file_id"`
	} `json:"result"`
}

func (s *SpeechService) UploadFile(ctx context.Context, file io.Reader) (string, error) {
	token, err := s.tokenStorage.GetToken()
	if err != nil {
		return "", fmt.Errorf("tokenStorage.GetToken: %w", err)
	}
	s.wg.Add(1)
	defer s.wg.Done()

	url := "https://smartspeech.sber.ru/rest/v1/data:upload"
	method := http.MethodPost

	client := &http.Client{Timeout: 10 * time.Second}
	// TODO: use client in struct instead of creating the new one
	req, err := http.NewRequestWithContext(ctx, method, url, file)

	if err != nil {
		return "", fmt.Errorf("http.NewRequestWithContext: %w", err)
	}
	req.Header.Add("Content-Type", "audio/mpeg")
	req.Header.Add("Accept", "application/json")
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
		return "", fmt.Errorf("salute speech request failed: %s, %s", res.Status, string(body))
	}
	var response responseApiUploadFile
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return "", fmt.Errorf("decoder.Decode: %w", err)
	}
	return response.Result.RequestFileID, nil
}
