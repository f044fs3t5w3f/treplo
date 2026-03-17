package salute

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type responseApiUploadFile struct {
	Status int `json:"status"`
	Result struct {
		RequestFileID string `json:"request_file_id"`
	} `json:"result"`
}

func (s *SpeachService) UploadFile(file io.Reader) (string, error) {
	token, err := s.tokenStorage.GetToken()
	if err != nil {
		return "", fmt.Errorf("tokenStorage.GetToken: %w", err)
	}

	url := "https://smartspeech.sber.ru/rest/v1/data:upload"
	method := http.MethodPost

	client := &http.Client{}
	req, err := http.NewRequest(method, url, file)

	if err != nil {
		return "", fmt.Errorf("http.NewRequest: %w", err)
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
		return "", fmt.Errorf("salute speach request failded: %s, %s", res.Status, string(body))
	}
	var response responseApiUploadFile
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return "", fmt.Errorf("decoder.Decode: %w", err)
	}
	return response.Result.RequestFileID, nil
}
