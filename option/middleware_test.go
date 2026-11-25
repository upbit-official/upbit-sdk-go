package option_test

import (
	"bytes"
	"context"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	upbit "github.com/upbit-official/upbit-sdk-go"
	"github.com/upbit-official/upbit-sdk-go/internal/requestconfig"
	"github.com/upbit-official/upbit-sdk-go/lib"
	"github.com/upbit-official/upbit-sdk-go/option"
)

func TestDebugLogRedactsHeadersAndPreservesWireContent(t *testing.T) {
	var logs bytes.Buffer
	req, _ := http.NewRequest(http.MethodPost, "https://example.invalid", strings.NewReader(`{"price":"1.0"}`))
	for _, name := range []string{"aUtHoRiZaTiOn", "aPi-KeY", "X-aPi-kEy", "cOoKiE", "sEt-CoOkIe"} {
		req.Header.Add(name, "request-secret-"+name)
		req.Header.Add(name, "request-secret-second-"+name)
	}
	req.Header.Set("X-Public", "public-request")
	original := req.Header.Clone()
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"X-Public": {"public-response"}},
		Body:       io.NopCloser(strings.NewReader("response-body")),
	}
	for _, name := range []string{"authorization", "api-key", "x-api-key", "cookie", "set-cookie"} {
		resp.Header.Add(name, "response-secret-"+name)
		resp.Header.Add(name, "response-secret-second-"+name)
	}
	originalResponse := resp.Header.Clone()
	cfg := &requestconfig.RequestConfig{}
	if err := option.WithDebugLog(log.New(&logs, "", 0)).Apply(cfg); err != nil {
		t.Fatal(err)
	}
	got, err := cfg.Middlewares[0](req, func(r *http.Request) (*http.Response, error) {
		if !reflect.DeepEqual(r.Header, original) {
			t.Fatal("wire request headers changed")
		}
		body, err := io.ReadAll(r.Body)
		if err != nil || string(body) != `{"price":"1.0"}` {
			t.Fatalf("wire body = %q, %v", body, err)
		}
		return resp, nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(req.Header, original) || !reflect.DeepEqual(got.Header, originalResponse) {
		t.Fatal("original request or response headers changed")
	}
	body, err := io.ReadAll(got.Body)
	if err != nil || string(body) != "response-body" {
		t.Fatalf("response body = %q, %v", body, err)
	}
	if strings.Contains(logs.String(), "request-secret-") || strings.Contains(logs.String(), "response-secret-") {
		t.Fatalf("sensitive headers leaked: %s", logs.String())
	}
	if strings.Count(logs.String(), "***") != 20 {
		t.Fatalf("multi-value redaction count mismatch: %s", logs.String())
	}
	for _, public := range []string{"public-request", "public-response", `{"price":"1.0"}`, "response-body"} {
		if !strings.Contains(logs.String(), public) {
			t.Fatalf("missing public content %q", public)
		}
	}
}

type debugHTTPClient func(*http.Request) (*http.Response, error)

func (f debugHTTPClient) Do(req *http.Request) (*http.Response, error) { return f(req) }

func TestDebugLogJWTOrdersRetriesAndPagination(t *testing.T) {
	for _, jwtFirst := range []bool{false, true} {
		name := "debug-before-jwt"
		if jwtFirst {
			name = "jwt-before-debug"
		}
		t.Run(name, func(t *testing.T) {
			var logs bytes.Buffer
			var tokens []string
			var queries []string
			doer := debugHTTPClient(func(req *http.Request) (*http.Response, error) {
				token := strings.TrimPrefix(req.Header.Get("Authorization"), "Bearer ")
				tokens = append(tokens, token)
				queries = append(queries, req.URL.RawQuery)
				for _, header := range []string{"Api-Key", "X-Api-Key", "Cookie", "Set-Cookie"} {
					if req.Header.Get(header) != "wire-secret" {
						t.Fatalf("wire header %s changed", header)
					}
				}
				body := `[{"uuid":"first-order"}]`
				status := http.StatusOK
				if len(tokens) == 1 {
					status = http.StatusInternalServerError
					body = `{"error":{"name":"temporary","message":"retry"}}`
				} else if len(tokens) == 3 {
					body = "[]"
				}
				return &http.Response{
					StatusCode: status,
					Header:     http.Header{"Content-Type": {"application/json"}, "Retry-After": {"0"}, "Set-Cookie": {"response-cookie-secret"}},
					Body:       io.NopCloser(strings.NewReader(body)),
					Request:    req,
				}, nil
			})
			debug := option.WithDebugLog(log.New(&logs, "", 0))
			jwt := lib.WithJWTAuth("dummy-access-secret", "dummy-signing-secret")
			opts := []option.RequestOption{
				option.WithHTTPClient(doer), option.WithBaseURL("https://example.invalid"),
				option.WithAccessKey("dummy-access-secret"), option.WithSecretKey(""),
				option.WithMaxRetries(1),
			}
			for _, header := range []string{"aPi-kEy", "X-aPi-kEy", "cOoKiE", "sEt-CoOkIe"} {
				opts = append(opts, option.WithHeader(header, "wire-secret"))
			}
			if jwtFirst {
				opts = append(opts, jwt, debug)
			} else {
				opts = append(opts, debug, jwt)
			}
			client := upbit.NewClient(opts...)
			page, err := client.Orders.ListOpen(context.Background(), upbit.OrderListOpenParams{Page: upbit.Int(1)})
			if err != nil {
				t.Fatal(err)
			}
			if _, err := page.GetNextPage(); err != nil {
				t.Fatal(err)
			}
			if len(tokens) != 3 || queries[0] != "page=1" || queries[1] != "page=1" || queries[2] != "page=2" {
				t.Fatalf("retry/page requests = %v", queries)
			}
			for i, token := range tokens {
				if strings.Contains(logs.String(), token) {
					t.Fatal("JWT leaked in debug log")
				}
				parts := strings.Split(token, ".")
				if len(parts) != 3 {
					t.Fatal("wire Authorization is not JWT")
				}
				payload, err := base64.RawURLEncoding.DecodeString(parts[1])
				if err != nil {
					t.Fatal(err)
				}
				var claims map[string]any
				if err := json.Unmarshal(payload, &claims); err != nil {
					t.Fatal(err)
				}
				sum := sha512.Sum512([]byte(queries[i]))
				if claims["query_hash"] != hex.EncodeToString(sum[:]) {
					t.Fatalf("request %d has stale query hash", i)
				}
			}
			if tokens[0] == tokens[1] || tokens[1] == tokens[2] {
				t.Fatal("JWT was not refreshed per request")
			}
			for _, secret := range []string{"dummy-access-secret", "wire-secret", "response-cookie-secret"} {
				if strings.Contains(logs.String(), secret) {
					t.Fatalf("%s leaked", secret)
				}
			}
			if strings.Count(logs.String(), "Request Content:") != 3 || strings.Count(logs.String(), "Response Content:") != 3 {
				t.Fatalf("missing retry/page debug logs: %s", logs.String())
			}
		})
	}
}

func TestDebugLogRawHeaderKeysOverHTTP(t *testing.T) {
	var logs bytes.Buffer
	sensitive := []string{"Authorization", "Api-Key", "X-Api-Key", "Cookie", "Set-Cookie"}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil || string(body) != "request-body" {
			t.Errorf("wire body = %q, %v", body, err)
		}
		for _, name := range sensitive {
			if len(r.Header.Values(name)) != 4 {
				t.Errorf("wire header %s = %v", name, r.Header.Values(name))
			}
			w.Header().Add(name, "Bearer response-fake-secret-"+name)
			w.Header().Add(name, "Basic response-fake-secret-second-"+name)
		}
		w.Header().Set("X-Public", "public-response")
		_, _ = w.Write([]byte("response-body"))
	}))
	defer server.Close()
	req, err := http.NewRequest(http.MethodPost, server.URL, strings.NewReader("request-body"))
	if err != nil {
		t.Fatal(err)
	}
	for _, name := range sensitive {
		// Raw map assignment retains case variants, unlike Header.Add/Set.
		req.Header[name] = []string{"Bearer request-fake-secret-" + name, "Basic request-fake-secret-second-" + name}
		req.Header[strings.ToLower(name)] = []string{"Bearer request-fake-secret-lower-" + name}
		req.Header[strings.ToUpper(name)] = []string{"Bearer request-fake-secret-upper-" + name}
	}
	req.Header.Set("X-Public", "public-request")
	original := req.Header.Clone()
	cfg := &requestconfig.RequestConfig{}
	if err := option.WithDebugLog(log.New(&logs, "", 0)).Apply(cfg); err != nil {
		t.Fatal(err)
	}
	middleware := cfg.Middlewares[0]
	var originalResponse http.Header
	got, err := middleware(req, func(r *http.Request) (*http.Response, error) {
		if !reflect.DeepEqual(r.Header, original) {
			t.Error("original request headers changed")
		}
		resp, err := http.DefaultClient.Do(r)
		if err != nil {
			return nil, err
		}
		// A wrapper transport may supply raw lowercase/mixed-case header keys.
		for _, name := range sensitive {
			values := resp.Header.Values(name)
			resp.Header.Del(name)
			resp.Header[strings.ToLower(name)] = values
			resp.Header[strings.ToUpper(name)] = []string{"Bearer response-fake-secret-upper-" + name}
		}
		originalResponse = resp.Header.Clone()
		return resp, nil
	})
	if err != nil {
		t.Fatal(err)
	}
	defer got.Body.Close()
	if !reflect.DeepEqual(req.Header, original) || !reflect.DeepEqual(got.Header, originalResponse) {
		t.Error("original headers changed")
	}
	body, err := io.ReadAll(got.Body)
	if err != nil || string(body) != "response-body" {
		t.Fatalf("response body = %q, %v", body, err)
	}
	if strings.Contains(logs.String(), "fake-secret") {
		t.Fatalf("sensitive headers leaked: %s", logs.String())
	}
	for _, public := range []string{"public-request", "public-response", "request-body", "response-body"} {
		if !strings.Contains(logs.String(), public) {
			t.Errorf("missing public content %q", public)
		}
	}

}
