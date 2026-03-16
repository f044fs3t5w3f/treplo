package tg

import "context"

type mechanicService interface {
	AddAudio(ctx context.Context, chatID int64, fileID string) error
}
