// Auto-generated example — do not edit manually.
// List Pair Trades (GET /v1/trades/ticks)
package main

import (
	"context"
	"fmt"

	"github.com/upbit-official/upbit-sdk-go"
)

// default
func exampleListPairTrades() {
	client := upbit.NewClient()
	params := upbit.TradeListParams{
		Market: "KRW-BTC",
	}
	res, err := client.Trades.List(context.TODO(), params)
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
}

func main() {
	exampleListPairTrades()
}
