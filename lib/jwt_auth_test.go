package lib_test

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"strings"
	"testing"

	"github.com/upbit-official/upbit-sdk-go/lib"
)

const (
	testAccessKey = "test-access-key"
	testSecretKey = "test-secret-key"
)

// decodeToken verifies the HS256 signature and returns the JWT claims.
func decodeToken(t *testing.T, tokenStr string) map[string]any {
	t.Helper()
	parts := strings.Split(tokenStr, ".")
	if len(parts) != 3 {
		t.Fatalf("invalid token: expected 3 parts, got %d", len(parts))
	}

	mac := hmac.New(sha256.New, []byte(testSecretKey))
	mac.Write([]byte(parts[0] + "." + parts[1]))
	expectedSig := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
	if expectedSig != parts[2] {
		t.Fatalf("signature verification failed")
	}

	payloadJSON, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		t.Fatalf("failed to decode payload: %v", err)
	}
	var claims map[string]any
	if err := json.Unmarshal(payloadJSON, &claims); err != nil {
		t.Fatalf("failed to unmarshal claims: %v", err)
	}
	return claims
}

func TestCreateJWTToken_BasicStructure(t *testing.T) {
	token, err := lib.CreateJWTToken(testAccessKey, testSecretKey, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	claims := decodeToken(t, token)

	if claims["access_key"] != testAccessKey {
		t.Errorf("access_key: got %v, want %v", claims["access_key"], testAccessKey)
	}
	if claims["nonce"] == nil || claims["nonce"] == "" {
		t.Errorf("nonce should be present and non-empty")
	}
	if _, ok := claims["query_hash"]; ok {
		t.Error("query_hash should not be present when queryString is empty")
	}
	if _, ok := claims["query_hash_alg"]; ok {
		t.Error("query_hash_alg should not be present when queryString is empty")
	}
}

func TestCreateJWTToken_WithQueryString(t *testing.T) {
	token, err := lib.CreateJWTToken(testAccessKey, testSecretKey, "market=KRW-BTC&side=bid")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	claims := decodeToken(t, token)

	if _, ok := claims["query_hash"]; !ok {
		t.Error("query_hash should be present when queryString is non-empty")
	}
	if claims["query_hash_alg"] != "SHA512" {
		t.Errorf("query_hash_alg: got %v, want SHA512", claims["query_hash_alg"])
	}
}

func TestCreateJWTToken_QueryHashIsSHA512(t *testing.T) {
	queryString := "market=KRW-BTC"
	token, err := lib.CreateJWTToken(testAccessKey, testSecretKey, queryString)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	sum := sha512.Sum512([]byte(queryString))
	expected := hex.EncodeToString(sum[:])

	claims := decodeToken(t, token)
	if claims["query_hash"] != expected {
		t.Errorf("query_hash: got %v, want %v", claims["query_hash"], expected)
	}
}

func TestCreateJWTToken_NonceUniqueness(t *testing.T) {
	token1, _ := lib.CreateJWTToken(testAccessKey, testSecretKey, "")
	token2, _ := lib.CreateJWTToken(testAccessKey, testSecretKey, "")

	claims1 := decodeToken(t, token1)
	claims2 := decodeToken(t, token2)

	if claims1["nonce"] == claims2["nonce"] {
		t.Error("nonce should be unique across calls")
	}
}

func TestCreateJWTToken_EmptyQueryString(t *testing.T) {
	token, err := lib.CreateJWTToken(testAccessKey, testSecretKey, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	claims := decodeToken(t, token)
	if _, ok := claims["query_hash"]; ok {
		t.Error("query_hash should not be present for empty queryString")
	}
}

func TestCreateJWTToken_AlgorithmHS256(t *testing.T) {
	token, err := lib.CreateJWTToken(testAccessKey, testSecretKey, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		t.Fatalf("invalid token format")
	}

	headerJSON, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		t.Fatalf("failed to decode header: %v", err)
	}
	var header map[string]string
	if err := json.Unmarshal(headerJSON, &header); err != nil {
		t.Fatalf("failed to unmarshal header: %v", err)
	}
	if header["alg"] != "HS256" {
		t.Errorf("algorithm: got %v, want HS256", header["alg"])
	}
}
