// Auto-generated example — do not edit manually.
// List Tickers by Pairs (GET /v1/ticker)
package main

import (
	"context"
	"fmt"

	"github.com/upbit-official/upbit-sdk-go"
)

// default
func exampleListTickers() {
	client := upbit.NewClient()
	params := upbit.TickerListByTradingPairsParams{
		Markets: "KRW-BTC,KRW-ETH",
	}
	res, err := client.Tickers.ListByTradingPairs(context.TODO(), params)
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
}

func main() {
	exampleListTickers()
}
