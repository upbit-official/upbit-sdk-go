package upbit_test

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/upbit-official/upbit-sdk-go"
	"github.com/upbit-official/upbit-sdk-go/internal/testutil"
	"github.com/upbit-official/upbit-sdk-go/option"
)

func TestOrderbookListWithOptionalParams(t *testing.T) {
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
	_, err := client.Orderbooks.List(context.TODO(), upbit.OrderbookListParams{
		Markets: "KRW-BTC,KRW-ETH",
		Count:   upbit.Int(10),
		Level:   upbit.String("1"),
	})
	if err != nil {
		var apierr *upbit.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestOrderbookListInstruments(t *testing.T) {
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
	_, err := client.Orderbooks.ListInstruments(context.TODO(), upbit.OrderbookListInstrumentsParams{
		Markets: "KRW-BTC,KRW-ETH",
	})
	if err != nil {
		var apierr *upbit.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}
