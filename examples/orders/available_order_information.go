// Auto-generated example — do not edit manually.
// Get Available Order Info (GET /v1/orders/chance)
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/upbit-official/upbit-sdk-go"
	"github.com/upbit-official/upbit-sdk-go/option"
)

// default
func exampleAvailableOrderInformation() {
	client := upbit.NewClient(
		option.WithAccessKey(os.Getenv("UPBIT_ACCESS_KEY")),
		option.WithSecretKey(os.Getenv("UPBIT_SECRET_KEY")),
	)
	params := upbit.OrderGetChanceParams{
		Market: "KRW-BTC",
	}
	res, err := client.Orders.GetChance(context.TODO(), params)
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
}

func main() {
	exampleAvailableOrderInformation()
}
