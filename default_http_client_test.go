package upbit

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/upbit-official/upbit-sdk-go/option"
)

type transportFunc func(*http.Request) (*http.Response, error)

func (f transportFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

type transportWrapper struct{ http.RoundTripper }

func TestDefaultHTTPClientClonesStandardTransport(t *testing.T) {
	original := http.DefaultTransport
	t.Cleanup(func() { http.DefaultTransport = original })
	transport := &http.Transport{ResponseHeaderTimeout: time.Second}
	http.DefaultTransport = transport

	client := defaultHTTPClient()
	cloned, ok := client.Transport.(*http.Transport)
	if !ok || cloned == transport {
		t.Fatal("default client must clone the standard transport")
	}
	if cloned.ResponseHeaderTimeout != defaultResponseHeaderTimeout {
		t.Fatalf("header timeout = %v", cloned.ResponseHeaderTimeout)
	}
	if transport.ResponseHeaderTimeout != time.Second {
		t.Fatal("original transport was modified")
	}
}

func TestDefaultHTTPClientPreservesWrappedTransport(t *testing.T) {
	original := http.DefaultTransport
	t.Cleanup(func() { http.DefaultTransport = original })
	calls := 0
	wrapped := &transportWrapper{RoundTripper: transportFunc(func(req *http.Request) (*http.Response, error) {
		calls++
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": {"application/json"}},
			Body:       io.NopCloser(strings.NewReader("[]")),
			Request:    req,
		}, nil
	})}
	http.DefaultTransport = wrapped
	client := defaultHTTPClient()
	if client.Transport != wrapped {
		t.Fatal("default client must preserve the wrapper")
	}
	req, _ := http.NewRequest(http.MethodGet, "https://example.invalid", nil)
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if calls != 1 {
		t.Fatalf("wrapper calls = %d", calls)
	}
}

func TestNewClientWithWrappedDefaultAndExplicitHTTPClient(t *testing.T) {
	original := http.DefaultTransport
	t.Cleanup(func() { http.DefaultTransport = original })
	http.DefaultTransport = transportFunc(func(*http.Request) (*http.Response, error) {
		t.Fatal("explicit client must take precedence over the default transport")
		return nil, nil
	})
	calls := 0
	explicit := &http.Client{Transport: transportFunc(func(req *http.Request) (*http.Response, error) {
		calls++
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": {"application/json"}},
			Body:       io.NopCloser(strings.NewReader("[]")),
			Request:    req,
		}, nil
	})}
	client := NewClient(option.WithHTTPClient(explicit), option.WithAccessKey("dummy-access"), option.WithMaxRetries(0))
	if _, err := client.Accounts.List(context.Background()); err != nil {
		t.Fatal(err)
	}
	if calls != 1 {
		t.Fatalf("explicit client calls = %d", calls)
	}
}
