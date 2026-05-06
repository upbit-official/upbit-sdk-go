// Auto-generated example — do not edit manually.
// List Trading Pairs (GET /v1/market/all)
package main

import (
	"context"
	"fmt"

	"github.com/upbit-official/upbit-sdk-go"
)

// default
func exampleListTradingPairs() {
	client := upbit.NewClient()
	params := upbit.TradingPairListParams{
		IsDetails: upbit.Bool(false),
	}
	res, err := client.TradingPairs.List(context.TODO(), params)
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
}

func main() {
	exampleListTradingPairs()
}
