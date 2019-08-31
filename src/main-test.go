package main

// import (
// 	"fmt"
// 	"log"
// 	"sync"
// 	api "telegrambottest/src/bipdev"

// 	_ "github.com/jinzhu/gorm/dialects/postgres"
// )

// var MinterAddress = "Mxc19bf5558d8b374ad02557fd87d57ade178fc14a"

// var wg sync.WaitGroup

// func main() {
// 	app := api.InitApp("https://mbank.dl-dev.ru/api/")

// 	addr, err := app.GetBTCDeposAddress(MinterAddress, "BIP", "xxx@yyy.ru")
// 	if err != nil {
// 		log.Fatal(err)
// 	}

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
// }
