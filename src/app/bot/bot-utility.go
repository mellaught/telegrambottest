package bot

import (
	"fmt"
	"strings"
	"time"

	vocab "github.com/mrKitikat/telegrambottest/src/app/bot/vocabulary"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
)

const (
	startCommand    = "start"
	priceCommand    = "price"
	buyCommand      = "buy"
	sellCommand     = "sell"
	salesCommand    = "lookat"
	getMainMenu     = "getmenu"
	checkcommand    = "check"
	settingsMenu    = "settings"
	language        = "language"
	engvocabCommand = "englanguage"
	rusvocabCommand = "ruslanguage"
	newBTC          = "newBTC"
	newMinter       = "newMinter"
	sendBTC         = "sendBTC"
	sendMinter      = "sendMinter"
	sendEmail       = "sendEmail"
	sendPrice       = "sendPrice"
	newEmail        = "newEmail"
	cancelComm      = "cancel"
	yescommand      = "yes"
	nocommand       = "not"
)

func (b *Bot) CancelHandler(ChatId int64) {
	fmt.Println("Cancel", UserHistory[ChatId])
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
			fmt.Println("HERE")
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

func (b *Bot) ChangeCurrency(ChatId int64) {
	timeout := time.After(10 * time.Minute)
	tick := time.Tick(1 * time.Minute)
	CallId := b.Dlg[ChatId].CallBackId
	MessageId := b.Dlg[ChatId].MessageId
	for {
		select {
		case <-timeout:
			return
		case <-tick:
			fmt.Println(b.Dlg[ChatId].MessageId, MessageId)
			if MessageId == b.Dlg[ChatId].MessageId && CallId == b.Dlg[ChatId].CallBackId {
				fmt.Println("HERE")
				fmt.Println(b.Dlg[ChatId].CallBackId, CallId)
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
		Text:      txt,
		ParseMode: "markdown",
	}
	b.Bot.Send(msg)
}

func (b *Bot) PrintAndSendError(err error, ChatId int64) {
	fmt.Println(err)
	msg := tgbotapi.NewMessage(b.Dlg[ChatId].ChatId, vocab.GetTranslate("Error", b.Dlg[ChatId].language))
	b.Bot.Send(msg)
}

func (b Bot) SendMessage(txt string, ChatId int64, kb interface{}) {
	msg := tgbotapi.NewMessage(b.Dlg[ChatId].ChatId, txt)
	msg.ParseMode = "markdown"
	msg.ReplyMarkup = kb
	b.Bot.Send(msg)
}
