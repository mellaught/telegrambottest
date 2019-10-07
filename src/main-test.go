package main

// import (
// 	"math/big"

// 	_ "github.com/jinzhu/gorm/dialects/postgres"
// 	api "github.com/mrKitikat/telegrambottest/src/app/bipdev"
// )

// var (
// 	str1 = "1000000000000000000000000"
// 	str2 = "100000"
// )

// func main() {

// 	app := api.InitApp("https://api.bip.dev/api/")

// 	app.GetBonus()
// 	//_ := StrToBig("100000000000000000000000")
// 	// r := new(big.Rat)
// 	// r.SetString(str1 + "/" + "1000000000000000000")
// 	// fmt.Println(str2, r.FloatString(1))
// 	// floatValue := new(big.Float).SetPrec(500).SetInt(value)
// 	// f, _ := new(big.Float).SetPrec(500).Quo(floatValue, pipInBip).Float64()

// 	// fmt.Println(humanize.FormatFloat("# ###.##", f))
// }

// func StrToBig(balance string) *big.Int {
// 	bigInt, success := big.NewInt(0).SetString(balance, 10)
// 	if success != true {
// 		panic("Failed to decode " + balance)
// 	}

// 	return bigInt
// }

// func PipToBip(value *big.Int) string {
// 	floatValue := new(big.Float).SetPrec(500).SetInt(value)
// 	f, _ := new(big.Float).SetPrec(500).Quo(floatValue, pipInBip).Float64()

// 	return humanize.FormatFloat("# ###.##", f)
// }
