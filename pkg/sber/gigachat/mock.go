package gigachat

import (
	"context"
)

type GigaChatServiceMock struct {
}

func (g GigaChatServiceMock) GetAnswer(ctx context.Context, messages []Message) (string, error) {
	return "Ответ", nil
}

var _ TextGenerator = (*GigaChatServiceMock)(nil)
