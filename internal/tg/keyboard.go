package tg

import (
	"fmt"
	"strconv"
	"time"

	"github.com/a-kuleshov/treplo/internal/models"
	tgBotApi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func makeAudioFilesKeyboard(files []*models.File, page int, hasPrevious, hasNext bool) tgBotApi.InlineKeyboardMarkup {
	rows := make([][]tgBotApi.InlineKeyboardButton, 0, len(files))

	for _, file := range files {
		rows = append(rows,
			tgBotApi.NewInlineKeyboardRow(
				tgBotApi.NewInlineKeyboardButtonData(
					file.CreatedAt.Format(time.DateTime),
					strconv.FormatInt(file.ID, 10),
				),
			),
		)
	}
	if hasPrevious || hasNext {
		paginationRow := tgBotApi.NewInlineKeyboardRow()
		if hasPrevious {
			paginationRow = append(paginationRow,
				tgBotApi.NewInlineKeyboardButtonData(
					"<",
					fmt.Sprintf("p%d", page-1),
				),
			)
		}
		if hasNext {
			paginationRow = append(paginationRow,
				tgBotApi.NewInlineKeyboardButtonData(
					">",
					fmt.Sprintf("p%d", page+1),
				),
			)
		}
		rows = append(rows, paginationRow)
	}

	keyboard := tgBotApi.NewInlineKeyboardMarkup(rows...)
	return keyboard
}
