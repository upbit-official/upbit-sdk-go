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

func TestOrderNewWithOptionalParams(t *testing.T) {
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
	_, err := client.Orders.New(context.TODO(), upbit.OrderNewParams{
		Market:      "KRW-BTC",
		OrdType:     upbit.OrdTypeEnumLimit,
		Side:        upbit.SideEnumBid,
		Identifier:  upbit.String("9ca023a5-851b-4fec-9f0a-48cd83c2eaae"),
		Price:       upbit.String("14000000"),
		SmpType:     upbit.SmpTypeEnumCancelMaker,
		TimeInForce: upbit.TimeInForceEnumFok,
		Volume:      upbit.String("1"),
	})
	if err != nil {
		var apierr *upbit.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestOrderGetWithOptionalParams(t *testing.T) {
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
	_, err := client.Orders.Get(context.TODO(), upbit.OrderGetParams{
		Identifier: upbit.String("9ca023a5-851b-4fec-9f0a-48cd83c2eaae"),
		Uuid:       upbit.String("9ca023a5-851b-4fec-9f0a-48cd83c2eaae"),
	})
	if err != nil {
		var apierr *upbit.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestOrderCancelWithOptionalParams(t *testing.T) {
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
	_, err := client.Orders.Cancel(context.TODO(), upbit.OrderCancelParams{
		Identifier: upbit.String("9ca023a5-851b-4fec-9f0a-48cd83c2eaae"),
		Uuid:       upbit.String("9ca023a5-851b-4fec-9f0a-48cd83c2eaae"),
	})
	if err != nil {
		var apierr *upbit.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestOrderCancelAndNewWithOptionalParams(t *testing.T) {
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
	_, err := client.Orders.CancelAndNew(context.TODO(), upbit.OrderCancelAndNewParams{
		NewOrdType:          upbit.OrdTypeEnumLimit,
		NewIdentifier:       upbit.String("9ca023a5-851b-4fec-9f0a-48cd83c2eaae"),
		NewPrice:            upbit.String("100000000"),
		NewSmpType:          upbit.SmpTypeEnumCancelMaker,
		NewTimeInForce:      upbit.TimeInForceEnumFok,
		NewVolume:           upbit.String("1"),
		PrevOrderIdentifier: upbit.String("9ca023a5-851b-4fec-9f0a-48cd83c2eaae"),
		PrevOrderUuid:       upbit.String("ad217e24-ed02-469c-9b30-c08dbbda6908"),
	})
	if err != nil {
		var apierr *upbit.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestOrderCancelByUuidsWithOptionalParams(t *testing.T) {
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
	_, err := client.Orders.CancelByUuids(context.TODO(), upbit.OrderCancelByUuidsParams{
		Identifiers: []string{"9ca023a5-851b-4fec-9f0a-48cd83c2eaae"},
		Uuids:       []string{"9ca023a5-851b-4fec-9f0a-48cd83c2eaae"},
	})
	if err != nil {
		var apierr *upbit.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestOrderCancelOpenWithOptionalParams(t *testing.T) {
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
	_, err := client.Orders.CancelOpen(context.TODO(), upbit.OrderCancelOpenParams{
		CancelSide:      upbit.OrderCancelOpenParamsCancelSideBid,
		Count:           upbit.Int(300),
		ExcludedPairs:   upbit.String("KRW-BTC,KRW-ETH"),
		OrderBy:         upbit.OrderCancelOpenParamsOrderByDesc,
		Pairs:           upbit.String("KRW-BTC,KRW-ETH"),
		QuoteCurrencies: upbit.String("KRW"),
	})
	if err != nil {
		var apierr *upbit.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestOrderListByUuidsWithOptionalParams(t *testing.T) {
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
	_, err := client.Orders.ListByUuids(context.TODO(), upbit.OrderListByUuidsParams{
		Identifiers: []string{"9ca023a5-851b-4fec-9f0a-48cd83c2eaae"},
		Market:      upbit.String("KRW-BTC"),
		OrderBy:     upbit.OrderListByUuidsParamsOrderByDesc,
		Uuids:       []string{"9ca023a5-851b-4fec-9f0a-48cd83c2eaae"},
	})
	if err != nil {
		var apierr *upbit.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestOrderListClosedWithOptionalParams(t *testing.T) {
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
	_, err := client.Orders.ListClosed(context.TODO(), upbit.OrderListClosedParams{
		EndTime:   upbit.String("2024-12-09T13:56:53+09:00"),
		Limit:     upbit.Int(1000),
		Market:    upbit.String("KRW-BTC"),
		OrderBy:   upbit.OrderListClosedParamsOrderByDesc,
		StartTime: upbit.String("2024-12-09T13:56:53+09:00"),
		State:     upbit.ClosedOrderStateDone,
		States:    []upbit.ClosedOrderState{upbit.ClosedOrderStateDone},
	})
	if err != nil {
		var apierr *upbit.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestOrderListOpenWithOptionalParams(t *testing.T) {
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
	_, err := client.Orders.ListOpen(context.TODO(), upbit.OrderListOpenParams{
		Limit:   upbit.Int(100),
		Market:  upbit.String("KRW-BTC"),
		OrderBy: upbit.OrderListOpenParamsOrderByDesc,
		Page:    upbit.Int(1),
		State:   upbit.OpenOrderStateWait,
		States:  []upbit.OpenOrderState{upbit.OpenOrderStateWait},
	})
	if err != nil {
		var apierr *upbit.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestOrderGetChance(t *testing.T) {
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
	_, err := client.Orders.GetChance(context.TODO(), upbit.OrderGetChanceParams{
		Market: "KRW-BTC",
	})
	if err != nil {
		var apierr *upbit.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestOrderTestNewWithOptionalParams(t *testing.T) {
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
	_, err := client.Orders.TestNew(context.TODO(), upbit.OrderTestNewParams{
		Market:      "KRW-BTC",
		OrdType:     upbit.OrdTypeEnumLimit,
		Side:        upbit.SideEnumBid,
		Identifier:  upbit.String("9ca023a5-851b-4fec-9f0a-48cd83c2eaae"),
		Price:       upbit.String("14000000"),
		SmpType:     upbit.SmpTypeEnumCancelMaker,
		TimeInForce: upbit.TimeInForceEnumFok,
		Volume:      upbit.String("1"),
	})
	if err != nil {
		var apierr *upbit.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}
