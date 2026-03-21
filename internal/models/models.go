package models

type File struct {
	ID                      int64
	FileID                  string
	ChatID                  int64
	MessageID               int
	Filepath                *string
	Encoding                string
	SaluteId                *string
	RecognizeTaskID         *string
	RecognizeStatus         *string
	ResponseFileID          *string
	Content                 *string
	ProcessNotificationSent bool
}
