// Auto-generated example — do not edit manually.
// List Day Candles (GET /v1/candles/days)
package main

import (
	"context"
	"fmt"

	"github.com/upbit-official/upbit-sdk-go"
)

// default
func exampleListCandlesDays() {
	client := upbit.NewClient()
	params := upbit.CandleListDaysParams{
		Market: "KRW-BTC",
	}
	res, err := client.Candles.ListDays(context.TODO(), params)
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
}

func main() {
	exampleListCandlesDays()
}
