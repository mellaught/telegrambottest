package main

// import (
// 	"fmt"
// 	"log"
// 	"sync"
// 	api "telegrambottest/src/bipdev"
// 	"time"

// 	_ "github.com/jinzhu/gorm/dialects/postgres"
// )

// var (
// 	MinterAddress = "Mxc19bf5558d8b374ad02557fd87d57ade178fc14a"
// 	BitcoinAddress = "n2x6Fu7ACk5BMUJiS75cLLAC3uFz6PgXyf"
// 	wg sync.WaitGroup
// )

// func main() {
// 	app := api.InitApp("https://mbank.dl-dev.ru/api/")

// 	addr, err := app.GetMinterDeposAddress(BitcoinAddress, "BIP", 0.1)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	if addr == nil {
// 		log.Fatalf("Empty addr")
// 	}

// 	fmt.Println(addr.Data.Address, addr.Data.Tag)
// 	time.Sleep(120 * time.Second)
// 	wg.Add(1)

// 	go app.CheckStatusSell(addr.Data.Tag, &wg)

// 	wg.Wait()
// 	fmt.Println("Test finish")
// }
