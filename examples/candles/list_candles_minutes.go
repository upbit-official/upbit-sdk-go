// Auto-generated example — do not edit manually.
// List Minute Candles (GET /v1/candles/minutes/{unit})
package main

import (
	"context"
	"fmt"

	"github.com/upbit-official/upbit-sdk-go"
)

// default
func exampleListCandlesMinutes() {
	client := upbit.NewClient()
	params := upbit.CandleListMinutesParams{
		Market: "KRW-BTC",
	}
	res, err := client.Candles.ListMinutes(context.TODO(), 15, params)
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
}

func main() {
	exampleListCandlesMinutes()
}
