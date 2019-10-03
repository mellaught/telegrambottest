package db

import (
	"database/sql"
	"time"

	stct "github.com/mrKitikat/telegrambottest/src/app/structs"

	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type DataBase struct {
	DB *sql.DB
}

// InitDB creates tables USERs or SALEs if tables not exists.
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

	_, err = db.Exec(CREATE_BITCOIN_DATA_IF_NOT_EXISTS)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(CREATE_MINTER_DATA_IF_NOT_EXISTS)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(CREATE_EMAIL_DATA_IF_NOT_EXISTS)
	if err != nil {
		return nil, err
	}

	return &d, nil
}

// PutUser adds user in database.
func (d *DataBase) PutUser(ChatId int64) error {

	_, err := d.DB.Exec("INSERT INTO USERS(id, chat_id, lang)"+
		"VALUES ($1,$2,$3)", int(ChatId), ChatId, "en")
	if err != nil {
		return err
	}

	return nil
}

// GetLanguage returns language for user by UserId.
func (d *DataBase) GetLanguage(ChatId int64) string {

	var lang string
	err := d.DB.QueryRow("SELECT lang FROM USERS WHERE id = $1 limit 1", int(ChatId)).Scan(&lang)
	if err != nil && err.Error() == "sql: no rows in result set" {
		d.PutUser(ChatId)
		return "en"
	}

	return lang
}

// SetLanguage is setting language for user by UserId.
func (d *DataBase) SetLanguage(UserId int, lang string) error {

	_, err := d.DB.Exec("UPDATE USERS SET lang = $1 where id = $2", lang, UserId)
	if err != nil {
		return err
	}

	return nil
}

// PutLoot puts new loot for sale.
func (d *DataBase) PutLoot(UserId int, tag string, taginfo *stct.TagInfo) error {
	_, err := d.DB.Exec("INSERT INTO LOOTS(user_id, tag, coin, price, amount, minter_address, created_at, last_sell_at)"+
		"VALUES ($1,$2,$3,$4,$5,$6,$7,$8)", UserId, tag, taginfo.Data.Coin, taginfo.Data.Price, taginfo.Data.Amount,
		taginfo.Data.MinterAddress, time.Now(), time.Time{})
	if err != nil {
		return err
	}

	return nil
}

// GetSales returns all sales for user by UserId.
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

// GetChatIDLang return user's chatID and Language.
func (d *DataBase) GetChatIDLang(UserId int) (int64, string, error) {

	var chatID int64
	var lang string

	err := d.DB.QueryRow("SELECT CHAT_ID, LANG FROM USERS WHERE ID = $1 LIMIT 1", UserId).Scan(&chatID, &lang)
	if err != nil {
		return 0, "", err
	}

	return chatID, lang, nil
}

// UpdateSales updates (insert new) sales for user by UserId.
func (d *DataBase) UpdateLoots(amount, tag string) (int64, string, error) {
	_, err := d.DB.Exec("UPDATE LOOTS SET last_sell_at = $1, amount = $2 where tag = $3", time.Now(), amount, tag)
	if err != nil {
		return 0, "", err
	}

	row := d.DB.QueryRow("SELECT USER_ID FROM LOOTS WHERE TAG = $1", tag)

	var userID int
	err = row.Scan(&userID)
	if err != nil {
		return 0, "", err
	}
	return d.GetChatIDLang(userID)
}

// GetBTCAddresses returns previously entered bitcoin addresses by UserID.
func (d *DataBase) GetBTCAddresses(userID int) ([]string, error) {

	rows, err := d.DB.Query("SELECT bitcoin_address FROM BITCOIN_DATA WHERE USER_ID = $1", userID)
	if err != nil {
		return nil, err
	}

	var addresses []string
	for rows.Next() {
		var addr string
		err := rows.Scan(&addr)
		if err != nil {
			return nil, err
		}

		addresses = append(addresses, addr)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return addresses, nil
}

// GetMinterAddresses returns previously entered minter addresses by UserID.
func (d *DataBase) GetMinterAddresses(userID int) ([]string, error) {

	rows, err := d.DB.Query("SELECT MINTER_ADDRESS FROM MINTER_DATA WHERE USER_ID = $1", userID)
	if err != nil {
		return nil, err
	}
	var addresses []string
	for rows.Next() {
		var addr string
		err := rows.Scan(&addr)
		if err != nil {
			return nil, err
		}

		addresses = append(addresses, addr)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return addresses, nil
}

// GetEmails returns previously entered emails by UserID.
func (d *DataBase) GetEmails(userID int) ([]string, error) {

	rows, err := d.DB.Query("SELECT EMAIL FROM EMAIL_DATA WHERE USER_ID = $1", userID)
	if err != nil {
		return nil, err
	}

	var emails []string
	for rows.Next() {
		var email string
		err := rows.Scan(&email)
		if err != nil {
			return nil, err
		}

		emails = append(emails, email)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return emails, nil
}

// PutBTCAddress puts new user's bitcoin address.
func (d *DataBase) PutBTCAddress(UserId int, bitcoinAddress string) error {
	_, err := d.DB.Exec("INSERT INTO bitcoin_data(user_id, bitcoin_address)  VALUES ($1, $2)"+
		"ON CONFLICT(bitcoin_address) DO NOTHING;", UserId, bitcoinAddress)
	if err != nil {
		return err
	}

	return nil
}

// PutMinterAddress puts new user's minter address.
func (d *DataBase) PutMinterAddress(UserId int, minterAddress string) error {
	_, err := d.DB.Exec("INSERT INTO minter_data(user_id, minter_address)  VALUES ($1, $2)"+
		"ON CONFLICT(minter_address) DO NOTHING;", UserId, minterAddress)
	if err != nil {
		return err
	}

	return nil
}

// PutEmail puts new user's email.
func (d *DataBase) PutEmail(UserId int, email string) error {
	_, err := d.DB.Exec("INSERT INTO email_data(user_id, email)  VALUES ($1, $2)"+
		"ON CONFLICT(email) DO NOTHING;", UserId, email)
	if err != nil {
		return err
	}

	return nil
}
