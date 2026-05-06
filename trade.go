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

// TradeService contains methods and other services that help with interacting with
// the upbit API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewTradeService] method instead.
type TradeService struct {
	options []option.RequestOption
}

// NewTradeService generates a new service that applies the given options to each
// request. These options are applied after the parent client's options (if there
// is one), and before any request-specific options.
func NewTradeService(opts ...option.RequestOption) (r TradeService) {
	r = TradeService{}
	r.options = opts
	return
}

// Recent Trades History
//
// GET /v1/trades/ticks
//
// - docs(kr): https://docs.upbit.com/kr/reference/recent-trades-history
// - docs(global): https://global-docs.upbit.com/reference/recent-trades-history
func (r *TradeService) List(ctx context.Context, query TradeListParams, opts ...option.RequestOption) (res *[]TradeListResponse, err error) {
	var preClientOpts = []option.RequestOption{requestconfig.WithSecurity(requestconfig.Security{})}
	opts = slices.Concat(preClientOpts, r.options, opts)
	path := "v1/trades/ticks"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, query, &res, opts...)
	return res, err
}

type TradeListResponse struct {
	// Trade direction — ASK (sell) or BID (buy).
	//
	// Any of "ASK", "BID".
	AskBid TradeListResponseAskBid `json:"ask_bid" api:"required"`
	// Price change from previous day's closing price. Calculated as "trade_price" -
	// "prev_closing_price".
	//
	// - Positive(+): current price rose from previous closing
	// - Negative(-): current price fell from previous closing
	// - 0: no change from previous closing
	ChangePrice string `json:"change_price" api:"required"`
	// Trading pair code.
	Market string `json:"market" api:"required"`
	// Previous day's closing price (UTC-based).
	PrevClosingPrice string `json:"prev_closing_price" api:"required"`
	// Unique identifier for the trade.
	//
	// This field does not guarantee the sequence of trades.
	SequentialID int64 `json:"sequential_id" api:"required"`
	// Milliseconds timestamp when the trade was executed.
	Timestamp int64 `json:"timestamp" api:"required"`
	// Trade date (UTC-based).
	//
	// [Format] yyyy-MM-dd
	TradeDateUtc string `json:"trade_date_utc" api:"required"`
	// Most recent trade price.
	TradePrice string `json:"trade_price" api:"required"`
	// Trade time (UTC-based).
	//
	// [Format] HH:mm:ss
	TradeTimeUtc string `json:"trade_time_utc" api:"required"`
	// Most recent trade volume.
	TradeVolume string `json:"trade_volume" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		AskBid           respjson.Field
		ChangePrice      respjson.Field
		Market           respjson.Field
		PrevClosingPrice respjson.Field
		SequentialID     respjson.Field
		Timestamp        respjson.Field
		TradeDateUtc     respjson.Field
		TradePrice       respjson.Field
		TradeTimeUtc     respjson.Field
		TradeVolume      respjson.Field
		ExtraFields      map[string]respjson.Field
		Raw              string
	} `json:"-"`
}

// RawJSON returns the unmodified JSON received from the API.
func (r TradeListResponse) RawJSON() string { return r.JSON.Raw }
func (r *TradeListResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Trade direction — ASK (sell) or BID (buy).
type TradeListResponseAskBid string

const (
	TradeListResponseAskBidAsk TradeListResponseAskBid = "ASK"
	TradeListResponseAskBidBid TradeListResponseAskBid = "BID"
)

type TradeListParams struct {
	// Trading pair code to query.
	Market string `query:"market" api:"required" json:"-"`
	// Number of trade records to retrieve.
	//
	// Up to 500 supported. Default: 1.
	Count param.Opt[int64] `query:"count,omitzero" json:"-"`
	// Cursor for pagination. Enter the "sequential_id" from the response to retrieve
	// the previous "count" trade records prior to that trade.
	Cursor param.Opt[string] `query:"cursor,omitzero" json:"-"`
	// Day offset between the query date and the request date. Must specify the target
	// date; up to 7 days of history is supported (UTC-based).
	//
	// Integer between 1 and 7. If omitted, returns trades for the current date. If 7,
	// returns trades from 7 days ago in reverse chronological order.
	DaysAgo param.Opt[int64] `query:"days_ago,omitzero" json:"-"`
	// End time within the query date range (UTC). Optional parameter to retrieve
	// trades within a specific time of the query date.
	//
	// Accepts HHmmss or HH:mm:ss format. Trade list is returned in reverse
	// chronological order from the specified time.
	To param.Opt[string] `query:"to,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [TradeListParams]'s query parameters as `url.Values`.
func (r TradeListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}
