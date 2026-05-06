// Auto-generated example — do not edit manually.
// Get Order (GET /v1/order)
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/upbit-official/upbit-sdk-go"
	"github.com/upbit-official/upbit-sdk-go/option"
)

// default
func exampleGetOrder() {
	client := upbit.NewClient(
		option.WithAccessKey(os.Getenv("UPBIT_ACCESS_KEY")),
		option.WithSecretKey(os.Getenv("UPBIT_SECRET_KEY")),
	)
	params := upbit.OrderGetParams{
		Uuid: upbit.String("9ca023a5-851b-4fec-9f0a-48cd83c2eaae"),
	}
	res, err := client.Orders.Get(context.TODO(), params)
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
}

// by_identifier
func exampleGetOrderByIdentifier() {
	client := upbit.NewClient(
		option.WithAccessKey(os.Getenv("UPBIT_ACCESS_KEY")),
		option.WithSecretKey(os.Getenv("UPBIT_SECRET_KEY")),
	)
	params := upbit.OrderGetParams{
		Identifier: upbit.String("9ca023a5-851b-4fec-9f0a-48cd83c2eaae"),
	}
	res, err := client.Orders.Get(context.TODO(), params)
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
}

func main() {
	exampleGetOrder()
	exampleGetOrderByIdentifier()
}
