package requestconfig_test

import (
	"bytes"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	upbit "github.com/upbit-official/upbit-sdk-go"
	"github.com/upbit-official/upbit-sdk-go/lib"
	"github.com/upbit-official/upbit-sdk-go/option"
)

func TestReaderBodySigningAndRetry(t *testing.T) {
	const body = `{"market":"KRW-BTC","price":"1.0"}`
	for _, kind := range []string{"buffer", "bytes", "strings"} {
		for _, partial := range []bool{false, true} {
			for _, override := range []bool{false, true} {
				name := kind
				if partial {
					name += "-partial"
				}
				if override {
					name += "-override"
				}
				t.Run(name, func(t *testing.T) {
					content := body
					if partial {
						content = "consumed-prefix" + content
					}
					var reader io.Reader
					switch kind {
					case "buffer":
						reader = bytes.NewBufferString(content)
					case "bytes":
						reader = bytes.NewReader([]byte(content))
					case "strings":
						reader = strings.NewReader(content)
					}
					if partial {
						_, _ = io.CopyN(io.Discard, reader, int64(len("consumed-prefix")))
					}
					query, err := lib.QueryStringFromJSONBody([]byte(body))
					if err != nil {
						t.Fatal(err)
					}
					sum := sha512.Sum512([]byte(query))
					attempts := 0
					server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						attempts++
						received, err := io.ReadAll(r.Body)
						if err != nil || string(received) != body || r.ContentLength != int64(len(body)) {
							t.Errorf("wire body = %q, length = %d, error = %v", received, r.ContentLength, err)
						}
						parts := strings.Split(strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer "), ".")
						if len(parts) != 3 {
							t.Error("missing JWT")
						} else {
							payload, err := base64.RawURLEncoding.DecodeString(parts[1])
							var claims map[string]any
							if err != nil || json.Unmarshal(payload, &claims) != nil || claims["query_hash"] != hex.EncodeToString(sum[:]) {
								t.Errorf("JWT body hash mismatch: %s, %v", payload, err)
							}
						}
						w.Header().Set("Content-Type", "application/json")
						if attempts == 1 {
							w.Header().Set("Retry-After", "0")
							w.WriteHeader(http.StatusInternalServerError)
						}
						_, _ = w.Write([]byte(`{}`))
					}))
					defer server.Close()
					options := append(lib.WithCredentials("test-access", "test-secret"),
						option.WithBaseURL(server.URL), option.WithMaxRetries(1),
						option.WithMiddleware(func(r *http.Request, next option.MiddlewareNext) (*http.Response, error) {
							if r.GetBody == nil {
								t.Error("missing independent body copy")
								return next(r)
							}
							first, err := r.GetBody()
							if err != nil {
								t.Fatal(err)
							}
							defer first.Close()
							_, _ = io.CopyN(io.Discard, first, 5)
							second, err := r.GetBody()
							if err != nil {
								t.Fatal(err)
							}
							defer second.Close()
							copied, err := io.ReadAll(second)
							if err != nil || string(copied) != body {
								t.Errorf("second copy = %q, %v", copied, err)
							}
							tail, err := io.ReadAll(first)
							if err != nil || string(tail) != body[5:] {
								t.Errorf("first copy cursor changed: %q, %v", tail, err)
							}
							return next(r)
						}))
					client := upbit.NewClient(options...)
					var result map[string]any
					params := any(reader)
					var requestOptions []option.RequestOption
					if override {
						params = nil
						requestOptions = []option.RequestOption{option.WithRequestBody("application/json", reader)}
					}
					err = client.Post(t.Context(), "/probe", params, &result, requestOptions...)
					if err != nil || attempts != 2 {
						t.Fatalf("request/retry failed: %v, attempts = %d", err, attempts)
					}
					// Preparing, signing and sending copies must not consume the caller's reader.
					remaining, err := io.ReadAll(reader)
					if err != nil || string(remaining) != body {
						t.Errorf("caller reader consumed: %q, %v", remaining, err)
					}
				})
			}
		}
	}
}

func TestDefaultClientHeaders(t *testing.T) {
	var captured http.Header

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		captured = r.Header.Clone()
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`[]`))
	}))
	defer server.Close()

	client := upbit.NewClient(
		option.WithAccessKey("test-access-key"),
		option.WithSecretKey("test-secret-key"),
		option.WithBaseURL(server.URL),
	)

	_, _ = client.TradingPairs.List(t.Context(), upbit.TradingPairListParams{})

	tests := []struct {
		header string
		want   string
	}{
		{"X-Upbit-Client-Type", "sdk"},
		{"X-Upbit-Client-Name", "upbit-go"},
	}

	for _, tc := range tests {
		got := captured.Get(tc.header)
		if got != tc.want {
			t.Errorf("%s: got %q, want %q", tc.header, got, tc.want)
		}
	}

	if captured.Get("X-Upbit-Client-Version") == "" {
		t.Error("X-Upbit-Client-Version header is missing")
	}
	if captured.Get("X-Upbit-Runtime") == "" {
		t.Error("X-Upbit-Runtime header is missing")
	}
	if captured.Get("X-Upbit-Runtime-Version") == "" {
		t.Error("X-Upbit-Runtime-Version header is missing")
	}
}
