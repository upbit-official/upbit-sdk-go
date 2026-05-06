// Auto-generated example — do not edit manually.
// Get Withdrawal (GET /v1/withdraw)
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/upbit-official/upbit-sdk-go"
	"github.com/upbit-official/upbit-sdk-go/option"
)

// default
func exampleGetWithdrawal() {
	client := upbit.NewClient(
		option.WithAccessKey(os.Getenv("UPBIT_ACCESS_KEY")),
		option.WithSecretKey(os.Getenv("UPBIT_SECRET_KEY")),
	)
	params := upbit.WithdrawGetParams{
		Currency: upbit.String("BTC"),
	}
	res, err := client.Withdraws.Get(context.TODO(), params)
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
}

func main() {
	exampleGetWithdrawal()
}
