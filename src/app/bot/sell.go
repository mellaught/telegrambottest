package bot

import (
	"fmt"
	"regexp"
	"strconv"
	"time"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
	vocab "github.com/mrKitikat/telegrambottest/src/app/bot/vocabulary"
)

var SellStatus = make(map[int64]string)

func (b *Bot) GetStatusSell(ChatId int64) string {
	return SellStatus[ChatId]
}

// CheckPrice ..
func (b *Bot) CheckPrice(chatId int64, price string) bool {
	if s, err := strconv.ParseFloat(price, 64); err == nil {
		if 0.01 <= s && s <= 0.32 {
			PriceToSell[chatId] = s
			return true
		} else {
			return false
		}
	} else {
		return false
	}
}

// CheckCoin ..
func (b *Bot) CheckCoin(coin string) bool {
	re := regexp.MustCompile("^[0-9-A-Z]{3,10}$")
	return re.MatchString(coin)
}

// CheckBitcoin ..
func (b *Bot) CheckBitcoin(address string) bool {
	re := regexp.MustCompile("^[13][a-km-zA-HJ-NP-Z1-9]{25,34}$")
	return re.MatchString(address)
}

//
func (b *Bot) SendBTCAddresses(ChatId int64) (tgbotapi.InlineKeyboardMarkup, string, error) {

	keyboard := tgbotapi.InlineKeyboardMarkup{}
	addresses, err := b.DB.GetBTCAddresses(b.Dlg[ChatId].UserId)
	if err != nil {
		return keyboard, "", err
	}
	if len(addresses) > 0 {
		txt := vocab.GetTranslate("Select bitcoin", b.Dlg[ChatId].language)
		for _, addr := range addresses {
			var row []tgbotapi.InlineKeyboardButton
			btn := tgbotapi.NewInlineKeyboardButtonData(addr, sendBTC+addr)
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
		txt := vocab.GetTranslate("New bitcoin", b.Dlg[ChatId].language)
		msg := tgbotapi.NewMessage(b.Dlg[ChatId].ChatId, txt)
		btn := tgbotapi.NewInlineKeyboardButtonData(vocab.GetTranslate("Cancel", b.Dlg[ChatId].language), cancelComm)
		var row []tgbotapi.InlineKeyboardButton
		row = append(row, btn)
		msg.ReplyMarkup = keyboard
		msg.ParseMode = "markdown"
		msg.ReplyMarkup = keyboard
		Message[ChatId] = msg
		return keyboard, txt, nil
	}
}

// SellFinal
func (b *Bot) SellFinal(ChatId int64) {
	fmt.Println("Sell data:", BitcoinAddress[b.Dlg[ChatId].ChatId], CoinToSell[b.Dlg[ChatId].ChatId], PriceToSell[b.Dlg[ChatId].ChatId])
	depos, err := b.Api.GetMinterDeposAddress(BitcoinAddress[b.Dlg[ChatId].ChatId], CoinToSell[b.Dlg[ChatId].ChatId], PriceToSell[b.Dlg[ChatId].ChatId])
	if err != nil {
		msg := tgbotapi.NewMessage(b.Dlg[ChatId].ChatId, err.Error())
		b.Bot.Send(msg)
		kb, txt, err := b.SendMenuMessage(ChatId)
		if err != nil {
			b.PrintAndSendError(err, ChatId)
			return
		}
		go b.ChangeCurrency(ChatId)
		b.SendMessage(txt, ChatId, kb)
		return
	}

	b.SendMenuChoose(ChatId)
	b.Dlg[ChatId].Command = ""
	txt := fmt.Sprintf(vocab.GetTranslate("Send your coins", b.Dlg[ChatId].language), CoinToSell[ChatId], CoinToSell[ChatId], "https://bip.dev/trade/"+depos.Data.Tag)
	msg := tgbotapi.NewMessage(b.Dlg[ChatId].ChatId, txt)
	msg.ParseMode = "markdown"
	msg.ReplyMarkup = b.ShareCancel(ChatId, "https://bip.dev/trade/"+depos.Data.Tag)
	b.Bot.Send(msg)
	newmsg := tgbotapi.NewMessage(b.Dlg[ChatId].ChatId, depos.Data.Address)
	b.Bot.Send(newmsg)
	go b.CheckStatusSell(depos.Data.Tag, ChatId)
	time.Sleep(3 * time.Second)
	kb, txt, err := b.SendMenuMessage(ChatId)
	if err != nil {
		b.PrintAndSendError(err, ChatId)
		return
	}
	go b.ChangeCurrency(ChatId)
	b.SendMessage(txt, ChatId, kb)
	return
}

// CheckStatusSell checks status of deposit for method Sell().
func (b *Bot) CheckStatusSell(tag string, ChatId int64) {
	timeout := time.After(30 * time.Minute)
	tick := time.Tick(5 * time.Second)
	amount := "0.0"
	SellStatus[ChatId] = vocab.GetTranslate("Wait deposit", b.Dlg[ChatId].language)
	for {
		select {
		case <-timeout:
			if amount == "0.0" {
				SellStatus[ChatId] = vocab.GetTranslate("No Sell", b.Dlg[ChatId].language)
				// msg := tgbotapi.NewMessage(b.Dlg[ChatId].ChatId, vocab.GetTranslate("timeout", b.Dlg[ChatId].language))
				// msg.ReplyMarkup = b.newMainKeyboard(ChatId)
				// b.Bot.Send(msg)
				return
			} else {
				continue
			}
		case <-tick:
			taginfo, err := b.Api.GetTagInfo(tag)
			if err != nil {
				fmt.Println(err)
				time.Sleep(10 * time.Second)
				continue
			}
			if taginfo.Data.Amount != amount {
				amount = taginfo.Data.Amount
				// Put in DB.
				b.DB.PutLoot(b.Dlg[ChatId].UserId, tag, taginfo)
				ans := fmt.Sprintf(vocab.GetTranslate("New deposit for sale", b.Dlg[ChatId].language),
					taginfo.Data.Amount, taginfo.Data.Coin, taginfo.Data.Price)
				msg := tgbotapi.NewMessage(b.Dlg[ChatId].ChatId, ans)
				b.Bot.Send(msg)
				//go a.CheckLootforSell(taginfo.Data.MinterAddress)
				return
			}

		}
	}
}
