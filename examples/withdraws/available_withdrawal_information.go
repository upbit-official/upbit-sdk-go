// Auto-generated example — do not edit manually.
// Get Available Withdrawal Information (GET /v1/withdraws/chance)
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/upbit-official/upbit-sdk-go"
	"github.com/upbit-official/upbit-sdk-go/option"
)

// default
func exampleAvailableWithdrawalInformation() {
	client := upbit.NewClient(
		option.WithAccessKey(os.Getenv("UPBIT_ACCESS_KEY")),
		option.WithSecretKey(os.Getenv("UPBIT_SECRET_KEY")),
	)
	params := upbit.WithdrawGetChanceParams{
		Currency: "BTC",
		NetType:  "BTC",
	}
	res, err := client.Withdraws.GetChance(context.TODO(), params)
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
}

func main() {
	exampleAvailableWithdrawalInformation()
}
