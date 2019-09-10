package stct

import (
	"math/big"
	"time"
)

//Config ..
type Config struct {
	DbName        string
	DbUser        string
	DbPassword    string
	DbDriver      string
	BotToken      string
	ServerPort    string
	BipdevApiHost string
}

// Price is a structure of resp method GetPrice()
type Price struct {
	Data struct {
		Delta float32 `json:"delta"`
		Price int     `json:"price"`
	} `json:"data"`
}

// Err is a struct if the request is erroneous
type Err struct {
	Error struct {
		Message string `json:"message"`
	} `json:"error"`
}

// DeposBTC is a deposit struct for func GetBTCDeposAddress
type DeposBTC struct {
	Data struct {
		Address string `json:"address"`
	} `json:"data"`
}

// DeposMNT is a deposit struct for func GetMinterDeposAddress
type DeposMNT struct {
	Data struct {
		Address string `json:"address"`
		Tag     string `json:"tag"`
	} `json:"data"`
}

// BTCStatus is a responce func GetBTCDepositStatus.
type BTCStatus struct {
	Data struct {
		Coin        string  `json:"coin"`
		WillReceive float64 `json:"will_receive"`
	} `json:"data"`
}

// TagInfo is a responce func GetTagInfo.
type TagInfo struct {
	Data struct {
		MinterAddress string  `json:"minter_address"`
		BTCPrice      big.Int `json:"btc_price"`
		Price         int     `json:"price"`
		Coin          string  `json:"coin"`
		Amount        string  `json:"amount"`
	} `json:"data"`
}

// AddrHistory is a responce funcs BTCAddressHistory and MinterAddressHistory
type AddrHistory struct {
	Data []*Data `json:"data"`
}

type Data struct {
	Amount string `json:"amount"`
}

type Available struct {
	Data []int `json:"data"`
}

// Loot is a responce GetLoots
type Loot struct {
	ID            int
	Tag           string
	Coin          string
	Price         int
	Amout         string
	MinterAddress string
	CreatedAt     time.Time
	LastSell      time.Time
}

// REST API request for update loot with current tag
type UPDLoot struct {
	Tag        string `json:"tag"`
	Amount     string `json:"amount"`
	SellAmount string `json:"sells"`
	Coin       string `json:"coin"`
	Price      int    `json:"price"`
}
