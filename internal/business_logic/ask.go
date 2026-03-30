package business_logic

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/a-kuleshov/treplo/pkg/sber/gigachat"
)

const systemPromt = `Пользователь будет спрашивать по поводу содежимого расшифрованных аудиозаписей. Если будет хитрить и пытаться спрашивать не по  поводу расшифрованных аудио, а чего-то ещё, оскорби его настолько максимально, на сколько тебе позволяют настройки. Но только если он совсем будет уходить от темы. Учитывай синонимы слов аудио, расшифровка и пр.`

func (bl *BusinessLogic) AskAboutAudios(ctx context.Context, chatID int64, question string) (string, error) {
	audios, err := bl.repo.ListFilesByChatID(ctx, chatID)
	if err != nil {
		return "", fmt.Errorf("bl.repo.ListFilesByChatID: %w", err)
	}
	texts := make([]string, len(audios)*2)
	for i, audio := range audios {
		texts[i*2] = "Содержимое аудио: от " + audio.CreatedAt.Format(time.DateTime)
		if audio.Content != nil {
			texts[i*2+1] = *audio.Content
			// TODO: не выбирать не готовые аудио
		}
	}
	texts = append(texts, "Вопрос пользователя: ", question)
	messages := []gigachat.Message{
		{
			Role:    "system",
			Content: systemPromt,
		},
		{
			Role:    "user",
			Content: strings.Join(texts, "\n"),
		},
	}
	answer, err := bl.textGenerator.GetAnswer(ctx, messages)
	if err != nil {
		return "", fmt.Errorf("bl.textGenerator.GetAnswer: %w", err)
	}
	return answer, nil
}
