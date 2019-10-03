package bot

import (
	"fmt"

	vocab "github.com/mrKitikat/telegrambottest/src/app/bot/vocabulary"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
)

func (b *Bot) SendMenuMessage(ChatId int64) (tgbotapi.InlineKeyboardMarkup, string, error) {
	UserHistory[ChatId] = ""
	kb := b.newMainMenuKeyboard(ChatId)
	price, diff, err := b.Api.GetPrice()
	if err != nil {
		return kb, "", err
	}

	txt := fmt.Sprintf(vocab.GetTranslate("Select", b.Dlg[ChatId].language), price, diff)
	msg := tgbotapi.NewMessage(b.Dlg[ChatId].ChatId, txt)
	msg.ReplyMarkup = kb
	msg.ParseMode = "markdown"
	return kb, txt, nil
}

// SendMenuChoose ..
func (b *Bot) SendMenuChoose(ChatId int64) {
	kb := b.GetChooseKb(ChatId)
	msg := tgbotapi.EditMessageTextConfig{
		BaseEdit: tgbotapi.BaseEdit{
			ChatID:      b.Dlg[ChatId].ChatId,
			MessageID:   b.Dlg[ChatId].MessageId,
			ReplyMarkup: &kb,
		},
		Text:      vocab.GetTranslate("Save", b.Dlg[ChatId].language),
		ParseMode: "markdown",
	}

	b.Bot.Send(msg)
}

// GetChooseKb ..
func (b *Bot) GetChooseKb(ChatId int64) tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(vocab.GetTranslate("Yes", b.Dlg[ChatId].language), yescommand),
			tgbotapi.NewInlineKeyboardButtonData(vocab.GetTranslate("No", b.Dlg[ChatId].language), nocommand),
		),
	)
}

// CheckKeyboard ..
func (b *Bot) CheckKeyboard(ChatId int64) tgbotapi.InlineKeyboardMarkup {
	_, ok := PreviousMessage[ChatId]
	if ok {
		delete(PreviousMessage, ChatId)
	}
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(vocab.GetTranslate("Check", b.Dlg[ChatId].language), checkcommand),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(vocab.GetTranslate("Cancel", b.Dlg[ChatId].language), cancelComm),
		),
	)
}

// newMainMenuKeyboard is main menu keyboar: price, buy, sell, sales.
func (b *Bot) newMainMenuKeyboard(ChatId int64) tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(vocab.GetTranslate("Buy", b.Dlg[ChatId].language), buyCommand),
			tgbotapi.NewInlineKeyboardButtonData(vocab.GetTranslate("Sell", b.Dlg[ChatId].language), sellCommand),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(vocab.GetTranslate("Loots", b.Dlg[ChatId].language), salesCommand),
		),
	)
}

// vocabuageKeybord is keybouad for select vocabuage.
func (b *Bot) newVocabuageKeybord() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Русский", rusvocabCommand),
			tgbotapi.NewInlineKeyboardButtonData("English", engvocabCommand),
		),
	)
}

// CancelKeyboard ..
func (b *Bot) CancelKeyboard(ChatId int64) tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(vocab.GetTranslate("Cancel", b.Dlg[ChatId].language), cancelComm),
		),
	)
}

func (b *Bot) ShareCancel(ChatId int64, link string) tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(vocab.GetTranslate("Cancel", b.Dlg[ChatId].language), cancelComm),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonSwitch(vocab.GetTranslate("Share", b.Dlg[ChatId].language), link),
		),
	)
}
