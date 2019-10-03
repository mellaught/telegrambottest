package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"

	stct "github.com/mrKitikat/telegrambottest/src/app/structs"
)

// App is main app for API Methods
type App struct {
	URL string
}

// InitApp inizilizations App
func InitApp(URL string) *App {

	a := App{
		URL: URL,
	}

	return &a
}

func GetMethod(req string) ([]byte, error) {
	//req := a.URL + "bitcoinDepositAddress?minterAddress=" + minterAddress + "&coin=" + coin + "&email=" + email
	response, err := http.Get(req)
	if err != nil {
		return nil, errors.New("http://bip.dev is not respond")
	}

	defer response.Body.Close()

	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("Something going wrong, sorry:(")
	}

	if response.StatusCode == 400 || response.StatusCode == 404 {
		data := &stct.Err{}
		err := json.Unmarshal([]byte(contents), data)
		if err != nil {
			return nil, err
		}
		return nil, errors.New(data.Error.Message)
	}

	return contents, nil
}

// GetPrice return current price BIP/USD
func (a *App) GetPrice() (float64, string, error) {

	req := a.URL + "price"
	contents, err := GetMethod(req)
	if err != nil {
		return 0, "", err
	}

	data := &stct.Price{}

	err = json.Unmarshal([]byte(contents), data)
	if err != nil {
		fmt.Println(err)
		return 0, "", err
	}

	currentPrice := float64(float64(data.Data.Price) / float64(10000))
	currnetDelta := float64(float64(data.Data.Delta) * float64(100))

	var delatastr string
	if currnetDelta > 0 {
		delatastr = fmt.Sprintf("+%.2f", currnetDelta)
	} else {
		delatastr = fmt.Sprintf("%.2f", currnetDelta)
	}

	return currentPrice, delatastr, nil
}

// --------------------------- Buy ----------------------------------
// -------------------------------- 1 --------------------------------

// GetBTCDeposAddress returns bitcoin address to deposit. (BUY coins)
func (a *App) GetBTCDeposAddress(minterAddress, coin, email string) (string, error) {

	req := a.URL + "bitcoinDepositAddress?minterAddress=" + minterAddress + "&coin=" + coin + "&email=" + email
	contents, err := GetMethod(req)
	if err != nil {
		return "", err
	}

	data := &stct.DeposBTC{}
	err = json.Unmarshal([]byte(contents), data)
	if err != nil {
		fmt.Println(err)
		return "", errors.New("Something going wrong, sorry:(")
	}

	return data.Data.Address, nil
}

// -------------------------------- 2 --------------------------------

// GetBTCDepositStatus returns the status for a given address.
func (a *App) GetBTCDepositStatus(bitcoinAddress string) (*stct.BTCStatus, error) {

	req := a.URL + "bitcoinAddressStatus?address=" + bitcoinAddress
	contents, err := GetMethod(req)
	if err != nil {
		return nil, err
	}
	data := &stct.BTCStatus{}

	err = json.Unmarshal([]byte(contents), data)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("Something going wrong, sorry:(")
	}

	return data, nil
}

// -------------------------------- Sell ----------------------------------
// -------------------------------- 1 --------------------------------

// GetMinterDeposAddress return deposit struct.
func (a *App) GetMinterDeposAddress(bitcoinAddress, coin string, price float64) (*stct.DeposMNT, error) {

	pricestr := fmt.Sprintf("%d", int(price*10000.))
	req := a.URL + "minterDepositAddress?bitcoinAddress=" + bitcoinAddress + "&price=" + pricestr + "&coin=" + coin
	contents, err := GetMethod(req)
	if err != nil {
		return nil, err
	}

	data := &stct.DeposMNT{}
	err = json.Unmarshal([]byte(contents), data)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("Something going wrong, sorry:(")
	}

	return data, nil
}

// -------------------------------- 2 --------------------------------

// GetTagInfo returns TagInfo struct.
func (a *App) GetTagInfo(tag string) (*stct.TagInfo, error) {

	req := a.URL + "tag?tag=" + tag
	contents, err := GetMethod(req)
	if err != nil {
		return nil, err
	}

	data := &stct.TagInfo{}

	err = json.Unmarshal([]byte(contents), data)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("Something going wrong, sorry:(")
	}

	r := new(big.Rat)
	r.SetString(data.Data.Amount + "/" + "1000000000000000000")

	data.Data.Amount = r.FloatString(1)

	data.Data.Price = float64(float64(data.Data.Price) / float64(10000))

	return data, nil
}

// BTCAddressHistory returns BTCAddress history
func (a *App) BTCAddressHistory(address string) (*stct.AddrHistory, error) {

	req := a.URL + "bitcoinAddressHistory?address=" + address
	return AddressHistory(req)
}

// MinterAddressHistory returns MinterAddress history
func (a *App) MinterAddressHistory(address string) (*stct.AddrHistory, error) {
	req := a.URL + "minterAddressHistory?address=" + address

	return AddressHistory(req)
}

// AddressHistory returns address history
func AddressHistory(req string) (*stct.AddrHistory, error) {

	response, err := http.Get(req)
	if err != nil {
		return nil, errors.New("http://bip.dev is not respond")
	}

	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("Something going wrong, sorry:(")
	}

	if response.StatusCode == 404 {
		data := &stct.Err{}
		err := json.Unmarshal([]byte(contents), data)
		if err != nil {
			fmt.Println(err)
			return nil, errors.New("Something going wrong, sorry:(")
		}
		return nil, errors.New(data.Error.Message)
	}

	data := &stct.AddrHistory{}
	err = json.Unmarshal([]byte(contents), data)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("Something going wrong, sorry:(")
	}

	return data, nil
}

// GetAvailablePrices
// func (a *App) GetAvailablePrices() ([]float64, error) {

// 	req := a.URL + "availablePrices"
// 	contents, err := GetMethod(req)
// 	if err != nil {
// 		return nil, err
// 	}

// 	data := &stct.Available{}

// 	err = json.Unmarshal([]byte(contents), data)
// 	if err != nil {
// 		fmt.Println(err)
// 		return nil, err
// 	}

// 	var currentData []float64

// 	for _, n := range data.Data {
// 		currentData = append(currentData, float64(float64(n)/10000))
// 	}

// 	return currentData, nil
// }
