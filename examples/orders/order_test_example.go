// Auto-generated example — do not edit manually.
// Test Order (POST /v1/orders/test)
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/upbit-official/upbit-sdk-go"
	"github.com/upbit-official/upbit-sdk-go/option"
)

// default
func exampleOrderTest() {
	client := upbit.NewClient(
		option.WithAccessKey(os.Getenv("UPBIT_ACCESS_KEY")),
		option.WithSecretKey(os.Getenv("UPBIT_SECRET_KEY")),
	)
	params := upbit.OrderTestNewParams{
		Market:  "KRW-BTC",
		Side:    upbit.SideEnumBid,
		Volume:  upbit.String("1"),
		Price:   upbit.String("14000000"),
		OrdType: upbit.OrdTypeEnumLimit,
	}
	res, err := client.Orders.TestNew(context.TODO(), params)
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
}

func main() {
	exampleOrderTest()
}
