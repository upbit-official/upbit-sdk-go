// Auto-generated example — do not edit manually.
// List Travel Rule Supporting VASPs (GET /v1/travel_rule/vasps)
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/upbit-official/upbit-sdk-go"
	"github.com/upbit-official/upbit-sdk-go/option"
)

// default
func exampleListTravelruleVasps() {
	client := upbit.NewClient(
		option.WithAccessKey(os.Getenv("UPBIT_ACCESS_KEY")),
		option.WithSecretKey(os.Getenv("UPBIT_SECRET_KEY")),
	)
	res, err := client.TravelRule.ListVasps(context.TODO())
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
}

func main() {
	exampleListTravelruleVasps()
}
