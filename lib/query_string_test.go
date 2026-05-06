package lib_test

import (
	"testing"

	"github.com/upbit-official/upbit-sdk-go/lib"
)

func TestQueryStringFromRawQuery_Simple(t *testing.T) {
	result := lib.QueryStringFromRawQuery("market=KRW-BTC")
	if result != "market=KRW-BTC" {
		t.Errorf("got %q, want %q", result, "market=KRW-BTC")
	}
}

func TestQueryStringFromRawQuery_WithBrackets(t *testing.T) {
	// Go's url.Values.Encode() encodes [] as %5B%5D; this should be decoded back.
	input := "uuids%5B%5D=uuid-1&uuids%5B%5D=uuid-2"
	expected := "uuids[]=uuid-1&uuids[]=uuid-2"
	result := lib.QueryStringFromRawQuery(input)
	if result != expected {
		t.Errorf("got %q, want %q", result, expected)
	}
}

func TestQueryStringFromRawQuery_Empty(t *testing.T) {
	result := lib.QueryStringFromRawQuery("")
	if result != "" {
		t.Errorf("got %q, want empty string", result)
	}
}

func TestQueryStringFromJSONBody_SimpleDict(t *testing.T) {
	body := []byte(`{"market":"KRW-BTC","side":"bid"}`)
	result, err := lib.QueryStringFromJSONBody(body)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Keys sorted alphabetically: market, side
	expected := "market=KRW-BTC&side=bid"
	if result != expected {
		t.Errorf("got %q, want %q", result, expected)
	}
}

func TestQueryStringFromJSONBody_WithArrayBrackets(t *testing.T) {
	body := []byte(`{"uuids":["uuid-1","uuid-2"]}`)
	result, err := lib.QueryStringFromJSONBody(body)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := "uuids[]=uuid-1&uuids[]=uuid-2"
	if result != expected {
		t.Errorf("got %q, want %q", result, expected)
	}
}

func TestQueryStringFromJSONBody_Empty(t *testing.T) {
	result, err := lib.QueryStringFromJSONBody([]byte(`{}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "" {
		t.Errorf("got %q, want empty string", result)
	}
}

func TestQueryStringFromJSONBody_Nil(t *testing.T) {
	result, err := lib.QueryStringFromJSONBody(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "" {
		t.Errorf("got %q, want empty string", result)
	}
}

func TestQueryStringFromJSONBody_NumericValue(t *testing.T) {
	// Numbers should not use scientific notation.
	body := []byte(`{"price":100}`)
	result, err := lib.QueryStringFromJSONBody(body)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := "price=100"
	if result != expected {
		t.Errorf("got %q, want %q", result, expected)
	}
}
