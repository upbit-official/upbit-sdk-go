package upbit_test

import (
	"context"
	"os"
	"testing"

	"github.com/upbit-official/upbit-sdk-go"
	"github.com/upbit-official/upbit-sdk-go/internal/testutil"
	"github.com/upbit-official/upbit-sdk-go/option"
)

func TestManualPagination(t *testing.T) {
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
	page, err := client.Orders.ListOpen(context.TODO(), upbit.OrderListOpenParams{})
	if err != nil {
		t.Fatalf("err should be nil: %s", err.Error())
	}
	for _, order := range page.Items {
		t.Logf("%+v\n", order.PaidFee)
	}
	// The mock server isn't going to give us real pagination
	page, err = page.GetNextPage()
	if err != nil {
		t.Fatalf("err should be nil: %s", err.Error())
	}
	if page != nil {
		for _, order := range page.Items {
			t.Logf("%+v\n", order.PaidFee)
		}
	}
}
