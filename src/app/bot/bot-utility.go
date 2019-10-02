package bot

import (
	"fmt"
	"regexp"
	"strconv"
	"time"

	stct "github.com/mrKitikat/telegrambottest/src/app/bipdev/structs"
	vocab "github.com/mrKitikat/telegrambottest/src/app/bot/vocabulary"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
)

const (
	lootlink        = "https://bip.dev/trade/"
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

//
var (
	N         float64 = 0
	BuyStatus         = make(map[int64]string)
)

func (b *Bot) ChangeCurrency(ChatId int64, MessageId int, Id string) {
	timeout := time.After(10 * time.Minute)
	tick := time.Tick(20 * time.Second)
	for {
		select {
		case <-timeout:
			return

		case <-tick:
			fmt.Println(b.Dlg[ChatId].MessageId, MessageId)
			price, diff, err := b.Api.GetPrice()
			if err != nil {
				fmt.Println(err)
				msg := tgbotapi.NewMessage(b.Dlg[ChatId].ChatId, vocab.GetTranslate("Error", b.Dlg[ChatId].language))
				b.Bot.Send(msg)
				return
			}
			//PreviousMessage[ChatId] = msg
			if MessageId == b.Dlg[ChatId].MessageId {
				N++
				fmt.Println(b.Dlg[ChatId].CallBackId, Id)
				kb := b.newMainMenuKeyboard(ChatId)
				msg := tgbotapi.EditMessageTextConfig{
					BaseEdit: tgbotapi.BaseEdit{
						//ChatID: ChatId,
						//MessageID:       MessageId,
						InlineMessageID: Id,
						ReplyMarkup:     &kb,
					},
					Text:      fmt.Sprintf(vocab.GetTranslate("Select", b.Dlg[ChatId].language), price+N, diff),
					ParseMode: "markdown",
				}

				b.Bot.Send(msg)
				continue
			} else {
				return
			}

		}
	}
}

//
func (b *Bot) StepsZero(ChatId int64) {
	SellSteps[ChatId] = 0
	BuySteps[ChatId] = 0
}

//
func (b *Bot) SendBTCAddresses(ChatId int64) error {

	keyboard := tgbotapi.InlineKeyboardMarkup{}
	addresses, err := b.DB.GetBTCAddresses(b.Dlg[ChatId].UserId)
	if err != nil {
		return err
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
		b.Bot.Send(msg)
		Message[ChatId] = msg

	} else {
		txt := vocab.GetTranslate("New bitcoin", b.Dlg[ChatId].language)
		msg := tgbotapi.NewMessage(b.Dlg[ChatId].ChatId, txt)
		btn := tgbotapi.NewInlineKeyboardButtonData(vocab.GetTranslate("Cancel", b.Dlg[ChatId].language), cancelComm)
		var row []tgbotapi.InlineKeyboardButton
		row = append(row, btn)
		msg.ReplyMarkup = keyboard
		msg.ParseMode = "markdown"
		msg.ReplyMarkup = keyboard
		b.Bot.Send(msg)
		Message[ChatId] = msg
	}

	return nil
}

// SendMinterAddresses --
func (b *Bot) SendMinterAddresses(ChatId int64) error {

	//PreviousMessage[ChatId] = Message[ChatId]
	keyboard := tgbotapi.InlineKeyboardMarkup{}
	addresses, err := b.DB.GetMinterAddresses(b.Dlg[ChatId].UserId)
	if err != nil {
		return err
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
		b.Bot.Send(msg)
		Message[ChatId] = msg

	} else {
		txt := vocab.GetTranslate("New minter", b.Dlg[ChatId].language)
		msg := tgbotapi.NewMessage(b.Dlg[ChatId].ChatId, txt)
		btn := tgbotapi.NewInlineKeyboardButtonData(vocab.GetTranslate("Cancel", b.Dlg[ChatId].language), cancelComm)
		var row []tgbotapi.InlineKeyboardButton
		row = append(row, btn)
		msg.ReplyMarkup = keyboard
		msg.ParseMode = "markdown"
		b.Bot.Send(msg)
		Message[ChatId] = msg
	}

	return nil
}

// SendEmail --
func (b *Bot) SendEmail(ChatId int64) error {

	PreviousMessage[ChatId] = Message[ChatId]
	keyboard := tgbotapi.InlineKeyboardMarkup{}
	addresses, err := b.DB.GetEmails(b.Dlg[ChatId].UserId)
	if err != nil {
		return err
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
		b.Bot.Send(msg)

	} else {
		txt := vocab.GetTranslate("New email", b.Dlg[ChatId].language)
		msg := tgbotapi.NewMessage(b.Dlg[ChatId].ChatId, txt)
		btn := tgbotapi.NewInlineKeyboardButtonData(vocab.GetTranslate("Cancel", b.Dlg[ChatId].language), cancelComm)
		var row []tgbotapi.InlineKeyboardButton
		row = append(row, btn)
		msg.ReplyMarkup = keyboard
		msg.ParseMode = "markdown"
		b.Bot.Send(msg)
		Message[ChatId] = msg
	}

	return nil
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

//
func (b *Bot) GetStatusBuy(ChatId int64) string {
	return BuyStatus[ChatId]
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
					msg.ReplyMarkup = b.newMainKeyboard(ChatId)
					BuyStatus[ChatId] = vocab.GetTranslate("No buy", b.Dlg[ChatId].language)
					b.Bot.Send(msg)
					b.SendMenuMessage(ChatId)
					return
				}
			}
		}
	}
}

// CheckStatusSell checks status of deposit for method Sell().
func (b *Bot) CheckStatusSell(tag string, ChatId int64) {
	timeout := time.After(30 * time.Minute)
	tick := time.Tick(5 * time.Second)
	amount := "0.0"
	for {
		select {
		case <-timeout:
			if amount == "0.0" {
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

// Method for sending loots in markdown style to user.
func (b *Bot) ComposeResp(loots []*stct.Loot, ChatId int64) {

	keyboard := tgbotapi.InlineKeyboardMarkup{}

	for _, loot := range loots {
		var row []tgbotapi.InlineKeyboardButton
		lText := fmt.Sprintf(vocab.GetTranslate("Loot", b.Dlg[ChatId].language), loot.Amout, loot.Coin, loot.Price)
		btn := tgbotapi.NewInlineKeyboardButtonURL(lText, lootlink+loot.Tag)
		row = append(row, btn)
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
	}

	var row []tgbotapi.InlineKeyboardButton
	msg := tgbotapi.NewMessage(b.Dlg[ChatId].ChatId, vocab.GetTranslate("Your loots", b.Dlg[ChatId].language))
	btn := tgbotapi.NewInlineKeyboardButtonData(vocab.GetTranslate("Cancel", b.Dlg[ChatId].language), cancelComm)
	row = append(row, btn)
	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
	msg.ReplyMarkup = keyboard
	b.Bot.Send(msg)
}

func (b *Bot) CheckMinter(address string) bool {
	return len(address) != 42 || address[:2] != "Mx" || address == "Mx00000000000000000000000000000000000001"
}

// CheckEmail ..
func (b *Bot) CheckEmail(email string) bool {
	re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	if !re.MatchString(email) || email == "mail@example.com" {
		return false
	}
	return true
}

// CheckBitcoin ..
func (b *Bot) CheckBitcoin(address string) bool {
	re := regexp.MustCompile("^[13][a-km-zA-HJ-NP-Z1-9]{25,34}$")
	return re.MatchString(address)
}

// CheckCoin ..
func (b *Bot) CheckCoin(coin string) bool {
	re := regexp.MustCompile("^[0-9-A-Z]{3,10}$")
	return re.MatchString(coin)
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

// func (b *Bot) GetPrice(ChatId int64) tgbotapi.InlineKeyboardMarkup {

// 	keyboard := tgbotapi.InlineKeyboardMarkup{}
// 	prices, err := b.Api.GetAvailablePrices()
// 	if err != nil {
// 		fmt.Println(err)
// 		return keyboard
// 	}
// 	if len(prices) > 0 {
// 		for _, price := range prices {
// 			var row []tgbotapi.InlineKeyboardButton
// 			btn := tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("%.4f $", price), sendPrice+fmt.Sprintf("%.4f", price))
// 			row = append(row, btn)
// 			keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
// 		}
// 	}

// 	return keyboard
// }
