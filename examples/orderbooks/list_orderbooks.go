// Auto-generated example — do not edit manually.
// Get Orderbook (GET /v1/orderbook)
package main

import (
	"context"
	"fmt"

	"github.com/upbit-official/upbit-sdk-go"
)

// default
func exampleListOrderbooks() {
	client := upbit.NewClient()
	params := upbit.OrderbookListParams{
		Markets: "KRW-BTC,KRW-ETH",
	}
	res, err := client.Orderbooks.List(context.TODO(), params)
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
}

func main() {
	exampleListOrderbooks()
}
