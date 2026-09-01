// Auto-generated example — do not edit manually.
// List Universal Transfers (GET /v1/pockets/universal_transfers)
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/upbit-official/upbit-sdk-go"
	"github.com/upbit-official/upbit-sdk-go/option"
)

// default
func exampleListUniversalTransfers() {
	client := upbit.NewClient(
		option.WithAccessKey(os.Getenv("UPBIT_ACCESS_KEY")),
		option.WithSecretKey(os.Getenv("UPBIT_SECRET_KEY")),
	)
	params := upbit.PocketListUniversalTransfersParams{
		States:   []upbit.PocketTransferState{"processing", "done"},
		Currency: upbit.String("XRP"),
		Limit:    upbit.Int(2),
	}
	res, err := client.Pockets.ListUniversalTransfers(context.TODO(), params)
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
}

func main() {
	exampleListUniversalTransfers()
}
