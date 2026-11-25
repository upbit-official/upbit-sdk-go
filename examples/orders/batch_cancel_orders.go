// Auto-generated example — do not edit manually.
// Batch Cancel Orders (DELETE /v1/orders/open)
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/upbit-official/upbit-sdk-go"
	"github.com/upbit-official/upbit-sdk-go/option"
)

// default
func exampleBatchCancelOrders() {
	client := upbit.NewClient(
		option.WithAccessKey(os.Getenv("UPBIT_ACCESS_KEY")),
		option.WithSecretKey(os.Getenv("UPBIT_SECRET_KEY")),
	)
	params := upbit.OrderCancelOpenParams{
		QuoteCurrencies: upbit.String("KRW,BTC"),
		ExcludedPairs:   upbit.String("KRW-ETH,BTC-XRP"),
	}
	res, err := client.Orders.CancelOpen(context.TODO(), params)
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
}

func main() {
	exampleBatchCancelOrders()
}
