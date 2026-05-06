// Auto-generated example — do not edit manually.
// List Year Candles (GET /v1/candles/years)
package main

import (
	"context"
	"fmt"

	"github.com/upbit-official/upbit-sdk-go"
)

// default
func exampleListCandlesYears() {
	client := upbit.NewClient()
	params := upbit.CandleListYearsParams{
		Market: "KRW-BTC",
	}
	res, err := client.Candles.ListYears(context.TODO(), params)
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
}

func main() {
	exampleListCandlesYears()
}
