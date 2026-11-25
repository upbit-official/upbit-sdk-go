// Auto-generated example — do not edit manually.
// Cancel and New Order (POST /v1/orders/cancel_and_new)
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/upbit-official/upbit-sdk-go"
	"github.com/upbit-official/upbit-sdk-go/option"
)

// default
func exampleCancelAndNewOrder() {
	client := upbit.NewClient(
		option.WithAccessKey(os.Getenv("UPBIT_ACCESS_KEY")),
		option.WithSecretKey(os.Getenv("UPBIT_SECRET_KEY")),
	)
	params := upbit.OrderCancelAndNewParams{
		PrevOrderUuid: upbit.String("ad217e24-ed02-469c-9b30-c08dbbda6908"),
		NewOrdType:    upbit.OrdTypeEnumLimit,
		NewVolume:     upbit.String("1"),
		NewPrice:      upbit.String("100000000"),
	}
	res, err := client.Orders.CancelAndNew(context.TODO(), params)
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
}

func main() {
	exampleCancelAndNewOrder()
}
