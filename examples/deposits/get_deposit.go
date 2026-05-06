// Auto-generated example — do not edit manually.
// Get Deposit (GET /v1/deposit)
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/upbit-official/upbit-sdk-go"
	"github.com/upbit-official/upbit-sdk-go/option"
)

// default
func exampleGetDeposit() {
	client := upbit.NewClient(
		option.WithAccessKey(os.Getenv("UPBIT_ACCESS_KEY")),
		option.WithSecretKey(os.Getenv("UPBIT_SECRET_KEY")),
	)
	params := upbit.DepositGetParams{
		Uuid: upbit.String("94332e99-3a87-4a35-ad98-28b0c969f830"),
	}
	res, err := client.Deposits.Get(context.TODO(), params)
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
}

func main() {
	exampleGetDeposit()
}
