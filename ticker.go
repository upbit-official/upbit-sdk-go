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
	"github.com/upbit-official/upbit-sdk-go/packages/respjson"
)

// TickerService contains methods and other services that help with interacting
// with the upbit API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewTickerService] method instead.
type TickerService struct {
	options []option.RequestOption
}

// NewTickerService generates a new service that applies the given options to each
// request. These options are applied after the parent client's options (if there
// is one), and before any request-specific options.
func NewTickerService(opts ...option.RequestOption) (r TickerService) {
	r = TickerService{}
	r.options = opts
	return
}

// List Tickers by Market
//
// GET /v1/ticker/all
//
// - docs(kr): https://docs.upbit.com/kr/reference/list-quote-tickers
// - docs(sg/th): https://global-docs.upbit.com/reference/list-quote-tickers
// - docs(id): https://docs.upbit.co.id/reference/list-quote-tickers
//
// Retrieves current price information of all pairs in the specified market(quote
// currency).
func (r *TickerService) ListByQuoteCurrencies(ctx context.Context, query TickerListByQuoteCurrenciesParams, opts ...option.RequestOption) (res *[]Ticker, err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithSecurity(requestconfig.Security{})}
	opts = slices.Concat(preClientOpts, r.options, opts)
	path := "v1/ticker/all"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, query, &res, opts...)
	return res, err
}

// List Tickers by Pairs
//
// GET /v1/ticker
//
// - docs(kr): https://docs.upbit.com/kr/reference/list-tickers
// - docs(sg/th): https://global-docs.upbit.com/reference/list-tickers
// - docs(id): https://docs.upbit.co.id/reference/list-tickers
//
// Retrieves the current price for the specified pair. Returns the ticker snapshot
// of the pair at the time of the request.
func (r *TickerService) ListByTradingPairs(ctx context.Context, query TickerListByTradingPairsParams, opts ...option.RequestOption) (res *[]Ticker, err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithSecurity(requestconfig.Security{})}
	opts = slices.Concat(preClientOpts, r.options, opts)
	path := "v1/ticker"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, query, &res, opts...)
	return res, err
}

type Ticker struct {
	// Accumulated trade amount since UTC midnight.
	AccTradePrice string `json:"acc_trade_price" api:"required"`
	// Accumulated trade amount over the past 24 hours.
	AccTradePrice24h string `json:"acc_trade_price_24h" api:"required"`
	// Accumulated trade volume since UTC midnight.
	AccTradeVolume string `json:"acc_trade_volume" api:"required"`
	// Accumulated trade volume over the past 24 hours.
	AccTradeVolume24h string `json:"acc_trade_volume_24h" api:"required"`
	// Price change status.
	//
	// - `EVEN`: No change
	// - `RISE`: Price rose
	// - `FALL`: Price fell
	//
	// Any of "EVEN", "RISE", "FALL".
	Change TickerChange `json:"change" api:"required"`
	// Absolute price change from previous day's closing price.
	//
	// Calculated as "trade_price" - "prev_closing_price".
	ChangePrice string `json:"change_price" api:"required"`
	// Absolute price change rate from previous day's closing price.
	//
	// Calculated as ("trade_price" - "prev_closing_price") / "prev_closing_price".
	ChangeRate string `json:"change_rate" api:"required"`
	// Highest trading price of the trading pair.
	HighPrice string `json:"high_price" api:"required"`
	// Date of 52-week highest price.
	//
	// [Format] yyyy-MM-dd
	Highest52WeekDate string `json:"highest_52_week_date" api:"required"`
	// 52-week highest price.
	Highest52WeekPrice string `json:"highest_52_week_price" api:"required"`
	// Lowest trading price of the trading pair.
	LowPrice string `json:"low_price" api:"required"`
	// Date of 52-week lowest price.
	//
	// [Format] yyyy-MM-dd
	Lowest52WeekDate string `json:"lowest_52_week_date" api:"required"`
	// 52-week lowest price.
	Lowest52WeekPrice string `json:"lowest_52_week_price" api:"required"`
	// Trading pair code.
	Market string `json:"market" api:"required"`
	// Opening price of the trading pair.
	OpeningPrice string `json:"opening_price" api:"required"`
	// Previous day's closing price (UTC-based).
	PrevClosingPrice string `json:"prev_closing_price" api:"required"`
	// Signed price change from previous day's closing price. Calculated as
	// "trade_price" - "prev_closing_price".
	//
	// - Positive(+): price rose
	// - Negative(-): price fell
	SignedChangePrice string `json:"signed_change_price" api:"required"`
	// Signed price change rate from previous day's closing price. Calculated as
	// ("trade_price" - "prev_closing_price") / "prev_closing_price".
	//
	// - Positive(+): price rose
	// - Negative(-): price fell
	//
	// [Example] 0.015 = 1.5% increase
	SignedChangeRate string `json:"signed_change_rate" api:"required"`
	// Timestamp (ms) when the ticker information was recorded.
	Timestamp int64 `json:"timestamp" api:"required"`
	// Most recent trade date (UTC-based).
	//
	// [Format] yyyyMMdd
	TradeDate string `json:"trade_date" api:"required"`
	// Current (closing) price of the trading pair.
	TradePrice string `json:"trade_price" api:"required"`
	// Most recent trade time (UTC-based).
	//
	// [Format] HHmmss
	TradeTime string `json:"trade_time" api:"required"`
	// Milliseconds timestamp when the trade was executed.
	TradeTimestamp int64 `json:"trade_timestamp" api:"required"`
	// Most recent trade volume.
	TradeVolume string `json:"trade_volume" api:"required"`
	// Most recent trade date (KST-based). Supported in KOREA only.
	//
	// [Format] yyyyMMdd
	TradeDateKst string `json:"trade_date_kst"`
	// Most recent trade time (KST-based). Supported in KOREA only.
	//
	// [Format] HHmmss
	TradeTimeKst string `json:"trade_time_kst"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		AccTradePrice      respjson.Field
		AccTradePrice24h   respjson.Field
		AccTradeVolume     respjson.Field
		AccTradeVolume24h  respjson.Field
		Change             respjson.Field
		ChangePrice        respjson.Field
		ChangeRate         respjson.Field
		HighPrice          respjson.Field
		Highest52WeekDate  respjson.Field
		Highest52WeekPrice respjson.Field
		LowPrice           respjson.Field
		Lowest52WeekDate   respjson.Field
		Lowest52WeekPrice  respjson.Field
		Market             respjson.Field
		OpeningPrice       respjson.Field
		PrevClosingPrice   respjson.Field
		SignedChangePrice  respjson.Field
		SignedChangeRate   respjson.Field
		Timestamp          respjson.Field
		TradeDate          respjson.Field
		TradePrice         respjson.Field
		TradeTime          respjson.Field
		TradeTimestamp     respjson.Field
		TradeVolume        respjson.Field
		TradeDateKst       respjson.Field
		TradeTimeKst       respjson.Field
		ExtraFields        map[string]respjson.Field
		raw                string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r Ticker) RawJSON() string { return r.JSON.raw }
func (r *Ticker) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Price change status.
//
// - `EVEN`: No change
// - `RISE`: Price rose
// - `FALL`: Price fell
type TickerChange string

const (
	TickerChangeEven TickerChange = "EVEN"
	TickerChangeRise TickerChange = "RISE"
	TickerChangeFall TickerChange = "FALL"
)

type TickerListByQuoteCurrenciesParams struct {
	// List of quote currencies to query. For multiple markets, use comma-separated
	// format.
	//
	// [Example] SGD,BTC
	QuoteCurrencies string `query:"quote_currencies" api:"required" json:"-"`
	paramObj
}

// URLQuery serializes [TickerListByQuoteCurrenciesParams]'s query parameters as
// `url.Values`.
func (r TickerListByQuoteCurrenciesParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type TickerListByTradingPairsParams struct {
	// List of trading pairs to query.
	//
	// For multiple pairs, use comma-separated format.
	//
	// [Example] SGD-BTC,SGD-ETH
	Markets string `query:"markets" api:"required" json:"-"`
	paramObj
}

// URLQuery serializes [TickerListByTradingPairsParams]'s query parameters as
// `url.Values`.
func (r TickerListByTradingPairsParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}
