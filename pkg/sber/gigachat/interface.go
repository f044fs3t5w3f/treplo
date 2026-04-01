package gigachat

import "context"

type TextGenerator interface {
	GetAnswer(ctx context.Context, messages []Message) (string, error)
}
