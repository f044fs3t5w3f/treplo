package business_logic

import (
	"github.com/a-kuleshov/treplo/internal/business_logic/pipe"
	"github.com/a-kuleshov/treplo/internal/db"
	"github.com/a-kuleshov/treplo/pkg/sber/gigachat"
)

type BusinessLogic struct {
	repo          db.Repository
	fileProcessor pipe.FileProcessor
	textGenerator gigachat.TextGenerator
}

func NewBusinessLogic(repo db.Repository, fileProcessor pipe.FileProcessor, textGenerator gigachat.TextGenerator) *BusinessLogic {
	return &BusinessLogic{
		repo:          repo,
		fileProcessor: fileProcessor,
		textGenerator: textGenerator,
	}
}
