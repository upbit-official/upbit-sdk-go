// Auto-generated example — do not edit manually.
// Get Sub-Pocket Balance (GET /v1/pockets/assets)
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/upbit-official/upbit-sdk-go"
	"github.com/upbit-official/upbit-sdk-go/option"
)

// default
func exampleGetSubPocketBalance() {
	client := upbit.NewClient(
		option.WithAccessKey(os.Getenv("UPBIT_ACCESS_KEY")),
		option.WithSecretKey(os.Getenv("UPBIT_SECRET_KEY")),
	)
	params := upbit.PocketGetBalanceParams{
		Uuid: "9ca023a5-851b-4fec-9f0a-48cd83c2eaae",
	}
	res, err := client.Pockets.GetBalance(context.TODO(), params)
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
}

func main() {
	exampleGetSubPocketBalance()
}
