package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var ctx context.Context

type DataBase struct {
	DB *sql.DB
}

// InitDB creates tables USERs or SALEs if tables not exists
func InitDB(db *sql.DB) error {

	_, err := db.Exec(CREATE_USERS_IF_NOT_EXISTS)
	if err != nil {
		return err
	}

	_, err = db.Exec(CREATE_SALES_IF_NOT_EXISTS)
	if err != nil {
		return err
	}

	return nil
}

// GetLanguage return language of user
func (d *DataBase) GetLanguage(Userid int) (string, error) {

	rows, err := d.DB.Query("Select lang from users where user_id = $1 limit 1", Userid)

	if err != nil {
		fmt.Println(err)
		return "", errors.New("Something going wrong with database, sorry:(")
	}

	var lang string

	defer rows.Close()

	err = rows.Scan(&lang)
	if err != nil {
		fmt.Println(err)
		return "", errors.New("Something going wrong with database, sorry:(")
	}

	return lang, nil
}
