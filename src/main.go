package main

import (
	stct "bipBot/src/bipdev/structs"
	"bipBot/src/bot"
	"database/sql"
	"fmt"

	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"

	//"bipbot/src/bot"
	"log"
	"os"
)

// MinterAddress is my minter testnet address
var MinterAddress = "Mxc19bf5558d8b374ad02557fd87d57ade178fc14a"

// BitcoinAddress is my bitcoin testnet address
var BitcoinAddress = ""

var email = "torres-dan@yandex.ru"

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	conf := stct.Config{
		URL:   os.Getenv("URL"),
		Token: os.Getenv("TOKEN"),
	}

	dbsql, err := sql.Open("postgres", "user=postgres dbname=gorm password=simsim sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer dbsql.Close()

	// Inizializaton users from DB, token for bot.
	bot := bot.InitBot(conf, dbsql)
	// Run bot
	fmt.Println("Bot started!")
	bot.Run()

}

// --------------------------------------------------GetBTCDeposAddress

// addr, err := app.GetBTCDeposAddress(MinterAddress, "BIP", email)
// if err != nil {
// 	log.Fatal(err)
// }
// fmt.Println(addr)

// ---------------------------------------------GetBTCDepositStatus

// stat, err := app.GetBTCDepositStatus("tb1qtfnwald5a667730yqrvdt67aslmgn3k7qykq5a")
// if err != nil {
// 	log.Fatal(err)
// }
// fmt.Println(stat.Data.Coin, stat.Data.WillReceive)

// --------------------------------------------- GetBTCDepositStatus
//Creating Table
// if os.Getenv("CREATE_TABLE") == "yes" {
// 	if err := createTable(); err != nil {
// 		panic(err)
// 	}
// }
