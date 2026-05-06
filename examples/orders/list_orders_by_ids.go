// Auto-generated example — do not edit manually.
// List Orders by IDs (GET /v1/orders/uuids)
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/upbit-official/upbit-sdk-go"
	"github.com/upbit-official/upbit-sdk-go/option"
)

// default
func exampleListOrdersByIds() {
	client := upbit.NewClient(
		option.WithAccessKey(os.Getenv("UPBIT_ACCESS_KEY")),
		option.WithSecretKey(os.Getenv("UPBIT_SECRET_KEY")),
	)
	params := upbit.OrderListByUuidsParams{
		Uuids: []string{"5d303952-8be9-41e6-915b-121a90026248", "3944c2c1-bd8c-441a-aa25-2370d08217a9", "5b95451b-971e-4e76-8f61-5ff441f078d5", "3b67e543-8ad3-48d0-8451-0dad315cae73"},
	}
	res, err := client.Orders.ListByUuids(context.TODO(), params)
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
}

func main() {
	exampleListOrdersByIds()
}
