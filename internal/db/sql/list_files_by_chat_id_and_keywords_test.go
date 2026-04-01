package sql

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetConditionsAndArgs(t *testing.T) {
	testCases := []struct {
		name       string
		keywords   []string
		chatId     int64
		conditions string
		args       []any
	}{
		{
			"no words",
			[]string{},
			1,
			"chat_id = $1",
			[]any{int64(1)},
		},
		{
			"with words",
			[]string{"word1", "word2"},
			1,
			"chat_id = $1 AND dialogue_content ILIKE $2 AND dialogue_content ILIKE $3",
			[]any{int64(1), "%word1%", "%word2%"},
		},
		{
			"escaping",
			[]string{"wo_%rd"},
			1,
			"chat_id = $1 AND dialogue_content ILIKE $2",
			[]any{int64(1), `%wo\_\%rd%`},
		}}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			conditions, args := getConditionsAndArgs(testCase.keywords, testCase.chatId)
			assert.Equal(t, testCase.conditions, conditions)
			assert.Equal(t, testCase.args, args)

		})
	}

}
