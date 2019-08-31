package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	stct "telegrambottest/src/bipdev/structs"
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

// BadReq func for statusCode == 400 , for GetBTCDeposAddress
func BadReq(contents []byte) (string, error) {
	data := &stct.Err{}
	err := json.Unmarshal([]byte(contents), data)
	if err != nil {
		return "", err
	}
	return data.Error.Message, errors.New(data.Error.Message)
}

// GetPrice return current price BIP/USD
func (a *App) GetPrice() (float64, error) {

	response, err := http.Get(a.URL + "price")
	if err != nil {
		return -1., errors.New("http://bip.dev is not respond")
	}

	defer response.Body.Close()

	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
		return -1., errors.New("Something going wrong, sorry:(")
	}

	data := &stct.Price{}

	err = json.Unmarshal([]byte(contents), data)
	if err != nil {
		fmt.Println(err)
		return -1., errors.New("Something going wrong, sorry:(")
	}

	currentPrice := float64(data.Data.Price / 1000)

	return currentPrice, nil

}

// --------------------------- Buy ----------------------------------
// -------------------------------- 1 --------------------------------

// GetBTCDeposAddress returns bitcoin address to deposit. (BUY coins)
func (a *App) GetBTCDeposAddress(minterAddress, coin, email string) (string, error) {

	req := a.URL + "bitcoinDepositAddress?minterAddress=" + minterAddress + "&coin=" + coin + "&email=" + email
	response, err := http.Get(req)
	if err != nil {
		return "", errors.New("http://bip.dev is not respond")
	}

	defer response.Body.Close()

	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
		return "", errors.New("Something going wrong, sorry:(")
	}

	if response.StatusCode == 400 {
		return BadReq(contents)
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
			return nil, err
		}
		return nil, errors.New(data.Error.Message)
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

	if response.StatusCode == 400 {
		data := &stct.Err{}
		err := json.Unmarshal([]byte(contents), data)
		if err != nil {
			return nil, err
		}
		return nil, errors.New(data.Error.Message)
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

	data := &stct.TagInfo{}
	err = json.Unmarshal([]byte(contents), data)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("Something going wrong, sorry:(")
	}

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

<<<<<<< HEAD
func (a *App) CheckStatus(address string, wg *sync.WaitGroup) {
	defer wg.Done()
	willcoin := 0.
	for {
		stat, err := a.GetBTCDepositStatus(address)
		if err != nil {
			log.Fatal(err)
			return
		}
		if stat.Data.WillReceive != willcoin {
			// Сообщить что придет столько то и
			if willcoin == 0. {
				// b.Bot.Send(msg)
				fmt.Println(stat.Data.WillReceive)
				willcoin = stat.Data.WillReceive
				time.Sleep(60 * time.Second)
			} else {
				// b.Bot.Send(msg)
				return
			}
		}

		fmt.Println(stat.Data.WillReceive)
		time.Sleep(3 * time.Second)
	}
}

// func (a *App) CheckStatus(address string) {

// 	for {
// 		stat, err := a.GetBTCDepositStatus(address)
// 		if err != nil {
// 			log.Fatal(err)
// 			return
// 		}
// 		if stat.Data.WillReceive != 0 {
// 			// Сообщить что придет столько то и
// 			fmt.Println(stat.Data.WillReceive)
// 			return
// 		}
// 		fmt.Println(stat.Data.WillReceive)
// 		time.Sleep(2 * time.Second)
// 	}
// }
=======
>>>>>>> 97af52583c4354e0e85352890f1f573f1701a764
