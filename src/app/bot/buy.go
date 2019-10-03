package bot

import (
	"fmt"
	"regexp"
	"time"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
	vocab "github.com/mrKitikat/telegrambottest/src/app/bot/vocabulary"
)

var BuyStatus = make(map[int64]string)

//
func (b *Bot) GetStatusBuy(ChatId int64) string {
	return BuyStatus[ChatId]
}

func (b *Bot) CheckMinter(address string) bool {
	return len(address) != 42 || address[:2] != "Mx" || address == "Mx00000000000000000000000000000000000001"
}

// SendMinterAddresses --
func (b *Bot) SendMinterAddresses(ChatId int64) (tgbotapi.InlineKeyboardMarkup, string, error) {

	//PreviousMessage[ChatId] = Message[ChatId]
	keyboard := tgbotapi.InlineKeyboardMarkup{}
	addresses, err := b.DB.GetMinterAddresses(b.Dlg[ChatId].UserId)
	if err != nil {
		return keyboard, "", err
	}

	if len(addresses) > 0 {
		txt := vocab.GetTranslate("Select minter", b.Dlg[ChatId].language)
		for _, addr := range addresses {
			var row []tgbotapi.InlineKeyboardButton
			btn := tgbotapi.NewInlineKeyboardButtonData(addr, sendMinter+addr)
			row = append(row, btn)
			keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
		}

		btn := tgbotapi.NewInlineKeyboardButtonData(vocab.GetTranslate("Cancel", b.Dlg[ChatId].language), cancelComm)
		var row []tgbotapi.InlineKeyboardButton
		row = append(row, btn)
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
		msg := tgbotapi.NewMessage(b.Dlg[ChatId].ChatId, txt)
		msg.ReplyMarkup = keyboard
		msg.ParseMode = "markdown"
		Message[ChatId] = msg
		Message[ChatId] = msg
		return keyboard, txt, nil
	} else {
		txt := vocab.GetTranslate("New minter", b.Dlg[ChatId].language)
		msg := tgbotapi.NewMessage(b.Dlg[ChatId].ChatId, txt)
		btn := tgbotapi.NewInlineKeyboardButtonData(vocab.GetTranslate("Cancel", b.Dlg[ChatId].language), cancelComm)
		var row []tgbotapi.InlineKeyboardButton
		row = append(row, btn)
		msg.ReplyMarkup = keyboard
		msg.ParseMode = "markdown"
		Message[ChatId] = msg
		return keyboard, txt, nil
	}

}

// CheckEmail ..
func (b *Bot) CheckEmail(email string) bool {
	re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	if !re.MatchString(email) || email == "mail@example.com" {
		return false
	}
	return true
}

// SendEmail --
func (b *Bot) SendEmail(ChatId int64) (tgbotapi.InlineKeyboardMarkup, string, error) {

	PreviousMessage[ChatId] = Message[ChatId]
	keyboard := tgbotapi.InlineKeyboardMarkup{}
	addresses, err := b.DB.GetEmails(b.Dlg[ChatId].UserId)
	if err != nil {
		return keyboard, "", err
	}

	if len(addresses) > 0 {
		txt := vocab.GetTranslate("Select email", b.Dlg[ChatId].language)
		for _, addr := range addresses {
			var row []tgbotapi.InlineKeyboardButton
			btn := tgbotapi.NewInlineKeyboardButtonData(addr, sendEmail+addr)
			row = append(row, btn)
			keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
		}
		btn := tgbotapi.NewInlineKeyboardButtonData(vocab.GetTranslate("Cancel", b.Dlg[ChatId].language), cancelComm)
		var row []tgbotapi.InlineKeyboardButton
		row = append(row, btn)
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
		msg := tgbotapi.NewMessage(b.Dlg[ChatId].ChatId, txt)
		msg.ReplyMarkup = keyboard
		msg.ParseMode = "markdown"
		Message[ChatId] = msg
		return keyboard, txt, nil

	} else {
		txt := vocab.GetTranslate("New email", b.Dlg[ChatId].language)
		msg := tgbotapi.NewMessage(b.Dlg[ChatId].ChatId, txt)
		btn := tgbotapi.NewInlineKeyboardButtonData(vocab.GetTranslate("Cancel", b.Dlg[ChatId].language), cancelComm)
		var row []tgbotapi.InlineKeyboardButton
		row = append(row, btn)
		msg.ReplyMarkup = keyboard
		msg.ParseMode = "markdown"
		Message[ChatId] = msg
		return keyboard, txt, nil
	}

}

// SendDepos --
func (b *Bot) SendDepos(ChatId int64) {
	price, diff, err := b.Api.GetPrice()
	if err != nil {
		fmt.Println(err)
		msg := tgbotapi.NewMessage(b.Dlg[ChatId].ChatId, vocab.GetTranslate("Error", b.Dlg[ChatId].language))
		b.Bot.Send(msg)
		return
	}
	txt := fmt.Sprintf(vocab.GetTranslate("Send deposit", b.Dlg[ChatId].language), price, diff)
	msg := tgbotapi.NewMessage(b.Dlg[ChatId].ChatId, txt)
	msg.ParseMode = "markdown"
	// msg.ReplyMarkup = b.CancelKeyboard(ChatId)
	b.Bot.Send(msg)
}

// CheckStatusBuy checks depos BTC and wait 2 confirmations
func (b *Bot) CheckStatusBuy(address string, ChatId int64) {
	timeout := time.After(60 * time.Minute)
	tick := time.Tick(5 * time.Second)
	willcoin := 0.
	BuyStatus[ChatId] = vocab.GetTranslate("Wait deposit", b.Dlg[ChatId].language)
	for {
		select {
		case <-timeout:
			if willcoin == 0. {
				BuyStatus[ChatId] = vocab.GetTranslate("No buy", b.Dlg[ChatId].language)
				return
			} else {
				continue
			}
		case <-tick:
			stat, err := b.Api.GetBTCDepositStatus(address)
			if err != nil {
				fmt.Println(err)
				time.Sleep(10 * time.Second)
				continue
			}
			if stat.Data.WillReceive != willcoin {
				if willcoin == 0. {
					willcoin = stat.Data.WillReceive
					BuyStatus[ChatId] = fmt.Sprintf(vocab.GetTranslate("New deposit", b.Dlg[ChatId].language), stat.Data.WillReceive)
					time.Sleep(60 * time.Second)
				} else {
					ans := fmt.Sprintf(vocab.GetTranslate("Exchange is successful", b.Dlg[ChatId].language), willcoin)
					msg := tgbotapi.NewMessage(b.Dlg[ChatId].ChatId, ans)
					msg.ReplyMarkup = b.newMainMenuKeyboard(ChatId)
					BuyStatus[ChatId] = vocab.GetTranslate("No buy", b.Dlg[ChatId].language)
					b.Bot.Send(msg)
					b.SendMenuMessage(ChatId)
					return
				}
			}
		}
	}
}
