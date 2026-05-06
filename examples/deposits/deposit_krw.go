// Auto-generated example — do not edit manually.
// Deposit KRW (POST /v1/deposits/krw)
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/upbit-official/upbit-sdk-go"
	"github.com/upbit-official/upbit-sdk-go/option"
)

// default
func exampleDepositKrw() {
	client := upbit.NewClient(
		option.WithAccessKey(os.Getenv("UPBIT_ACCESS_KEY")),
		option.WithSecretKey(os.Getenv("UPBIT_SECRET_KEY")),
	)
	params := upbit.DepositDepositKrwParams{
		Amount:        "10000",
		TwoFactorType: "naver",
	}
	res, err := client.Deposits.DepositKrw(context.TODO(), params)
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
}

func main() {
	exampleDepositKrw()
}
