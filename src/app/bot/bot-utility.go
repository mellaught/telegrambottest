package bot

import (
	"fmt"
	"strings"
	"time"

	vocab "github.com/mrKitikat/telegrambottest/src/app/bot/vocabulary"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
)

const (
	startCommand     = "start"
	priceCommand     = "price"
	buyCommand       = "buy"
	sellCommand      = "sell"
	salesCommand     = "orders"
	getMainMenu      = "home"
	checkcommandBuy  = "checkBuy"
	checkcommandSell = "checkSell"
	settingsMenu     = "settings"
	language         = "language"
	engvocabCommand  = "englanguage"
	rusvocabCommand  = "ruslanguage"
	newBTC           = "newBTC"
	newMinter        = "newMinter"
	sendBTC          = "sendBTC"
	sendMinter       = "sendMinter"
	sendEmail        = "sendEmail"
	sendPrice        = "sendPrice"
	newEmail         = "newEmail"
	cancelComm       = "cancel"
	yescommand       = "yes"
	nocommand        = "not"
)

func (b *Bot) CancelHandler(ChatId int64) {
	if strings.Contains(UserHistory[ChatId], "buy") {
		if UserHistory[ChatId][4:] == "1" {
			kb, txt, err := b.SendMenuMessage(ChatId)
			if err != nil {
				b.PrintAndSendError(err, ChatId)
				return
			}
			go b.ChangeCurrency(ChatId)
			b.EditAndSend(&kb, txt, ChatId)
		} else if UserHistory[ChatId][4:] == "2" {
			UserHistory[ChatId] = "buy_1"
			kb, txt, err := b.SendMinterAddresses(ChatId)
			if err != nil {
				fmt.Println(err)
			}

			b.EditAndSend(&kb, txt, ChatId)
		}

	} else if strings.Contains(UserHistory[ChatId], "sell") {
		fmt.Println(UserHistory[ChatId])
		if UserHistory[ChatId][5:] == "1" {
			kb, txt, err := b.SendMenuMessage(ChatId)
			if err != nil {
				b.PrintAndSendError(err, ChatId)
				return
			}
			go b.ChangeCurrency(ChatId)
			b.EditAndSend(&kb, txt, ChatId)
		} else if UserHistory[ChatId][5:] == "2" {
			UserHistory[ChatId] = "sell_1"
			kb := b.CancelKeyboard(ChatId)
			txt := vocab.GetTranslate("Coin", b.Dlg[ChatId].language)
			b.EditAndSend(&kb, txt, ChatId)
		} else if UserHistory[ChatId][5:] == "3" {
			UserHistory[ChatId] = "sell_2"
			txt := vocab.GetTranslate("Select price", b.Dlg[ChatId].language)
			kb := b.CancelKeyboard(ChatId)
			b.EditAndSend(&kb, txt, ChatId)
		} else if UserHistory[ChatId][5:] == "4" {
			kb, txt, err := b.SendBTCAddresses(ChatId)
			if err != nil {
				b.PrintAndSendError(err, ChatId)
				return
			}
			b.SendMessage(txt, ChatId, kb)
		}
	} else if strings.Contains(UserHistory[ChatId], "loot") {
		kb, txt, err := b.SendMenuMessage(ChatId)
		if err != nil {
			b.PrintAndSendError(err, ChatId)
			return
		}
		go b.ChangeCurrency(ChatId)
		b.EditAndSend(&kb, txt, ChatId)
	}

}

func (b *Bot) UpdatePrice() {
	for {

		price, diff, err := b.Api.GetPrice()
		if err != nil {
			fmt.Println(err)
			time.Sleep(5 * time.Second)
			continue
		}
		CurrentPrice = price
		CurrnetMarkup = diff
		time.Sleep(30 * time.Minute)
	}
}

func (b *Bot) ChangeCurrency(ChatId int64) {
	tick := time.Tick(15 * time.Minute)
	CallId := b.Dlg[ChatId].CallBackId
	MessageId := b.Dlg[ChatId].MessageId
	for {
		select {
		case <-tick:
			if MessageId == b.Dlg[ChatId].MessageId || CallId == b.Dlg[ChatId].CallBackId {
				kb, txt, err := b.SendMenuMessage(ChatId)
				if err != nil {
					fmt.Println(err)
					return
				}
				b.EditAndSend(&kb, txt, ChatId)
				continue
			} else {
				return
			}

		}
	}
}

func (b *Bot) EditAndSend(kb *tgbotapi.InlineKeyboardMarkup, txt string, ChatId int64) {
	msg := tgbotapi.EditMessageTextConfig{
		BaseEdit: tgbotapi.BaseEdit{
			ChatID:      b.Dlg[ChatId].ChatId,
			MessageID:   b.Dlg[ChatId].MessageId,
			ReplyMarkup: kb,
		},
		DisableWebPagePreview: true,
		Text:                  txt,
		ParseMode:             "markdown",
	}
	b.Bot.Send(msg)
}

func (b *Bot) PrintAndSendError(err error, ChatId int64) {
	fmt.Println(err)
	b.SendMessage(vocab.GetTranslate("Error", b.Dlg[ChatId].language), ChatId, nil)
}

func (b Bot) SendMessage(txt string, ChatId int64, kb interface{}) {
	msg := tgbotapi.NewMessage(b.Dlg[ChatId].ChatId, txt)
	msg.ParseMode = "markdown"
	msg.ReplyMarkup = kb
	msg.DisableWebPagePreview = true
	b.Bot.Send(msg)
	b.Dlg[ChatId].MessageId++
}
