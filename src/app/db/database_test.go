package db

import (
	"database/sql"
	api "telegrambottest/src/app/bipdev"
	"testing"
)

// Test for InitDB ( create tables ) and PutUser.
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

// Test For GetLanguage.
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
	lang := db.GetLanguage(11)
	if lang != "ru" {
		t.Fatalf("I want see ru, but see %s", lang)
	}
}

// Test for set language.
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

// Test for put user's loot to sell.
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
	taginfo, err := app.GetTagInfo("kxCUNnIQdJkFfVPTuE4V")
	if err != nil {
		t.Fatal(err)
	}

	err = db.PutLoot(344178872, "kxCUNnIQdJkFfVPTuE4V", taginfo)
	if err != nil {
		t.Fatal(err)
	}
}

// Test for get user's loots.
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

// Test for update user's loot: last_sell_at, amount.
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

	chatid, lang, err := db.UpdateLoots("2000", "FFY37X9kRfvfGeDT8hZv")
	if err != nil {
		t.Fatal(err)
	}

	if chatid != 344178872 || lang != "ru" {
		t.Errorf("I want see chatid , but %d and lang want see ,but lang = %s ", chatid, lang)
	}
}

// Test for get user's entered bitcoin addresses.
// Result: Success: Tests passed.
func TestGetBTCAddresses(t *testing.T) {

	dbsql, err := sql.Open("postgres", "user=postgres dbname=gorm password=simsim sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}

	db, err := InitDB(dbsql)
	if err != nil {
		t.Fatal(err)
	}

	addresses, err := db.GetBTCAddresses(344178872)
	if err != nil {
		t.Fatal(err)
	}

	if addresses != nil {
		t.Errorf("Addresses must be empty %d!", len(addresses))
	}

}

// Test for get user's entered minter addresses.
// Result: Success: Tests passed.
func TestGetMinterAddresses(t *testing.T) {

	dbsql, err := sql.Open("postgres", "user=postgres dbname=gorm password=simsim sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}

	db, err := InitDB(dbsql)
	if err != nil {
		t.Fatal(err)
	}

	addresses, err := db.GetMinterAddresses(344178872)
	if err != nil {
		t.Fatal(err)
	}

	if addresses != nil {
		t.Errorf("Addresses must be empty %d!", len(addresses))
	}

}

// Test for get user's entered email addresses.
// Result: Success: Tests passed.
func TestGetEmails(t *testing.T) {

	dbsql, err := sql.Open("postgres", "user=postgres dbname=gorm password=simsim sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}

	db, err := InitDB(dbsql)
	if err != nil {
		t.Fatal(err)
	}

	addresses, err := db.GetEmails(344178872)
	if err != nil {
		t.Fatal(err)
	}

	if addresses != nil {
		t.Errorf("Addresses must be empty %d!", len(addresses))
	}

}

// Test for put new user's email.
// Result: Success: Tests passed.
func TestPut(t *testing.T) {
	dbsql, err := sql.Open("postgres", "user=postgres dbname=gorm password=simsim sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}

	db, err := InitDB(dbsql)
	if err != nil {
		t.Fatal(err)
	}

	err = db.PutEmail(344178872, "torres-dan@yandex.ru")
	if err != nil {
		t.Fatal(err)
	}

	err = db.PutEmail(344178872, "torres-dan@yandex.ru")
	if err != nil {
		t.Fatal(err)
	}
	emails, err := db.GetEmails(344178872)
	if err != nil {
		t.Fatal(err)
	}
	if len(emails) != 1 {
		t.Errorf("I want see 1, but see %d", len(emails))
	}
	
	err = db.PutMinterAddress(344178872, "Mxc19bf5558d8b374ad02557fd87d57ade178fc14a")
	if err != nil {
		t.Fatal(err)
	}

	err = db.PutBTCAddress(344178872, "mkWZZPqd1FebZM1MNFfZBQoYFqA4EpE8vD")
	if err != nil {
		t.Fatal(err)
	}

}
