// Auto-generated example — do not edit manually.
// Verify Travel Rule by Deposit UUID (POST /v1/travel_rule/deposit/uuid)
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/upbit-official/upbit-sdk-go"
	"github.com/upbit-official/upbit-sdk-go/option"
)

// default
func exampleVerifyTravelruleByUuid() {
	client := upbit.NewClient(
		option.WithAccessKey(os.Getenv("UPBIT_ACCESS_KEY")),
		option.WithSecretKey(os.Getenv("UPBIT_SECRET_KEY")),
	)
	params := upbit.TravelRuleVerifyDepositByUuidParams{
		DepositUuid: "5b871d34-fe38-4025-8f5c-9b22028f85d3",
		VaspUuid:    "8d4fe968-82b2-42e5-822f-3840a245f802",
	}
	res, err := client.TravelRule.VerifyDepositByUuid(context.TODO(), params)
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
}

func main() {
	exampleVerifyTravelruleByUuid()
}
