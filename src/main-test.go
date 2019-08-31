package main

// import (
// 	"fmt"
// 	"log"
// 	"sync"
// 	api "telegrambottest/src/bipdev"

// 	_ "github.com/jinzhu/gorm/dialects/postgres"
// )

<<<<<<< HEAD
// var MinterAddress = "Mxc19bf5558d8b374ad02557fd87d57ade178fc14a"

// var wg sync.WaitGroup
=======
// var (
// 	MinterAddress  = "Mxc19bf5558d8b374ad02557fd87d57ade178fc14a"
// 	BitcoinAddress = "n2x6Fu7ACk5BMUJiS75cLLAC3uFz6PgXyf"
// 	wg             sync.WaitGroup
// )
>>>>>>> 97af52583c4354e0e85352890f1f573f1701a764

// func main() {
// 	app := api.InitApp("https://mbank.dl-dev.ru/api/")

<<<<<<< HEAD
// 	addr, err := app.GetBTCDeposAddress(MinterAddress, "BIP", "xxx@yyy.ru")
=======
// 	addr, err := app.GetMinterDeposAddress(BitcoinAddress, "MNT", 0.1)
>>>>>>> 97af52583c4354e0e85352890f1f573f1701a764
// 	if err != nil {
// 		log.Fatal(err)
// 	}

<<<<<<< HEAD
// 	if addr == "" {
// 		log.Fatalf("Empty addr")
// 	}
// 	fmt.Println(addr)
// 	//time.Sleep(1 * time.Minute)
// 	wg.Add(1)

// 	go app.CheckStatus(addr, &wg)
// 	//time.Sleep(10 * time.Second)
// 	wg.Wait()
// 	fmt.Println("+ Test over")
=======
// 	if addr == nil {
// 		log.Fatalf("Empty addr")
// 	}

// 	fmt.Println(addr.Data.Address, "     ", addr.Data.Tag)
// 	//time.Sleep(120 * time.Second)
// 	wg.Add(1)

// 	go app.CheckStatusSell("FFY37X9kRfvfGeDT8hZv", &wg)

// 	wg.Wait()
// 	fmt.Println("Test finish")
>>>>>>> 97af52583c4354e0e85352890f1f573f1701a764
// }
