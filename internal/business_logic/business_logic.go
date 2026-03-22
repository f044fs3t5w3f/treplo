package business_logic

import (
	"github.com/a-kuleshov/treplo/internal/business_logic/pipe"
	"github.com/a-kuleshov/treplo/internal/db"
)

type BusinessLogic struct {
	repo          db.Repository
	fileProcessor pipe.FileProcessor
}

func NewBusinessLogic(repo db.Repository, fileProcessor pipe.FileProcessor) *BusinessLogic {
	return &BusinessLogic{
		repo:          repo,
		fileProcessor: fileProcessor,
	}
}
