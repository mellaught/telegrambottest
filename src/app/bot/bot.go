package bot

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	api "github.com/mrKitikat/telegrambottest/src/app/bipdev"
	vocab "github.com/mrKitikat/telegrambottest/src/app/bot/vocabulary"
	"github.com/mrKitikat/telegrambottest/src/app/db"
	stct "github.com/mrKitikat/telegrambottest/src/app/structs"

	//strt "bipbot/src/bipdev/structs"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
)

var (
	commands        = make(map[int64]string)
	UserHistory     = make(map[int64]string)
	PreviousMessage = make(map[int64]tgbotapi.MessageConfig)
	Message         = make(map[int64]tgbotapi.MessageConfig)
	MinterAddress   = make(map[int64]string)
	BitcoinAddress  = make(map[int64]string)
	CoinToSell      = make(map[int64]string)
	EmailAddress    = make(map[int64]string)
	PriceToSell     = make(map[int64]float64)
	//SellSteps       = make(map[int64]int)
	//BuySteps        = make(map[int64]int)
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
	fmt.Printf("UserHistory: %s \n", UserHistory[ChatId])
	// Проверка команды <<купить>>.
	if strings.Contains(UserHistory[ChatId], "buy") {
		// команда не выбрана.
		if UserHistory[ChatId][4:] == "1" {
			// Проверка минтер адреса.
			if !b.CheckMinter(text) {
				b.SendMessage(vocab.GetTranslate("Wrong minter", b.Dlg[ChatId].language), ChatId, nil)
				return
			} else {
				MinterAddress[ChatId] = text
				// Отправьте почту.
				UserHistory[ChatId] = "buy_2"
				kb, txt, err := b.SendEmail(ChatId)
				if err != nil {
					fmt.Println(err)
				}

				b.EditAndSend(&kb, txt, ChatId)
				return
			}
		} else if UserHistory[ChatId][4:] == "2" {
			// Проверка почты.
			if !b.CheckEmail(text) {
				b.SendMessage(vocab.GetTranslate("Wrong email", b.Dlg[ChatId].language), ChatId, nil)
				return
			} else {
				EmailAddress[ChatId] = text
				// Отправьте депозит на биткоин адрес.
				b.SendMenuChoose(ChatId)
				b.SendDepos(ChatId)
				b.BuyFinal(ChatId)
				return
			}
		}
	} else {
		if UserHistory[ChatId][5:] == "1" {
			// Проверка названия монеты.
			if !b.CheckCoin(text) {
				b.SendMessage(vocab.GetTranslate("Wrong coin name", b.Dlg[ChatId].language), ChatId, nil)
				return
			} else {
				// Отравьте цены за монеты.
				UserHistory[ChatId] = "sell_2"
				CoinToSell[ChatId] = text
				b.SendMessage(vocab.GetTranslate("Select price", b.Dlg[ChatId].language), ChatId, b.CancelKeyboard(ChatId))
				// PreviousMessage[ChatId] = Message[ChatId]
				// Message[ChatId] = msg
				return
			}
		} else if UserHistory[ChatId][5:] == "2" {
			// Проверка цены за монеты.
			if !b.CheckPrice(ChatId, text) {
				b.SendMessage(vocab.GetTranslate("Wrong price", b.Dlg[ChatId].language), ChatId, nil)
				return
			} else {
				// Отправьте биткоин адрес.
				UserHistory[ChatId] = "sell_3"
				kb, txt, err := b.SendBTCAddresses(ChatId)
				if err != nil {
					b.PrintAndSendError(err, ChatId)
					return
				}
				b.SendMessage(txt, ChatId, kb)
				return
			}
		} else if UserHistory[ChatId][5:] == "3" {
			// Проверка биткоин адреса.
			if !b.CheckBitcoin(text) {
				b.SendMessage(vocab.GetTranslate("Wrong bitcoin", b.Dlg[ChatId].language), ChatId, nil)
				return
			} else {
				// Сохранить ли введенные данные?
				// Отправьте монеты на адрес.
				BitcoinAddress[ChatId] = text
				b.SellFinal(ChatId)
				return
			}
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
		UserHistory[ChatId] = "start"
		b.SendMessage(vocab.GetTranslate("Hello", b.Dlg[ChatId].language), ChatId, b.newVocabuageKeybord())
		return

	// engvocabCommand sets english lang for user.
	case engvocabCommand:
		commands[ChatId] = ""

		b.DB.SetLanguage(b.Dlg[ChatId].UserId, "en")
		b.Dlg[ChatId].language = "en"
		kb, txt, err := b.SendMenuMessage(ChatId)
		if err != nil {
			b.PrintAndSendError(err, ChatId)
			return
		}
		b.EditAndSend(&kb, txt, ChatId)
		go b.ChangeCurrency(ChatId)
		//PreviousMessage[ChatId] = msg
		//fmt.Println("Message id:", b.Dlg[ChatId].MessageId)
		//go b.ChangeCurrency(ChatId, b.Dlg[ChatId].MessageId, b.Dlg[ChatId].CallBackId)
		return

	// rusvocabCommand sets russian lang for user.
	case rusvocabCommand:
		commands[ChatId] = ""
		b.DB.SetLanguage(b.Dlg[ChatId].UserId, "ru")
		b.Dlg[ChatId].language = "ru"
		kb, txt, err := b.SendMenuMessage(ChatId)
		if err != nil {
			b.PrintAndSendError(err, ChatId)
			return
		}

		b.EditAndSend(&kb, txt, ChatId)
		go b.ChangeCurrency(ChatId)
		return

	case cancelComm:
		b.CancelHandler(ChatId)

	// Save user's data in DataBase.
	case yescommand:
		err := b.DB.PutMinterAddress(b.Dlg[ChatId].UserId, MinterAddress[ChatId])
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
		UserHistory[ChatId] = "buy_1"
		_, ok := PreviousMessage[ChatId]
		if ok {
			delete(PreviousMessage, ChatId)
		}

		kb, txt, err := b.SendMinterAddresses(ChatId)
		if err != nil {
			fmt.Println(err)
		}

		b.EditAndSend(&kb, txt, ChatId)
		return

	//sendMinter after the user has selected the minter address from the proposed. ( BUY )
	case sendMinter:
		UserHistory[ChatId] = "buy_2"
		MinterAddress[ChatId] = b.Dlg[ChatId].Text
		kb, txt, err := b.SendEmail(ChatId)
		if err != nil {
			fmt.Println(err)
		}

		b.EditAndSend(&kb, txt, ChatId)
		return

	// sendEmail after the user has selected email from the proposed. ( BUY )
	case sendEmail:
		EmailAddress[ChatId] = b.Dlg[ChatId].Text
		b.SendMenuChoose(ChatId)
		b.SendDepos(ChatId)
		b.BuyFinal(ChatId)

	// sellCommand collects data from the user to transmit their request. ( SELL )
	case sellCommand:
		UserHistory[ChatId] = "sell_1"
		_, ok := PreviousMessage[ChatId]
		if ok {
			delete(PreviousMessage, ChatId)
		}

		kb := b.CancelKeyboard(ChatId)
		txt := vocab.GetTranslate("Coin", b.Dlg[ChatId].language)
		b.EditAndSend(&kb, txt, ChatId)
		return
		// Message[ChatId] = msg

	// sendBTC after the user has sent his bitcoin address. ( SELL )
	case sendBTC:
		BitcoinAddress[ChatId] = b.Dlg[ChatId].Text
		b.SellFinal(ChatId)

	// salesCommand sends a request to the database to get user's loots. ( LOOTS )
	case salesCommand:
		UserHistory[ChatId] = "loots"
		loots, err := b.DB.GetLoots(b.Dlg[ChatId].UserId)
		if err != nil {
			b.PrintAndSendError(err, ChatId)
			return
		} else if len(loots) == 0 {
			b.SendMessage(vocab.GetTranslate("Empty loots", b.Dlg[ChatId].language), ChatId, b.newMainMenuKeyboard(ChatId))
			return
		}

		b.SendLoots(loots, ChatId)
		return
	}
}
