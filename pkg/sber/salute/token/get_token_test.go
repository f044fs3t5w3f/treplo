package token

import (
	"bytes"
	"io"
	"testing"

	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetAccessTokenFromResponse_ok(t *testing.T) {
	goodResponse := `{
    "access_token": "token",
    "expires_at": 1617814516729
}
`
	reader := bytes.NewReader([]byte(goodResponse))
	closer := io.NopCloser(reader)
	token, expiresAt, err := getAccessTokenFromResponse(closer)
	assert.NoError(t, err)
	assert.Equal(t, "token", token)
	assert.Equal(t, time.Unix(1617814516729, 0), expiresAt)
}
func TestGetAccessTokenFromResponse_error(t *testing.T) {
	badResponse := `{"error": {"status_code": 401, "message": "Authorization error: header is incorrect"}}`
	reader := bytes.NewReader([]byte(badResponse))
	closer := io.NopCloser(reader)
	_, _, err := getAccessTokenFromResponse(closer)
	assert.Error(t, err)
}
