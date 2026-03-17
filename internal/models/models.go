package models

type File struct {
	ID              int64
	FileID          string
	ChatID          int64
	Filepath        *string
	SaluteId        *string
	RecognizeTaskID *string
	RecognizeStatus *string
}
