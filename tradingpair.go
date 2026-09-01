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

// TradingPairService contains methods and other services that help with
// interacting with the upbit API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewTradingPairService] method instead.
type TradingPairService struct {
	options []option.RequestOption
}

// NewTradingPairService generates a new service that applies the given options to
// each request. These options are applied after the parent client's options (if
// there is one), and before any request-specific options.
func NewTradingPairService(opts ...option.RequestOption) (r TradingPairService) {
	r = TradingPairService{}
	r.options = opts
	return
}

// List Trading Pairs
//
// GET /v1/market/all
//
// - docs(kr): https://docs.upbit.com/kr/reference/list-trading-pairs
// - docs(global): https://global-docs.upbit.com/reference/list-trading-pairs
//
// Retrieves the list of all trading pairs supported by Upbit.
func (r *TradingPairService) List(ctx context.Context, query TradingPairListParams, opts ...option.RequestOption) (res *[]TradingPairListResponse, err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithSecurity(requestconfig.Security{})}
	opts = slices.Concat(preClientOpts, r.options, opts)
	path := "v1/market/all"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, query, &res, opts...)
	return res, err
}

type TradingPairListResponse struct {
	// English name of the digital asset.
	EnglishName string `json:"english_name" api:"required"`
	// Korean name of the digital asset.
	KoreanName string `json:"korean_name" api:"required"`
	// Trading pair code representing the market.
	//
	// [Example]: "SGD-BTC"
	Market string `json:"market" api:"required"`
	// Market alert information.
	MarketEvent TradingPairListResponseMarketEvent `json:"market_event"`
	// Trading pair warning status.
	//
	// Any of "NONE", "CAUTION".
	MarketWarning TradingPairListResponseMarketWarning `json:"market_warning"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		EnglishName   respjson.Field
		KoreanName    respjson.Field
		Market        respjson.Field
		MarketEvent   respjson.Field
		MarketWarning respjson.Field
		ExtraFields   map[string]respjson.Field
		raw           string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r TradingPairListResponse) RawJSON() string { return r.JSON.raw }
func (r *TradingPairListResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Market alert information.
type TradingPairListResponseMarketEvent struct {
	// Caution item status. If flagged, one or more of the following alert types may
	// apply.
	Caution TradingPairListResponseMarketEventCaution `json:"caution"`
	// Whether the trading pair is flagged as a warning item.
	Warning bool `json:"warning"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Caution     respjson.Field
		Warning     respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r TradingPairListResponseMarketEvent) RawJSON() string { return r.JSON.raw }
func (r *TradingPairListResponseMarketEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Caution item status. If flagged, one or more of the following alert types may
// apply.
type TradingPairListResponseMarketEventCaution struct {
	// Concentration of small accounts alert.
	ConcentrationOfSmallAccounts bool `json:"CONCENTRATION_OF_SMALL_ACCOUNTS"`
	// Deposit amount surge alert.
	DepositAmountSoaring bool `json:"DEPOSIT_AMOUNT_SOARING"`
	// Global price difference alert.
	GlobalPriceDifferences bool `json:"GLOBAL_PRICE_DIFFERENCES"`
	// Price fluctuation alert.
	PriceFluctuations bool `json:"PRICE_FLUCTUATIONS"`
	// Trading volume surge alert.
	TradingVolumeSoaring bool `json:"TRADING_VOLUME_SOARING"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ConcentrationOfSmallAccounts respjson.Field
		DepositAmountSoaring         respjson.Field
		GlobalPriceDifferences       respjson.Field
		PriceFluctuations            respjson.Field
		TradingVolumeSoaring         respjson.Field
		ExtraFields                  map[string]respjson.Field
		raw                          string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r TradingPairListResponseMarketEventCaution) RawJSON() string { return r.JSON.raw }
func (r *TradingPairListResponseMarketEventCaution) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Trading pair warning status.
type TradingPairListResponseMarketWarning string

const (
	TradingPairListResponseMarketWarningNone    TradingPairListResponseMarketWarning = "NONE"
	TradingPairListResponseMarketWarningCaution TradingPairListResponseMarketWarning = "CAUTION"
)

type TradingPairListParams struct {
	// Whether to include detailed information in the query. If true, the response
	// includes detail fields such as caution or warning designation. Default: false.
	IsDetails param.Opt[bool] `query:"is_details,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [TradingPairListParams]'s query parameters as `url.Values`.
func (r TradingPairListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}
