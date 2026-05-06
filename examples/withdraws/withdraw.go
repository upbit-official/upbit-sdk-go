// Auto-generated example — do not edit manually.
// Withdraw Digital Asset (POST /v1/withdraws/coin)
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/upbit-official/upbit-sdk-go"
	"github.com/upbit-official/upbit-sdk-go/option"
)

// default
func exampleWithdraw() {
	client := upbit.NewClient(
		option.WithAccessKey(os.Getenv("UPBIT_ACCESS_KEY")),
		option.WithSecretKey(os.Getenv("UPBIT_SECRET_KEY")),
	)
	params := upbit.WithdrawNewWithdrawalParams{
		Currency: "BTC",
		NetType:  "BTC",
		Amount:   "0.01",
		Address:  "1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa",
	}
	res, err := client.Withdraws.NewWithdrawal(context.TODO(), params)
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
}

func main() {
	exampleWithdraw()
}
