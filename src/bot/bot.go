package bot

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	api "telegrambottest/src/bipdev"
	stct "telegrambottest/src/bipdev/structs"
	vocab "telegrambottest/src/bot/vocabulary"
	"telegrambottest/src/db"
	"time"

	//strt "bipbot/src/bipdev/structs"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
)

const (
	startCommand    = "start"
	priceCommand    = "price"
	buyCommand      = "buy"
	sellCommand     = "sell"
	salesCommand    = "lookat"
	getMainMenu     = "getmainmenu"
	engvocabCommand = "englanguage"
	rusvocabCommand = "ruslanguage"
)

var (
	commands    = make(map[int]string)
	CommandInfo = make(map[int]string)
	CoinToSell  = make(map[int]string)
	PriceToSell = make(map[int]float64)
)

type Dialog struct {
	ChatId   int64
	UserId   int
	Text     string
	language string
	Command  string
}

type Bot struct {
	Token string
	Api   *api.App
	DB    *db.DataBase
	Bot   *tgbotapi.BotAPI
}

func InitBot(config stct.Config, dbsql *sql.DB) *Bot {

	b := Bot{
		Token: config.Token,
		DB:    &db.DataBase{},
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

		if update.Message != nil && update.Message.ReplyToMessage != nil {
			if dialog.Command == "buy" {
				b.Buy(&dialog)
				continue
			} else if dialog.Command == "sell" {
				b.Sell(dialog)
				continue
			}
		}

		if botCommand := b.getCommand(update); botCommand != "" {
			b.RunCommand(botCommand, dialog)
			continue
		}

		msg := tgbotapi.NewMessage(dialog.ChatId, vocab.GetTranslate("Select", dialog.language))
		msg.ReplyMarkup = newMainMenuKeyboard(&dialog)
		b.Bot.Send(msg)

	}
}

// assembleUpdate
func (b *Bot) assembleUpdate(update tgbotapi.Update) (Dialog, bool) {
	dialog := Dialog{}
	if update.Message != nil {
		fmt.Println("111")
		dialog.language = b.DB.GetLanguage(update.Message.Chat.ID)
		dialog.ChatId = update.Message.Chat.ID
		dialog.UserId = int(update.Message.Chat.ID)
		dialog.Text = update.Message.Text
	} else if update.CallbackQuery != nil {
		fmt.Println("222")
		dialog.language = b.DB.GetLanguage(update.CallbackQuery.Message.Chat.ID)
		dialog.ChatId = update.CallbackQuery.Message.Chat.ID
		dialog.UserId = int(update.CallbackQuery.Message.Chat.ID)
		dialog.Text = ""
	} else {
		fmt.Println("333")
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
			fmt.Println("command: ", update.Message.Command())
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
		msg := tgbotapi.NewMessage(dialog.ChatId, vocab.GetTranslate("Hello", dialog.language))
		msg.ReplyMarkup = newvocabuageKeybord()
		b.Bot.Send(msg)
	case engvocabCommand:
		b.DB.SetLanguage(dialog.UserId, "en")
		dialog.language = "en"
		msg := tgbotapi.NewMessage(dialog.ChatId, vocab.GetTranslate("Installed", dialog.language)+" "+
			vocab.GetTranslate("english", dialog.language))
		msg.ReplyMarkup = newMainMenuKeyboard(&dialog)
		b.Bot.Send(msg)
	case rusvocabCommand:
		b.DB.SetLanguage(dialog.UserId, "ru")
		dialog.language = "ru"
		msg := tgbotapi.NewMessage(dialog.ChatId, vocab.GetTranslate("Installed", dialog.language)+" "+
			vocab.GetTranslate("russian", dialog.language))
		msg.ReplyMarkup = newMainMenuKeyboard(&dialog)
		b.Bot.Send(msg)
	case priceCommand:
		price, err := b.Api.GetPrice()
		if err != nil {
			msg := tgbotapi.NewMessage(dialog.ChatId, err.Error())
			b.Bot.Send(msg)
		}
		ans := fmt.Sprintf(vocab.GetTranslate("Now", dialog.language), price)
		msg := tgbotapi.NewMessage(dialog.ChatId, ans)
		b.Bot.Send(msg)
	case buyCommand:
		msg := tgbotapi.NewMessage(dialog.ChatId, vocab.GetTranslate("Send minter", dialog.language))
		msg.ReplyMarkup = tgbotapi.ForceReply{
			ForceReply: true,
			Selective:  true,
		}
		b.Bot.Send(msg)
	case sellCommand:
		msg := tgbotapi.NewMessage(dialog.ChatId, vocab.GetTranslate("Development", dialog.language))
		b.Bot.Send(msg)
	case salesCommand:
		msg := tgbotapi.NewMessage(dialog.ChatId, vocab.GetTranslate("Development", dialog.language))
		// loots, err := b.DB.GetLoots(dialog.UserId)
		// if err != nil {
		// 	fmt.Println(err)
		// 	msg := tgbotapi.NewMessage(dialog.ChatId, vocab.GetTranslate("Error", dialog.language))
		// 	b.Bot.Send(msg)
		// }
		b.Bot.Send(msg)
		// case getMainMenu:
		// 	msg := tgbotapi.NewMessage(dialog.ChatId, "You can get current price BIP/USD\n"+
		// 		"Also buy or sell your coins for BTC\n"+
		// 		"My service give your chance to see your sales")
		// 	msg.ReplyMarkup = newMainMenuKeyboard()
		// 	b.Bot.Send(msg)
	}
}

// Buy is function for method Buy
func (b *Bot) Buy(dialog *Dialog) {
	if strings.Contains(dialog.Text, "@") {
		addr, err := b.Api.GetBTCDeposAddress(CommandInfo[dialog.UserId], "BIP",
			dialog.Text)
		if err != nil {
			msg := tgbotapi.NewMessage(dialog.ChatId, err.Error())
			b.Bot.Send(msg)
			return
		}
		ans := fmt.Sprintf(vocab.GetTranslate("BTC deposit", dialog.language), addr)
		msg := tgbotapi.NewMessage(dialog.ChatId, ans)
		msg.ReplyMarkup = tgbotapi.ForceReply{
			ForceReply: false,
			Selective:  false,
		}
		dialog.Command = ""
		b.Bot.Send(msg)
		go b.CheckStatusBuy(dialog, addr)
		return
	} else {
		CommandInfo[dialog.UserId] = dialog.Text
		msg := tgbotapi.NewMessage(dialog.ChatId, vocab.GetTranslate("Email", dialog.language))
		msg.ReplyMarkup = tgbotapi.ForceReply{
			ForceReply: true,
			Selective:  true,
		}
		b.Bot.Send(msg)
		return
	}
}

// CheckStatusBuy checks depos BTC and 2 confirme
func (b *Bot) CheckStatusBuy(dialog *Dialog, address string) {
	timeout := time.After(2 * time.Minute)
	tick := time.Tick(3 * time.Second)
	willcoin := 0.
	for {
		select {
		case <-timeout:
			if willcoin == 0. {
				msg := tgbotapi.NewMessage(dialog.ChatId, vocab.GetTranslate("timeout", dialog.language))
				msg.ReplyMarkup = newMainMenuKeyboard(dialog)
				b.Bot.Send(msg)
				return
			} else {
				continue
			}

		case <-tick:
			stat, err := b.Api.GetBTCDepositStatus(address)
			if err != nil {
				fmt.Println(err)
				msg := tgbotapi.NewMessage(dialog.ChatId, vocab.GetTranslate("Error", dialog.language))
				b.Bot.Send(msg)
				return
			}
			if stat.Data.WillReceive != willcoin {
				if willcoin == 0. {
					willcoin = stat.Data.WillReceive
					ans := fmt.Sprintf(vocab.GetTranslate("New deposit", dialog.language), stat.Data.WillReceive)
					msg := tgbotapi.NewMessage(dialog.ChatId, ans)
					b.Bot.Send(msg)
					time.Sleep(60 * time.Second)
				} else {
					ans := fmt.Sprintf(vocab.GetTranslate("Exchange is successful", dialog.language), willcoin)
					msg := tgbotapi.NewMessage(dialog.ChatId, ans)
					b.Bot.Send(msg)
					return
				}
			}
		}
	}
}

//
// Sell is function for method Sell
func (b *Bot) Sell(dialog Dialog) {
	if len(dialog.Text) > 24 {
		// checkvalidbitcoin
		depos, err := b.Api.GetMinterDeposAddress(dialog.Text, CoinToSell[dialog.UserId], PriceToSell[dialog.UserId])
		if err != nil {
			msg := tgbotapi.NewMessage(dialog.ChatId, err.Error())
			b.Bot.Send(msg)
			return
		}

		ans := fmt.Sprintf(vocab.GetTranslate("Minter deposit and tag", dialog.language), depos.Data.Address, depos.Data.Tag)
		msg := tgbotapi.NewMessage(dialog.ChatId, ans)
		dialog.Command = ""
		b.Bot.Send(msg)
		return
		// Проверка статуса пошла
		//go b.CheckStatusSell(dialog, addr)
	} else {
		CoinToSell[dialog.UserId] = dialog.Text
		msg := tgbotapi.NewMessage(dialog.ChatId, vocab.GetTranslate("Send BTC", dialog.language))
		msg.ReplyMarkup = tgbotapi.ForceReply{
			ForceReply: true,
			Selective:  true,
		}
		b.Bot.Send(msg)
		return
	}
}

// CheckStatusSell checks status of deposit for method Sell
func (b *Bot) CheckStatusSell(tag string, dialog *Dialog) {
	timeout := time.After(2 * time.Minute)
	tick := time.Tick(3 * time.Second)
	amount := "0"
	for {
		select {
		case <-timeout:
			if amount == "0" {
				msg := tgbotapi.NewMessage(dialog.ChatId, vocab.GetTranslate("timeout", dialog.language))
				msg.ReplyMarkup = newMainMenuKeyboard(dialog)
				b.Bot.Send(msg)
				return
			} else {
				continue
			}
		case <-tick:
			taginfo, err := b.Api.GetTagInfo(tag)
			if err != nil {
				fmt.Println(err)
				msg := tgbotapi.NewMessage(dialog.ChatId, vocab.GetTranslate("Error", dialog.language))
				b.Bot.Send(msg)
				return
			}
			if taginfo.Data.Amount != amount {
				amount = taginfo.Data.Amount
				fmt.Printf("Новый депозит на продажу %s %s по %d $\n", taginfo.Data.Amount, taginfo.Data.Coin, taginfo.Data.Price)
				// Добавить в БД
				b.DB.PutLoot(dialog.UserId, tag, taginfo)
				//go a.CheckLootforSell(taginfo.Data.MinterAddress)
				return
			}

		}
	}
}

// func (a *App) CheckLootforSell(addr string) {
// 	tick := time.Tick(1 * time.Hour)
// 	lenght := 0
// 	for {
// 		select {
// 		case <-tick:
// 			history, err := a.MinterAddressHistory(addr)
// 			if err != nil {
// 				log.Fatal(err)
// 				return
// 			}
// 			if len(history.Data) > lenght {

// 			}

// 		}
// 	}
// }

// newMainMenuKeyboard is main menu keyboar : price, buy, sell, sales
func newMainMenuKeyboard(dialog *Dialog) tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(vocab.GetTranslate("Price", dialog.language), priceCommand),
			tgbotapi.NewInlineKeyboardButtonData(vocab.GetTranslate("Buy", dialog.language), buyCommand),
			tgbotapi.NewInlineKeyboardButtonData(vocab.GetTranslate("Sell", dialog.language), sellCommand),
			tgbotapi.NewInlineKeyboardButtonData(vocab.GetTranslate("Sales", dialog.language), salesCommand),
		),
	)
}

// vocabuageKeybord is keybouad for select vocabuage
func newvocabuageKeybord() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🇷🇺 Русский", rusvocabCommand),
			tgbotapi.NewInlineKeyboardButtonData("🇬🇧 English", engvocabCommand),
		),
	)
}
