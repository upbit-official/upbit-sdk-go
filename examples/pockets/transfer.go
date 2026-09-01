// Auto-generated example — do not edit manually.
// Transfer (POST /v1/pockets/transfers)
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/upbit-official/upbit-sdk-go"
	"github.com/upbit-official/upbit-sdk-go/option"
)

// default
func exampleTransfer() {
	client := upbit.NewClient(
		option.WithAccessKey(os.Getenv("UPBIT_ACCESS_KEY")),
		option.WithSecretKey(os.Getenv("UPBIT_SECRET_KEY")),
	)
	params := upbit.PocketTransferParams{
		To:         "9ca023a5-851b-4fec-9f0a-48cd83c2eaae",
		Currency:   "XRP",
		Amount:     "10.00",
		Identifier: upbit.String("unique_transfer_id_20260531_01"),
	}
	res, err := client.Pockets.Transfer(context.TODO(), params)
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
}

func main() {
	exampleTransfer()
}
