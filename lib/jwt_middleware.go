package lib

import (
	"fmt"
	"io"
	"net/http"

	"github.com/upbit-official/upbit-sdk-go/option"
)

// WithJWTAuth returns a RequestOption that installs a middleware intercepting
// every outgoing request that carries an Authorization header and replaces the
// simple Bearer token with a properly signed JWT.
//
// The middleware fires per-request, so pagination requests (which carry a
// stale query hash) are automatically re-signed with the current query params.
//
// Usage:
//
//	client := upbit.NewClient(
//	    option.WithAccessKey(accessKey),
//	    lib.WithJWTAuth(accessKey, secretKey),
//	)
func WithJWTAuth(accessKey, secretKey string) option.RequestOption {
	return option.WithMiddleware(func(req *http.Request, next option.MiddlewareNext) (*http.Response, error) {
		// Only intercept authenticated requests.
		if req.Header.Get("Authorization") == "" {
			return next(req)
		}

		queryString, err := extractQueryString(req)
		if err != nil {
			return nil, err
		}
		token, err := CreateJWTToken(accessKey, secretKey, queryString)
		if err != nil {
			return nil, fmt.Errorf("could not create JWT: %w", err)
		}

		reqClone := req.Clone(req.Context())
		reqClone.Header.Set("authorization", "Bearer "+token)
		return next(reqClone)
	})
}

// extractQueryString returns the URL-decoded query string suitable for JWT
// query hashing. Prefers URL query params; falls back to JSON body.
func extractQueryString(req *http.Request) (string, error) {
	if rawQuery := req.URL.RawQuery; rawQuery != "" {
		return QueryStringFromRawQuery(rawQuery), nil
	}
	if req.GetBody != nil {
		bodyReader, err := req.GetBody()
		if err != nil {
			return "", fmt.Errorf("could not copy body for JWT: %w", err)
		}
		if bodyReader == nil {
			return "", fmt.Errorf("could not copy body for JWT: nil reader")
		}
		body, readErr := io.ReadAll(bodyReader)
		closeErr := bodyReader.Close()
		if readErr != nil {
			return "", fmt.Errorf("could not read body for JWT: %w", readErr)
		}
		if closeErr != nil {
			return "", fmt.Errorf("could not close body copy for JWT: %w", closeErr)
		}
		qs, err := QueryStringFromJSONBody(body)
		if err != nil {
			return "", fmt.Errorf("could not convert body for JWT: %w", err)
		}
		return qs, nil
	}
	return "", nil
}

// WithCredentials is a convenience helper that sets both the access key and
// installs the JWT auth middleware in a single call.
//
// Usage:
//
//	client := upbit.NewClient(lib.WithCredentials(accessKey, secretKey)...)
func WithCredentials(accessKey, secretKey string) []option.RequestOption {
	return []option.RequestOption{
		option.WithAccessKey(accessKey),
		option.WithSecretKey(secretKey),
		WithJWTAuth(accessKey, secretKey),
	}
}
