// Auto-generated example — do not edit manually.
// Cancel Order (DELETE /v1/order)
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/upbit-official/upbit-sdk-go"
	"github.com/upbit-official/upbit-sdk-go/option"
)

// default
func exampleCancelOrder() {
	client := upbit.NewClient(
		option.WithAccessKey(os.Getenv("UPBIT_ACCESS_KEY")),
		option.WithSecretKey(os.Getenv("UPBIT_SECRET_KEY")),
	)
	params := upbit.OrderCancelParams{
		Uuid: upbit.String("cdd92199-2897-4e14-9448-f923320408ad"),
	}
	res, err := client.Orders.Cancel(context.TODO(), params)
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
}

func main() {
	exampleCancelOrder()
}
