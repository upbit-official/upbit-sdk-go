package lib

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/upbit-official/upbit-sdk-go/internal/requestconfig"
)

type failingBodyCopy struct {
	readErr  error
	closeErr error
}

func (b failingBodyCopy) Read([]byte) (int, error) { return 0, b.readErr }
func (b failingBodyCopy) Close() error             { return b.closeErr }

func TestJWTBodyPreparationFailureDoesNotSend(t *testing.T) {
	sentinel := errors.New("body copy failed")
	cases := []struct {
		name    string
		body    string
		getBody func() (io.ReadCloser, error)
		wantErr error
	}{
		{name: "invalid JSON", body: `{"x":`},
		{name: "non-object JSON", body: `[1]`},
		{name: "copy failure", body: `{}`, getBody: func() (io.ReadCloser, error) { return nil, sentinel }, wantErr: sentinel},
		{name: "nil copy", body: `{}`, getBody: func() (io.ReadCloser, error) { return nil, nil }},
		{name: "read failure", body: `{}`, getBody: func() (io.ReadCloser, error) {
			return failingBodyCopy{readErr: sentinel}, nil
		}, wantErr: sentinel},
		{name: "close failure", body: `{}`, getBody: func() (io.ReadCloser, error) {
			return failingBodyCopy{readErr: io.EOF, closeErr: sentinel}, nil
		}, wantErr: sentinel},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodPost, "https://example.invalid/orders", strings.NewReader(tc.body))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Authorization", "Bearer test-access")
			if tc.getBody != nil {
				req.GetBody = tc.getBody
			}
			cfg := &requestconfig.RequestConfig{Request: req}
			if err := WithJWTAuth("test-access", "test-secret").Apply(cfg); err != nil {
				t.Fatal(err)
			}
			sent := 0
			resp, err := cfg.Middlewares[0](req, func(*http.Request) (*http.Response, error) {
				sent++
				return &http.Response{StatusCode: http.StatusOK}, nil
			})
			if err == nil || resp != nil || sent != 0 {
				t.Fatalf("expected local error and no request: response=%v error=%v sent=%d", resp, err, sent)
			}
			if tc.wantErr != nil && !errors.Is(err, tc.wantErr) {
				t.Fatalf("error did not preserve cause: %v", err)
			}
			body, err := io.ReadAll(req.Body)
			if err != nil || string(body) != tc.body {
				t.Fatalf("original body changed: %q, %v", body, err)
			}
			if req.Header.Get("Authorization") != "Bearer test-access" {
				t.Fatal("original authorization changed")
			}
		})
	}
}

func TestJWTBodyPreparationPreservesSentBody(t *testing.T) {
	const body = `{"market":"KRW-BTC","price":"9007199254740993"}`
	req, err := http.NewRequest(http.MethodPost, "https://example.invalid/orders", strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer test-access")
	cfg := &requestconfig.RequestConfig{Request: req}
	if err := WithJWTAuth("test-access", "test-secret").Apply(cfg); err != nil {
		t.Fatal(err)
	}
	sent := 0
	_, err = cfg.Middlewares[0](req, func(out *http.Request) (*http.Response, error) {
		sent++
		if out.Header.Get("Authorization") == req.Header.Get("Authorization") {
			t.Fatal("JWT was not attached to outgoing request")
		}
		actual, err := io.ReadAll(out.Body)
		if err != nil || string(actual) != body {
			t.Fatalf("sent body changed: %q, %v", actual, err)
		}
		return &http.Response{StatusCode: http.StatusOK}, nil
	})
	if err != nil || sent != 1 {
		t.Fatalf("normal request failed: %v, sent=%d", err, sent)
	}
}
