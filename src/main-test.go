package main

// import (
// 	"fmt"
// 	"log"
// 	"sync"
// 	api "telegrambottest/src/bipdev"

// 	_ "github.com/jinzhu/gorm/dialects/postgres"
// )

// var (
// 	MinterAddress  = "Mxc19bf5558d8b374ad02557fd87d57ade178fc14a"
// 	BitcoinAddress = "n2x6Fu7ACk5BMUJiS75cLLAC3uFz6PgXyf"
// 	wg             sync.WaitGroup
// )

// func main() {
// 	app := api.InitApp("https://mbank.dl-dev.ru/api/")

// 	addr, err := app.GetMinterDeposAddress(BitcoinAddress, "MNT", 0.1)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	if addr == nil {
// 		log.Fatalf("Empty addr")
// 	}

// 	fmt.Println(addr.Data.Address, "     ", addr.Data.Tag)
// 	//time.Sleep(120 * time.Second)
// 	wg.Add(1)

// 	go app.CheckStatusSell("FFY37X9kRfvfGeDT8hZv", &wg)

// 	wg.Wait()
// 	fmt.Println("Test finish")
// }
