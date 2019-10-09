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
	return kb, txt, nil
}

// SendMenuChoose ..
func (b *Bot) SendMenuChoose(ChatId int64) {
	kb := b.GetChooseKb(ChatId)
	b.SendMessage(vocab.GetTranslate("Save", b.Dlg[ChatId].language), ChatId, kb)
}

func (b *Bot) EditMenuChoose(ChatId int64) {
	kb := b.GetChooseKb(ChatId)
	b.EditAndSend(&kb, vocab.GetTranslate("Save", b.Dlg[ChatId].language), ChatId)
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

// CheckKeyboardBuy ..
func (b *Bot) CheckKeyboardBuy(ChatId int64) tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(vocab.GetTranslate("Check", b.Dlg[ChatId].language), checkcommandBuy),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(vocab.GetTranslate("Cancel", b.Dlg[ChatId].language), cancelComm),
		),
	)
}

// CheckKeyboardSell ..
func (b *Bot) CheckKeyboardSell(ChatId int64) tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(vocab.GetTranslate("Check", b.Dlg[ChatId].language), checkcommandSell),
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

func (b *Bot) Share(ChatId int64, link string) tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonSwitch(vocab.GetTranslate("Share", b.Dlg[ChatId].language), link),
		),
	)
}
