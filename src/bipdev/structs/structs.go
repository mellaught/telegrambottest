package stct

import (
	"math/big"
)

//Config ..
type Config struct {
	Token      string
	URL        string
	Driver     string
	DataSource string
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
		Address string `json:"message"`
		Tag     string `json:"tag"`
	} `json:"data"`
}

// BTCStatus is a responce func GetBTCDepositStatus.
type BTCStatus struct {
	Data struct {
		Coin        string `json:"coin"`
		WillReceive string `json:"will_receive"`
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
	Data struct {
		Amount int64 `json:"amount"`
	} `json:"data"`
}
