// Auto-generated example — do not edit manually.
// Get Available Deposit Information (GET /v1/deposits/chance/coin)
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/upbit-official/upbit-sdk-go"
	"github.com/upbit-official/upbit-sdk-go/option"
)

// default
func exampleAvailableDepositInformation() {
	client := upbit.NewClient(
		option.WithAccessKey(os.Getenv("UPBIT_ACCESS_KEY")),
		option.WithSecretKey(os.Getenv("UPBIT_SECRET_KEY")),
	)
	params := upbit.DepositGetChanceParams{
		Currency: "BTC",
		NetType:  "BTC",
	}
	res, err := client.Deposits.GetChance(context.TODO(), params)
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
}

func main() {
	exampleAvailableDepositInformation()
}
