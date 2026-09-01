// Auto-generated example — do not edit manually.
// List Transfers (GET /v1/pockets/transfers)
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/upbit-official/upbit-sdk-go"
	"github.com/upbit-official/upbit-sdk-go/option"
)

// default
func exampleListTransfers() {
	client := upbit.NewClient(
		option.WithAccessKey(os.Getenv("UPBIT_ACCESS_KEY")),
		option.WithSecretKey(os.Getenv("UPBIT_SECRET_KEY")),
	)
	params := upbit.PocketListTransfersParams{
		Direction: upbit.PocketTransferDirectionOut,
		Currency:  upbit.String("XRP"),
		Limit:     upbit.Int(2),
	}
	res, err := client.Pockets.ListTransfers(context.TODO(), params)
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
}

func main() {
	exampleListTransfers()
}
