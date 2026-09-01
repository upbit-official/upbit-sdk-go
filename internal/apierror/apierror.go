package apierror

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"

	"github.com/upbit-official/upbit-sdk-go/packages/respjson"
)

// APIError represents an error that originates from the API, i.e. when a request is
// made and the API returns a response with a HTTP status code. Other errors are
// not wrapped by this SDK.
type APIError struct {
	// Name is the Upbit API error name (e.g. "invalid_parameter", "order_not_found").
	Name string
	// Message is the Upbit API error message.
	Message string
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Name        respjson.Field
		Message     respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
	StatusCode int
	Request    *http.Request
	Response   *http.Response
}

// Error is an alias for APIError for backward compatibility.
type Error = APIError

// Returns the unmodified JSON received from the API
func (r APIError) RawJSON() string { return r.JSON.raw }

// UnmarshalJSON parses the Upbit API error body.
// Handles shape B ({"error":{"name":...,"message":...}}) and
// shape A ({"name":...,"message":...}).
func (r *APIError) UnmarshalJSON(data []byte) error {
	r.JSON.raw = string(data)

	// Populate ExtraFields with all top-level keys not handled as named fields,
	// preserving backward-compat with callers that inspect r.JSON.ExtraFields.
	var topLevel map[string]json.RawMessage
	if json.Unmarshal(data, &topLevel) == nil {
		extra := make(map[string]respjson.Field, len(topLevel))
		for k, v := range topLevel {
			// Shape A known fields are handled below; omit them from ExtraFields.
			if k != "name" && k != "message" {
				extra[k] = respjson.NewField(string(v))
			}
		}
		if len(extra) > 0 {
			r.JSON.ExtraFields = extra
		}
	}

	// Shape B: {"error": {"name": "...", "message": "..."}} — Upbit standard format
	var shapeB struct {
		Error struct {
			Name    json.RawMessage `json:"name"`
			Message json.RawMessage `json:"message"`
		} `json:"error"`
	}
	if json.Unmarshal(data, &shapeB) == nil && len(shapeB.Error.Name) > 0 {
		// "error" is a shape B envelope key, not a domain field — remove from ExtraFields.
		delete(r.JSON.ExtraFields, "error")
		r.applyNameMessage(shapeB.Error.Name, shapeB.Error.Message)
		return nil
	}

	// Shape A: {"name": "...", "message": "..."}
	var shapeA struct {
		Name    json.RawMessage `json:"name"`
		Message json.RawMessage `json:"message"`
	}
	if json.Unmarshal(data, &shapeA) == nil && len(shapeA.Name) > 0 {
		r.applyNameMessage(shapeA.Name, shapeA.Message)
	}

	return nil
}

func (r *APIError) applyNameMessage(nameRaw, msgRaw json.RawMessage) {
	if len(nameRaw) > 0 {
		var s string
		if json.Unmarshal(nameRaw, &s) == nil {
			r.Name = s
			r.JSON.Name = respjson.NewField(string(nameRaw))
		} else {
			r.JSON.Name = respjson.NewInvalidField(string(nameRaw))
		}
	}
	if len(msgRaw) > 0 {
		var s string
		if json.Unmarshal(msgRaw, &s) == nil {
			r.Message = s
			r.JSON.Message = respjson.NewField(string(msgRaw))
		} else {
			r.JSON.Message = respjson.NewInvalidField(string(msgRaw))
		}
	}
}

func (r *APIError) Error() string {
	if r.Name != "" {
		base := fmt.Sprintf("%s %q: %d %s [%s]", r.Request.Method, r.Request.URL, r.Response.StatusCode, http.StatusText(r.Response.StatusCode), r.Name)
		if r.Message != "" {
			return base + " " + r.Message
		}
		return base
	}
	return fmt.Sprintf("%s %q: %d %s %s", r.Request.Method, r.Request.URL, r.Response.StatusCode, http.StatusText(r.Response.StatusCode), r.JSON.raw)
}

func (r *APIError) DumpRequest(body bool) []byte {
	if r.Request.GetBody != nil {
		r.Request.Body, _ = r.Request.GetBody()
	}
	out, _ := httputil.DumpRequestOut(r.Request, body)
	return out
}

func (r *APIError) DumpResponse(body bool) []byte {
	out, _ := httputil.DumpResponse(r.Response, body)
	return out
}
