// Auto-generated example — do not edit manually.
// List Week Candles (GET /v1/candles/weeks)
package main

import (
	"context"
	"fmt"

	"github.com/upbit-official/upbit-sdk-go"
)

// default
func exampleListCandlesWeeks() {
	client := upbit.NewClient()
	params := upbit.CandleListWeeksParams{
		Market: "KRW-BTC",
	}
	res, err := client.Candles.ListWeeks(context.TODO(), params)
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
}

func main() {
	exampleListCandlesWeeks()
}
