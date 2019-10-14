package bot

import (
	"fmt"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
	vocab "github.com/mrKitikat/telegrambottest/src/app/bot/vocabulary"
	stct "github.com/mrKitikat/telegrambottest/src/app/structs"
)

const lootlink = "https://bip.dev/trade/"

// Method for sending loots in markdown style to user.
func (b *Bot) SendLoots(loots []*stct.Loot, ChatId int64) (tgbotapi.InlineKeyboardMarkup, string) {

	keyboard := tgbotapi.InlineKeyboardMarkup{}

	for _, loot := range loots {
		var row []tgbotapi.InlineKeyboardButton
		lText := fmt.Sprintf(vocab.GetTranslate("Loot", b.Dlg[ChatId].language), loot.Amout, loot.Coin, loot.Price)
		btn := tgbotapi.NewInlineKeyboardButtonURL(lText, lootlink+loot.Tag)
		row = append(row, btn)
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
	}

	var row []tgbotapi.InlineKeyboardButton
	btn := tgbotapi.NewInlineKeyboardButtonData(vocab.GetTranslate("Cancel", b.Dlg[ChatId].language), cancelComm)
	row = append(row, btn)
	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
	b.EditAndSend(&keyboard, vocab.GetTranslate("Your loots", b.Dlg[ChatId].language), ChatId)
	return keyboard, vocab.GetTranslate("Your loots", b.Dlg[ChatId].language)
}
