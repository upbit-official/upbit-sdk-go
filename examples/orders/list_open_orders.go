// Auto-generated example — do not edit manually.
// List Open Orders (GET /v1/orders/open)
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/upbit-official/upbit-sdk-go"
	"github.com/upbit-official/upbit-sdk-go/option"
)

// default
func exampleListOpenOrders() {
	client := upbit.NewClient(
		option.WithAccessKey(os.Getenv("UPBIT_ACCESS_KEY")),
		option.WithSecretKey(os.Getenv("UPBIT_SECRET_KEY")),
	)
	params := upbit.OrderListOpenParams{
		Market: upbit.String("KRW-BTC"),
		States: []upbit.OpenOrderState{"wait", "watch"},
	}
	res, err := client.Orders.ListOpen(context.TODO(), params)
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
}

func main() {
	exampleListOpenOrders()
}
