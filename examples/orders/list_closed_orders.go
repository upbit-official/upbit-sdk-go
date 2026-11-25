// Auto-generated example — do not edit manually.
// List Closed Orders (GET /v1/orders/closed)
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/upbit-official/upbit-sdk-go"
	"github.com/upbit-official/upbit-sdk-go/option"
)

// default
func exampleListClosedOrders() {
	client := upbit.NewClient(
		option.WithAccessKey(os.Getenv("UPBIT_ACCESS_KEY")),
		option.WithSecretKey(os.Getenv("UPBIT_SECRET_KEY")),
	)
	params := upbit.OrderListClosedParams{
		Market: upbit.String("KRW-BTC"),
	}
	res, err := client.Orders.ListClosed(context.TODO(), params)
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
}

func main() {
	exampleListClosedOrders()
}
