// Auto-generated example — do not edit manually.
// Cancel Withdrawal (DELETE /v1/withdraws/coin)
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/upbit-official/upbit-sdk-go"
	"github.com/upbit-official/upbit-sdk-go/option"
)

// default
func exampleCancelWithdrawal() {
	client := upbit.NewClient(
		option.WithAccessKey(os.Getenv("UPBIT_ACCESS_KEY")),
		option.WithSecretKey(os.Getenv("UPBIT_SECRET_KEY")),
	)
	params := upbit.WithdrawCancelWithdrawalParams{
		Uuid: "9f432943-54e0-40b7-825f-b6fec8b42b79",
	}
	res, err := client.Withdraws.CancelWithdrawal(context.TODO(), params)
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
}

func main() {
	exampleCancelWithdrawal()
}
