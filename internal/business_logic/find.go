package business_logic

import (
	"context"
	"fmt"
	"strings"
	"unicode"

	"github.com/a-kuleshov/treplo/internal/models"
)

const shift = 200

type entry struct {
	Entry  string
	Before string
	After  string
}
type result struct {
	File  *models.File
	Entry []entry
}

func (bl *BusinessLogic) FindFiles(ctx context.Context, words string, chatID int64) ([]result, error) {
	// TODO: make inline keyboard
	keyWords := splitText(words)
	if len(keyWords) == 0 {
		return nil, ErrNoFiles
	}
	files, err := bl.repo.ListFilesByChatIDAndKeywords(ctx, keyWords, chatID)
	if err != nil {
		return nil, fmt.Errorf("bl.repo.ListFilesByChatIDAndKeywords: %w", err)
	}
	if len(files) == 0 {
		return nil, ErrNoFiles
	}
	results := make([]result, 0, len(files))
	for _, file := range files {
		results = append(results, result{
			File:  file,
			Entry: nil,
		})
		// TODO: find text entry, text before, text after
	}
	return results, nil
}

func splitText(s string) []string {
	return strings.FieldsFunc(s, func(r rune) bool {
		if unicode.IsSpace(r) {
			return true
		}

		if unicode.IsPunct(r) && r != '-' {
			return true
		}

		return false
	})
}
