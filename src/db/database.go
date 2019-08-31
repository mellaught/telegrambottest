package db

import (
	"database/sql"
<<<<<<< HEAD

	//"log"
=======
	stct "telegrambottest/src/bipdev/structs"
	"time"
>>>>>>> 97af52583c4354e0e85352890f1f573f1701a764

	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type DataBase struct {
	DB *sql.DB
}

// InitDB creates tables USERs or SALEs if tables not exists
func InitDB(db *sql.DB) (*DataBase, error) {

	d := DataBase{
		DB: db,
	}
	_, err := db.Exec(CREATE_USERS_IF_NOT_EXISTS)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(CREATE_LOOTS_IF_NOT_EXISTS)
	if err != nil {
		return nil, err
	}

	return &d, nil
}

// PutUser adds user in database
func (d *DataBase) PutUser(ChatId int64) error {

	_, err := d.DB.Exec("INSERT INTO USERS(id, chat_id, lang)"+
		"VALUES ($1,$2,$3)", int(ChatId), ChatId, "en")
	if err != nil {
		return err
	}

	return nil
}

// GetLanguage returns language for user by UserId
func (d *DataBase) GetLanguage(ChatId int64) string {

	rows := d.DB.QueryRow("SELECT lang FROM USERS WHERE id = $1 limit 1", int(ChatId))
	var lang string
	err := rows.Scan(&lang)
	if err != nil && err.Error() == "sql: no rows in result set" {
		d.PutUser(ChatId)
		return "en"
	}

	return lang
}

// SetLanguage is setting language for user by UserId
func (d *DataBase) SetLanguage(UserId int, lang string) error {

	_, err := d.DB.Exec("UPDATE USERS SET lang = $1 where id = $2", lang, UserId)
	if err != nil {
		return err
	}

	return nil
}

// PutLoot puts new loot for sale
func (d *DataBase) PutLoot(UserId int, tag string, taginfo *stct.TagInfo) error {
	_, err := d.DB.Exec("INSERT INTO LOOTS(user_id, tag, coin, price, amount, minter_address, created_at)"+
		"VALUES ($1,$2,$3,$4,$5,$6,$7)", UserId, tag, taginfo.Data.Coin, taginfo.Data.Price, taginfo.Data.Amount, taginfo.Data.MinterAddress, time.Now())
	if err != nil {
		return err
	}

	return nil
}

// GetSales returns all sales for user by UserId
<<<<<<< HEAD
// func (d *DataBase) GetSales(UserId int) []string {

// 	rows, err := d.DB.Query("SELECT * FROM SALES WHERE user_id = $1", UserId)

// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	defer rows.Close()
// 	var lang string
// 	for rows.Next() {
// 		err := rows.Scan(&lang)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 	}

// 	if err = rows.Err(); err != nil {
// 		log.Fatal(err)
// 	}

// 	return lang
// }

// UpdateSales updates (insert new) sales for user by UserId
// func (d *DataBase) UpdateSales() error {
// 	_, err := d.DB.Exec("INSERT INTO SALES(chat_id, lang)"+
// 		"VALUES ($1,$2,$3)", int(ChatId), "en")
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }
=======
func (d *DataBase) GetLoots(UserId int) ([]*stct.Loot, error) {
	rows, err := d.DB.Query("SELECT * FROM LOOTS WHERE user_id = $1 ", UserId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	loots := []*stct.Loot{}
	for rows.Next() {
		var u int
		loot := new(stct.Loot)
		err := rows.Scan(&loot.ID, &u, &loot.Tag, &loot.Coin, &loot.Price, &loot.Amout, &loot.MinterAddress, &loot.CreatedAt, &loot.LastSell)
		if err != nil {
			return nil, err
		}

		loots = append(loots, loot)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return loots, nil
}

// UpdateSales updates (insert new) sales for user by UserId
func (d *DataBase) UpdateLoots(amount, tag string) error {
	_, err := d.DB.Exec("UPDATE LOOTS SET last_sell_at = $1, amount = $2 where tag = $3", time.Now(), amount, tag)
	if err != nil {
		return err
	}

	return nil
}
>>>>>>> 97af52583c4354e0e85352890f1f573f1701a764
