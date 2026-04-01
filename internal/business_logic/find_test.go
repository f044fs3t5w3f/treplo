package business_logic

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSplitText(t *testing.T) {
	testCases := []struct {
		name     string
		text     string
		expected []string
	}{
		{
			"Multi spaces",
			"кот    пёс",
			[]string{"кот", "пёс"},
		},
		{
			"dash",
			"кот-пёс",
			[]string{"кот-пёс"},
		},
		{
			"punctuation",
			"кот, пёс",
			[]string{"кот", "пёс"},
		},
	}
	for _, testCase := range testCases {
		results := splitText(testCase.text)
		t.Run(testCase.name, func(t *testing.T) {
			assert.Equal(t, testCase.expected, results)
		})
	}
}
