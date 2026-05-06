// Auto-generated example — do not edit manually.
// List Second Candles (GET /v1/candles/seconds)
package main

import (
	"context"
	"fmt"

	"github.com/upbit-official/upbit-sdk-go"
)

// default
func exampleListCandlesSeconds() {
	client := upbit.NewClient()
	params := upbit.CandleListSecondsParams{
		Market: "KRW-BTC",
	}
	res, err := client.Candles.ListSeconds(context.TODO(), params)
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
}

func main() {
	exampleListCandlesSeconds()
}
