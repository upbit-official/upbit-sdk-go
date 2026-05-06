package requestconfig_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	upbit "github.com/upbit-official/upbit-sdk-go"
	"github.com/upbit-official/upbit-sdk-go/option"
)

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
