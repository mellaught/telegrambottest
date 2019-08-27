package db

import (
	"database/sql"
	api "telegrambottest/src/bipdev"
	"testing"
)

// Test for InitDB ( create tables ) and PutUser
// Result: Success: Tests passed.
func TestCreateTablePutUser(t *testing.T) {

	dbsql, err := sql.Open("postgres", "user=postgres dbname=gorm password=simsim sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}

	db, err := InitDB(dbsql)
	if err != nil {
		t.Fatal(err)
	}
	err = db.PutUser(12312)
	if err != nil {
		t.Fatal(err)
	}
}

// Test For GetLanguage
// Result: Success: Tests passed.
func TestGetLanguage(t *testing.T) {
	dbsql, err := sql.Open("postgres", "user=postgres dbname=gorm password=simsim sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}

	db, err := InitDB(dbsql)
	if err != nil {
		t.Fatal(err)
	}
	lang := db.GetLanguage(12312)
	if lang != "en" {
		t.Fatalf("I want see en, but see %s", lang)
	}
}

// Test for set language
// Result: Success: Tests passed.
func TestSetLanguage(t *testing.T) {
	dbsql, err := sql.Open("postgres", "user=postgres dbname=gorm password=simsim sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}

	db, err := InitDB(dbsql)
	if err != nil {
		t.Fatal(err)
	}
	err = db.SetLanguage(12312, "ru")
	if err != nil {
		t.Fatal(err)
	}
	lang := db.GetLanguage(12312)
	if lang != "ru" {
		t.Fatalf("I want see en, but see %s", lang)
	}
}

// Test for put user's loot to sell
// Result: Success: Tests passed.
func TestPutLoot(t *testing.T) {
	dbsql, err := sql.Open("postgres", "user=postgres dbname=gorm password=simsim sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}

	db, err := InitDB(dbsql)
	if err != nil {
		t.Fatal(err)
	}
	app := api.InitApp("https://mbank.dl-dev.ru/api/")
	taginfo, err := app.GetTagInfo("FFY37X9kRfvfGeDT8hZv")
	if err != nil {
		t.Fatal(err)
	}

	err = db.PutLoot(344178872, "FFY37X9kRfvfGeDT8hZv", taginfo)
	if err != nil {
		t.Fatal(err)
	}
}

// Test for get user's loots
// Result: Success: Tests passed.
func TestGetLoots(t *testing.T) {
	dbsql, err := sql.Open("postgres", "user=postgres dbname=gorm password=simsim sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}

	db, err := InitDB(dbsql)
	if err != nil {
		t.Fatal(err)
	}

	loots, err := db.GetLoots(344178872)
	if err != nil {
		t.Fatal(err)
	}
	if loots[0].Price != 1000 || loots[0].Amout != "1000" || loots[0].Coin != "MNT" {
		t.Errorf("I wanna see price = 1000, but %d, amount = 1000, but %s, coin = MNT, but %s", loots[0].Price, loots[0].Amout, loots[0].Coin)
	}
}

// Test for update user's loot: last_sell_at, amount
// Result: Success: Tests passed.
func TestUpdateLoots(t *testing.T) {
	dbsql, err := sql.Open("postgres", "user=postgres dbname=gorm password=simsim sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}

	db, err := InitDB(dbsql)
	if err != nil {
		t.Fatal(err)
	}

	err = db.UpdateLoots("1000", "FFY37X9kRfvfGeDT8hZv")
	if err != nil {
		t.Fatal(err)
	}
}
