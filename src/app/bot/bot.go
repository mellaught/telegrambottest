package bot

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"
	api "telegrambottest/src/app/bipdev"
	stct "telegrambottest/src/app/bipdev/structs"
	vocab "telegrambottest/src/app/bot/vocabulary"
	"telegrambottest/src/app/db"

	//strt "bipbot/src/bipdev/structs"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
)

var (
	commands    = make(map[int]string)
	CommandInfo = make(map[int]string)
	CoinToSell  = make(map[int]string)
	PriceToSell = make(map[int]string)
	Actions     = make(map[int]map[string]int)
)

// Dialog is struct for dialog with user:   - ChatId: User's ChatID
//											- UserId:   Struct App for Rest Api methods
//											- MessageId: Last Message id
//											- Text:   	Text of the last message from the user
//											- language: User's current language
//											- Command: Last command from user
type Dialog struct {
	ChatId    int64
	UserId    int
	MessageId int
	Text      string
	language  string
	Command   string
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
		Token: config.Token,
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
	b.Api = api.InitApp(config.URL)
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

		if update.Message != nil && update.Message.ReplyToMessage != nil {
			// Buy command
			// SECOND STEP
			if dialog.Command == "newMinter" {
				b.BuySecondStep(update.Message.Chat.ID)
				continue
			} else if dialog.Command == "newEmail" {
				b.BuyFinal(update.Message.Chat.ID)
				continue
			} else if dialog.Command == "sell" {
				b.CoinName(update.Message.Chat.ID)
				continue
			} else if dialog.Command == "newBTC" {
				b.SellFinal(update.Message.Chat.ID)
				continue
			}
		}

		if botCommand := b.getCommand(update); botCommand != "" {
			b.RunCommand(botCommand, dialog.ChatId)
			continue
		}

	}
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
		fmt.Println(update.CallbackQuery.Data)
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
			fmt.Println(update.CallbackQuery.Data[10:])
			fmt.Println(update.CallbackQuery.Data[:10])
			dialog.Text = update.CallbackQuery.Data[10:]
			update.CallbackQuery.Data = update.CallbackQuery.Data[:10]
		} else if strings.Contains(update.CallbackQuery.Data, sendEmail) {
			fmt.Println(update.CallbackQuery.Data[9:])
			fmt.Println(update.CallbackQuery.Data[:9])
			dialog.Text = update.CallbackQuery.Data[9:]
			update.CallbackQuery.Data = update.CallbackQuery.Data[:9]
		}
	} else {
		dialog.language = "en"
		return dialog, false
	}

	command, isset := commands[dialog.UserId]
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
	commands[b.Dlg[ChatId].UserId] = command
	switch command {

	// "/Start" interacting with the bot, bot description and available commands.
	case startCommand:
		msg := tgbotapi.NewMessage(b.Dlg[ChatId].ChatId, vocab.GetTranslate("Hello", b.Dlg[ChatId].language))
		msg.ReplyMarkup = b.newVocabuageKeybord()
		b.Bot.Send(msg)

	// settingsMenu return Inline KeyBoard newVocabuageKeybord to select a language.
	case settingsMenu:
		kb := b.newVocabuageKeybord()
		msg := tgbotapi.EditMessageTextConfig{
			BaseEdit: tgbotapi.BaseEdit{
				ChatID:      b.Dlg[ChatId].ChatId,
				MessageID:   b.Dlg[ChatId].MessageId,
				ReplyMarkup: &kb,
			},
			Text: vocab.GetTranslate("Settings", b.Dlg[ChatId].language),
		}

		b.Bot.Send(msg)

	// engvocabCommand sets english lang for user.
	case engvocabCommand:
		b.DB.SetLanguage(b.Dlg[ChatId].UserId, "en")
		b.Dlg[ChatId].language = "en"
		b.SendMenu(ChatId)

	// rusvocabCommand sets russian lang for user.
	case rusvocabCommand:
		b.DB.SetLanguage(b.Dlg[ChatId].UserId, "ru")
		b.Dlg[ChatId].language = "ru"
		b.SendMenu(ChatId)

	// priceCommand requests the server for the current BIP / USD rate and sends a message to user with the server responce. ( PRICE )
	case priceCommand:
		price, err := b.Api.GetPrice()
		if err != nil {
			fmt.Println(err)
			msg := tgbotapi.NewMessage(b.Dlg[ChatId].ChatId, vocab.GetTranslate("Error", b.Dlg[ChatId].language))
			b.Bot.Send(msg)
			return
		}
		ans := fmt.Sprintf(vocab.GetTranslate("Now", b.Dlg[ChatId].language), price)
		msg := tgbotapi.NewMessage(b.Dlg[ChatId].ChatId, ans)
		msg.ReplyMarkup = b.newMainKeyboard()
		b.Bot.Send(msg)

	// buyCommand collects data from the user to transmit their request.
	// The user will receive the address for the deposit.
	// After he sends the money he will receive a notification from bot.
	// After the money is confirmed, he will receive another notification from bot.
	// ( BUY )
	case buyCommand:
		kb := b.GetMinterAddresses(ChatId)
		msg := tgbotapi.EditMessageTextConfig{
			BaseEdit: tgbotapi.BaseEdit{
				ChatID:      b.Dlg[ChatId].ChatId,
				MessageID:   b.Dlg[ChatId].MessageId,
				ReplyMarkup: &kb,
			},
			Text: vocab.GetTranslate("Select minter", b.Dlg[ChatId].language),
		}
		b.Bot.Send(msg)

	// newMinter after the user decided to enter a new minter address. ( BUY )
	case newMinter:
		msg := tgbotapi.NewMessage(b.Dlg[ChatId].ChatId, vocab.GetTranslate("Send minter", b.Dlg[ChatId].language))
		msg.ReplyMarkup = tgbotapi.ForceReply{
			ForceReply: true,
			Selective:  true,
		}
		b.Bot.Send(msg)

	// sendMinter after the user has selected the minter address from the proposed. ( BUY )
	case sendMinter:
		CommandInfo[b.Dlg[ChatId].UserId] = b.Dlg[ChatId].Text
		kb := b.GetEmail(ChatId)
		msg := tgbotapi.EditMessageTextConfig{
			BaseEdit: tgbotapi.BaseEdit{
				ChatID:      b.Dlg[ChatId].ChatId,
				MessageID:   b.Dlg[ChatId].MessageId,
				ReplyMarkup: &kb,
			},
			Text: vocab.GetTranslate("Select email", b.Dlg[ChatId].language),
		}
		b.Bot.Send(msg)

	// newEmail after the user decided to enter a new email. ( BUY )
	case newEmail:
		msg := tgbotapi.NewMessage(b.Dlg[ChatId].ChatId, vocab.GetTranslate("Email", b.Dlg[ChatId].language))
		msg.ReplyMarkup = tgbotapi.ForceReply{
			ForceReply: true,
			Selective:  true,
		}
		b.Bot.Send(msg)

	// sendEmail after the user has selected email from the proposed. ( BUY )
	case sendEmail:
		b.BuyFinal(ChatId)

	// sellCommand collects data from the user to transmit their request. ( SELL )
	case sellCommand:
		msg := tgbotapi.NewMessage(b.Dlg[ChatId].ChatId, vocab.GetTranslate("Coin", b.Dlg[ChatId].language))
		msg.ReplyMarkup = tgbotapi.ForceReply{
			ForceReply: true,
			Selective:  true,
		}
		b.Bot.Send(msg)

	// sendPrice after the user has chosen a price for his coin. ( SELL )
	case sendPrice:
		PriceToSell[b.Dlg[ChatId].UserId] = b.Dlg[ChatId].Text
		kb := b.GetBTCAddresses(ChatId)
		msg := tgbotapi.EditMessageTextConfig{
			BaseEdit: tgbotapi.BaseEdit{
				ChatID:      b.Dlg[ChatId].ChatId,
				MessageID:   b.Dlg[ChatId].MessageId,
				ReplyMarkup: &kb,
			},
			Text: vocab.GetTranslate("Select bitcoin", b.Dlg[ChatId].language),
		}
		b.Bot.Send(msg)
	// newBTC after the user decided to enter a new bitcoin address. ( SELL )
	case newBTC:
		msg := tgbotapi.NewMessage(b.Dlg[ChatId].ChatId, vocab.GetTranslate("Send BTC", b.Dlg[ChatId].language))
		msg.ReplyMarkup = tgbotapi.ForceReply{
			ForceReply: true,
			Selective:  true,
		}
		b.Bot.Send(msg)

	// sendBTC after the user has sent his bitcoin address. ( SELL )
	case sendBTC:
		b.SellFinal(ChatId)
		// sendPrice

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
			msg.ReplyMarkup = b.newMainKeyboard()
			b.Bot.Send(msg)
			return
		}
		b.ComposeResp(loots, ChatId)

	// getMainMenu return Inline Keyboard newMainMenuKeyboard()
	case getMainMenu:
		msg := tgbotapi.NewMessage(b.Dlg[ChatId].ChatId, vocab.GetTranslate("Select", b.Dlg[ChatId].language))
		msg.ReplyMarkup = b.newMainMenuKeyboard(ChatId)
		b.Bot.Send(msg)
	}
}

// BuySecondStep
func (b *Bot) BuySecondStep(ChatId int64) {
	CommandInfo[b.Dlg[ChatId].UserId] = b.Dlg[ChatId].Text
	msg := tgbotapi.NewMessage(b.Dlg[ChatId].ChatId, vocab.GetTranslate("Select email", b.Dlg[ChatId].language))
	msg.ReplyMarkup = b.GetEmail(ChatId)
	b.Bot.Send(msg)
}

// BuyFinal is function for command "/buy".
// Requests an email from the user and Minter deposit address.
// Requests the "bitcoinDepositAddress" method with the received data.
func (b *Bot) BuyFinal(ChatId int64) {
	dialog := b.Dlg[ChatId]
	fmt.Println("In FINAL", CommandInfo[dialog.UserId], dialog.Text)
	addr, err := b.Api.GetBTCDeposAddress(CommandInfo[dialog.UserId], "BIP", dialog.Text)
	if err != nil {
		dialog.Command = ""
		msg := tgbotapi.NewMessage(dialog.ChatId, err.Error())
		msg.ReplyMarkup = b.newMainKeyboard()
		b.Bot.Send(msg)
		return
	}
	err = b.DB.PutMinterAddress(b.Dlg[ChatId].UserId, CommandInfo[b.Dlg[ChatId].UserId])
	if err != nil {
		fmt.Println(err)
	}
	err = b.DB.PutEmail(b.Dlg[ChatId].UserId, b.Dlg[ChatId].Text)
	if err != nil {
		fmt.Println(err)
	}
	b.Dlg[ChatId].Command = ""
	msg := tgbotapi.NewMessage(b.Dlg[ChatId].ChatId, vocab.GetTranslate("BTC deposit", b.Dlg[ChatId].language))
	b.Bot.Send(msg)
	newmsg := tgbotapi.NewMessage(b.Dlg[ChatId].ChatId, addr)
	newmsg.ReplyMarkup = b.newMainKeyboard()
	b.Bot.Send(newmsg)
	go b.CheckStatusBuy(addr, ChatId)
	return
}

func (b *Bot) CoinName(ChatId int64) {

	CoinToSell[b.Dlg[ChatId].UserId] = b.Dlg[ChatId].Text
	msg := tgbotapi.NewMessage(b.Dlg[ChatId].ChatId, vocab.GetTranslate("Select price", b.Dlg[ChatId].language))
	msg.ReplyMarkup = b.GetPrice(ChatId)

	b.Bot.Send(msg)
}

// SellFinal
func (b *Bot) SellFinal(ChatId int64) {

	price, err := strconv.ParseFloat(PriceToSell[b.Dlg[ChatId].UserId], 64)
	fmt.Println("Final sell:", price, b.Dlg[ChatId].Text, CoinToSell[b.Dlg[ChatId].UserId])
	if err != nil {
		fmt.Println(err)
		msg := tgbotapi.NewMessage(b.Dlg[ChatId].ChatId, vocab.GetTranslate("Error", b.Dlg[ChatId].language))
		msg.ReplyMarkup = b.newMainKeyboard()
		b.Bot.Send(msg)
		return
	}

	depos, err := b.Api.GetMinterDeposAddress(b.Dlg[ChatId].Text, CoinToSell[b.Dlg[ChatId].UserId], price)
	if err != nil {
		msg := tgbotapi.NewMessage(b.Dlg[ChatId].ChatId, err.Error())
		msg.ReplyMarkup = b.newMainKeyboard()
		b.Bot.Send(msg)
		return
	}
	err = b.DB.PutBTCAddress(b.Dlg[ChatId].UserId, b.Dlg[ChatId].Text)
	if err != nil {
		fmt.Println(err)
	}

	b.Dlg[ChatId].Command = ""
	msg := tgbotapi.NewMessage(b.Dlg[ChatId].ChatId, vocab.GetTranslate("Minter deposit and tag", b.Dlg[ChatId].language))
	b.Bot.Send(msg)
	newmsg := tgbotapi.NewMessage(b.Dlg[ChatId].ChatId, depos.Data.Address)
	b.Bot.Send(newmsg)
	newmsg2 := tgbotapi.NewMessage(b.Dlg[ChatId].ChatId, depos.Data.Tag)
	newmsg2.ReplyMarkup = b.newMainKeyboard()
	b.Bot.Send(newmsg2)
	go b.CheckStatusSell(depos.Data.Tag, ChatId)
	return
}
