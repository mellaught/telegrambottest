package bot

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	api "github.com/mrKitikat/telegrambottest/src/app/bipdev"
	stct "github.com/mrKitikat/telegrambottest/src/app/bipdev/structs"
	vocab "github.com/mrKitikat/telegrambottest/src/app/bot/vocabulary"
	"github.com/mrKitikat/telegrambottest/src/app/db"

	//strt "bipbot/src/bipdev/structs"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
)

var (
	commands        = make(map[int64]string)
	CommandInfo     = make(map[int64]string)
	PreviousMessage = make(map[int64]tgbotapi.MessageConfig)
	Message         = make(map[int64]tgbotapi.MessageConfig)
	MinterAddress   = make(map[int64]string)
	CoinToSell      = make(map[int64]string)
	EmailAddress    = make(map[int64]string)
	PriceToSell     = make(map[int64]float64)
	BitcoinAddress  = make(map[int64]string)
	SellSteps       = make(map[int64]int)
	BuySteps        = make(map[int64]int)
)

// Dialog is struct for dialog with user:   - ChatId: User's ChatID
//											- UserId:   Struct App for Rest Api methods
//											- MessageId: Last Message id
//											- Text:   	Text of the last message from the user
//											- language: User's current language
//											- Command: Last command from user
type Dialog struct {
	ChatId     int64
	UserId     int
	CallBackId string
	MessageId  int
	Text       string
	language   string
	Command    string
}

// Bot is struct for Bot:   - Token: secret token from .env
//							- Api:   Struct App for Rest Api methods
//							- DB:    Postgres DB fro users and user's loots.
//							- Bot:	 tgbotapi Bot(token)
//							- Dlg:   For dialog struct
type Bot struct {
	Token string
	Api   *api.App
	DB    *db.DataBase
	Bot   *tgbotapi.BotAPI
	Dlg   map[int64]*Dialog
}

//InitBot initialization: loading the database, creating a bot by token from the config.
func InitBot(config *stct.Config, dbsql *sql.DB) *Bot {

	b := Bot{
		Token: config.BotToken,
		DB:    &db.DataBase{},
		Dlg:   map[int64]*Dialog{},
	}

	// Create table if not exists
	db, err := db.InitDB(dbsql)
	if err != nil {
		log.Fatal(err)
	}

	b.DB = db
	// Define URL
	b.Api = api.InitApp(config.BipdevApiHost)
	// Create new bot
	bot, err := tgbotapi.NewBotAPI(b.Token)
	if err != nil {
		log.Fatal(err)
	}

	b.Bot = bot
	go b.Run()

	return &b
}

// Run is starting bot.
func (b *Bot) Run() {

	//Set update timeout
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	//Get updates from bot
	updates, _ := b.Bot.GetUpdatesChan(u)
	for update := range updates {

		if update.Message == nil && update.CallbackQuery == nil {
			continue
		}

		dialog, exist := b.assembleUpdate(update)
		if !exist {
			continue
		}

		b.Dlg[dialog.ChatId] = dialog

		if botCommand := b.getCommand(update); botCommand != "" {
			b.RunCommand(botCommand, dialog.ChatId)
			continue
		}

		b.RunMessageText(update.Message.Text, dialog.ChatId)
		continue
	}
}

// RunMessageText
func (b *Bot) RunMessageText(text string, ChatId int64) {
	fmt.Printf("Buy step: %d, Sell step: %d\n", BuySteps[ChatId], SellSteps[ChatId])
	// Проверка команды <<купить>>.
	switch BuySteps[ChatId] {
	// команда не выбрана.
	case 1:
		// Проверка минтер адреса.
		if b.CheckMinter(text) {
			msg := tgbotapi.NewMessage(b.Dlg[ChatId].ChatId, vocab.GetTranslate("Wrong minter", b.Dlg[ChatId].language))
			msg.ParseMode = "markdown"
			b.Bot.Send(msg)
			return
		} else {
			MinterAddress[ChatId] = text
			// Отправьте почту.
			BuySteps[ChatId] = 2
			err := b.SendEmail(ChatId)
			if err != nil {
				fmt.Println(err)
			}
			return
		}
	case 2:
		// Проверка почты.
		if !b.CheckEmail(text) {
			msg := tgbotapi.NewMessage(b.Dlg[ChatId].ChatId, vocab.GetTranslate("Wrong email", b.Dlg[ChatId].language))
			msg.ParseMode = "markdown"
			b.Bot.Send(msg)
			return
		} else {
			EmailAddress[ChatId] = text
			// Отправьте депозит на биткоин адрес.
			b.StepsZero(ChatId)
			b.SendDepos(ChatId)
			b.BuyFinal(ChatId)
		}
	}

	// Проверка команды <<продать>>.
	switch SellSteps[ChatId] {
	// команда не выбрана.
	case 1:
		// Проверка названия монеты.
		if !b.CheckCoin(text) {
			msg := tgbotapi.NewMessage(b.Dlg[ChatId].ChatId, vocab.GetTranslate("Wrong coin name", b.Dlg[ChatId].language))
			msg.ParseMode = "markdown"
			b.Bot.Send(msg)
			return
		} else {
			// Отравьте цены за монеты.
			SellSteps[ChatId] = 2
			CoinToSell[ChatId] = text
			msg := tgbotapi.NewMessage(b.Dlg[ChatId].ChatId, vocab.GetTranslate("Select price", b.Dlg[ChatId].language))
			msg.ReplyMarkup = b.CancelKeyboard(ChatId)
			msg.ParseMode = "markdown"
			PreviousMessage[ChatId] = Message[ChatId]
			b.Bot.Send(msg)
			Message[ChatId] = msg
			return
		}
	case 2:
		// Проверка цены за монеты.
		if !b.CheckPrice(ChatId, text) {
			msg := tgbotapi.NewMessage(b.Dlg[ChatId].ChatId, vocab.GetTranslate("Wrong price", b.Dlg[ChatId].language))
			msg.ParseMode = "markdown"
			b.Bot.Send(msg)
			return
		} else {
			// Отправьте биткоин адрес.
			SellSteps[ChatId] = 3
			err := b.SendBTCAddresses(ChatId)
			if err != nil {
				fmt.Println(err)
			}
			return
		}
	case 3:
		// Проверка биткоин адреса.
		if !b.CheckBitcoin(text) {
			msg := tgbotapi.NewMessage(b.Dlg[ChatId].ChatId, vocab.GetTranslate("Wrong bitcoin", b.Dlg[ChatId].language))
			msg.ParseMode = "markdown"
			b.Bot.Send(msg)
			return
		} else {
			// Сохранить ли введенные данные?
			// Отправьте монеты на адрес.
			BitcoinAddress[ChatId] = text
			b.SellFinal(ChatId)
			return
		}
	}
	return
}

// assembleUpdate
func (b *Bot) assembleUpdate(update tgbotapi.Update) (*Dialog, bool) {
	dialog := &Dialog{}
	if update.Message != nil {
		dialog.language = b.DB.GetLanguage(update.Message.Chat.ID)
		dialog.ChatId = update.Message.Chat.ID
		dialog.MessageId = update.Message.MessageID
		dialog.UserId = int(update.Message.Chat.ID)
		dialog.Text = update.Message.Text
	} else if update.CallbackQuery != nil {
		dialog.CallBackId = update.CallbackQuery.ID
		dialog.language = b.DB.GetLanguage(update.CallbackQuery.Message.Chat.ID)
		dialog.ChatId = update.CallbackQuery.Message.Chat.ID
		dialog.MessageId = update.CallbackQuery.Message.MessageID
		dialog.UserId = int(update.CallbackQuery.Message.Chat.ID)
		if strings.Contains(update.CallbackQuery.Data, sendPrice) {
			dialog.Text = update.CallbackQuery.Data[9:]
			update.CallbackQuery.Data = update.CallbackQuery.Data[:9]
		} else if strings.Contains(update.CallbackQuery.Data, sendBTC) {
			dialog.Text = update.CallbackQuery.Data[7:]
			update.CallbackQuery.Data = update.CallbackQuery.Data[:7]
		} else if strings.Contains(update.CallbackQuery.Data, sendMinter) {
			dialog.Text = update.CallbackQuery.Data[10:]
			update.CallbackQuery.Data = update.CallbackQuery.Data[:10]
		} else if strings.Contains(update.CallbackQuery.Data, sendEmail) {
			dialog.Text = update.CallbackQuery.Data[9:]
			update.CallbackQuery.Data = update.CallbackQuery.Data[:9]
		}
	} else {
		dialog.language = "en"
		return dialog, false
	}

	command, isset := commands[dialog.ChatId]
	if isset {
		dialog.Command = command
	} else {
		dialog.Command = ""
	}

	return dialog, true
}

// getCommand returns command from telegram update
func (b *Bot) getCommand(update tgbotapi.Update) string {
	if update.Message != nil {
		if update.Message.IsCommand() {
			return update.Message.Command()
		}
	} else if update.CallbackQuery != nil {
		return update.CallbackQuery.Data
	}

	return ""
}

// RunCommand executes the input command.
func (b *Bot) RunCommand(command string, ChatId int64) {
	commands[ChatId] = command
	switch command {
	// "/Start" interacting with the bot, bot description and available commands.
	case startCommand:
		b.StepsZero(ChatId)
		msg := tgbotapi.NewMessage(b.Dlg[ChatId].ChatId, vocab.GetTranslate("Hello", b.Dlg[ChatId].language))
		msg.ParseMode = "markdown"
		msg.ReplyMarkup = b.newVocabuageKeybord()
		b.Bot.Send(msg)
		return

	// engvocabCommand sets english lang for user.
	case engvocabCommand:
		b.StepsZero(ChatId)
		b.DB.SetLanguage(b.Dlg[ChatId].UserId, "en")
		b.Dlg[ChatId].language = "en"
		b.SendMenuMessage(ChatId)
		commands[ChatId] = ""
		return

	// rusvocabCommand sets russian lang for user.
	case rusvocabCommand:
		b.StepsZero(ChatId)
		b.DB.SetLanguage(b.Dlg[ChatId].UserId, "ru")
		b.Dlg[ChatId].language = "ru"
		b.SendMenuMessage(ChatId)
		commands[ChatId] = ""
		return

	case cancelComm:
		if msg, ok := PreviousMessage[ChatId]; ok {
			if BuySteps[ChatId] == 1 || SellSteps[ChatId] == 1 {
				b.StepsZero(ChatId)
				b.SendMenuMessage(ChatId)
			} else if BuySteps[ChatId] != 0 {
				BuySteps[ChatId]--
				b.Bot.Send(msg)
			} else if SellSteps[ChatId] != 0 {
				SellSteps[ChatId]--
				b.Bot.Send(msg)
			} else {
				b.SendMenuMessage(ChatId)
			}
		} else {
			if SellSteps[ChatId] != 0 {
				SellSteps[ChatId]--
			} else if BuySteps[ChatId] != 0 {
				BuySteps[ChatId]--
			}

			b.SendMenuMessage(ChatId)
		}

	// Save user's data in DataBase.
	case yescommand:
		err := b.DB.PutMinterAddress(b.Dlg[ChatId].UserId, CommandInfo[ChatId])
		if err != nil {
			fmt.Println(err)
		}
		err = b.DB.PutEmail(b.Dlg[ChatId].UserId, b.Dlg[ChatId].Text)
		if err != nil {
			fmt.Println(err)
		}
		err = b.DB.PutBTCAddress(b.Dlg[ChatId].UserId, BitcoinAddress[ChatId])
		if err != nil {
			fmt.Println(err)
		}
		return

	// Don't save data in DataBase.
	case nocommand:
		return

	// Returns status of buy operation:
	// 1. Ожидание транзакции BTC…
	// 2. BTC уже в пути, вы получите как минимум xxx.xx BIP.
	case checkcommand:
		b.Bot.AnswerCallbackQuery(tgbotapi.NewCallbackWithAlert(b.Dlg[ChatId].CallBackId, b.GetStatusBuy(ChatId)))
	// buyCommand collects data from the user to transmit their request.
	// The user will receive the address for the deposit.
	// After he sends the money he will receive a notification from bot.
	// After the money is confirmed, he will receive another notification from bot.
	// ( BUY )
	case buyCommand:
		b.StepsZero(ChatId)
		_, ok := PreviousMessage[ChatId]
		if ok {
			delete(PreviousMessage, ChatId)
		}
		err := b.SendMinterAddresses(ChatId)
		if err != nil {
			fmt.Println(err)
		}
		BuySteps[ChatId] = 1
		return

	//sendMinter after the user has selected the minter address from the proposed. ( BUY )
	case sendMinter:
		BuySteps[ChatId] = 2
		MinterAddress[ChatId] = b.Dlg[ChatId].Text
		err := b.SendEmail(ChatId)
		if err != nil {
			fmt.Println(err)
		}
		return

	// sendEmail after the user has selected email from the proposed. ( BUY )
	case sendEmail:
		b.StepsZero(ChatId)
		EmailAddress[ChatId] = b.Dlg[ChatId].Text
		b.SendMenuChoose(ChatId)
		b.SendDepos(ChatId)
		b.BuyFinal(ChatId)

	// sellCommand collects data from the user to transmit their request. ( SELL )
	case sellCommand:
		b.StepsZero(ChatId)
		_, ok := PreviousMessage[ChatId]
		if ok {
			delete(PreviousMessage, ChatId)
		}
		SellSteps[ChatId] = 1
		msg := tgbotapi.NewMessage(b.Dlg[ChatId].ChatId, vocab.GetTranslate("Coin", b.Dlg[ChatId].language))
		msg.ReplyMarkup = b.CancelKeyboard(ChatId)
		msg.ParseMode = "markdown"
		Message[ChatId] = msg
		b.Bot.Send(msg)

	// sendBTC after the user has sent his bitcoin address. ( SELL )
	case sendBTC:
		b.StepsZero(ChatId)
		BitcoinAddress[ChatId] = b.Dlg[ChatId].Text
		b.SendMenuChoose(ChatId)
		b.SellFinal(ChatId)

	// salesCommand sends a request to the database to get user's loots. ( LOOTS )
	case salesCommand:
		loots, err := b.DB.GetLoots(b.Dlg[ChatId].UserId)
		if err != nil {
			fmt.Println(err)
			msg := tgbotapi.NewMessage(b.Dlg[ChatId].ChatId, vocab.GetTranslate("Error", b.Dlg[ChatId].language))
			b.Bot.Send(msg)
			return
		} else if len(loots) == 0 {
			msg := tgbotapi.NewMessage(b.Dlg[ChatId].ChatId, vocab.GetTranslate("Empty loots", b.Dlg[ChatId].language))
			msg.ReplyMarkup = b.newMainKeyboard(ChatId)
			b.Bot.Send(msg)
			return
		}

		b.ComposeResp(loots, ChatId)
		return
	}
}

// BuyFinal is function for command "/buy".
// Requests an email from the user and Minter deposit address.
// Requests the "bitcoinDepositAddress" method with the received data.
func (b *Bot) BuyFinal(ChatId int64) {
	fmt.Println("Buy data:", MinterAddress[b.Dlg[ChatId].ChatId], EmailAddress[b.Dlg[ChatId].ChatId])
	addr, err := b.Api.GetBTCDeposAddress(MinterAddress[b.Dlg[ChatId].ChatId], "BIP", EmailAddress[b.Dlg[ChatId].ChatId])
	if err != nil {
		b.Dlg[ChatId].Command = ""
		msg := tgbotapi.NewMessage(b.Dlg[ChatId].ChatId, err.Error())
		msg.ReplyMarkup = b.newMainKeyboard(ChatId)
		b.Bot.Send(msg)
		return
	}
	b.SendMenuChoose(ChatId)
	b.Dlg[ChatId].Command = ""
	newmsg := tgbotapi.NewMessage(b.Dlg[ChatId].ChatId, addr)
	newmsg.ReplyMarkup = b.CheckKeyboard(ChatId)
	b.Bot.Send(newmsg)
	go b.CheckStatusBuy(addr, ChatId)
	return
}

// SellFinal
func (b *Bot) SellFinal(ChatId int64) {
	fmt.Println("Sell data:", BitcoinAddress[b.Dlg[ChatId].ChatId], CoinToSell[b.Dlg[ChatId].ChatId], PriceToSell[b.Dlg[ChatId].ChatId])
	depos, err := b.Api.GetMinterDeposAddress(BitcoinAddress[b.Dlg[ChatId].ChatId], CoinToSell[b.Dlg[ChatId].ChatId], PriceToSell[b.Dlg[ChatId].ChatId])
	if err != nil {
		msg := tgbotapi.NewMessage(b.Dlg[ChatId].ChatId, err.Error())
		b.Bot.Send(msg)
		b.SendMenuMessage(ChatId)
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
	time.Sleep(5 * time.Second)
	b.SendMenuMessage(ChatId)
	return
}

// settingsMenu return Inline KeyBoard newVocabuageKeybord to select a language.
// case settingsMenu:
// 	kb := b.newVocabuageKeybord()
// 	msg := tgbotapi.EditMessageTextConfig{
// 		BaseEdit: tgbotapi.BaseEdit{
// 			ChatID:      b.Dlg[ChatId].ChatId,
// 			MessageID:   b.Dlg[ChatId].MessageId,
// 			ReplyMarkup: &kb,
// 		},
// 		Text: vocab.GetTranslate("Settings", b.Dlg[ChatId].language),
// 	}

// 	b.Bot.Send(msg)

// priceCommand requests the server for the current BIP / USD rate and sends a message to user with the server responce. ( PRICE )
// case priceCommand:
// 	price, err := b.Api.GetPrice()
// 	if err != nil {
// 		fmt.Println(err)
// 		msg := tgbotapi.NewMessage(b.Dlg[ChatId].ChatId, vocab.GetTranslate("Error", b.Dlg[ChatId].language))
// 		b.Bot.Send(msg)
// 		return
// 	}
// 	msg := tgbotapi.EditMessageTextConfig{
// 		BaseEdit: tgbotapi.BaseEdit{
// 			ChatID:    b.Dlg[ChatId].ChatId,
// 			MessageID: b.Dlg[ChatId].MessageId,
// 		},
// 		Text: vocab.GetTranslate("Now", b.Dlg[ChatId].language),
// 	}
// 	b.Bot.Send(msg)
// 	newmsg := tgbotapi.NewMessage(ChatId, fmt.Sprintf("%.4f $", price))
// 	newmsg.ReplyMarkup = b.newMainKeyboard(ChatId)
// 	b.Bot.Send(newmsg)
// ans := fmt.Sprintf(vocab.GetTranslate("Now", b.Dlg[ChatId].language), price)
// msg := tgbotapi.NewMessage(b.Dlg[ChatId].ChatId, ans)
// msg.ReplyMarkup = b.newMainKeyboard()
// b.Bot.Send(msg)

// buyCommand collects data from the user to transmit their request.
// The user will receive the address for the deposit.
// After he sends the money he will receive a notification from bot.
// After the money is confirmed, he will receive another notification from bot.
// ( BUY )

// newMinter after the user decided to enter a new minter address. ( BUY )
// case newMinter:
// 	msg := tgbotapi.NewMessage(b.Dlg[ChatId].ChatId, vocab.GetTranslate("Send minter", b.Dlg[ChatId].language))
// 	msg.ReplyMarkup = tgbotapi.ForceReply{
// 		ForceReply: true,
// 		Selective:  true,
// 	}
// 	b.Bot.Send(msg)

// newEmail after the user decided to enter a new email. ( BUY )
// case newEmail:
// 	msg := tgbotapi.NewMessage(b.Dlg[ChatId].ChatId, vocab.GetTranslate("Email", b.Dlg[ChatId].language))
// 	msg.ReplyMarkup = tgbotapi.ForceReply{
// 		ForceReply: true,
// 		Selective:  true,
// 	}
// 	b.Bot.Send(msg)

// getMainMenu return Inline Keyboard newMainMenuKeyboard()
// case getMainMenu:
// 	msg := tgbotapi.NewMessage(b.Dlg[ChatId].ChatId, vocab.GetTranslate("Select", b.Dlg[ChatId].language))
// 	msg.ReplyMarkup = b.newMainMenuKeyboard(ChatId)
// 	b.Bot.Send(msg)
// 	return
//fmt.Println("Messasge id: ", b.Dlg[ChatId].MessageId)

// msg := tgbotapi.NewMessage(b.Dlg[ChatId].ChatId, vocab.GetTranslate("Select", b.Dlg[ChatId].language))
// msg.ReplyMarkup = b.newMainMenuKeyboard(ChatId)
// b.Bot.Send(msg)

// newBTC after the user decided to enter a new bitcoin address. ( SELL )
// case newBTC:
// 	msg := tgbotapi.NewMessage(b.Dlg[ChatId].ChatId, vocab.GetTranslate("Send BTC", b.Dlg[ChatId].language))
// 	msg.ReplyMarkup = tgbotapi.ForceReply{
// 		ForceReply: true,
// 		Selective:  true,
// 	}
// 	b.Bot.Send(msg)

// func (b *Bot) CoinName(ChatId int64) {
// 	re := regexp.MustCompile("^[0-9-A-Z]{3,10}$")
// 	if !re.MatchString(b.Dlg[ChatId].Text) {
// 		SellStep[ChatId] = 0
// 		msg := tgbotapi.NewMessage(b.Dlg[ChatId].ChatId, vocab.GetTranslate("Coin name", b.Dlg[ChatId].language))
// 		msg.ReplyMarkup = b.newMainKeyboard(ChatId)
// 		b.Bot.Send(msg)
// 		return
// 	}
// 	CoinToSell[b.Dlg[ChatId].UserId] = b.Dlg[ChatId].Text
// 	SellStep[ChatId] = 2
// 	fmt.Println(b.Dlg[ChatId].Text)
// 	msg := tgbotapi.NewMessage(b.Dlg[ChatId].ChatId, vocab.GetTranslate("Coin price", b.Dlg[ChatId].language))
// 	msg.ReplyMarkup = tgbotapi.ForceReply{
// 		ForceReply: true,
// 		Selective:  true,
// 	}
// 	b.Bot.Send(msg)
// }

// msg := tgbotapi.EditMessageTextConfig{
// 	BaseEdit: tgbotapi.BaseEdit{
// 		ChatID:      b.Dlg[ChatId].ChatId,
// 		MessageID:   b.Dlg[ChatId].MessageId,
// 		ReplyMarkup: &kb,
// 	},
// 	Text: vocab.GetTranslate("Select minter", b.Dlg[ChatId].language),
// }

///b.Bot.Send(msg)

// // BuySecondStep
// func (b *Bot) BuySecondStep(ChatId int64) {
// 	CommandInfo[b.Dlg[ChatId].UserId] = b.Dlg[ChatId].Text
// 	msg := tgbotapi.NewMessage(b.Dlg[ChatId].ChatId, vocab.GetTranslate("Select email", b.Dlg[ChatId].language))
// 	msg.ReplyMarkup = b.GetEmail(ChatId)
// 	b.Bot.Send(msg)
// }

// if update.Message != nil && update.Message.ReplyToMessage != nil {
// 	// Buy command
// 	// SECOND STEP
// 	if dialog.Command == "newMinter" {
// 		b.BuySecondStep(update.Message.Chat.ID)
// 		continue
// 	} else if dialog.Command == "newEmail" {
// 		b.BuyFinal(update.Message.Chat.ID)
// 		continue
// 	} else if dialog.Command == "sell" {
// 		if SellStep[update.Message.Chat.ID] == 1 {
// 			b.CoinName(update.Message.Chat.ID)
// 			continue
// 		} else {
// 			b.SellSecondStep(update.Message.Chat.ID)
// 			continue
// 		}
// 	} else if dialog.Command == "newBTC" {
// 		b.SellFinal(update.Message.Chat.ID)
// 		continue
// 	}
// }

// SellSecondStep after the user has sended a price for his coin. ( SELL )
// func (b *Bot) SellSecondStep(ChatId int64) {
// 	PriceToSell[b.Dlg[ChatId].ChatId] = b.Dlg[ChatId].Text
// 	msg := tgbotapi.NewMessage(b.Dlg[ChatId].ChatId, vocab.GetTranslate("Select bitcoin", b.Dlg[ChatId].language))
// 	msg.ReplyMarkup = b.SendBTCAddresses(ChatId)
// 	b.Bot.Send(msg)
// }
