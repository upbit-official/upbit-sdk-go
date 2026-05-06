package upbit

import (
	"context"
	"net/http"
	"slices"

	"github.com/upbit-official/upbit-sdk-go/internal/apijson"
	"github.com/upbit-official/upbit-sdk-go/internal/requestconfig"
	"github.com/upbit-official/upbit-sdk-go/option"
	"github.com/upbit-official/upbit-sdk-go/packages/respjson"
)

// APIKeyService contains methods and other services that help with interacting
// with the upbit API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewAPIKeyService] method instead.
type APIKeyService struct {
	options []option.RequestOption
}

// NewAPIKeyService generates a new service that applies the given options to each
// request. These options are applied after the parent client's options (if there
// is one), and before any request-specific options.
func NewAPIKeyService(opts ...option.RequestOption) (r APIKeyService) {
	r = APIKeyService{}
	r.options = opts
	return
}

// List API Keys
//
// GET /v1/api_keys
//
// - docs(kr): https://docs.upbit.com/kr/reference/list-api-keys
// - docs(global): https://global-docs.upbit.com/reference/list-api-keys
func (r *APIKeyService) List(ctx context.Context, opts ...option.RequestOption) (res *[]APIKeyListResponse, err error) {
	opts = slices.Concat(r.options, opts)
	path := "v1/api_keys"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

type APIKeyListResponse struct {
	// Access key of the API key.
	AccessKey string `json:"access_key" api:"required"`
	// Expiration timestamp of the access key.
	//
	// [Format] yyyy-MM-dd'T'HH:mm:ss+09:00
	ExpireAt string `json:"expire_at" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		AccessKey   respjson.Field
		ExpireAt    respjson.Field
		ExtraFields map[string]respjson.Field
		Raw         string
	} `json:"-"`
}

// RawJSON returns the unmodified JSON received from the API.
func (r APIKeyListResponse) RawJSON() string { return r.JSON.Raw }
func (r *APIKeyListResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}
