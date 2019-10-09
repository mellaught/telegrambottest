package main

// import (
// 	"fmt"

// 	_ "github.com/jinzhu/gorm/dialects/postgres"
// )

// var (
// 	str1 = "1000000000000000000000000"
// 	str2 = "100000"
// )

// func main() {
// 	address := "Mx3fdca30cb0ffb8ae1096acf60f5d0e0d908a7f75"
// 	fmt.Println(address[:2],)
// 	fmt.Println(len(address) != 42 || address[:2] != "Mx")
// 	//_ := StrToBig("100000000000000000000000")
// 	// r := new(big.Rat)
// 	// r.SetString(str1 + "/" + "1000000000000000000")
// 	// fmt.Println(str2, r.FloatString(1))
// 	// floatValue := new(big.Float).SetPrec(500).SetInt(value)
// 	// f, _ := new(big.Float).SetPrec(500).Quo(floatValue, pipInBip).Float64()

// 	// fmt.Println(humanize.FormatFloat("# ###.##", f))
// }

// // func StrToBig(balance string) *big.Int {
// // 	bigInt, success := big.NewInt(0).SetString(balance, 10)
// // 	if success != true {
// // 		panic("Failed to decode " + balance)
// // 	}

// // 	return bigInt
// // }

// // func PipToBip(value *big.Int) string {
// // 	floatValue := new(big.Float).SetPrec(500).SetInt(value)
// // 	f, _ := new(big.Float).SetPrec(500).Quo(floatValue, pipInBip).Float64()

// // 	return humanize.FormatFloat("# ###.##", f)
// // }
