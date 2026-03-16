package mechanics

import "context"

type Repository interface {
	AddAudio(ctx context.Context, chatID int64, fileID string) error
}
