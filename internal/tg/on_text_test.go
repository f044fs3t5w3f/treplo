package tg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractCommandAndPayload(t *testing.T) {
	cases := []struct {
		name, input, command, payload string
	}{
		{"Just command", "/start", "start", ""},
		{"Command with spaces", "/start ", "start", ""},
		{"Extra spaces", "/start     Hi!", "start", "Hi!"},
		{"Extra spaces in payload", "/start     Hi!", "start", "Hi!"},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			command, payload := extractCommandAndPayload(testCase.input)
			assert.Equal(t, testCase.command, command)
			assert.Equal(t, testCase.payload, payload)
		})
	}
}
