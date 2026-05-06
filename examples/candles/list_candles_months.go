// Auto-generated example — do not edit manually.
// List Month Candles (GET /v1/candles/months)
package main

import (
	"context"
	"fmt"

	"github.com/upbit-official/upbit-sdk-go"
)

// default
func exampleListCandlesMonths() {
	client := upbit.NewClient()
	params := upbit.CandleListMonthsParams{
		Market: "KRW-BTC",
	}
	res, err := client.Candles.ListMonths(context.TODO(), params)
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
}

func main() {
	exampleListCandlesMonths()
}
