package upbit_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/upbit-official/upbit-sdk-go"
	"github.com/upbit-official/upbit-sdk-go/internal"
	"github.com/upbit-official/upbit-sdk-go/option"
)

type closureTransport struct {
	fn func(req *http.Request) (*http.Response, error)
}

func (t *closureTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return t.fn(req)
}

func TestUserAgentHeader(t *testing.T) {
	var userAgent string
	client := upbit.NewClient(
		option.WithAccessKey("My Access Key"),
		option.WithHTTPClient(&http.Client{
			Transport: &closureTransport{
				fn: func(req *http.Request) (*http.Response, error) {
					userAgent = req.Header.Get("User-Agent")
					return &http.Response{
						StatusCode: http.StatusOK,
					}, nil
				},
			},
		}),
	)
	_, _ = client.Accounts.List(context.Background())
	if userAgent != fmt.Sprintf("Upbit/Go %s", internal.PackageVersion) {
		t.Errorf("Expected User-Agent to be correct, but got: %#v", userAgent)
	}
}

func TestDeprecatedCustomHeadersEnvIgnored(t *testing.T) {
	t.Setenv("UPBIT_CUSTOM_HEADERS", "X-Env-Header: should-not-be-sent")

	var customHeader string
	client := upbit.NewClient(
		option.WithHTTPClient(&http.Client{
			Transport: &closureTransport{
				fn: func(req *http.Request) (*http.Response, error) {
					customHeader = req.Header.Get("X-Env-Header")
					return &http.Response{
						StatusCode: http.StatusOK,
					}, nil
				},
			},
		}),
	)

	_, _ = client.TradingPairs.List(context.Background(), upbit.TradingPairListParams{})
	if customHeader != "" {
		t.Errorf("Expected UPBIT_CUSTOM_HEADERS to be ignored, got header value: %#v", customHeader)
	}
}

func TestClientLevelHeaderOptions(t *testing.T) {
	var customValues []string
	var removedHeader string
	client := upbit.NewClient(
		option.WithHeader("X-Custom-Header", "one"),
		option.WithHeaderAdd("X-Custom-Header", "two"),
		option.WithHeader("X-Removed-Header", "remove-me"),
		option.WithHeaderDel("X-Removed-Header"),
		option.WithHTTPClient(&http.Client{
			Transport: &closureTransport{
				fn: func(req *http.Request) (*http.Response, error) {
					customValues = req.Header.Values("X-Custom-Header")
					removedHeader = req.Header.Get("X-Removed-Header")
					return &http.Response{
						StatusCode: http.StatusOK,
					}, nil
				},
			},
		}),
	)

	_, _ = client.TradingPairs.List(context.Background(), upbit.TradingPairListParams{})
	if len(customValues) != 2 || customValues[0] != "one" || customValues[1] != "two" {
		t.Errorf("Expected client-level header options to apply, got: %#v", customValues)
	}
	if removedHeader != "" {
		t.Errorf("Expected client-level header delete to apply, got: %#v", removedHeader)
	}
}

func TestRetryAfter(t *testing.T) {
	attempts := 0
	client := upbit.NewClient(
		option.WithAccessKey("My Access Key"),
		option.WithHTTPClient(&http.Client{
			Transport: &closureTransport{
				fn: func(req *http.Request) (*http.Response, error) {
					attempts += 1
					return &http.Response{
						StatusCode: http.StatusTooManyRequests,
						Header: http.Header{
							http.CanonicalHeaderKey("Retry-After"): []string{"0.1"},
						},
					}, nil
				},
			},
		}),
	)
	_, err := client.Accounts.List(context.Background())
	if err == nil {
		t.Error("Expected there to be a cancel error")
	}

	if attempts != 3 {
		t.Errorf("Expected %d attempts, got %d", 3, attempts)
	}
}

func TestRetryAfterMs(t *testing.T) {
	attempts := 0
	client := upbit.NewClient(
		option.WithAccessKey("My Access Key"),
		option.WithHTTPClient(&http.Client{
			Transport: &closureTransport{
				fn: func(req *http.Request) (*http.Response, error) {
					attempts++
					return &http.Response{
						StatusCode: http.StatusTooManyRequests,
						Header: http.Header{
							http.CanonicalHeaderKey("Retry-After-Ms"): []string{"100"},
						},
					}, nil
				},
			},
		}),
	)
	_, err := client.Accounts.List(context.Background())
	if err == nil {
		t.Error("Expected there to be a cancel error")
	}
	if want := 3; attempts != want {
		t.Errorf("Expected %d attempts, got %d", want, attempts)
	}
}

func TestContextCancel(t *testing.T) {
	client := upbit.NewClient(
		option.WithAccessKey("My Access Key"),
		option.WithHTTPClient(&http.Client{
			Transport: &closureTransport{
				fn: func(req *http.Request) (*http.Response, error) {
					<-req.Context().Done()
					return nil, req.Context().Err()
				},
			},
		}),
	)
	cancelCtx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err := client.Accounts.List(cancelCtx)
	if err == nil {
		t.Error("Expected there to be a cancel error")
	}
}

func TestContextCancelDelay(t *testing.T) {
	client := upbit.NewClient(
		option.WithAccessKey("My Access Key"),
		option.WithHTTPClient(&http.Client{
			Transport: &closureTransport{
				fn: func(req *http.Request) (*http.Response, error) {
					<-req.Context().Done()
					return nil, req.Context().Err()
				},
			},
		}),
	)
	cancelCtx, cancel := context.WithTimeout(context.Background(), 2*time.Millisecond)
	defer cancel()
	_, err := client.Accounts.List(cancelCtx)
	if err == nil {
		t.Error("expected there to be a cancel error")
	}
}

func TestContextDeadline(t *testing.T) {
	testTimeout := time.After(3 * time.Second)
	testDone := make(chan struct{})

	deadline := time.Now().Add(100 * time.Millisecond)
	deadlineCtx, cancel := context.WithDeadline(context.Background(), deadline)
	defer cancel()

	go func() {
		client := upbit.NewClient(
			option.WithAccessKey("My Access Key"),
			option.WithHTTPClient(&http.Client{
				Transport: &closureTransport{
					fn: func(req *http.Request) (*http.Response, error) {
						<-req.Context().Done()
						return nil, req.Context().Err()
					},
				},
			}),
		)
		_, err := client.Accounts.List(deadlineCtx)
		if err == nil {
			t.Error("expected there to be a deadline error")
		}
		close(testDone)
	}()

	select {
	case <-testTimeout:
		t.Fatal("client didn't finish in time")
	case <-testDone:
		if diff := time.Since(deadline); diff < -30*time.Millisecond || 30*time.Millisecond < diff {
			t.Fatalf("client did not return within 30ms of context deadline, got %s", diff)
		}
	}
}
