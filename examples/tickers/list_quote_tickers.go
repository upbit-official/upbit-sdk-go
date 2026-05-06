// Auto-generated example — do not edit manually.
// List Tickers by Market (GET /v1/ticker/all)
package main

import (
	"context"
	"fmt"

	"github.com/upbit-official/upbit-sdk-go"
)

// default
func exampleListQuoteTickers() {
	client := upbit.NewClient()
	params := upbit.TickerListByQuoteCurrenciesParams{
		QuoteCurrencies: "KRW,BTC,USDT",
	}
	res, err := client.Tickers.ListByQuoteCurrencies(context.TODO(), params)
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
}

func main() {
	exampleListQuoteTickers()
}
