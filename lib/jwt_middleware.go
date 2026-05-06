package lib

import (
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

		queryString := extractQueryString(req)
		token, err := CreateJWTToken(accessKey, secretKey, queryString)
		if err != nil {
			return next(req)
		}

		reqClone := req.Clone(req.Context())
		reqClone.Header.Set("authorization", "Bearer "+token)
		return next(reqClone)
	})
}

// extractQueryString returns the URL-decoded query string suitable for JWT
// query hashing. Prefers URL query params; falls back to JSON body.
func extractQueryString(req *http.Request) string {
	if rawQuery := req.URL.RawQuery; rawQuery != "" {
		return QueryStringFromRawQuery(rawQuery)
	}
	if req.GetBody != nil {
		bodyReader, err := req.GetBody()
		if err == nil {
			defer bodyReader.Close()
			body, err := io.ReadAll(bodyReader)
			if err == nil && len(body) > 0 {
				qs, err := QueryStringFromJSONBody(body)
				if err == nil {
					return qs
				}
			}
		}
	}
	return ""
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
