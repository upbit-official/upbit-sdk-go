// Auto-generated example — do not edit manually.
// Cancel Orders by IDs (DELETE /v1/orders/uuids)
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/upbit-official/upbit-sdk-go"
	"github.com/upbit-official/upbit-sdk-go/option"
)

// default
func exampleCancelOrdersByIds() {
	client := upbit.NewClient(
		option.WithAccessKey(os.Getenv("UPBIT_ACCESS_KEY")),
		option.WithSecretKey(os.Getenv("UPBIT_SECRET_KEY")),
	)
	params := upbit.OrderCancelByUuidsParams{
		Uuids: []string{"bbbb8e07-1689-4769-af3e-a117016623f8", "4312ba49-5f1a-4a01-9f3b-2d2bce17267e", "bdb49a54-de36-4eb4-a963-9c8d4337a9da"},
	}
	res, err := client.Orders.CancelByUuids(context.TODO(), params)
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
}

func main() {
	exampleCancelOrdersByIds()
}
