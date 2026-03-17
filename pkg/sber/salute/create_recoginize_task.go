package salute

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const VoiceAudioEncoding = "OPUS"

const createRecognizeTaskPayloadTemplate = `{
  "options": {
    "model": "general",
    "audio_encoding": "%s",
    "sample_rate": 16000,
    "language": "ru-RU",
    "enable_profanity_filter": false,
    "hypotheses_count": 1,
    "hints": {
      "words": [ 
      ],
      "enable_letters": false
    },
    "channels_count": 1,
    "speaker_separation_options": {
      "enable": true,
      "enable_only_main_speaker": false,
      "count": 10
    },
    "insight_models": [
    ]
  },
  "request_file_id": "%s"
}`

type responseApiCreateRecognizeTask struct {
	Status int `json:"status"`
	Result struct {
		ID        string `json:"id"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
		Status    string `json:"status"`
	} `json:"result"`
}

func (s *SpeachService) CreateRecognizeTask(saluteFileId string) (string, string, error) {
	token, err := s.tokenStorage.GetToken()
	if err != nil {
		return "", "", fmt.Errorf("tokenStorage.GetToken: %w", err)
	}

	url := "https://smartspeech.sber.ru/rest/v1/speech:async_recognize"
	method := http.MethodPost

	encoding := VoiceAudioEncoding
	payload := strings.NewReader(
		fmt.Sprintf(createRecognizeTaskPayloadTemplate, encoding, saluteFileId),
	)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		return "", "", fmt.Errorf("http.NewRequest: %w", err)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+token)

	res, err := client.Do(req)
	if err != nil {
		return "", "", fmt.Errorf("client.Do: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		body, err := io.ReadAll(res.Body)
		fmt.Println(string(body))
		if err != nil {
			return "", "", fmt.Errorf("io.ReadAll: %w", err)
		}
		return "", "", fmt.Errorf("salute speach request failded: %s, %s", res.Status, string(body))
	}
	var response responseApiCreateRecognizeTask
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return "", "", fmt.Errorf("decoder.Decode: %w", err)
	}
	return response.Result.ID, response.Result.Status, nil
}
