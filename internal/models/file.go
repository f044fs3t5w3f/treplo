package models

import "time"

const (
	FileStatusNew   = "NEW"
	FileStatusDone  = "DONE"
	FileStatusError = "ERROR"
)

type File struct {
	ID                      int64
	FileID                  string
	ChatID                  int64
	MessageID               int
	Filepath                *string
	Encoding                *string
	SaluteId                *string
	RecognizeTaskID         *string
	RecognizeStatus         *string
	ResponseFileID          *string
	Content                 *string
	ProcessNotificationSent bool
	CreatedAt               time.Time
	Status                  string
}
