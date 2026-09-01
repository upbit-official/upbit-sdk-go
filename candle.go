package upbit

import (
	"context"
	"fmt"
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

// CandleService contains methods and other services that help with interacting
// with the upbit API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewCandleService] method instead.
type CandleService struct {
	options []option.RequestOption
}

// NewCandleService generates a new service that applies the given options to each
// request. These options are applied after the parent client's options (if there
// is one), and before any request-specific options.
func NewCandleService(opts ...option.RequestOption) (r CandleService) {
	r = CandleService{}
	r.options = opts
	return
}

// List Day Candles
//
// GET /v1/candles/days
//
// - docs(kr): https://docs.upbit.com/kr/reference/list-candles-days
// - docs(global): https://global-docs.upbit.com/reference/list-candles-days
//
// Retrieves the list of candles per day.
func (r *CandleService) ListDays(ctx context.Context, query CandleListDaysParams, opts ...option.RequestOption) (res *[]CandleListDaysResponse, err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithSecurity(requestconfig.Security{})}
	opts = slices.Concat(preClientOpts, r.options, opts)
	path := "v1/candles/days"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, query, &res, opts...)
	return res, err
}

// List Minute Candles
//
// GET /v1/candles/minutes/{unit}
//
// - docs(kr): https://docs.upbit.com/kr/reference/list-candles-minutes
// - docs(global): https://global-docs.upbit.com/reference/list-candles-minutes
//
// Retrieves the list of candles per minute.
func (r *CandleService) ListMinutes(ctx context.Context, unit int64, query CandleListMinutesParams, opts ...option.RequestOption) (res *[]CandleListMinutesResponse, err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithSecurity(requestconfig.Security{})}
	opts = slices.Concat(preClientOpts, r.options, opts)
	path := fmt.Sprintf("v1/candles/minutes/%v", unit)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, query, &res, opts...)
	return res, err
}

// List Month Candles
//
// GET /v1/candles/months
//
// - docs(kr): https://docs.upbit.com/kr/reference/list-candles-months
// - docs(global): https://global-docs.upbit.com/reference/list-candles-months
//
// Retrieves the list of candles per month.
func (r *CandleService) ListMonths(ctx context.Context, query CandleListMonthsParams, opts ...option.RequestOption) (res *[]CandlePeriod, err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithSecurity(requestconfig.Security{})}
	opts = slices.Concat(preClientOpts, r.options, opts)
	path := "v1/candles/months"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, query, &res, opts...)
	return res, err
}

// List Second Candles
//
// GET /v1/candles/seconds
//
// - docs(kr): https://docs.upbit.com/kr/reference/list-candles-seconds
// - docs(global): https://global-docs.upbit.com/reference/list-candles-seconds
//
// Retrieves the list of candles per second.
func (r *CandleService) ListSeconds(ctx context.Context, query CandleListSecondsParams, opts ...option.RequestOption) (res *[]Candle, err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithSecurity(requestconfig.Security{})}
	opts = slices.Concat(preClientOpts, r.options, opts)
	path := "v1/candles/seconds"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, query, &res, opts...)
	return res, err
}

// List Week Candles
//
// GET /v1/candles/weeks
//
// - docs(kr): https://docs.upbit.com/kr/reference/list-candles-weeks
// - docs(global): https://global-docs.upbit.com/reference/list-candles-weeks
//
// Retrieves the list of candles per week.
func (r *CandleService) ListWeeks(ctx context.Context, query CandleListWeeksParams, opts ...option.RequestOption) (res *[]CandlePeriod, err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithSecurity(requestconfig.Security{})}
	opts = slices.Concat(preClientOpts, r.options, opts)
	path := "v1/candles/weeks"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, query, &res, opts...)
	return res, err
}

// List Year Candles
//
// GET /v1/candles/years
//
// - docs(kr): https://docs.upbit.com/kr/reference/list-candles-years
// - docs(global): https://global-docs.upbit.com/reference/list-candles-years
//
// Retrieves the list of candles per year.
func (r *CandleService) ListYears(ctx context.Context, query CandleListYearsParams, opts ...option.RequestOption) (res *[]CandlePeriod, err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithSecurity(requestconfig.Security{})}
	opts = slices.Concat(preClientOpts, r.options, opts)
	path := "v1/candles/years"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, query, &res, opts...)
	return res, err
}

type Candle struct {
	// The total trade amount (in the quoted currency) accumulated during the candle
	// period.
	CandleAccTradePrice string `json:"candle_acc_trade_price" api:"required"`
	// The total traded volume (in the base asset) accumulated during the candle
	// period.
	CandleAccTradeVolume string `json:"candle_acc_trade_volume" api:"required"`
	// Start time of the candle period in KST. Supported in KOREA only.
	//
	// [Format]: yyyy-MM-dd'T'HH:mm:ss
	CandleDateTimeKst string `json:"candle_date_time_kst" api:"required"`
	// Start time of the candle period in UTC.
	//
	// [Format]: yyyy-MM-dd'T'HH:mm:ss
	CandleDateTimeUtc string `json:"candle_date_time_utc" api:"required"`
	// The highest trading price, recorded during the candle period.
	HighPrice string `json:"high_price" api:"required"`
	// The lowest trading price, recorded during the candle period.
	LowPrice string `json:"low_price" api:"required"`
	// Trading pair code representing the market.
	//
	// [Example]: "SGD-BTC"
	Market string `json:"market" api:"required"`
	// The opening price of the candle, representing the first trading price during the
	// candle period.
	OpeningPrice string `json:"opening_price" api:"required"`
	// The timestamp (in milliseconds) when the last tick of the candle was recorded.
	Timestamp int64 `json:"timestamp" api:"required"`
	// The closing price of the candle, representing the last trading price during the
	// candle period.
	TradePrice string `json:"trade_price" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		CandleAccTradePrice  respjson.Field
		CandleAccTradeVolume respjson.Field
		CandleDateTimeKst    respjson.Field
		CandleDateTimeUtc    respjson.Field
		HighPrice            respjson.Field
		LowPrice             respjson.Field
		Market               respjson.Field
		OpeningPrice         respjson.Field
		Timestamp            respjson.Field
		TradePrice           respjson.Field
		ExtraFields          map[string]respjson.Field
		raw                  string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r Candle) RawJSON() string { return r.JSON.raw }
func (r *Candle) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type CandlePeriod struct {
	// Start date of candle aggregation period (yyyy-MM-dd).
	FirstDayOfPeriod string `json:"first_day_of_period" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		FirstDayOfPeriod respjson.Field
		ExtraFields      map[string]respjson.Field
		raw              string
	} `json:"-"`
	Candle
}

// Returns the unmodified JSON received from the API
func (r CandlePeriod) RawJSON() string { return r.JSON.raw }
func (r *CandlePeriod) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type CandleListDaysResponse struct {
	// Price change = trade_price - prev_closing_price.
	ChangePrice string `json:"change_price" api:"required"`
	// Change rate = (trade_price - prev_closing_price) / prev_closing_price.
	ChangeRate string `json:"change_rate" api:"required"`
	// Previous day's closing price (UTC-based).
	PrevClosingPrice string `json:"prev_closing_price" api:"required"`
	// Closing price converted per converting_price_unit. Null for non-digital asset
	// markets.
	ConvertedTradePrice string `json:"converted_trade_price"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ChangePrice         respjson.Field
		ChangeRate          respjson.Field
		PrevClosingPrice    respjson.Field
		ConvertedTradePrice respjson.Field
		ExtraFields         map[string]respjson.Field
		raw                 string
	} `json:"-"`
	Candle
}

// Returns the unmodified JSON received from the API
func (r CandleListDaysResponse) RawJSON() string { return r.JSON.raw }
func (r *CandleListDaysResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type CandleListMinutesResponse struct {
	// Candle unit in minutes.
	//
	// Any of 1, 3, 5, 10, 15, 30, 60, 240.
	Unit int64 `json:"unit" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Unit        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
	Candle
}

// Returns the unmodified JSON received from the API
func (r CandleListMinutesResponse) RawJSON() string { return r.JSON.raw }
func (r *CandleListMinutesResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type CandleListDaysParams struct {
	// Trading pair code to query.
	Market string `query:"market" api:"required" json:"-"`
	// Currency to convert the closing price into. When specified, the response
	// includes a "converted_trade_price" field with the closing price converted to the
	// given currency.
	//
	// [Example] Specifying "KRW" returns the closing price converted to KRW.
	ConvertingPriceUnit param.Opt[string] `query:"converting_price_unit,omitzero" json:"-"`
	// Number of candles to retrieve.
	//
	// Up to 200 candles are supported. Default: 1.
	Count param.Opt[int64] `query:"count,omitzero" json:"-"`
	// End time of the query range. Retrieves candles before the specified time. If not
	// specified, the most recent candles are returned.
	//
	// Accepts ISO 8601 datetime format.
	To param.Opt[string] `query:"to,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [CandleListDaysParams]'s query parameters as `url.Values`.
func (r CandleListDaysParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type CandleListMinutesParams struct {
	// Trading pair code to query.
	Market string `query:"market" api:"required" json:"-"`
	// Number of candles to retrieve.
	//
	// Up to 200 candles are supported. Default: 1.
	Count param.Opt[int64] `query:"count,omitzero" json:"-"`
	// End time of the query range. Retrieves candles before the specified time. If not
	// specified, the most recent candles are returned.
	//
	// Accepts ISO 8601 datetime format.
	//
	// [Example] 2025-06-24T04:56:53Z 2025-06-24 04:56:53 2025-06-24T13:56:53+08:00
	To param.Opt[string] `query:"to,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [CandleListMinutesParams]'s query parameters as
// `url.Values`.
func (r CandleListMinutesParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type CandleListMonthsParams struct {
	// Trading pair code to query.
	Market string `query:"market" api:"required" json:"-"`
	// Number of candles to retrieve. Up to 200, default: 1.
	Count param.Opt[int64] `query:"count,omitzero" json:"-"`
	// End time of the query range. If not specified, the most recent candles are
	// returned.
	//
	// Accepts ISO 8601 datetime format.
	To param.Opt[string] `query:"to,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [CandleListMonthsParams]'s query parameters as `url.Values`.
func (r CandleListMonthsParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type CandleListSecondsParams struct {
	// Trading pair code to query.
	Market string `query:"market" api:"required" json:"-"`
	// Number of candles to retrieve.
	//
	// Up to 200 candles are supported. Default: 1.
	Count param.Opt[int64] `query:"count,omitzero" json:"-"`
	// End time of the query range. Retrieves candles before the specified time. If not
	// specified, the most recent candles are returned.
	//
	// Accepts ISO 8601 datetime format. URL encoding is required for spaces and
	// special characters.
	//
	// [Example] 2025-06-24T04:56:53Z 2025-06-24 04:56:53 2025-06-24T13:56:53+08:00
	//
	// Second candles only support data up to 3 months prior. An empty array is
	// returned for older times.
	To param.Opt[string] `query:"to,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [CandleListSecondsParams]'s query parameters as
// `url.Values`.
func (r CandleListSecondsParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type CandleListWeeksParams struct {
	// Trading pair code to query.
	Market string `query:"market" api:"required" json:"-"`
	// Number of candles to retrieve. Up to 200, default: 1.
	Count param.Opt[int64] `query:"count,omitzero" json:"-"`
	// End time of the query range. Retrieves candles before the specified time. If not
	// specified, the most recent candles are returned.
	To param.Opt[string] `query:"to,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [CandleListWeeksParams]'s query parameters as `url.Values`.
func (r CandleListWeeksParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type CandleListYearsParams struct {
	// Trading pair code to query.
	Market string `query:"market" api:"required" json:"-"`
	// Number of candles to retrieve. Up to 200, default: 1.
	Count param.Opt[int64] `query:"count,omitzero" json:"-"`
	// End time of the query range. ISO 8601 format.
	To param.Opt[string] `query:"to,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [CandleListYearsParams]'s query parameters as `url.Values`.
func (r CandleListYearsParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}
