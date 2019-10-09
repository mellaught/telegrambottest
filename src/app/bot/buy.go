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
	if address == "Mx00000000000000000000000000000000000001" {
		return false
	}

	return len(address) == 42 || address[:2] != "Mx"
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
		return keyboard, txt, nil
	} else {
		txt := vocab.GetTranslate("New minter", b.Dlg[ChatId].language)
		btn := tgbotapi.NewInlineKeyboardButtonData(vocab.GetTranslate("Cancel", b.Dlg[ChatId].language), cancelComm)
		var row []tgbotapi.InlineKeyboardButton
		row = append(row, btn)
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
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
		return keyboard, txt, nil

	} else {
		txt := vocab.GetTranslate("New email", b.Dlg[ChatId].language)
		btn := tgbotapi.NewInlineKeyboardButtonData(vocab.GetTranslate("Cancel", b.Dlg[ChatId].language), cancelComm)
		var row []tgbotapi.InlineKeyboardButton
		row = append(row, btn)
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
		return keyboard, txt, nil
	}

}

// SendDepos --
func (b *Bot) SendDepos(ChatId int64) {
	price, diff, err := b.Api.GetPrice()
	if err != nil {
		b.PrintAndSendError(err, ChatId)
		return
	}
	amount, bonus, err := b.Api.GetBonus()
	if err != nil {
		fmt.Println(err)
		return
	}

	txt := fmt.Sprintf(vocab.GetTranslate("Send deposit", b.Dlg[ChatId].language), price, diff, amount, bonus)
	b.SendMessage(txt, ChatId, nil)
}

// SendDepos --
func (b *Bot) EditDepos(ChatId int64) {
	price, diff, err := b.Api.GetPrice()
	if err != nil {
		b.PrintAndSendError(err, ChatId)
		return
	}
	amount, bonus, err := b.Api.GetBonus()
	if err != nil {
		fmt.Println(err)
		return
	}
	txt := fmt.Sprintf(vocab.GetTranslate("Send deposit", b.Dlg[ChatId].language), price, diff, amount, bonus)
	b.EditAndSend(nil, txt, ChatId)
}

// BuyFinal is function for command "/buy".
// Requests an email from the user and Minter deposit address.
// Requests the "bitcoinDepositAddress" method with the received data.
func (b *Bot) BuyFinal(ChatId int64) {
	fmt.Println("Buy data:", MinterAddress[b.Dlg[ChatId].ChatId], EmailAddress[b.Dlg[ChatId].ChatId])
	addr, err := b.Api.GetBTCDeposAddress(MinterAddress[b.Dlg[ChatId].ChatId], "BIP", EmailAddress[b.Dlg[ChatId].ChatId])
	if err != nil {
		b.Dlg[ChatId].Command = ""
		b.SendMessage(err.Error(), ChatId, b.newMainMenuKeyboard(ChatId))
		return
	}
	b.Dlg[ChatId].Command = ""
	b.SendMessage(addr, ChatId, b.CheckKeyboardBuy(ChatId))
	go b.CheckStatusBuy(addr, ChatId)
	return
}

// CheckStatusBuy checks depos BTC and wait 2 confirmations
func (b *Bot) CheckStatusBuy(address string, ChatId int64) {
	timeout := time.After(60 * time.Minute)
	tick := time.Tick(5 * time.Second)
	stat, err := b.Api.GetBTCDepositStatus(address)
	if err != nil {
		fmt.Println("Buy api error:", err)
		time.Sleep(30 * time.Second)
	}
	willcoin := stat.Data.WillReceive
	start := willcoin
	BuyStatus[ChatId] = vocab.GetTranslate("Wait deposit", b.Dlg[ChatId].language)
	for {
		select {
		case <-timeout:
			if willcoin == start {
				BuyStatus[ChatId] = vocab.GetTranslate("No buy", b.Dlg[ChatId].language)
				return
			} else {
				continue
			}
		case <-tick:
			stat, err := b.Api.GetBTCDepositStatus(address)
			if err != nil {
				fmt.Println("Buy api error:", err)
				time.Sleep(30 * time.Second)
				continue
			}
			if stat.Data.WillReceive != willcoin {
				if willcoin == start {
					willcoin = stat.Data.WillReceive - start
					coinSend := stat.Data.WillReceive - start
					BuyStatus[ChatId] = fmt.Sprintf(vocab.GetTranslate("New deposit", b.Dlg[ChatId].language), coinSend)
					time.Sleep(60 * time.Second)
				} else {
					ans := fmt.Sprintf(vocab.GetTranslate("Exchange is successful", b.Dlg[ChatId].language), willcoin)
					b.SendMessage(ans, ChatId, b.newMainMenuKeyboard(ChatId))
					BuyStatus[ChatId] = vocab.GetTranslate("No buy", b.Dlg[ChatId].language)
					time.Sleep(5 * time.Second)
					kb, txt, err := b.SendMenuMessage(ChatId)
					if err != nil {
						b.PrintAndSendError(err, ChatId)
						return
					}
					go b.ChangeCurrency(ChatId)
					b.SendMessage(txt, ChatId, kb)
					return
				}
			}
		}
	}
}
