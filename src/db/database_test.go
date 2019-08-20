package db

import (
	"database/sql"
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
