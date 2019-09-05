package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/mrKitikat/telegrambottest/src/app"
	stct "github.com/mrKitikat/telegrambottest/src/app/bipdev/structs"

	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"

	"log"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	conf := &stct.Config{
		URL:        os.Getenv("URL"),
		Token:      os.Getenv("TOKEN"),
		Driver:     os.Getenv("DRIVER"),
		DataSource: os.Getenv("DATASOURCE"),
	}

	dbsql, err := sql.Open(conf.Driver, conf.DataSource)
	if err != nil {
		log.Fatal(err)
	}

	defer dbsql.Close()

	app := app.NewApp(conf, dbsql)
	fmt.Println("APP Started!")
	app.Run(":8000")
}
