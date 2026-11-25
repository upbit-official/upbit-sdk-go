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

// AccountService contains methods and other services that help with interacting
// with the upbit API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewAccountService] method instead.
type AccountService struct {
	options []option.RequestOption
}

// NewAccountService generates a new service that applies the given options to each
// request. These options are applied after the parent client's options (if there
// is one), and before any request-specific options.
func NewAccountService(opts ...option.RequestOption) (r AccountService) {
	r = AccountService{}
	r.options = opts
	return
}

// Get Account Balances
//
// GET /v1/accounts
//
// - docs(kr): https://docs.upbit.com/kr/reference/get-balance
// - docs(sg/th): https://global-docs.upbit.com/reference/get-balance
// - docs(id): https://docs.upbit.co.id/reference/get-balance
//
// Retrieves the list of assets and balances owned by the account.
func (r *AccountService) List(ctx context.Context, opts ...option.RequestOption) (res *[]AccountListResponse, err error) {
	opts = slices.Concat(r.options, opts)
	path := "v1/accounts"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

type AccountListResponse struct {
	// Average buy price of the asset.
	AvgBuyPrice string `json:"avg_buy_price" api:"required"`
	// Indicates whether the average buy price has been modified.
	AvgBuyPriceModified bool `json:"avg_buy_price_modified" api:"required"`
	// Available amount or volume for orders. For digital assets, this represents the
	// available quantity. For fiat currency, this represents the available amount.
	Balance string `json:"balance" api:"required"`
	// Currency code to be queried.
	Currency string `json:"currency" api:"required"`
	// Amount or quantity locked by pending orders or withdrawals.
	Locked string `json:"locked" api:"required"`
	// Currency unit used as the basis for avg_buy_price.
	UnitCurrency string `json:"unit_currency" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		AvgBuyPrice         respjson.Field
		AvgBuyPriceModified respjson.Field
		Balance             respjson.Field
		Currency            respjson.Field
		Locked              respjson.Field
		UnitCurrency        respjson.Field
		ExtraFields         map[string]respjson.Field
		raw                 string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AccountListResponse) RawJSON() string { return r.JSON.raw }
func (r *AccountListResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}
