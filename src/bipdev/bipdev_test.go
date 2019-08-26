package api

import (
	"testing"
)

// MinterAddress is my minter testnet address
var MinterAddress = "Mxc19bf5558d8b374ad02557fd87d57ade178fc14a"

// BitcoinAddress is my bitcoin testnet address
var BitcoinAddress = "mkWZZPqd1FebZM1MNFfZBQoYFqA4EpE8vD"

// Test for GetPrice
// Result: Success: Tests passed
func TestGetPrice(t *testing.T) {

	a := InitApp("https://mbank.dl-dev.ru/api/")

	p, err := a.GetPrice()
	if err != nil {
		t.Fatal(err)
	}

	if p != 1 {
		t.Errorf("Error price %f, want 1", p)
	}
}

// Test for GetBTCDeposAddress
// Result: Success: Tests passed.
func TestGetBTCDeposAddress(t *testing.T) {

	a := InitApp("https://mbank.dl-dev.ru/api/")

	addr, err := a.GetBTCDeposAddress(MinterAddress, "BIP", "xxx@yyy.ru")
	if err != nil {
		t.Fatal(err)
	}

	if addr == "" {
		t.Errorf("Empty address %s", addr)
	}

}

// Test for GetBTCDepositStatus
// Result: Success: Tests passed.
func TestGetBTCDepositStatus(t *testing.T) {

	a := InitApp("https://mbank.dl-dev.ru/api/")

	stat, err := a.GetBTCDepositStatus("tb1qtfnwald5a667730yqrvdt67aslmgn3k7qykq5a")
	if err != nil {
		t.Fatal(err)
	}

	if stat == nil {
		t.Errorf("Empty stat")
	}

	stat, err = a.GetBTCDepositStatus("saawdadadw")
	if err.Error() != "Address not found" {
		t.Errorf("Cannot found err")
	}

}

func TestTagInfo(t *testing.T) {

	a := InitApp("https://mbank.dl-dev.ru/api/")

	tag, err := a.GetTagInfo("")
	if err != nil {
		t.Fatal(err)
	}
	if tag == nil {
		t.Fatalf("Empty responce")
	}

	t.Log(tag.Data)
}

// Test for GetBTCDepositStatus
// Result: Success: Tests passed.
func TestGetMinterDeposAddress(t *testing.T) {

	a := InitApp("https://mbank.dl-dev.ru/api/")

	addr, err := a.GetMinterDeposAddress(BitcoinAddress, "BIP", 0.1)
	if err != nil {
		t.Fatal(err)
	}
	if addr == nil {
		t.Errorf("Empty responce")
	}

	if addr.Data.Tag == "" || addr.Data.Address == ""{
		t.Errorf("Empty tag or address: %s and %s ", addr.Data.Tag, addr.Data.Address)
	}

}

// Test for MinterAddressHistory and also BTCAddressHistory  
// Result: Success: Tests passed.
func TestAddressHistory(t *testing.T) {

	a := InitApp("https://mbank.dl-dev.ru/api/")
	h, err := a.MinterAddressHistory("Mx6a55fa3a81fc55124b46a3c36101d11a39c27bbe")
	if err != nil {
		t.Fatal(err)
	}
	if h == nil{
		t.Errorf("Empty responce")
	}
	if len(h.Data) != 2 {
		t.Errorf("Wrong len of Data %d, want : %d", len(h.Data), 2)
	}

	if h.Data[1].Amount != "200000000000000000000" {
		t.Errorf("Wrong amout of Data[1] %s, want: %s", h.Data[1].Amount, "200000000000000000000")
	}
}

//Creating Table
// if os.Getenv("CREATE_TABLE") == "yes" {
// 	if err := createTable(); err != nil {
// 		panic(err)
// 	}
// }
