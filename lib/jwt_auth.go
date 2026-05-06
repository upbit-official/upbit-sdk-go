package lib

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
)

// CreateJWTToken creates a JWT token for Upbit API authentication.
//
// queryString should be the URL-decoded query string to hash (e.g. "market=KRW-BTC&side=bid").
// Pass an empty string when there are no request parameters.
func CreateJWTToken(accessKey, secretKey, queryString string) (string, error) {
	nonce, err := newUUIDv4()
	if err != nil {
		return "", err
	}

	claims := map[string]any{
		"access_key": accessKey,
		"nonce":      nonce,
	}

	if queryString != "" {
		sum := sha512.Sum512([]byte(queryString))
		claims["query_hash"] = hex.EncodeToString(sum[:])
		claims["query_hash_alg"] = "SHA512"
	}

	return signHS256(claims, []byte(secretKey))
}

func newUUIDv4() (string, error) {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "", err
	}
	b[6] = (b[6] & 0x0f) | 0x40 // version 4
	b[8] = (b[8] & 0x3f) | 0x80 // variant bits (RFC 4122)
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:16]), nil
}

func signHS256(claims map[string]any, secret []byte) (string, error) {
	enc := base64.RawURLEncoding

	headerJSON, err := json.Marshal(map[string]string{"alg": "HS256", "typ": "JWT"})
	if err != nil {
		return "", err
	}
	payloadJSON, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}

	signingInput := enc.EncodeToString(headerJSON) + "." + enc.EncodeToString(payloadJSON)

	mac := hmac.New(sha256.New, secret)
	mac.Write([]byte(signingInput))

	return signingInput + "." + enc.EncodeToString(mac.Sum(nil)), nil
}
