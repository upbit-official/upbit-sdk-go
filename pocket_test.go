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

func TestPocketList(t *testing.T) {
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
	_, err := client.Pockets.List(context.TODO())
	if err != nil {
		var apierr *upbit.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestPocketListAPIKeysWithOptionalParams(t *testing.T) {
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
	_, err := client.Pockets.ListAPIKeys(context.TODO(), upbit.PocketListAPIKeysParams{
		IncludeExpired: upbit.Bool(true),
		Uuids:          []string{"string"},
	})
	if err != nil {
		var apierr *upbit.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestPocketListTransfersWithOptionalParams(t *testing.T) {
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
	_, err := client.Pockets.ListTransfers(context.TODO(), upbit.PocketListTransfersParams{
		Currency:    upbit.String("XRP"),
		Direction:   upbit.PocketTransferDirectionOut,
		EndTime:     upbit.String("end_time"),
		Identifiers: []string{"string"},
		Limit:       upbit.Int(2),
		OrderBy:     upbit.PocketListTransfersParamsOrderByDesc,
		StartTime:   upbit.String("start_time"),
		States:      []upbit.PocketTransferState{upbit.PocketTransferStateProcessing},
		Uuids:       []string{"string"},
	})
	if err != nil {
		var apierr *upbit.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestPocketListUniversalTransfersWithOptionalParams(t *testing.T) {
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
	_, err := client.Pockets.ListUniversalTransfers(context.TODO(), upbit.PocketListUniversalTransfersParams{
		Currency:    upbit.String("XRP"),
		EndTime:     upbit.String("end_time"),
		From:        upbit.String("from"),
		Identifiers: []string{"string"},
		Limit:       upbit.Int(2),
		OrderBy:     upbit.PocketListUniversalTransfersParamsOrderByDesc,
		StartTime:   upbit.String("start_time"),
		States:      []upbit.PocketTransferState{upbit.PocketTransferStateProcessing},
		To:          upbit.String("to"),
		Uuids:       []string{"string"},
	})
	if err != nil {
		var apierr *upbit.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestPocketGetBalance(t *testing.T) {
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
	_, err := client.Pockets.GetBalance(context.TODO(), upbit.PocketGetBalanceParams{
		Uuid: "uuid",
	})
	if err != nil {
		var apierr *upbit.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestPocketTransferWithOptionalParams(t *testing.T) {
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
	_, err := client.Pockets.Transfer(context.TODO(), upbit.PocketTransferParams{
		Amount:     "10.00",
		Currency:   "XRP",
		To:         "9ca023a5-851b-4fec-9f0a-48cd83c2eaae",
		Identifier: upbit.String("unique_transfer_id_20260531_01"),
	})
	if err != nil {
		var apierr *upbit.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestPocketUniversalTransferWithOptionalParams(t *testing.T) {
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
	_, err := client.Pockets.UniversalTransfer(context.TODO(), upbit.PocketUniversalTransferParams{
		Amount:     "10.00",
		Currency:   "XRP",
		To:         "9ca023a5-851b-4fec-9f0a-48cd83c2eaae",
		From:       upbit.String("9ca023a5-851b-4fec-9f0a-48cd83c2eaae"),
		Identifier: upbit.String("unique_transfer_id_20260531_01"),
	})
	if err != nil {
		var apierr *upbit.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}
