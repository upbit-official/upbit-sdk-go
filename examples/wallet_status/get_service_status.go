// Auto-generated example — do not edit manually.
// Get Deposit/Withdrawal Service Status (GET /v1/status/wallet)
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/upbit-official/upbit-sdk-go"
	"github.com/upbit-official/upbit-sdk-go/option"
)

// default
func exampleGetServiceStatus() {
	client := upbit.NewClient(
		option.WithAccessKey(os.Getenv("UPBIT_ACCESS_KEY")),
		option.WithSecretKey(os.Getenv("UPBIT_SECRET_KEY")),
	)
	res, err := client.WalletStatus.List(context.TODO())
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
}

func main() {
	exampleGetServiceStatus()
}
