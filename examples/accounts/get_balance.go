// Auto-generated example — do not edit manually.
// Get Account Balances (GET /v1/accounts)
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/upbit-official/upbit-sdk-go"
	"github.com/upbit-official/upbit-sdk-go/option"
)

// default
func exampleGetBalance() {
	client := upbit.NewClient(
		option.WithAccessKey(os.Getenv("UPBIT_ACCESS_KEY")),
		option.WithSecretKey(os.Getenv("UPBIT_SECRET_KEY")),
	)
	res, err := client.Accounts.List(context.TODO())
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
}

func main() {
	exampleGetBalance()
}
