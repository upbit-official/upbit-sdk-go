package upbit

import (
	"context"
	"net/http"
	"net/url"
	"slices"

	"github.com/upbit-official/upbit-sdk-go/internal/apijson"
	"github.com/upbit-official/upbit-sdk-go/internal/apiquery"
	"github.com/upbit-official/upbit-sdk-go/internal/requestconfig"
	"github.com/upbit-official/upbit-sdk-go/option"
	"github.com/upbit-official/upbit-sdk-go/packages/param"
	"github.com/upbit-official/upbit-sdk-go/packages/respjson"
)

// OrderbookService contains methods and other services that help with interacting
// with the upbit API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewOrderbookService] method instead.
type OrderbookService struct {
	options []option.RequestOption
}

// NewOrderbookService generates a new service that applies the given options to
// each request. These options are applied after the parent client's options (if
// there is one), and before any request-specific options.
func NewOrderbookService(opts ...option.RequestOption) (r OrderbookService) {
	r = OrderbookService{}
	r.options = opts
	return
}

// Get Orderbook
//
// GET /v1/orderbook
//
// - docs(kr): https://docs.upbit.com/kr/reference/list-orderbooks
// - docs(sg/th): https://global-docs.upbit.com/reference/list-orderbooks
// - docs(id): https://docs.upbit.co.id/reference/list-orderbooks
//
// Retrieves the real-time orderbook data for the specified markets.
func (r *OrderbookService) List(ctx context.Context, query OrderbookListParams, opts ...option.RequestOption) (res *[]OrderbookListResponse, err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithSecurity(requestconfig.Security{})}
	opts = slices.Concat(preClientOpts, r.options, opts)
	path := "v1/orderbook"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, query, &res, opts...)
	return res, err
}

// List Orderbook Instruments
//
// GET /v1/orderbook/instruments
//
//   - docs(kr): https://docs.upbit.com/kr/reference/list-orderbook-instruments
//   - docs(sg/th):
//     https://global-docs.upbit.com/reference/list-orderbook-instruments
//   - docs(id): https://docs.upbit.co.id/reference/list-orderbook-instruments
//
// Retrieves the tick size (order price unit) for the specified trading pairs.
func (r *OrderbookService) ListInstruments(ctx context.Context, query OrderbookListInstrumentsParams, opts ...option.RequestOption) (res *[]OrderbookListInstrumentsResponse, err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithSecurity(requestconfig.Security{})}
	opts = slices.Concat(preClientOpts, r.options, opts)
	path := "v1/orderbook/instruments"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, query, &res, opts...)
	return res, err
}

type OrderbookListResponse struct {
	// Price unit applied for this orderbook level.
	Level string `json:"level" api:"required"`
	// Trading pair code.
	Market string `json:"market" api:"required"`
	// List of orderbook entries ordered from level 1 to level 30.
	OrderbookUnits []OrderbookListResponseOrderbookUnit `json:"orderbook_units" api:"required"`
	// Timestamp of the query request (ms).
	Timestamp int64 `json:"timestamp" api:"required"`
	// Total ask volume in the current orderbook.
	TotalAskSize string `json:"total_ask_size" api:"required"`
	// Total bid volume in the current orderbook.
	TotalBidSize string `json:"total_bid_size" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Level          respjson.Field
		Market         respjson.Field
		OrderbookUnits respjson.Field
		Timestamp      respjson.Field
		TotalAskSize   respjson.Field
		TotalBidSize   respjson.Field
		ExtraFields    map[string]respjson.Field
		raw            string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r OrderbookListResponse) RawJSON() string { return r.JSON.raw }
func (r *OrderbookListResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Individual orderbook entry.
type OrderbookListResponseOrderbookUnit struct {
	// Lowest sell (ask) price.
	AskPrice string `json:"ask_price" api:"required"`
	// Sell quantity at ask price.
	AskSize string `json:"ask_size" api:"required"`
	// Highest buy (bid) price.
	BidPrice string `json:"bid_price" api:"required"`
	// Buy quantity at bid price.
	BidSize string `json:"bid_size" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		AskPrice    respjson.Field
		AskSize     respjson.Field
		BidPrice    respjson.Field
		BidSize     respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r OrderbookListResponseOrderbookUnit) RawJSON() string { return r.JSON.raw }
func (r *OrderbookListResponseOrderbookUnit) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type OrderbookListInstrumentsResponse struct {
	// Trading pair code.
	Market string `json:"market" api:"required"`
	// Quote currency code for the trading pair (e.g., SGD, BTC, USDT).
	QuoteCurrency string `json:"quote_currency" api:"required"`
	// Price unit applied to the orderbook level.
	TickSize string `json:"tick_size" api:"required"`
	// Supported orderbook aggregation levels for this trading pair.
	//
	// - 0: Default unit
	// - Aggregation is only supported for KRW markets. (BTC, USDT markets only have 0)
	SupportedLevels []string `json:"supported_levels"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Market          respjson.Field
		QuoteCurrency   respjson.Field
		TickSize        respjson.Field
		SupportedLevels respjson.Field
		ExtraFields     map[string]respjson.Field
		raw             string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r OrderbookListInstrumentsResponse) RawJSON() string { return r.JSON.raw }
func (r *OrderbookListInstrumentsResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type OrderbookListParams struct {
	// List of trading pairs to query.
	//
	// For multiple pairs, use comma-separated format.
	Markets string `query:"markets" api:"required" json:"-"`
	// Number of orderbook entries to retrieve.
	//
	// Based on the best bid-ask pair, returns the specified number of pairs.
	// Default: 30.
	Count param.Opt[int64] `query:"count,omitzero" json:"-"`
	// Orderbook aggregation level. Only supported for KRW markets. Groups ask/bid
	// price and size by the specified unit. Provide as a numeric string. Use an
	// integer string for units >= 1, or a double string for fractional units. Defaults
	// to 0 if not specified.
	Level param.Opt[string] `query:"level,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [OrderbookListParams]'s query parameters as `url.Values`.
func (r OrderbookListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type OrderbookListInstrumentsParams struct {
	// List of trading pairs to query.
	//
	// For multiple pairs, use comma-separated format.
	Markets string `query:"markets" api:"required" json:"-"`
	paramObj
}

// URLQuery serializes [OrderbookListInstrumentsParams]'s query parameters as
// `url.Values`.
func (r OrderbookListInstrumentsParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}
