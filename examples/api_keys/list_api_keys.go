// Auto-generated example — do not edit manually.
// List API Keys (GET /v1/api_keys)
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/upbit-official/upbit-sdk-go"
	"github.com/upbit-official/upbit-sdk-go/option"
)

// default
func exampleListAPIKeys() {
	client := upbit.NewClient(
		option.WithAccessKey(os.Getenv("UPBIT_ACCESS_KEY")),
		option.WithSecretKey(os.Getenv("UPBIT_SECRET_KEY")),
	)
	res, err := client.APIKeys.List(context.TODO())
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
}

func main() {
	exampleListAPIKeys()
}
