package salute

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type responseApiDownloadContent []struct {
	Results []RecognizeChunk `json:"results"`

	Eou bool `json:"eou"`

	EmotionsResult struct {
		Positive float64 `json:"positive"`
		Neutral  float64 `json:"neutral"`
		Negative float64 `json:"negative"`
	} `json:"emotions_result"`

	ProcessedAudioStart string `json:"processed_audio_start"`
	ProcessedAudioEnd   string `json:"processed_audio_end"`

	BackendInfo struct {
		ModelName     string `json:"model_name"`
		ModelVersion  string `json:"model_version"`
		ServerVersion string `json:"server_version"`
	} `json:"backend_info"`

	Channel int `json:"channel"`

	SpeakerInfo struct {
		SpeakerID             int     `json:"speaker_id"`
		MainSpeakerConfidence float64 `json:"main_speaker_confidence"`
	} `json:"speaker_info"`

	EouReason string `json:"eou_reason"`
	Insight   string `json:"insight"`

	PersonIdentity struct {
		Age         string  `json:"age"`
		Gender      string  `json:"gender"`
		AgeScore    float64 `json:"age_score"`
		GenderScore float64 `json:"gender_score"`
	} `json:"person_identity"`
}

func (s *SpeachService) DownloadContent(ctx context.Context, saluteFileId string) (string, error) {
	token, err := s.tokenStorage.GetToken()
	if err != nil {
		return "", fmt.Errorf("tokenStorage.GetToken: %w", err)
	}

	url := "https://smartspeech.sber.ru/rest/v1/data:download?response_file_id=" + saluteFileId

	client := &http.Client{}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)

	if err != nil {
		return "", fmt.Errorf("http.NewRequestWithContext: %w", err)
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
	var response responseApiDownloadContent
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return "", fmt.Errorf("decoder.Decode: %w", err)
	}

	partsCount := 0
	for _, result := range response {
		partsCount += len(result.Results)
	}

	parts := make([]string, 0, partsCount)
	for _, result := range response {
		for _, chunk := range result.Results {
			parts = append(parts, chunk.NormalizedText)
		}
	}
	return strings.Join(parts, "\n"), nil
}
