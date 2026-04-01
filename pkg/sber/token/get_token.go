package token

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

func generateRqUID() string {
	u4 := uuid.New()
	return u4.String()
}

func getAccessToken(clientSecret string, scope string) (accessToken string, expiresAt time.Time, err error) {
	url := "https://ngw.devices.sberbank.ru:9443/api/v2/oauth"

	payload := strings.NewReader("scope=" + scope)

	req, err := http.NewRequest("POST", url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Basic "+clientSecret)
	req.Header.Add("RqUID", generateRqUID())

	client := http.DefaultClient
	res, err := client.Do(req)
	if err != nil {
		return "", time.Now(), fmt.Errorf("client.Do: %w", err)
	}
	return getAccessTokenFromResponse(res.Body)
}

func getAccessTokenFromResponse(respBody io.ReadCloser) (accessToken string, expiresAt time.Time, err error) {
	defer respBody.Close()

	responseBody, err := io.ReadAll(respBody)
	if err != nil {
		return
	}
	var response struct {
		AccessToken string `json:"access_token"`
		ExpiresAt   int64  `json:"expires_at"`
		Error       *struct {
			Message string `json:"message"`
		} `json:"error"`
	}
	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		return "", time.Now(), fmt.Errorf("json.Unmarshal: %w", err)
	}

	if response.Error != nil {
		return "", time.Now(), fmt.Errorf("error from server: %s", response.Error.Message)
	}

	expiresAt = time.Unix(response.ExpiresAt/time.Second.Milliseconds(), 0)

	return response.AccessToken, expiresAt, nil
}
