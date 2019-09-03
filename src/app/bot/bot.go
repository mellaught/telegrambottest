package bot

import (
	"database/sql"
	"fmt"
	"log"
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
	PriceToSell = make(map[int]float64)
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
	Dlg   *Dialog
}

//InitBot initialization: loading the database, creating a bot by token from the config.
func InitBot(config *stct.Config, dbsql *sql.DB) *Bot {

	b := Bot{
		Token: config.Token,
		DB:    &db.DataBase{},
		Dlg:   &Dialog{},
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

		b.Dlg = dialog

		if update.Message != nil && update.Message.ReplyToMessage != nil {
			// Buy command
			// SECOND STEP
			if dialog.Command == "newMinter" {
				b.BuySecondStep()
				continue
			} else if dialog.Command == "newEmail" {
				b.BuyFinal()
				continue
			} else if dialog.Command == "newBTC" {
				b.SellSecondStep()
				continue
			} else if dialog.Command == "sell" {
				b.SellFinal()
				continue
			}
		}

		if botCommand := b.getCommand(update); botCommand != "" {
			b.RunCommand(botCommand)
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
		dialog.language = b.DB.GetLanguage(update.CallbackQuery.Message.Chat.ID)
		dialog.ChatId = update.CallbackQuery.Message.Chat.ID
		dialog.MessageId = update.CallbackQuery.Message.MessageID
		dialog.UserId = int(update.CallbackQuery.Message.Chat.ID)
		dialog.Text = update.CallbackQuery.Message.Text
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
func (b *Bot) RunCommand(command string) {
	commands[b.Dlg.UserId] = command
	switch command {

	// "/Start" interacting with the bot, bot description and available commands.
	case startCommand:
		msg := tgbotapi.NewMessage(b.Dlg.ChatId, vocab.GetTranslate("Hello", b.Dlg.language))
		msg.ReplyMarkup = b.newVocabuageKeybord()
		b.Bot.Send(msg)

	// settingsMenu return Inline KeyBoard newVocabuageKeybord to select a language.
	case settingsMenu:
		kb := b.newVocabuageKeybord()
		msg := tgbotapi.EditMessageTextConfig{
			BaseEdit: tgbotapi.BaseEdit{
				ChatID:      b.Dlg.ChatId,
				MessageID:   b.Dlg.MessageId,
				ReplyMarkup: &kb,
			},
			Text: vocab.GetTranslate("Settings", b.Dlg.language),
		}

		b.Bot.Send(msg)

	// engvocabCommand sets english lang for user.
	case engvocabCommand:
		b.DB.SetLanguage(b.Dlg.UserId, "en")
		b.Dlg.language = "en"
		b.SendMenu()

	// rusvocabCommand sets russian lang for user.
	case rusvocabCommand:
		b.DB.SetLanguage(b.Dlg.UserId, "ru")
		b.Dlg.language = "ru"
		b.SendMenu()

	// priceCommand requests the server for the current BIP / USD rate and sends a message to user with the server responce.
	case priceCommand:
		price, err := b.Api.GetPrice()
		if err != nil {
			fmt.Println(err)
			msg := tgbotapi.NewMessage(b.Dlg.ChatId, vocab.GetTranslate("Error", b.Dlg.language))
			b.Bot.Send(msg)
			return
		}
		ans := fmt.Sprintf(vocab.GetTranslate("Now", b.Dlg.language), price)
		msg := tgbotapi.NewMessage(b.Dlg.ChatId, ans)
		msg.ReplyMarkup = b.newMainKeyboard()
		b.Bot.Send(msg)

	// buyCommand collects data from the user to transmit their request.
	// The user will receive the address for the deposit.
	// After he sends the money he will receive a notification from bot.
	// After the money is confirmed, he will receive another notification from bot.
	// FIRST STEP
	case buyCommand:
		msg := tgbotapi.NewMessage(b.Dlg.ChatId, vocab.GetTranslate("Select minter", b.Dlg.language))
		msg.ReplyMarkup = b.GetMinterAddresses()
		b.Bot.Send(msg)

	// Buy command
	// FIRST STEP(IF new minter address)
	case newMinter:
		msg := tgbotapi.NewMessage(b.Dlg.ChatId, vocab.GetTranslate("Send minter", b.Dlg.language))
		msg.ReplyMarkup = tgbotapi.ForceReply{
			ForceReply: true,
			Selective:  true,
		}
		b.Bot.Send(msg)

	// Buy command
	// SECOND STEP
	case sendMinter:
		kb := b.GetEmail()
		msg := tgbotapi.EditMessageTextConfig{
			BaseEdit: tgbotapi.BaseEdit{
				ChatID:      b.Dlg.ChatId,
				MessageID:   b.Dlg.MessageId,
				ReplyMarkup: &kb,
			},
			Text: vocab.GetTranslate("Select email", b.Dlg.language),
		}
		b.Bot.Send(msg)
	// newEmail
	case newEmail:
		msg := tgbotapi.NewMessage(b.Dlg.ChatId, vocab.GetTranslate("Email", b.Dlg.language))
		msg.ReplyMarkup = tgbotapi.ForceReply{
			ForceReply: true,
			Selective:  true,
		}
		b.Bot.Send(msg)

	// sendEmail
	case sendEmail:
		b.BuyFinal()

	// sellCommand collects data from the user to transmit their request.
	case sellCommand:
		msg := tgbotapi.NewMessage(b.Dlg.ChatId, vocab.GetTranslate("Select bitcoin", b.Dlg.language))
		msg.ReplyMarkup = b.GetBTCAddresses()
		b.Bot.Send(msg)

	// salesCommand sends a request to the database to get user's loots.
	case newBTC:
		msg := tgbotapi.NewMessage(b.Dlg.ChatId, vocab.GetTranslate("Send BTC", b.Dlg.language))
		msg.ReplyMarkup = tgbotapi.ForceReply{
			ForceReply: true,
			Selective:  true,
		}
		b.Bot.Send(msg)

	// sendBTC
	case sendBTC:
		kb := b.GetPrice()
		msg := tgbotapi.EditMessageTextConfig{
			BaseEdit: tgbotapi.BaseEdit{
				ChatID:      b.Dlg.ChatId,
				MessageID:   b.Dlg.MessageId,
				ReplyMarkup: &kb,
			},
			Text: vocab.GetTranslate("Select price", b.Dlg.language),
		}
		b.Bot.Send(msg)

	// sendPrice
	case sendPrice:
		msg := tgbotapi.NewMessage(b.Dlg.ChatId, vocab.GetTranslate("Coin price", b.Dlg.language))
		msg.ReplyMarkup = tgbotapi.ForceReply{
			ForceReply: true,
			Selective:  true,
		}
		b.Bot.Send(msg)

	// salesCommand
	case salesCommand:
		loots, err := b.DB.GetLoots(b.Dlg.UserId)
		if err != nil {
			fmt.Println(err)
			msg := tgbotapi.NewMessage(b.Dlg.ChatId, vocab.GetTranslate("Error", b.Dlg.language))
			b.Bot.Send(msg)
			return
		} else if len(loots) == 0 {
			msg := tgbotapi.NewMessage(b.Dlg.ChatId, vocab.GetTranslate("Empty loots", b.Dlg.language))
			msg.ReplyMarkup = b.newMainKeyboard()
			b.Bot.Send(msg)
			return
		}
		b.ComposeResp(loots)

	// getMainMenu return Inline Keyboard newMainMenuKeyboard()
	case getMainMenu:
		msg := tgbotapi.NewMessage(b.Dlg.ChatId, vocab.GetTranslate("Select", b.Dlg.language))
		msg.ReplyMarkup = b.newMainMenuKeyboard()
		b.Bot.Send(msg)
	}
}

// BuySecondStep
func (b *Bot) BuySecondStep() {
	msg := tgbotapi.NewMessage(b.Dlg.ChatId, vocab.GetTranslate("Select email", b.Dlg.language))
	msg.ReplyMarkup = b.GetEmail()
	b.Bot.Send(msg)
}

// BuyFinal
func (b *Bot) BuyFinal() {

	addr, err := b.Api.GetBTCDeposAddress(CommandInfo[b.Dlg.UserId], "BIP", b.Dlg.Text)
	if err != nil {
		b.Dlg.Command = ""
		msg := tgbotapi.NewMessage(b.Dlg.ChatId, err.Error())
		msg.ReplyMarkup = b.newMainKeyboard()
		b.Bot.Send(msg)
		return
	}

	b.Dlg.Command = ""
	msg := tgbotapi.NewMessage(b.Dlg.ChatId, vocab.GetTranslate("BTC deposit", b.Dlg.language))
	b.Bot.Send(msg)
	msg.ReplyMarkup = b.newMainKeyboard()
	msg = tgbotapi.NewMessage(b.Dlg.ChatId, addr)
	b.Bot.Send(msg)
	go b.CheckStatusBuy(addr)
	return
}

// SellSecondStep
func (b *Bot) SellSecondStep() {
	msg := tgbotapi.NewMessage(b.Dlg.ChatId, vocab.GetTranslate("Select price", b.Dlg.language))
	msg.ReplyMarkup = b.GetPrice()
	b.Bot.Send(msg)
}

// SellFinal
func (b *Bot) SellFinal() {
	CoinToSell[b.Dlg.UserId] = "MNT"
	depos, err := b.Api.GetMinterDeposAddress(b.Dlg.Text, CoinToSell[b.Dlg.UserId], PriceToSell[b.Dlg.UserId])
	if err != nil {
		msg := tgbotapi.NewMessage(b.Dlg.ChatId, err.Error())
		msg.ReplyMarkup = b.newMainKeyboard()
		b.Bot.Send(msg)
		return
	}
	ans := fmt.Sprintf(vocab.GetTranslate("Minter deposit and tag", b.Dlg.language), depos.Data.Address, depos.Data.Tag)
	msg := tgbotapi.NewMessage(b.Dlg.ChatId, ans)
	b.Dlg.Command = ""
	msg.ReplyMarkup = b.newMainKeyboard()
	b.Bot.Send(msg)
	go b.CheckStatusSell(depos.Data.Tag)
	return
}