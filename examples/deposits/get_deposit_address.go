// Auto-generated example — do not edit manually.
// Get Deposit Address (GET /v1/deposits/coin_address)
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/upbit-official/upbit-sdk-go"
	"github.com/upbit-official/upbit-sdk-go/option"
)

// default
func exampleGetDepositAddress() {
	client := upbit.NewClient(
		option.WithAccessKey(os.Getenv("UPBIT_ACCESS_KEY")),
		option.WithSecretKey(os.Getenv("UPBIT_SECRET_KEY")),
	)
	params := upbit.DepositGetCoinAddressParams{
		Currency: "BTC",
		NetType:  "BTC",
	}
	res, err := client.Deposits.GetCoinAddress(context.TODO(), params)
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
}

func main() {
	exampleGetDepositAddress()
}
