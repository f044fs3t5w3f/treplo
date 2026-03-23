package pipe

import tgBotApi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type tgNotifier struct {
	tgbotapi *tgBotApi.BotAPI
}

func (tn tgNotifier) Notify(replyToMessageId int, chatId int64, message string) error {
	tgMessage := tgBotApi.NewMessage(chatId, message)
	tgMessage.ReplyToMessageID = replyToMessageId
	_, err := tn.tgbotapi.Send(tgMessage)
	return err
}
