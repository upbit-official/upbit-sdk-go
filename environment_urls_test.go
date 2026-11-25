package upbit_test

import (
	"context"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"

	upbit "github.com/upbit-official/upbit-sdk-go"
	"github.com/upbit-official/upbit-sdk-go/option"
)

func TestEnvironmentRequestURLs(t *testing.T) {
	t.Setenv("UPBIT_BASE_URL", "")
	if err := os.Unsetenv("UPBIT_BASE_URL"); err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		name string
		env  string
		opts []option.RequestOption
		host string
	}{
		{"default", "", nil, "api.upbit.com"},
		{"kr", "", []option.RequestOption{option.WithEnvironmentKr()}, "api.upbit.com"},
		{"sg", "", []option.RequestOption{option.WithEnvironmentSg()}, "sg-api.upbit.com"},
		{"id", "", []option.RequestOption{option.WithEnvironmentId()}, "api.upbit.co.id"},
		{"th", "", []option.RequestOption{option.WithEnvironmentTh()}, "th-api.upbit.com"},
		{"env overrides country", "https://custom.example", []option.RequestOption{option.WithEnvironmentId()}, "custom.example"},
		{"explicit overrides env", "https://custom.example", []option.RequestOption{option.WithEnvironmentId(), option.WithBaseURL("https://api.upbit.co.id")}, "api.upbit.co.id"},
		{"country cannot override explicit", "", []option.RequestOption{option.WithBaseURL("https://custom.example"), option.WithEnvironmentId()}, "custom.example"},
		{"explicit old address preserved", "", []option.RequestOption{option.WithEnvironmentId(), option.WithBaseURL("https://id-api.upbit.com")}, "id-api.upbit.com"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if tc.env != "" {
				t.Setenv("UPBIT_BASE_URL", tc.env)
			}
			calls := 0
			opts := append([]option.RequestOption{}, tc.opts...)
			opts = append(opts, option.WithHTTPClient(&http.Client{Transport: &closureTransport{
				fn: func(req *http.Request) (*http.Response, error) {
					if req.URL.Host != tc.host || req.URL.Path != "/v1/market/all" {
						t.Errorf("unexpected request URL: %s", req.URL)
					}
					calls++
					return &http.Response{StatusCode: 200, Header: http.Header{"Content-Type": {"application/json"}}, Body: io.NopCloser(strings.NewReader("[]"))}, nil
				},
			}}))
			client := upbit.NewClient(opts...)
			if _, err := client.TradingPairs.List(context.Background(), upbit.TradingPairListParams{}); err != nil {
				t.Fatal(err)
			}
			if calls != 1 {
				t.Fatalf("expected one request, got %d", calls)
			}
		})
	}
}

func TestIDPaginationRequestURLs(t *testing.T) {
	t.Setenv("UPBIT_BASE_URL", "https://api.upbit.co.id")
	var pages []string
	client := upbit.NewClient(option.WithEnvironmentId(), option.WithAccessKey("key"), option.WithSecretKey(strings.Repeat("x", 32)),
		option.WithHTTPClient(&http.Client{Transport: &closureTransport{
			fn: func(req *http.Request) (*http.Response, error) {
				if req.URL.Host != "api.upbit.co.id" || req.URL.Path != "/v1/orders/open" {
					t.Errorf("unexpected pagination URL: %s", req.URL)
				}
				pages = append(pages, req.URL.Query().Get("page"))
				body := "[]"
				if len(pages) == 1 {
					body = `[{"uuid":"order-1"}]`
				}
				return &http.Response{StatusCode: 200, Header: http.Header{"Content-Type": {"application/json"}}, Body: io.NopCloser(strings.NewReader(body))}, nil
			},
		}}))
	iter := client.Orders.ListOpenAutoPaging(context.Background(), upbit.OrderListOpenParams{Page: upbit.Int(1), Limit: upbit.Int(1)})
	var orders []string
	for iter.Next() {
		orders = append(orders, iter.Current().Uuid)
	}
	if err := iter.Err(); err != nil {
		t.Fatal(err)
	}
	if strings.Join(pages, ",") != "1,2" || strings.Join(orders, ",") != "order-1" {
		t.Fatalf("unexpected pagination: pages=%v orders=%v", pages, orders)
	}
}
