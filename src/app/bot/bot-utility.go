package bot

import (
	"fmt"
	stct "telegrambottest/src/app/bipdev/structs"
	vocab "telegrambottest/src/app/bot/vocabulary"
	"time"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
)

const (
	startCommand    = "start"
	priceCommand    = "price"
	buyCommand      = "buy"
	sellCommand     = "sell"
	salesCommand    = "lookat"
	getMainMenu     = "getmenu"
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
)

//
func (b *Bot) GetBTCAddresses() tgbotapi.InlineKeyboardMarkup {

	keyboard := tgbotapi.InlineKeyboardMarkup{}
	var row []tgbotapi.InlineKeyboardButton
	btn := tgbotapi.NewInlineKeyboardButtonData(vocab.GetTranslate("New BTC", b.Dlg.language), newBTC)
	row = append(row, btn)
	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
	addresses, err := b.DB.GetBTCAddresses(b.Dlg.UserId)
	if err != nil {
		fmt.Println(err)
		return keyboard
	}
	if len(addresses) > 0 {
		for _, addr := range addresses {
			var row []tgbotapi.InlineKeyboardButton
			btn := tgbotapi.NewInlineKeyboardButtonData(addr, sendBTC+addr)
			row = append(row, btn)
			keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
		}
	}

	return keyboard
}

//
func (b *Bot) GetMinterAddresses() tgbotapi.InlineKeyboardMarkup {

	keyboard := tgbotapi.InlineKeyboardMarkup{}
	var row []tgbotapi.InlineKeyboardButton
	btn := tgbotapi.NewInlineKeyboardButtonData(vocab.GetTranslate("New minter", b.Dlg.language), newMinter)
	row = append(row, btn)
	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
	addresses, err := b.DB.GetBTCAddresses(b.Dlg.UserId)
	if err != nil {
		fmt.Println(err)
		return keyboard
	}
	if len(addresses) > 0 {
		for _, addr := range addresses {
			var row []tgbotapi.InlineKeyboardButton
			btn := tgbotapi.NewInlineKeyboardButtonData(addr, sendMinter+addr)
			row = append(row, btn)
			keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
		}
	}

	return keyboard

}

//
func (b *Bot) GetEmail() tgbotapi.InlineKeyboardMarkup {

	keyboard := tgbotapi.InlineKeyboardMarkup{}
	var row []tgbotapi.InlineKeyboardButton
	btn := tgbotapi.NewInlineKeyboardButtonData(vocab.GetTranslate("New email", b.Dlg.language), newEmail)
	row = append(row, btn)
	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
	addresses, err := b.DB.GetBTCAddresses(b.Dlg.UserId)
	if err != nil {
		fmt.Println(err)
		return keyboard
	}
	if len(addresses) > 0 {
		for _, addr := range addresses {
			var row []tgbotapi.InlineKeyboardButton
			btn := tgbotapi.NewInlineKeyboardButtonData(addr, sendEmail+addr)
			row = append(row, btn)
			keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
		}
	}

	return keyboard

}

func (b *Bot) GetPrice() tgbotapi.InlineKeyboardMarkup {

	keyboard := tgbotapi.InlineKeyboardMarkup{}
	prices, err := b.Api.GetAvailablePrices()
	if err != nil {
		fmt.Println(err)
		return keyboard
	}
	if len(prices) > 0 {
		for _, price := range prices {
			var row []tgbotapi.InlineKeyboardButton
			btn := tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("%d", price), sendPrice+fmt.Sprintf("%d", price))
			row = append(row, btn)
			keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
		}
	}

	return keyboard
}

// SendMenu edit message and send Inline Keyboard newMainMenuKeyboard()
func (b *Bot) SendMenu() {

	kb := b.newMainMenuKeyboard()
	newmsg := tgbotapi.EditMessageTextConfig{
		BaseEdit: tgbotapi.BaseEdit{
			ChatID:      b.Dlg.ChatId,
			MessageID:   b.Dlg.MessageId,
			ReplyMarkup: &kb,
		},
		Text: vocab.GetTranslate("Select", b.Dlg.language),
	}

	b.Bot.Send(newmsg)
}

// CheckStatusBuy checks depos BTC and wait 2 confirmations
func (b *Bot) CheckStatusBuy(address string) {
	timeout := time.After(2 * time.Minute)
	tick := time.Tick(3 * time.Second)
	willcoin := 0.
	for {
		select {
		case <-timeout:
			if willcoin == 0. {
				msg := tgbotapi.NewMessage(b.Dlg.ChatId, vocab.GetTranslate("timeout", b.Dlg.language))
				msg.ReplyMarkup = b.newMainKeyboard()
				b.Bot.Send(msg)
				return
			} else {
				continue
			}

		case <-tick:
			stat, err := b.Api.GetBTCDepositStatus(address)
			if err != nil {
				fmt.Println(err)
				msg := tgbotapi.NewMessage(b.Dlg.ChatId, vocab.GetTranslate("Error", b.Dlg.language))
				b.Bot.Send(msg)
				return
			}
			if stat.Data.WillReceive != willcoin {
				if willcoin == 0. {
					willcoin = stat.Data.WillReceive
					ans := fmt.Sprintf(vocab.GetTranslate("New deposit", b.Dlg.language), stat.Data.WillReceive)
					msg := tgbotapi.NewMessage(b.Dlg.ChatId, ans)
					b.Bot.Send(msg)
					time.Sleep(60 * time.Second)
				} else {
					ans := fmt.Sprintf(vocab.GetTranslate("Exchange is successful", b.Dlg.language), willcoin)
					msg := tgbotapi.NewMessage(b.Dlg.ChatId, ans)
					msg.ReplyMarkup = b.newMainKeyboard()
					b.Bot.Send(msg)
					return
				}
			}
		}
	}
}

// Sell is function for command /sell.
// func (b *Bot) Sell() {
// 	if len(b.Dlg.Text) > 3 {
// 		// checkvalidbitcoin
// 		CoinToSell[b.Dlg.UserId] = "MNT"
// 		depos, err := b.Api.GetMinterDeposAddress(b.Dlg.Text, CoinToSell[b.Dlg.UserId], PriceToSell[b.Dlg.UserId])
// 		if err != nil {
// 			msg := tgbotapi.NewMessage(b.Dlg.ChatId, err.Error())
// 			msg.ReplyMarkup = b.newMainKeyboard()
// 			b.Bot.Send(msg)
// 			return
// 		}
// 		ans := fmt.Sprintf(vocab.GetTranslate("Minter deposit and tag", b.Dlg.language), depos.Data.Address, depos.Data.Tag)
// 		msg := tgbotapi.NewMessage(b.Dlg.ChatId, ans)
// 		b.Dlg.Command = ""
// 		msg.ReplyMarkup = b.newMainKeyboard()
// 		b.Bot.Send(msg)
// 		go b.CheckStatusSell(depos.Data.Tag)
// 		return
// 	} else {
// 		price, err := strconv.ParseFloat(b.Dlg.Text, 64)
// 		if err != nil {
// 			fmt.Println(err)
// 			msg := tgbotapi.NewMessage(b.Dlg.ChatId, vocab.GetTranslate("Wrong price", b.Dlg.language))
// 			b.Bot.Send(msg)
// 			return
// 		}

// 		PriceToSell[b.Dlg.UserId] = price
// 		msg := tgbotapi.NewMessage(b.Dlg.ChatId, vocab.GetTranslate("Send BTC", b.Dlg.language))
// 		msg.ReplyMarkup = tgbotapi.ForceReply{
// 			ForceReply: true,
// 			Selective:  true,
// 		}
// 		b.Bot.Send(msg)
// 		return
// 	}
// }

// CheckStatusSell checks status of deposit for method Sell().
func (b *Bot) CheckStatusSell(tag string) {
	timeout := time.After(2 * time.Minute)
	tick := time.Tick(3 * time.Second)
	amount := "0"
	for {
		select {
		case <-timeout:
			if amount == "0" {
				msg := tgbotapi.NewMessage(b.Dlg.ChatId, vocab.GetTranslate("timeout", b.Dlg.language))
				msg.ReplyMarkup = b.newMainMenuKeyboard()
				b.Bot.Send(msg)
				return
			} else {
				continue
			}
		case <-tick:
			taginfo, err := b.Api.GetTagInfo(tag)
			if err != nil {
				fmt.Println(err)
				msg := tgbotapi.NewMessage(b.Dlg.ChatId, vocab.GetTranslate("Error", b.Dlg.language))
				b.Bot.Send(msg)
				return
			}
			if taginfo.Data.Amount != amount {
				amount = taginfo.Data.Amount
				// Put in DB.
				b.DB.PutLoot(b.Dlg.UserId, tag, taginfo)
				ans := fmt.Sprintf(vocab.GetTranslate("New deposit for sale", b.Dlg.language), taginfo.Data.Amount, taginfo.Data.Coin, taginfo.Data.Price)
				msg := tgbotapi.NewMessage(b.Dlg.ChatId, ans)
				b.Bot.Send(msg)
				//go a.CheckLootforSell(taginfo.Data.MinterAddress)
				return
			}

		}
	}
}

// Method for sending loots in markdown style to user.
func (b *Bot) ComposeResp(loots []*stct.Loot) {
	for _, loot := range loots {
		text := fmt.Sprintf(
			"*Tag:*  %s\n"+
				"*Coin:*  %s  "+
				"   *Price:*  %v\n"+
				"*Amount:*  %s\n"+
				"*Minted address:*  %s\n"+
				"*Created at:*  %s\n"+
				"*Last sell at:*  %s",
			loot.Tag,
			loot.Coin,
			loot.Price,
			loot.Amout,
			loot.MinterAddress,
			loot.CreatedAt.Format("2006-01-02 15:04:05"),
			loot.LastSell.Format("2006-01-02 15:04:05"))

		msg := tgbotapi.NewMessage(b.Dlg.ChatId, text)
		msg.ParseMode = "markdown"
		msg.ReplyMarkup = b.newMainKeyboard()
		b.Bot.Send(msg)
	}
}

// newMainMenuKeyboard is main menu keyboar: price, buy, sell, sales.
func (b *Bot) newMainMenuKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(vocab.GetTranslate("Price", b.Dlg.language), priceCommand),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(vocab.GetTranslate("Buy", b.Dlg.language), buyCommand),
			tgbotapi.NewInlineKeyboardButtonData(vocab.GetTranslate("Sell", b.Dlg.language), sellCommand),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(vocab.GetTranslate("Loots", b.Dlg.language), salesCommand),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(vocab.GetTranslate("Settings", b.Dlg.language), settingsMenu),
		),
	)
}

// vocabuageKeybord is keybouad for select vocabuage.
func (b *Bot) newVocabuageKeybord() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🇷🇺 Русский", rusvocabCommand),
			tgbotapi.NewInlineKeyboardButtonData("🇬🇧 English", engvocabCommand),
		),
	)
}

// newMainKeyboard is keyboard for main menu.
func (b *Bot) newMainKeyboard() tgbotapi.ReplyKeyboardMarkup {
	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("/getmenu"),
		),
	)

	keyboard.OneTimeKeyboard = true
	return keyboard
}
