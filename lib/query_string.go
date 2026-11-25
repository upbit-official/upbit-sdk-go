package lib

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/tidwall/gjson"
)

// QueryStringFromRawQuery takes a URL-encoded query string (as found in url.RawQuery)
// and returns the URL-decoded version suitable for JWT query hashing.
//
// Go's url.Values.Encode() percent-encodes brackets as %5B%5D; this function
// restores them to [] to match Python's unquote(qs.stringify(params)) behavior.
func QueryStringFromRawQuery(rawQuery string) string {
	if rawQuery == "" {
		return ""
	}
	decoded, err := url.QueryUnescape(rawQuery)
	if err != nil {
		return rawQuery
	}
	return decoded
}

// QueryStringFromJSONBody converts a JSON-encoded request body into a URL-decoded
// query string using brackets array format, matching the Python SDK's behavior of
// unquote(Querystring(array_format="brackets").stringify(params)).
//
// JSON object key insertion order is preserved (using gjson) so that the hash
// matches the server's expectation. Arrays are serialized as key[]=val1&key[]=val2.
// Returns empty string for an empty or null body.
func QueryStringFromJSONBody(body []byte) (string, error) {
	if len(body) == 0 {
		return "", nil
	}
	if !gjson.ValidBytes(body) {
		return "", fmt.Errorf("query_string: invalid JSON body")
	}

	parsed := gjson.ParseBytes(body)
	if !parsed.IsObject() {
		return "", fmt.Errorf("query_string: JSON body is not an object")
	}

	var parts []string
	var iterErr error

	parsed.ForEach(func(key, value gjson.Result) bool {
		k := key.String()
		if value.IsArray() {
			bracketKey := url.QueryEscape(k + "[]")
			value.ForEach(func(_, elem gjson.Result) bool {
				parts = append(parts, bracketKey+"="+url.QueryEscape(gjsonToString(elem)))
				return true
			})
		} else {
			parts = append(parts, url.QueryEscape(k)+"="+url.QueryEscape(gjsonToString(value)))
		}
		return iterErr == nil
	})

	if iterErr != nil {
		return "", iterErr
	}
	if len(parts) == 0 {
		return "", nil
	}

	encoded := strings.Join(parts, "&")
	decoded, err := url.QueryUnescape(encoded)
	if err != nil {
		return encoded, nil
	}
	return decoded, nil
}

// gjsonToString converts a gjson.Result to its string representation.
// Floats that are whole numbers are rendered without a decimal point.
func gjsonToString(v gjson.Result) string {
	switch v.Type {
	case gjson.String:
		return v.String()
	case gjson.Number:
		f := v.Float()
		return strconv.FormatFloat(f, 'f', -1, 64)
	case gjson.True:
		return "true"
	case gjson.False:
		return "false"
	default:
		return v.String()
	}
}
