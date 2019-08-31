package api

import (
	"testing"
)

// MinterAddress is my minter testnet address
var MinterAddress = "Mxc19bf5558d8b374ad02557fd87d57ade178fc14a"

// BitcoinAddress is my bitcoin testnet address
var BitcoinAddress = "mrByFzaeHgwnayPznApmNiLJdcGbvehJqw"

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

// Test for GetTagInfo
// Result: Success: Tests passed.
func TestTagInfo(t *testing.T) {

	a := InitApp("https://mbank.dl-dev.ru/api/")

	tag, err := a.GetTagInfo("PCSmQDFTt2EOmBNSSQtF")
	if err != nil {
		t.Fatal(err)
	}

	if tag == nil {
		t.Fatalf("Empty responce")
	}

	tag, err = a.GetTagInfo("")
	if err.Error() != "Tag not found"{
		t.Errorf("Dont found err, want %s", err.Error())
	}

}

// Test for GetMinterDeposAddress
// Result: Success: Tests passed.
func TestGetMinterDeposAddress(t *testing.T) {

	a := InitApp("https://mbank.dl-dev.ru/api/")

	addr, err := a.GetMinterDeposAddress(BitcoinAddress, "BI", 0.2)
	if err != nil {
		t.Fatal(err)
	}
	if addr == nil {
		t.Errorf("Empty responce")
	}

}

func TestAddressHistory(t *testing.T) {

	a := InitApp("https://mbank.dl-dev.ru/api/")

	a.BTCAddressHistory("")
	a.MinterAddressHistory("")
}
