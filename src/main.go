package main

import (
	"database/sql"
	"fmt"

	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/mrKitikat/telegrambottest/src/app"
	stct "github.com/mrKitikat/telegrambottest/src/app/bipdev/structs"

	"log"
)

func main() {

	conf := stct.Config{
		DBName:        cfg.GetString("database.name"),
		DbUser:        cfg.GetString("database.user"),
		DbPassword:    cfg.GetString("database.password"),
		DbDriver:      cfg.GetString("database.driver"),
		BotToken       cfg.GetString("bot.token"),
		ServerPort:    cfg.GetString("server.port"),
		BipdevApiHost: cfg.GetString("bipdev.api"),
	}

	DatasourseName := "user" + conf.DbUser + "dbname=" + conf.DbName + "password=" + conf.DbPassword + "sslmode=disable"
	dbsql, err := sql.Open(conf.DbDriver, DatasourseName)
	if err != nil {
		log.Fatal(err)
	}

	defer dbsql.Close()

	app := app.NewApp(conf, dbsql)
	fmt.Println("APP Started!")
	app.Run(conf.ServerPort)
}
