package apierror_test

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/upbit-official/upbit-sdk-go/internal/apierror"
)

func makeBase(status int, body string) *apierror.APIError {
	req := &http.Request{
		Method: "GET",
		URL:    &url.URL{Scheme: "https", Host: "api.upbit.com", Path: "/v1/order"},
	}
	resp := &http.Response{StatusCode: status, Request: req}
	aerr := &apierror.APIError{Request: req, Response: resp, StatusCode: status}
	aerr.UnmarshalJSON([]byte(body)) //nolint:errcheck
	return aerr
}

func TestUnmarshalShapeB(t *testing.T) {
	body := `{"error":{"name":"invalid_parameter","message":"잘못된 파라미터"}}`
	aerr := makeBase(400, body)

	if aerr.Name != "invalid_parameter" {
		t.Errorf("Name = %q, want %q", aerr.Name, "invalid_parameter")
	}
	if aerr.Message != "잘못된 파라미터" {
		t.Errorf("Message = %q, want %q", aerr.Message, "잘못된 파라미터")
	}
	if !aerr.JSON.Name.Valid() {
		t.Error("JSON.Name should be valid")
	}
	if !aerr.JSON.Message.Valid() {
		t.Error("JSON.Message should be valid")
	}
}

func TestUnmarshalShapeA(t *testing.T) {
	body := `{"name":"order_not_found","message":"주문을 찾을 수 없습니다."}`
	aerr := makeBase(404, body)

	if aerr.Name != "order_not_found" {
		t.Errorf("Name = %q, want %q", aerr.Name, "order_not_found")
	}
	if aerr.Message != "주문을 찾을 수 없습니다." {
		t.Errorf("Message = %q, want %q", aerr.Message, "주문을 찾을 수 없습니다.")
	}
}

func TestRawJSONPreserved(t *testing.T) {
	body := `{"error":{"name":"test","message":"msg"}}`
	aerr := makeBase(400, body)

	if aerr.RawJSON() != body {
		t.Errorf("RawJSON = %q, want %q", aerr.RawJSON(), body)
	}
}

func TestNonJSONBody(t *testing.T) {
	aerr := makeBase(500, "Internal Server Error")

	if aerr.Name != "" {
		t.Errorf("Name should be empty for non-JSON body, got %q", aerr.Name)
	}
	if aerr.RawJSON() != "Internal Server Error" {
		t.Error("RawJSON should preserve raw body")
	}
}

func TestExtraFieldsShapeBNoEnvelopeKey(t *testing.T) {
	body := `{"error":{"name":"invalid_parameter","message":"bad"}}`
	aerr := makeBase(400, body)

	if _, ok := aerr.JSON.ExtraFields["error"]; ok {
		t.Error(`ExtraFields must not contain "error" envelope key for shape B response`)
	}
}

func TestExtraFieldsShapeAPreservedUnknownKeys(t *testing.T) {
	body := `{"name":"x","message":"y","extra_field":"val"}`
	aerr := makeBase(400, body)

	if _, ok := aerr.JSON.ExtraFields["extra_field"]; !ok {
		t.Error(`ExtraFields should contain unknown keys from shape A`)
	}
}

func TestErrorStringNoTrailingSpace(t *testing.T) {
	req := &http.Request{Method: "GET", URL: &url.URL{Scheme: "https", Host: "api.upbit.com", Path: "/v1/order"}}
	resp := &http.Response{StatusCode: 400, Request: req}
	aerr := &apierror.APIError{Request: req, Response: resp, StatusCode: 400, Name: "test_error"}

	s := aerr.Error()
	if len(s) > 0 && s[len(s)-1] == ' ' {
		t.Errorf("Error() has trailing space: %q", s)
	}
}

func TestNonStringNameInvalid(t *testing.T) {
	body := `{"error":{"name":123,"message":"msg"}}`
	aerr := makeBase(400, body)

	if aerr.Name != "" {
		t.Errorf("Name should be empty for non-string, got %q", aerr.Name)
	}
	if aerr.JSON.Name.Valid() {
		t.Error("JSON.Name should be invalid for non-string value")
	}
}

func TestErrorStringWithName(t *testing.T) {
	body := `{"error":{"name":"invalid_parameter","message":"잘못됨"}}`
	aerr := makeBase(400, body)
	s := aerr.Error()

	if s == "" {
		t.Error("Error() should not be empty")
	}
	want := "[invalid_parameter]"
	if len(s) < len(want) {
		t.Errorf("Error() too short, missing %q", want)
	}
}
