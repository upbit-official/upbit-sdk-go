// Auto-generated example — do not edit manually.
// List Orderbook Instruments (GET /v1/orderbook/instruments)
package main

import (
	"context"
	"fmt"

	"github.com/upbit-official/upbit-sdk-go"
)

// default
func exampleListOrderbookInstruments() {
	client := upbit.NewClient()
	params := upbit.OrderbookListInstrumentsParams{
		Markets: "KRW-BTC,KRW-ETH",
	}
	res, err := client.Orderbooks.ListInstruments(context.TODO(), params)
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
}

func main() {
	exampleListOrderbookInstruments()
}
