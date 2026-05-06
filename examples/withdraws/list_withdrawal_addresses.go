// Auto-generated example — do not edit manually.
// List Withdrawal Allowed Addresses (GET /v1/withdraws/coin_addresses)
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/upbit-official/upbit-sdk-go"
	"github.com/upbit-official/upbit-sdk-go/option"
)

// default
func exampleListWithdrawalAddresses() {
	client := upbit.NewClient(
		option.WithAccessKey(os.Getenv("UPBIT_ACCESS_KEY")),
		option.WithSecretKey(os.Getenv("UPBIT_SECRET_KEY")),
	)
	res, err := client.Withdraws.ListCoinAddresses(context.TODO())
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
}

func main() {
	exampleListWithdrawalAddresses()
}
