// Auto-generated example — do not edit manually.
// Create Deposit Address (POST /v1/deposits/generate_coin_address)
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/upbit-official/upbit-sdk-go"
	"github.com/upbit-official/upbit-sdk-go/option"
)

// default
func exampleCreateDepositAddress() {
	client := upbit.NewClient(
		option.WithAccessKey(os.Getenv("UPBIT_ACCESS_KEY")),
		option.WithSecretKey(os.Getenv("UPBIT_SECRET_KEY")),
	)
	params := upbit.DepositNewCoinAddressParams{
		Currency: "BTC",
		NetType:  "BTC",
	}
	res, err := client.Deposits.NewCoinAddress(context.TODO(), params)
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
}

func main() {
	exampleCreateDepositAddress()
}
