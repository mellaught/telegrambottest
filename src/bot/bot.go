package bot

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	api "telegrambottest/src/bipdev"
	stct "telegrambottest/src/bipdev/structs"
	"telegrambottest/src/db"

	//strt "bipbot/src/bipdev/structs"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
)

const (
	startCommand   = "start"
	priceCommand   = "price"
	buyCommand     = "buy"
	sellCommand    = "sell"
	salesCommand   = "lookat"
	getMainMenu    = "getmainmenu"
	engLangCommand = "englanguage"
	rusLangCommand = "ruslanguage"
)

var (
	commands    = make(map[int]string)
	CommandInfo = make(map[int]string)
)

type Dialog struct {
	ChatId   int64
	UserId   int
	Text     string
	Language string
	Command  string
}

type BuySell struct {
	Address string
	Email   string
	Price   float32
	Coin    string
}

type Bot struct {
	Token string
	Api   *api.App
	DB    db.DataBase
	Bot   *tgbotapi.BotAPI
}

func InitBot(config stct.Config, dbsql *sql.DB) *Bot {

	b := Bot{
		Token: config.Token,
		DB: db.DataBase{
			DB: dbsql,
		},
	}

	b.Api = api.InitApp(config.URL)
	//db.InitDB(dbsql)
	bot, err := tgbotapi.NewBotAPI(b.Token)
	if err != nil {
		log.Fatal(err)
	}
	b.Bot = bot
	//	b.initDB()

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

		dialog, exist := assembleUpdate(update)
		if !exist {
			continue
		}

		if update.Message != nil && update.Message.ReplyToMessage != nil {
			if dialog.Command == "buy" {
				b.Buy(dialog)
				continue
			}
		}
		if botCommand := b.getCommand(update); botCommand != "" {
			b.RunCommand(botCommand, dialog)
			continue
		}
		msg := tgbotapi.NewMessage(dialog.ChatId, "Select, please, what you want:)")
		msg.ReplyMarkup = newMainMenuKeyboard()
		b.Bot.Send(msg)

	}
}

// assembleUpdate
func assembleUpdate(update tgbotapi.Update) (Dialog, bool) {
	dialog := Dialog{}

	if update.Message != nil {
		dialog.ChatId = update.Message.Chat.ID
		dialog.UserId = int(update.Message.Chat.ID)
		dialog.Text = update.Message.Text
	} else if update.CallbackQuery != nil {
		dialog.ChatId = update.CallbackQuery.Message.Chat.ID
		dialog.UserId = int(update.CallbackQuery.Message.Chat.ID)
		dialog.Text = ""
	} else {
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

// RunCommand executes the input command
func (b *Bot) RunCommand(command string, dialog Dialog) {
	commands[dialog.UserId] = command
	switch command {
	case startCommand:
		msg := tgbotapi.NewMessage(dialog.ChatId, "Privet, i'm a exchange BIP/BTC or BTC/BIP bot")
		msg.ReplyMarkup = newLanguageKeybord()
		b.Bot.Send(msg)
	case getMainMenu:
		msg := tgbotapi.NewMessage(dialog.ChatId, "You can get current price BIP/USD\n"+
			"Also buy or sell your coins for BTC\n"+
			"My service give your chance to see your sales")
		msg.ReplyMarkup = newMainMenuKeyboard()
		b.Bot.Send(msg)
	case priceCommand:
		price, err := b.Api.GetPrice()
		if err != nil {
			msg := tgbotapi.NewMessage(dialog.ChatId, err.Error())
			b.Bot.Send(msg)
		}
		ans := fmt.Sprintf("📈 Now BIP/USD %f $", price)
		msg := tgbotapi.NewMessage(dialog.ChatId, ans)
		b.Bot.Send(msg)
	case buyCommand:
		msg := tgbotapi.NewMessage(dialog.ChatId, "Send me your Minter Address:)")
		msg.ReplyMarkup = tgbotapi.ForceReply{
			ForceReply: true,
			Selective:  true,
		}
		b.Bot.Send(msg)
	case sellCommand:
		msg := tgbotapi.NewMessage(dialog.ChatId, "Delepment")
		b.Bot.Send(msg)
	case salesCommand:
		msg := tgbotapi.NewMessage(dialog.ChatId, "Delepment")
		b.Bot.Send(msg)
	}
}

// Buy is function if method Buy.
func (b *Bot) Buy(dialog Dialog) {
	if strings.Contains(dialog.Text, "@") {
		addr, err := b.Api.GetBTCDeposAddress(CommandInfo[dialog.UserId], "BIP",
			dialog.Text)
		if err != nil {
			msg := tgbotapi.NewMessage(dialog.ChatId, err.Error())
			b.Bot.Send(msg)
			return
		}
		ans := fmt.Sprintf("Your BTC deposit address %s", addr)
		msg := tgbotapi.NewMessage(dialog.ChatId, ans)
		dialog.Command = ""
		b.Bot.Send(msg)
		return
		// Проверка статуса пошла
		//go b.CheckStatus(dialog, addr)
	} else {
		CommandInfo[dialog.UserId] = dialog.Text
		msg := tgbotapi.NewMessage(dialog.ChatId, "Send me your email!\n Example: myfriend@bipbest.com")
		msg.ReplyMarkup = tgbotapi.ForceReply{
			ForceReply: true,
			Selective:  true,
		}
		b.Bot.Send(msg)
		return
	}
}

// func (b *Bot) CheckStatus(dialog *Dialog, address string) {
// 	for {
// 		stat, err := b.Api.GetBTCDepositStatus(address)
// 		if err != nil{
// 			msg := tgbotapi.NewMessage(dialog.ChatId, err.Error())
// 			b.Bot.Send(msg)
// 			return
// 		}
// 		if stat.Data.
// 	}
// }

// newMainMenuKeyboard is main menu keyboar : price, buy, sell, sales
func newMainMenuKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("💹Price", priceCommand),
			tgbotapi.NewInlineKeyboardButtonData("💰Buy", buyCommand),
			tgbotapi.NewInlineKeyboardButtonData("💰Sell", sellCommand),
			tgbotapi.NewInlineKeyboardButtonData("📃My sales", salesCommand),
		),
	)
}

// LanguageKeybord is keybouad for select language
func newLanguageKeybord() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🇷🇺 Русский", rusLangCommand),
			tgbotapi.NewInlineKeyboardButtonData("🇬🇧 English", engLangCommand),
		),
	)
}

// if update.Message.IsCommand() {

// 	switch update.Message.Command() {
// 	case "start":
// 		msg.Text = "Privet, i'm a exchange BIP/BTC bot or BTC/BIP."
// 		bot.Send(msg)

//

// 	case "sell":
// 		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Delepment")
// 		bot.Send(msg)

// 	case "buy":
// 		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Delepment")
// 		bot.Send(msg)

// 	case "lookAt":
// 		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Delepment")
// 		bot.Send(msg)
// 	}

// } else {
// 	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Please, send command!:)")
// 	bot.Send(msg)
// }

// msg.ReplyMarkup = newMainMenuKeyboard()
// bot.Send(msg)

// }
