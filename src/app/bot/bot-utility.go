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
func (b *Bot) GetBTCAddresses(ChatId int64) tgbotapi.InlineKeyboardMarkup {

	keyboard := tgbotapi.InlineKeyboardMarkup{}
	var row []tgbotapi.InlineKeyboardButton
	btn := tgbotapi.NewInlineKeyboardButtonData(vocab.GetTranslate("New BTC", b.Dlg[ChatId].language), newBTC)
	row = append(row, btn)
	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
	addresses, err := b.DB.GetBTCAddresses(b.Dlg[ChatId].UserId)
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
func (b *Bot) GetMinterAddresses(ChatId int64) tgbotapi.InlineKeyboardMarkup {

	keyboard := tgbotapi.InlineKeyboardMarkup{}
	var row []tgbotapi.InlineKeyboardButton
	btn := tgbotapi.NewInlineKeyboardButtonData(vocab.GetTranslate("New minter", b.Dlg[ChatId].language), newMinter)
	row = append(row, btn)
	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
	addresses, err := b.DB.GetMinterAddresses(b.Dlg[ChatId].UserId)
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
func (b *Bot) GetEmail(ChatId int64) tgbotapi.InlineKeyboardMarkup {

	keyboard := tgbotapi.InlineKeyboardMarkup{}
	var row []tgbotapi.InlineKeyboardButton
	btn := tgbotapi.NewInlineKeyboardButtonData(vocab.GetTranslate("New email", b.Dlg[ChatId].language), newEmail)
	row = append(row, btn)
	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
	addresses, err := b.DB.GetEmails(b.Dlg[ChatId].UserId)
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

func (b *Bot) GetPrice(ChatId int64) tgbotapi.InlineKeyboardMarkup {

	keyboard := tgbotapi.InlineKeyboardMarkup{}
	prices, err := b.Api.GetAvailablePrices()
	if err != nil {
		fmt.Println(err)
		return keyboard
	}
	if len(prices) > 0 {
		for _, price := range prices {
			var row []tgbotapi.InlineKeyboardButton
			btn := tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("%.4f $", price), sendPrice+fmt.Sprintf("%.4f", price))
			row = append(row, btn)
			keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
		}
	}

	return keyboard
}

// SendMenu edit message and send Inline Keyboard newMainMenuKeyboard()
func (b *Bot) SendMenu(ChatId int64) {

	kb := b.newMainMenuKeyboard(ChatId)
	newmsg := tgbotapi.EditMessageTextConfig{
		BaseEdit: tgbotapi.BaseEdit{
			ChatID:      b.Dlg[ChatId].ChatId,
			MessageID:   b.Dlg[ChatId].MessageId,
			ReplyMarkup: &kb,
		},
		Text: vocab.GetTranslate("Select", b.Dlg[ChatId].language),
	}

	b.Bot.Send(newmsg)
}

// CheckStatusBuy checks depos BTC and wait 2 confirmations
func (b *Bot) CheckStatusBuy(address string, ChatId int64) {
	timeout := time.After(2 * time.Minute)
	tick := time.Tick(5 * time.Second)
	willcoin := 0.
	for {
		select {
		case <-timeout:
			if willcoin == 0. {
				msg := tgbotapi.NewMessage(b.Dlg[ChatId].ChatId, vocab.GetTranslate("timeout", b.Dlg[ChatId].language))
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
				time.Sleep(10 * time.Second)
				continue
			}
			if stat.Data.WillReceive != willcoin {
				if willcoin == 0. {
					willcoin = stat.Data.WillReceive
					ans := fmt.Sprintf(vocab.GetTranslate("New deposit", b.Dlg[ChatId].language), stat.Data.WillReceive)
					msg := tgbotapi.NewMessage(b.Dlg[ChatId].ChatId, ans)
					b.Bot.Send(msg)
					time.Sleep(60 * time.Second)
				} else {
					ans := fmt.Sprintf(vocab.GetTranslate("Exchange is successful", b.Dlg[ChatId].language), willcoin)
					msg := tgbotapi.NewMessage(b.Dlg[ChatId].ChatId, ans)
					msg.ReplyMarkup = b.newMainKeyboard()
					b.Bot.Send(msg)
					return
				}
			}
		}
	}
}

// CheckStatusSell checks status of deposit for method Sell().
func (b *Bot) CheckStatusSell(tag string, ChatId int64) {
	timeout := time.After(2 * time.Minute)
	tick := time.Tick(5 * time.Second)
	amount := "0"
	for {
		select {
		case <-timeout:
			if amount == "0" {
				msg := tgbotapi.NewMessage(b.Dlg[ChatId].ChatId, vocab.GetTranslate("timeout", b.Dlg[ChatId].language))
				msg.ReplyMarkup = b.newMainKeyboard()
				b.Bot.Send(msg)
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
					taginfo.Data.Amount, taginfo.Data.Coin, float64(float64(taginfo.Data.Price)/10000.))
				msg := tgbotapi.NewMessage(b.Dlg[ChatId].ChatId, ans)
				b.Bot.Send(msg)
				//go a.CheckLootforSell(taginfo.Data.MinterAddress)
				return
			}

		}
	}
}

// Method for sending loots in markdown style to user.
func (b *Bot) ComposeResp(loots []*stct.Loot, ChatId int64) {
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

		msg := tgbotapi.NewMessage(b.Dlg[ChatId].ChatId, text)
		msg.ParseMode = "markdown"
		msg.ReplyMarkup = b.newMainKeyboard()
		b.Bot.Send(msg)
	}
}

// newMainMenuKeyboard is main menu keyboar: price, buy, sell, sales.
func (b *Bot) newMainMenuKeyboard(ChatId int64) tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(vocab.GetTranslate("Price", b.Dlg[ChatId].language), priceCommand),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(vocab.GetTranslate("Buy", b.Dlg[ChatId].language), buyCommand),
			tgbotapi.NewInlineKeyboardButtonData(vocab.GetTranslate("Sell", b.Dlg[ChatId].language), sellCommand),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(vocab.GetTranslate("Loots", b.Dlg[ChatId].language), salesCommand),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(vocab.GetTranslate("Settings", b.Dlg[ChatId].language), settingsMenu),
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
