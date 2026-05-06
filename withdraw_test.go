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

func TestWithdrawGetWithOptionalParams(t *testing.T) {
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
	_, err := client.Withdraws.Get(context.TODO(), upbit.WithdrawGetParams{
		Currency: upbit.String("KRW"),
		Txid:     upbit.String("98c15999f0bdc4ae0e8a-ed35868bb0c204fe6ec29e4058a3451e-88636d1040f4baddf943274ce37cf9cc"),
		Uuid:     upbit.String("9ca023a5-851b-4fec-9f0a-48cd83c2eaae"),
	})
	if err != nil {
		var apierr *upbit.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestWithdrawListWithOptionalParams(t *testing.T) {
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
	_, err := client.Withdraws.List(context.TODO(), upbit.WithdrawListParams{
		Currency: upbit.String("KRW"),
		From:     upbit.String("from"),
		Limit:    upbit.Int(100),
		OrderBy:  upbit.WithdrawListParamsOrderByDesc,
		Page:     upbit.Int(1),
		State:    upbit.WithdrawListParamsStateWaiting,
		To:       upbit.String("to"),
		Txids:    []string{"98c15999f0bdc4ae0e8a-ed35868bb0c204fe6ec29e4058a3451e-88636d1040f4baddf943274ce37cf9cc"},
		Uuids:    []string{"9ca023a5-851b-4fec-9f0a-48cd83c2eaae"},
	})
	if err != nil {
		var apierr *upbit.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestWithdrawCancelWithdrawal(t *testing.T) {
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
	_, err := client.Withdraws.CancelWithdrawal(context.TODO(), upbit.WithdrawCancelWithdrawalParams{
		Uuid: "9ca023a5-851b-4fec-9f0a-48cd83c2eaae",
	})
	if err != nil {
		var apierr *upbit.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestWithdrawNewKrwWithdrawal(t *testing.T) {
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
	_, err := client.Withdraws.NewKrwWithdrawal(context.TODO(), upbit.WithdrawNewKrwWithdrawalParams{
		Amount:        "10000",
		TwoFactorType: upbit.WithdrawNewKrwWithdrawalParamsTwoFactorTypeNaver,
	})
	if err != nil {
		var apierr *upbit.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestWithdrawNewWithdrawalWithOptionalParams(t *testing.T) {
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
	_, err := client.Withdraws.NewWithdrawal(context.TODO(), upbit.WithdrawNewWithdrawalParams{
		Address:          "1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa",
		Amount:           "0.01",
		Currency:         "BTC",
		NetType:          "BTC",
		SecondaryAddress: upbit.String("secondary_address"),
		TransactionType:  upbit.WithdrawNewWithdrawalParamsTransactionTypeDefault,
	})
	if err != nil {
		var apierr *upbit.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestWithdrawListCoinAddresses(t *testing.T) {
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
	_, err := client.Withdraws.ListCoinAddresses(context.TODO())
	if err != nil {
		var apierr *upbit.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestWithdrawGetChance(t *testing.T) {
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
	_, err := client.Withdraws.GetChance(context.TODO(), upbit.WithdrawGetChanceParams{
		Currency: "BTC",
		NetType:  "BTC",
	})
	if err != nil {
		var apierr *upbit.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}
