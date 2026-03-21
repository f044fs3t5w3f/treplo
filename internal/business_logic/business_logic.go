package business_logic

import (
	"github.com/a-kuleshov/treplo/internal/pipe"
)

type BusinessLogic struct {
	repo          Repository
	fileProcessor pipe.FileProcessor
}

func NewBusinessLogic(repo Repository, fileProcessor pipe.FileProcessor) *BusinessLogic {
	return &BusinessLogic{
		repo:          repo,
		fileProcessor: fileProcessor,
	}
}
