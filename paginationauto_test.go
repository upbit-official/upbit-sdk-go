package upbit_test

import (
	"context"
	"os"
	"testing"

	"github.com/upbit-official/upbit-sdk-go"
	"github.com/upbit-official/upbit-sdk-go/internal/testutil"
	"github.com/upbit-official/upbit-sdk-go/option"
)

func TestAutoPagination(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	baseURL := "http://localhost:4010"
	if envURL, ok := os.LookupEnv("TEST_API_BASE_URL"); ok {
		baseURL = envURL
	}
	if !testutil.CheckTestServer(t, baseURL) {
		return
	}
	client := upbit.NewClient(
		option.WithBaseURL(baseURL),
		option.WithAccessKey("My Access Key"),
	)
	iter := client.Orders.ListOpenAutoPaging(context.TODO(), upbit.OrderListOpenParams{})
	// The mock server isn't going to give us real pagination
	for i := 0; i < 3 && iter.Next(); i++ {
		order := iter.Current()
		t.Logf("%+v\n", order.PaidFee)
	}
	if err := iter.Err(); err != nil {
		t.Fatalf("err should be nil: %s", err.Error())
	}
}
