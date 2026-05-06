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
	"github.com/upbit-official/upbit-sdk-go/packages/pagination"
	"github.com/upbit-official/upbit-sdk-go/packages/param"
	"github.com/upbit-official/upbit-sdk-go/packages/respjson"
)

// OrderService contains methods and other services that help with interacting with
// the upbit API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewOrderService] method instead.
type OrderService struct {
	options []option.RequestOption
}

// NewOrderService generates a new service that applies the given options to each
// request. These options are applied after the parent client's options (if there
// is one), and before any request-specific options.
func NewOrderService(opts ...option.RequestOption) (r OrderService) {
	r = OrderService{}
	r.options = opts
	return
}

// Create Order
//
// POST /v1/orders
//
// - docs(kr): https://docs.upbit.com/kr/reference/new-order
// - docs(global): https://global-docs.upbit.com/reference/new-order
//
// Order parameter rules by `ord_type`:
//
//   - `limit`: requires both `price` and `volume`. Optional: `time_in_force` (`ioc`,
//     `fok`, `post_only`) and `smp_type`. Note: `post_only` cannot be used with
//     `smp_type`.
//
// - `price` (market buy): requires `price`, and `volume` must be omitted.
//
// - `market` (market sell): requires `volume`, and `price` must be omitted.
//
// - `best`: `time_in_force` is required (`ioc` or `fok`).
//   - `side=bid`: requires `price`, and `volume` must be omitted.
//   - `side=ask`: requires `volume`, and `price` must be omitted.
func (r *OrderService) New(ctx context.Context, body OrderNewParams, opts ...option.RequestOption) (res *Order, err error) {
	opts = slices.Concat(r.options, opts)
	path := "v1/orders"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// Get Order
//
// GET /v1/order
//
// - docs(kr): https://docs.upbit.com/kr/reference/get-order
// - docs(global): https://global-docs.upbit.com/reference/get-order
//
// **Required (one of):** `uuid` or `identifier` must be provided.
func (r *OrderService) Get(ctx context.Context, query OrderGetParams, opts ...option.RequestOption) (res *Order, err error) {
	opts = slices.Concat(r.options, opts)
	path := "v1/order"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, query, &res, opts...)
	return res, err
}

// Cancel Order
//
// DELETE /v1/order
//
// - docs(kr): https://docs.upbit.com/kr/reference/cancel-order
// - docs(global): https://global-docs.upbit.com/reference/cancel-order
//
// **Required (one of):** `uuid` or `identifier` must be provided.
func (r *OrderService) Cancel(ctx context.Context, body OrderCancelParams, opts ...option.RequestOption) (res *Order, err error) {
	opts = slices.Concat(r.options, opts)
	path := "v1/order"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, body, &res, opts...)
	return res, err
}

// Cancel and New Order
//
// POST /v1/orders/cancel_and_new
//
// - docs(kr): https://docs.upbit.com/kr/reference/cancel-and-new-order
// - docs(global): https://global-docs.upbit.com/reference/cancel-and-new-order
func (r *OrderService) CancelAndNew(ctx context.Context, body OrderCancelAndNewParams, opts ...option.RequestOption) (res *OrderCancelAndNewResponse, err error) {
	opts = slices.Concat(r.options, opts)
	path := "v1/orders/cancel_and_new"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// Cancel Orders by IDs
//
// DELETE /v1/orders/uuids
//
// - docs(kr): https://docs.upbit.com/kr/reference/cancel-orders-by-ids
// - docs(global): https://global-docs.upbit.com/reference/cancel-orders-by-ids
//
// **Required (one of):** `uuids` or `identifiers` must be provided.
func (r *OrderService) CancelByUuids(ctx context.Context, body OrderCancelByUuidsParams, opts ...option.RequestOption) (res *BatchCancelResult, err error) {
	opts = slices.Concat(r.options, opts)
	path := "v1/orders/uuids"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, body, &res, opts...)
	return res, err
}

// Batch Cancel Orders
//
// DELETE /v1/orders/open
//
// - docs(kr): https://docs.upbit.com/kr/reference/batch-cancel-orders
// - docs(global): https://global-docs.upbit.com/reference/batch-cancel-orders
func (r *OrderService) CancelOpen(ctx context.Context, body OrderCancelOpenParams, opts ...option.RequestOption) (res *BatchCancelResult, err error) {
	opts = slices.Concat(r.options, opts)
	path := "v1/orders/open"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, body, &res, opts...)
	return res, err
}

// List Orders by IDs
//
// GET /v1/orders/uuids
//
// - docs(kr): https://docs.upbit.com/kr/reference/list-orders-by-ids
// - docs(global): https://global-docs.upbit.com/reference/list-orders-by-ids
//
// **Required (one of):** `uuids` or `identifiers` must be provided.
func (r *OrderService) ListByUuids(ctx context.Context, query OrderListByUuidsParams, opts ...option.RequestOption) (res *[]Order, err error) {
	opts = slices.Concat(r.options, opts)
	path := "v1/orders/uuids"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, query, &res, opts...)
	return res, err
}

// List Closed Orders
//
// GET /v1/orders/closed
//
// - docs(kr): https://docs.upbit.com/kr/reference/list-closed-orders
// - docs(global): https://global-docs.upbit.com/reference/list-closed-orders
func (r *OrderService) ListClosed(ctx context.Context, query OrderListClosedParams, opts ...option.RequestOption) (res *[]Order, err error) {
	opts = slices.Concat(r.options, opts)
	path := "v1/orders/closed"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, query, &res, opts...)
	return res, err
}

// List Open Orders
//
// GET /v1/orders/open
//
// - docs(kr): https://docs.upbit.com/kr/reference/list-open-orders
// - docs(global): https://global-docs.upbit.com/reference/list-open-orders
func (r *OrderService) ListOpen(ctx context.Context, query OrderListOpenParams, opts ...option.RequestOption) (res *pagination.PageNumberPage[Order], err error) {
	var raw *http.Response
	opts = slices.Concat(r.options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	path := "v1/orders/open"
	cfg, err := requestconfig.NewRequestConfig(ctx, http.MethodGet, path, query, &res, opts...)
	if err != nil {
		return nil, err
	}
	err = cfg.Execute()
	if err != nil {
		return nil, err
	}
	res.SetPageConfig(cfg, raw)
	return res, nil
}

// List Open Orders
//
// GET /v1/orders/open
//
// - docs(kr): https://docs.upbit.com/kr/reference/list-open-orders
// - docs(global): https://global-docs.upbit.com/reference/list-open-orders
func (r *OrderService) ListOpenAutoPaging(ctx context.Context, query OrderListOpenParams, opts ...option.RequestOption) *pagination.PageNumberPageAutoPager[Order] {
	return pagination.NewPageNumberPageAutoPager(r.ListOpen(ctx, query, opts...))
}

// Get Available Order Info
//
// GET /v1/orders/chance
//
//   - docs(kr): https://docs.upbit.com/kr/reference/available-order-information
//   - docs(global):
//     https://global-docs.upbit.com/reference/available-order-information
func (r *OrderService) GetChance(ctx context.Context, query OrderGetChanceParams, opts ...option.RequestOption) (res *OrderGetChanceResponse, err error) {
	opts = slices.Concat(r.options, opts)
	path := "v1/orders/chance"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, query, &res, opts...)
	return res, err
}

// Test Order Creation
//
// POST /v1/orders/test
//
// - docs(kr): https://docs.upbit.com/kr/reference/test-order
// - docs(global): https://global-docs.upbit.com/reference/test-order
//
// Order parameter rules by `ord_type`:
//
//   - `limit`: requires both `price` and `volume`. Optional: `time_in_force` (`ioc`,
//     `fok`, `post_only`) and `smp_type`. Note: `post_only` cannot be used with
//     `smp_type`.
//
// - `price` (market buy): requires `price`, and `volume` must be omitted.
//
// - `market` (market sell): requires `volume`, and `price` must be omitted.
//
// - `best`: `time_in_force` is required (`ioc` or `fok`).
//   - `side=bid`: requires `price`, and `volume` must be omitted.
//   - `side=ask`: requires `volume`, and `price` must be omitted.
func (r *OrderService) TestNew(ctx context.Context, body OrderTestNewParams, opts ...option.RequestOption) (res *Order, err error) {
	opts = slices.Concat(r.options, opts)
	path := "v1/orders/test"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

type BatchCancelResult struct {
	// Information about orders that failed to cancel.
	Failed BatchCancelResultFailed `json:"failed" api:"required"`
	// Information about successfully cancelled orders.
	Success BatchCancelResultSuccess `json:"success" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Failed      respjson.Field
		Success     respjson.Field
		ExtraFields map[string]respjson.Field
		Raw         string
	} `json:"-"`
}

// RawJSON returns the unmodified JSON received from the API.
func (r BatchCancelResult) RawJSON() string { return r.JSON.Raw }
func (r *BatchCancelResult) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Information about orders that failed to cancel.
type BatchCancelResultFailed struct {
	// Number of orders that failed to cancel.
	Count int64 `json:"count" api:"required"`
	// List of orders that failed to cancel.
	Orders []BatchCancelResultFailedOrder `json:"orders" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Count       respjson.Field
		Orders      respjson.Field
		ExtraFields map[string]respjson.Field
		Raw         string
	} `json:"-"`
}

// RawJSON returns the unmodified JSON received from the API.
func (r BatchCancelResultFailed) RawJSON() string { return r.JSON.Raw }
func (r *BatchCancelResultFailed) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BatchCancelResultFailedOrder struct {
	// Trading pair code.
	Market string `json:"market" api:"required"`
	// Unique identifier of the order.
	Uuid string `json:"uuid" api:"required"`
	// Client-assigned order identifier.
	//
	// - Only available for orders created after October 18, 2024.
	Identifier string `json:"identifier"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Market      respjson.Field
		Uuid        respjson.Field
		Identifier  respjson.Field
		ExtraFields map[string]respjson.Field
		Raw         string
	} `json:"-"`
}

// RawJSON returns the unmodified JSON received from the API.
func (r BatchCancelResultFailedOrder) RawJSON() string { return r.JSON.Raw }
func (r *BatchCancelResultFailedOrder) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Information about successfully cancelled orders.
type BatchCancelResultSuccess struct {
	// Number of successfully cancelled orders.
	Count int64 `json:"count" api:"required"`
	// List of successfully cancelled orders.
	Orders []BatchCancelResultSuccessOrder `json:"orders" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Count       respjson.Field
		Orders      respjson.Field
		ExtraFields map[string]respjson.Field
		Raw         string
	} `json:"-"`
}

// RawJSON returns the unmodified JSON received from the API.
func (r BatchCancelResultSuccess) RawJSON() string { return r.JSON.Raw }
func (r *BatchCancelResultSuccess) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BatchCancelResultSuccessOrder struct {
	// Trading pair code.
	Market string `json:"market" api:"required"`
	// Unique identifier of the order.
	Uuid string `json:"uuid" api:"required"`
	// Client-assigned order identifier.
	//
	// - Only available for orders created after October 18, 2024.
	Identifier string `json:"identifier"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Market      respjson.Field
		Uuid        respjson.Field
		Identifier  respjson.Field
		ExtraFields map[string]respjson.Field
		Raw         string
	} `json:"-"`
}

// RawJSON returns the unmodified JSON received from the API.
func (r BatchCancelResultSuccessOrder) RawJSON() string { return r.JSON.Raw }
func (r *BatchCancelResultSuccessOrder) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Closed order state.
//
// - `done`: Fully executed
// - `cancel`: Cancelled
type ClosedOrderState string

const (
	ClosedOrderStateDone   ClosedOrderState = "done"
	ClosedOrderStateCancel ClosedOrderState = "cancel"
)

// Open order state.
//
// - `wait`: Pending execution
// - `watch`: Pending reservation (stop order)
type OpenOrderState string

const (
	OpenOrderStateWait  OpenOrderState = "wait"
	OpenOrderStateWatch OpenOrderState = "watch"
)

// Order type.
//
// - `limit`: Limit order
// - `price`: Market buy order (by total price)
// - `market`: Market sell order (by volume)
// - `best`: Best available price order
type OrdTypeEnum string

const (
	OrdTypeEnumLimit  OrdTypeEnum = "limit"
	OrdTypeEnumPrice  OrdTypeEnum = "price"
	OrdTypeEnumMarket OrdTypeEnum = "market"
	OrdTypeEnumBest   OrdTypeEnum = "best"
)

type Order struct {
	// Timestamp when the order was created.
	//
	// [Format] yyyy-MM-ddTHH:mm:ss+00:00
	CreatedAt string `json:"created_at" api:"required"`
	// Executed volume.
	ExecutedVolume string `json:"executed_volume" api:"required"`
	// Amount locked for this order.
	Locked string `json:"locked" api:"required"`
	// Trading pair code.
	Market string `json:"market" api:"required"`
	// Order type.
	//
	// Any of "limit", "price", "market", "best".
	OrdType OrderOrdType `json:"ord_type" api:"required"`
	// Fee paid so far.
	PaidFee string `json:"paid_fee" api:"required"`
	// Asset unlocked due to Self-Match Prevention cancellation.
	//
	// - For buy orders: cancelled amount
	// - For sell orders: cancelled volume
	PreventedLocked string `json:"prevented_locked" api:"required"`
	// Volume cancelled due to Self-Match Prevention.
	PreventedVolume string `json:"prevented_volume" api:"required"`
	// Remaining fee.
	RemainingFee string `json:"remaining_fee" api:"required"`
	// Remaining volume after execution.
	RemainingVolume string `json:"remaining_volume" api:"required"`
	// Fee reserved for this order.
	ReservedFee string `json:"reserved_fee" api:"required"`
	// Order direction (ask=sell, bid=buy).
	//
	// Any of "ask", "bid".
	Side OrderSide `json:"side" api:"required"`
	// Order state.
	//
	// - `wait`: Pending execution
	// - `watch`: Waiting for reservation order
	// - `done`: Fully executed
	// - `cancel`: Cancelled
	//
	// Any of "wait", "watch", "done", "cancel".
	State OrderState `json:"state" api:"required"`
	// Number of trades executed for this order.
	TradesCount int64 `json:"trades_count" api:"required"`
	// Unique identifier of the order.
	Uuid string `json:"uuid" api:"required"`
	// Client-assigned order identifier specified at order creation.
	//
	//   - The identifier field is only available for orders created after October
	//     18, 2024.
	Identifier string `json:"identifier"`
	// Order price or total amount. For limit orders, this is the unit price. For
	// market buy orders, this is the total buy amount.
	Price string `json:"price"`
	// Self-Match Prevention (SMP) mode.
	//
	// Any of "reduce", "cancel_maker", "cancel_taker".
	SmpType OrderSmpType `json:"smp_type"`
	// Time in force condition.
	//
	// Any of "fok", "ioc", "post_only".
	TimeInForce OrderTimeInForce `json:"time_in_force"`
	// List of trades executed for this order (returned in individual order query).
	Trades []OrderTrade `json:"trades"`
	// Requested order volume.
	Volume string `json:"volume"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		CreatedAt       respjson.Field
		ExecutedVolume  respjson.Field
		Locked          respjson.Field
		Market          respjson.Field
		OrdType         respjson.Field
		PaidFee         respjson.Field
		PreventedLocked respjson.Field
		PreventedVolume respjson.Field
		RemainingFee    respjson.Field
		RemainingVolume respjson.Field
		ReservedFee     respjson.Field
		Side            respjson.Field
		State           respjson.Field
		TradesCount     respjson.Field
		Uuid            respjson.Field
		Identifier      respjson.Field
		Price           respjson.Field
		SmpType         respjson.Field
		TimeInForce     respjson.Field
		Trades          respjson.Field
		Volume          respjson.Field
		ExtraFields     map[string]respjson.Field
		Raw             string
	} `json:"-"`
}

// RawJSON returns the unmodified JSON received from the API.
func (r Order) RawJSON() string { return r.JSON.Raw }
func (r *Order) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Order type.
type OrderOrdType string

const (
	OrderOrdTypeLimit  OrderOrdType = "limit"
	OrderOrdTypePrice  OrderOrdType = "price"
	OrderOrdTypeMarket OrderOrdType = "market"
	OrderOrdTypeBest   OrderOrdType = "best"
)

// Order direction (ask=sell, bid=buy).
type OrderSide string

const (
	OrderSideAsk OrderSide = "ask"
	OrderSideBid OrderSide = "bid"
)

// Order state.
//
// - `wait`: Pending execution
// - `watch`: Waiting for reservation order
// - `done`: Fully executed
// - `cancel`: Cancelled
type OrderState string

const (
	OrderStateWait   OrderState = "wait"
	OrderStateWatch  OrderState = "watch"
	OrderStateDone   OrderState = "done"
	OrderStateCancel OrderState = "cancel"
)

// Self-Match Prevention (SMP) mode.
type OrderSmpType string

const (
	OrderSmpTypeReduce      OrderSmpType = "reduce"
	OrderSmpTypeCancelMaker OrderSmpType = "cancel_maker"
	OrderSmpTypeCancelTaker OrderSmpType = "cancel_taker"
)

// Time in force condition.
type OrderTimeInForce string

const (
	OrderTimeInForceFok      OrderTimeInForce = "fok"
	OrderTimeInForceIoc      OrderTimeInForce = "ioc"
	OrderTimeInForcePostOnly OrderTimeInForce = "post_only"
)

type OrderTrade struct {
	// Timestamp of the trade.
	CreatedAt string `json:"created_at"`
	// Trade funds (price × volume).
	Funds string `json:"funds"`
	// Trading pair code.
	Market string `json:"market"`
	// Trade price.
	Price string `json:"price"`
	// Order side of the trade.
	Side string `json:"side"`
	// Price trend of the trade.
	Trend string `json:"trend"`
	// Unique identifier of the trade.
	Uuid string `json:"uuid"`
	// Trade volume.
	Volume string `json:"volume"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		CreatedAt   respjson.Field
		Funds       respjson.Field
		Market      respjson.Field
		Price       respjson.Field
		Side        respjson.Field
		Trend       respjson.Field
		Uuid        respjson.Field
		Volume      respjson.Field
		ExtraFields map[string]respjson.Field
		Raw         string
	} `json:"-"`
}

// RawJSON returns the unmodified JSON received from the API.
func (r OrderTrade) RawJSON() string { return r.JSON.Raw }
func (r *OrderTrade) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Order direction (ask=sell, bid=buy)
type SideEnum string

const (
	SideEnumAsk SideEnum = "ask"
	SideEnumBid SideEnum = "bid"
)

// Self-Match Prevention (SMP) mode.
//
// - `cancel_maker`: Cancel maker order
// - `cancel_taker`: Cancel taker order
// - `reduce`: Reduce order quantity
type SmpTypeEnum string

const (
	SmpTypeEnumCancelMaker SmpTypeEnum = "cancel_maker"
	SmpTypeEnumCancelTaker SmpTypeEnum = "cancel_taker"
	SmpTypeEnumReduce      SmpTypeEnum = "reduce"
)

// Time in force condition.
//
// - `fok`: Fill or Kill
// - `ioc`: Immediate or Cancel
// - `post_only`: Post only (maker only)
type TimeInForceEnum string

const (
	TimeInForceEnumFok      TimeInForceEnum = "fok"
	TimeInForceEnumIoc      TimeInForceEnum = "ioc"
	TimeInForceEnumPostOnly TimeInForceEnum = "post_only"
)

type OrderCancelAndNewResponse struct {
	// Unique identifier of the newly created order.
	NewOrderUuid string `json:"new_order_uuid" api:"required"`
	// Client-assigned identifier of the newly created order.
	NewOrderIdentifier string `json:"new_order_identifier"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		NewOrderUuid       respjson.Field
		NewOrderIdentifier respjson.Field
		ExtraFields        map[string]respjson.Field
		Raw                string
	} `json:"-"`
	Order
}

// RawJSON returns the unmodified JSON received from the API.
func (r OrderCancelAndNewResponse) RawJSON() string { return r.JSON.Raw }
func (r *OrderCancelAndNewResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type OrderGetChanceResponse struct {
	// Base asset account information.
	AskAccount OrderGetChanceResponseAskAccount `json:"ask_account" api:"required"`
	// Fee rate applied when placing a sell order.
	AskFee string `json:"ask_fee" api:"required"`
	// Quote asset account information.
	BidAccount OrderGetChanceResponseBidAccount `json:"bid_account" api:"required"`
	// Fee rate applied when placing a buy order.
	BidFee string `json:"bid_fee" api:"required"`
	// Maker fee rate for sell orders.
	MakerAskFee string `json:"maker_ask_fee" api:"required"`
	// Maker fee rate for buy orders.
	MakerBidFee string                       `json:"maker_bid_fee" api:"required"`
	Market      OrderGetChanceResponseMarket `json:"market" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		AskAccount  respjson.Field
		AskFee      respjson.Field
		BidAccount  respjson.Field
		BidFee      respjson.Field
		MakerAskFee respjson.Field
		MakerBidFee respjson.Field
		Market      respjson.Field
		ExtraFields map[string]respjson.Field
		Raw         string
	} `json:"-"`
}

// RawJSON returns the unmodified JSON received from the API.
func (r OrderGetChanceResponse) RawJSON() string { return r.JSON.Raw }
func (r *OrderGetChanceResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Base asset account information.
type OrderGetChanceResponseAskAccount struct {
	// Average buy price of the asset.
	AvgBuyPrice string `json:"avg_buy_price" api:"required"`
	// Indicates whether the average buy price has been modified.
	AvgBuyPriceModified bool `json:"avg_buy_price_modified" api:"required"`
	// Available amount or volume for orders.
	Balance string `json:"balance" api:"required"`
	// Currency code to be queried.
	Currency string `json:"currency" api:"required"`
	// Amount or quantity locked by pending orders or withdrawals.
	Locked string `json:"locked" api:"required"`
	// Currency unit used as the basis for avg_buy_price.
	UnitCurrency string `json:"unit_currency"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		AvgBuyPrice         respjson.Field
		AvgBuyPriceModified respjson.Field
		Balance             respjson.Field
		Currency            respjson.Field
		Locked              respjson.Field
		UnitCurrency        respjson.Field
		ExtraFields         map[string]respjson.Field
		Raw                 string
	} `json:"-"`
}

// RawJSON returns the unmodified JSON received from the API.
func (r OrderGetChanceResponseAskAccount) RawJSON() string { return r.JSON.Raw }
func (r *OrderGetChanceResponseAskAccount) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Quote asset account information.
type OrderGetChanceResponseBidAccount struct {
	// Average buy price of the asset.
	AvgBuyPrice string `json:"avg_buy_price" api:"required"`
	// Indicates whether the average buy price has been modified.
	AvgBuyPriceModified bool `json:"avg_buy_price_modified" api:"required"`
	// Available amount or volume for orders.
	Balance string `json:"balance" api:"required"`
	// Currency code to be queried.
	Currency string `json:"currency" api:"required"`
	// Amount or quantity locked by pending orders or withdrawals.
	Locked string `json:"locked" api:"required"`
	// Currency unit used as the basis for avg_buy_price.
	UnitCurrency string `json:"unit_currency"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		AvgBuyPrice         respjson.Field
		AvgBuyPriceModified respjson.Field
		Balance             respjson.Field
		Currency            respjson.Field
		Locked              respjson.Field
		UnitCurrency        respjson.Field
		ExtraFields         map[string]respjson.Field
		Raw                 string
	} `json:"-"`
}

// RawJSON returns the unmodified JSON received from the API.
func (r OrderGetChanceResponseBidAccount) RawJSON() string { return r.JSON.Raw }
func (r *OrderGetChanceResponseBidAccount) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type OrderGetChanceResponseMarket struct {
	// Trading pair code.
	ID string `json:"id" api:"required"`
	// Sell order constraints.
	Ask OrderGetChanceResponseMarketAsk `json:"ask" api:"required"`
	// Supported sell order types.
	//
	// [Example] ["market", "limit", "best_fok", ...]
	AskTypes []string `json:"ask_types" api:"required"`
	// Buy order constraints.
	Bid OrderGetChanceResponseMarketBid `json:"bid" api:"required"`
	// Supported buy order types.
	//
	// [Example] ["limit", "price", "best_fok", ...]
	BidTypes []string `json:"bid_types" api:"required"`
	// Maximum order amount.
	MaxTotal string `json:"max_total" api:"required"`
	// Trading pair name (base asset / quote currency).
	Name string `json:"name" api:"required"`
	// Supported order sides.
	OrderSides []string `json:"order_sides" api:"required"`
	// Market operation state.
	//
	// Any of "active".
	State string `json:"state" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		Ask         respjson.Field
		AskTypes    respjson.Field
		Bid         respjson.Field
		BidTypes    respjson.Field
		MaxTotal    respjson.Field
		Name        respjson.Field
		OrderSides  respjson.Field
		State       respjson.Field
		ExtraFields map[string]respjson.Field
		Raw         string
	} `json:"-"`
}

// RawJSON returns the unmodified JSON received from the API.
func (r OrderGetChanceResponseMarket) RawJSON() string { return r.JSON.Raw }
func (r *OrderGetChanceResponseMarket) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Sell order constraints.
type OrderGetChanceResponseMarketAsk struct {
	// Currency of the asset being sold.
	//
	// [Example] "BTC", "ETH"
	Currency string `json:"currency" api:"required"`
	// Minimum order amount for sell orders (in quote currency).
	MinTotal string `json:"min_total" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Currency    respjson.Field
		MinTotal    respjson.Field
		ExtraFields map[string]respjson.Field
		Raw         string
	} `json:"-"`
}

// RawJSON returns the unmodified JSON received from the API.
func (r OrderGetChanceResponseMarketAsk) RawJSON() string { return r.JSON.Raw }
func (r *OrderGetChanceResponseMarketAsk) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Buy order constraints.
type OrderGetChanceResponseMarketBid struct {
	// Quote currency used to purchase digital assets (SGD, BTC, USDT).
	Currency string `json:"currency" api:"required"`
	// Minimum order amount for buy orders (in quote currency).
	MinTotal string `json:"min_total" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Currency    respjson.Field
		MinTotal    respjson.Field
		ExtraFields map[string]respjson.Field
		Raw         string
	} `json:"-"`
}

// RawJSON returns the unmodified JSON received from the API.
func (r OrderGetChanceResponseMarketBid) RawJSON() string { return r.JSON.Raw }
func (r *OrderGetChanceResponseMarketBid) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type OrderNewParams struct {
	// Target trading pair for the order. (required)
	Market string `json:"market" api:"required"`
	// Order type.
	//
	// - `limit`: Limit order
	// - `price`: Market buy order (by total price)
	// - `market`: Market sell order (by volume)
	// - `best`: Best available price order
	//
	// Any of "limit", "price", "market", "best".
	OrdType OrdTypeEnum `json:"ord_type,omitzero" api:"required"`
	// Order direction (ask=sell, bid=buy)
	//
	// Any of "ask", "bid".
	Side SideEnum `json:"side,omitzero" api:"required"`
	// Client-assigned order identifier.
	Identifier param.Opt[string] `json:"identifier,omitzero"`
	// Order price or total amount.
	Price param.Opt[string] `json:"price,omitzero"`
	// Order volume.
	Volume param.Opt[string] `json:"volume,omitzero"`
	// Self-Match Prevention (SMP) mode.
	//
	// - `cancel_maker`: Cancel maker order
	// - `cancel_taker`: Cancel taker order
	// - `reduce`: Reduce order quantity
	//
	// Any of "cancel_maker", "cancel_taker", "reduce".
	SmpType SmpTypeEnum `json:"smp_type,omitzero"`
	// Time in force condition.
	//
	// - `fok`: Fill or Kill
	// - `ioc`: Immediate or Cancel
	// - `post_only`: Post only (maker only)
	//
	// Any of "fok", "ioc", "post_only".
	TimeInForce TimeInForceEnum `json:"time_in_force,omitzero"`
	paramObj
}

func (r OrderNewParams) MarshalJSON() (data []byte, err error) {
	type shadow OrderNewParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *OrderNewParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type OrderGetParams struct {
	// Client-assigned identifier of the order to query.
	//
	// Used when querying by the identifier assigned at order creation.
	Identifier param.Opt[string] `query:"identifier,omitzero" json:"-"`
	// UUID of the order to query.
	Uuid param.Opt[string] `query:"uuid,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [OrderGetParams]'s query parameters as `url.Values`.
func (r OrderGetParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type OrderCancelParams struct {
	// Client-assigned identifier of the order to cancel.
	//
	// Used when cancelling by the identifier assigned at order creation.
	Identifier param.Opt[string] `query:"identifier,omitzero" json:"-"`
	// UUID of the order to cancel.
	Uuid param.Opt[string] `query:"uuid,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [OrderCancelParams]'s query parameters as `url.Values`.
func (r OrderCancelParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type OrderCancelAndNewParams struct {
	// Order type.
	//
	// - `limit`: Limit order
	// - `price`: Market buy order (by total price)
	// - `market`: Market sell order (by volume)
	// - `best`: Best available price order
	//
	// Any of "limit", "price", "market", "best".
	NewOrdType OrdTypeEnum `json:"new_ord_type,omitzero" api:"required"`
	// Client-assigned identifier for the new order.
	NewIdentifier param.Opt[string] `json:"new_identifier,omitzero"`
	// Price or total amount for the new order.
	NewPrice param.Opt[string] `json:"new_price,omitzero"`
	// Volume for the new order.
	NewVolume param.Opt[string] `json:"new_volume,omitzero"`
	// Client-assigned identifier of the order to cancel.
	PrevOrderIdentifier param.Opt[string] `json:"prev_order_identifier,omitzero"`
	// UUID of the order to cancel.
	PrevOrderUuid param.Opt[string] `json:"prev_order_uuid,omitzero"`
	// Self-Match Prevention (SMP) mode.
	//
	// - `cancel_maker`: Cancel maker order
	// - `cancel_taker`: Cancel taker order
	// - `reduce`: Reduce order quantity
	//
	// Any of "cancel_maker", "cancel_taker", "reduce".
	NewSmpType SmpTypeEnum `json:"new_smp_type,omitzero"`
	// Time in force condition.
	//
	// - `fok`: Fill or Kill
	// - `ioc`: Immediate or Cancel
	// - `post_only`: Post only (maker only)
	//
	// Any of "fok", "ioc", "post_only".
	NewTimeInForce TimeInForceEnum `json:"new_time_in_force,omitzero"`
	paramObj
}

func (r OrderCancelAndNewParams) MarshalJSON() (data []byte, err error) {
	type shadow OrderCancelAndNewParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *OrderCancelAndNewParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type OrderCancelByUuidsParams struct {
	// List of client-assigned identifiers of orders to cancel. Maximum 20 orders.
	//
	// [Example] identifiers[]=id1&identifiers[]=id2…
	Identifiers []string `query:"identifiers,omitzero" json:"-"`
	// List of UUIDs of orders to cancel. Maximum 20 orders.
	//
	// [Example] uuids[]=uuid1&uuids[]=uuid2…
	Uuids []string `query:"uuids,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [OrderCancelByUuidsParams]'s query parameters as
// `url.Values`.
func (r OrderCancelByUuidsParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type OrderCancelOpenParams struct {
	// Maximum number of orders to cancel. Max 300, default 20.
	Count param.Opt[int64] `query:"count,omitzero" json:"-"`
	// Trading pair exclusion filter. Cancels all open orders except those for the
	// specified pairs. Up to 20 pairs, comma-separated.
	//
	// [Example] excluded_pairs=KRW-BTC,KRW-ETH
	ExcludedPairs param.Opt[string] `query:"excluded_pairs,omitzero" json:"-"`
	// Trading pair filter. Cancels open orders only for the specified pairs. Up to 20
	// pairs, comma-separated.
	//
	// [Example] pairs=KRW-BTC,KRW-ETH
	Pairs param.Opt[string] `query:"pairs,omitzero" json:"-"`
	// Quote currency filter (KRW, BTC, USDT). Cancels all open orders in markets with
	// the specified quote currency.
	//
	// [Example] "KRW" cancels all open orders in the KRW market.
	QuoteCurrencies param.Opt[string] `query:"quote_currencies,omitzero" json:"-"`
	// Side filter. "all" (both), "ask" (sell only), "bid" (buy only).
	//
	// Any of "bid", "ask", "all".
	CancelSide OrderCancelOpenParamsCancelSide `query:"cancel_side,omitzero" json:"-"`
	// Sort order for determining which orders to cancel. "desc" (newest first) or
	// "asc" (oldest first). Default is "desc".
	//
	// Any of "asc", "desc".
	OrderBy OrderCancelOpenParamsOrderBy `query:"order_by,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [OrderCancelOpenParams]'s query parameters as `url.Values`.
func (r OrderCancelOpenParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

// Side filter. "all" (both), "ask" (sell only), "bid" (buy only).
type OrderCancelOpenParamsCancelSide string

const (
	OrderCancelOpenParamsCancelSideBid OrderCancelOpenParamsCancelSide = "bid"
	OrderCancelOpenParamsCancelSideAsk OrderCancelOpenParamsCancelSide = "ask"
	OrderCancelOpenParamsCancelSideAll OrderCancelOpenParamsCancelSide = "all"
)

// Sort order for determining which orders to cancel. "desc" (newest first) or
// "asc" (oldest first). Default is "desc".
type OrderCancelOpenParamsOrderBy string

const (
	OrderCancelOpenParamsOrderByAsc  OrderCancelOpenParamsOrderBy = "asc"
	OrderCancelOpenParamsOrderByDesc OrderCancelOpenParamsOrderBy = "desc"
)

type OrderListByUuidsParams struct {
	// Trading pair to filter by.
	Market param.Opt[string] `query:"market,omitzero" json:"-"`
	// List of client-assigned identifiers of orders to query. Maximum 100 orders.
	//
	// [Example] identifiers[]=id1&identifiers[]=id2…
	Identifiers []string `query:"identifiers,omitzero" json:"-"`
	// Sort order. "desc" (newest first) or "asc" (oldest first). Default is "desc".
	//
	// Any of "asc", "desc".
	OrderBy OrderListByUuidsParamsOrderBy `query:"order_by,omitzero" json:"-"`
	// List of UUIDs of orders to query. Maximum 100 orders.
	//
	// [Example] uuids[]=uuid1&uuids[]=uuid2…
	Uuids []string `query:"uuids,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [OrderListByUuidsParams]'s query parameters as `url.Values`.
func (r OrderListByUuidsParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

// Sort order. "desc" (newest first) or "asc" (oldest first). Default is "desc".
type OrderListByUuidsParamsOrderBy string

const (
	OrderListByUuidsParamsOrderByAsc  OrderListByUuidsParamsOrderBy = "asc"
	OrderListByUuidsParamsOrderByDesc OrderListByUuidsParamsOrderBy = "desc"
)

type OrderListClosedParams struct {
	// End time of the query range. Max range is 7 days.
	//
	// - If only "end_time" is given, the range is 7 days before it.
	//
	// Format: ISO 8601 (2025-06-24T13:56:53+09:00) or millisecond timestamp
	// (1750741013000)
	EndTime param.Opt[string] `query:"end_time,omitzero" json:"-"`
	// Number of results per request (default: 100, max: 1,000).
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// Trading pair to filter by.
	Market param.Opt[string] `query:"market,omitzero" json:"-"`
	// Start time of the query range. Max range is 7 days.
	//
	// - If only "start_time" is given, the range is 7 days after it.
	// - If neither is given, the default is the past 7 days.
	//
	// Format: ISO 8601 (2025-06-24T13:56:53+09:00) or millisecond timestamp
	// (1750741013000)
	StartTime param.Opt[string] `query:"start_time,omitzero" json:"-"`
	// Sort order. "desc" (newest first) or "asc" (oldest first). Default is "desc".
	//
	// Any of "asc", "desc".
	OrderBy OrderListClosedParamsOrderBy `query:"order_by,omitzero" json:"-"`
	// Closed order state.
	//
	// - `done`: Fully executed
	// - `cancel`: Cancelled
	//
	// Any of "done", "cancel".
	State ClosedOrderState `query:"state,omitzero" json:"-"`
	// Order state filter (array form). "done" or "cancel". Default returns all states.
	//
	// [Example] states[]=done&states[]=cancel
	States []ClosedOrderState `query:"states,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [OrderListClosedParams]'s query parameters as `url.Values`.
func (r OrderListClosedParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

// Sort order. "desc" (newest first) or "asc" (oldest first). Default is "desc".
type OrderListClosedParamsOrderBy string

const (
	OrderListClosedParamsOrderByAsc  OrderListClosedParamsOrderBy = "asc"
	OrderListClosedParamsOrderByDesc OrderListClosedParamsOrderBy = "desc"
)

type OrderListOpenParams struct {
	// Number of results per request (default: 100, max: 100).
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// Trading pair to filter by.
	Market param.Opt[string] `query:"market,omitzero" json:"-"`
	// Page number for pagination. Default is 1.
	Page param.Opt[int64] `query:"page,omitzero" json:"-"`
	// Sort order. "desc" (newest first) or "asc" (oldest first). Default is "desc".
	//
	// Any of "asc", "desc".
	OrderBy OrderListOpenParamsOrderBy `query:"order_by,omitzero" json:"-"`
	// Open order state.
	//
	// - `wait`: Pending execution
	// - `watch`: Pending reservation (stop order)
	//
	// Any of "wait", "watch".
	State OpenOrderState `query:"state,omitzero" json:"-"`
	// Order state filter (array form). "wait" or "watch". Default is ["wait"].
	//
	// [Example] states[]=wait&states[]=watch
	States []OpenOrderState `query:"states,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [OrderListOpenParams]'s query parameters as `url.Values`.
func (r OrderListOpenParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

// Sort order. "desc" (newest first) or "asc" (oldest first). Default is "desc".
type OrderListOpenParamsOrderBy string

const (
	OrderListOpenParamsOrderByAsc  OrderListOpenParamsOrderBy = "asc"
	OrderListOpenParamsOrderByDesc OrderListOpenParamsOrderBy = "desc"
)

type OrderGetChanceParams struct {
	// Trading pair to query.
	Market string `query:"market" api:"required" json:"-"`
	paramObj
}

// URLQuery serializes [OrderGetChanceParams]'s query parameters as `url.Values`.
func (r OrderGetChanceParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type OrderTestNewParams struct {
	// Target trading pair for the order. (required)
	Market string `json:"market" api:"required"`
	// Order type.
	//
	// - `limit`: Limit order
	// - `price`: Market buy order (by total price)
	// - `market`: Market sell order (by volume)
	// - `best`: Best available price order
	//
	// Any of "limit", "price", "market", "best".
	OrdType OrdTypeEnum `json:"ord_type,omitzero" api:"required"`
	// Order direction (ask=sell, bid=buy)
	//
	// Any of "ask", "bid".
	Side SideEnum `json:"side,omitzero" api:"required"`
	// Client-assigned order identifier.
	Identifier param.Opt[string] `json:"identifier,omitzero"`
	// Order price or total amount.
	Price param.Opt[string] `json:"price,omitzero"`
	// Order volume.
	Volume param.Opt[string] `json:"volume,omitzero"`
	// Self-Match Prevention (SMP) mode.
	//
	// - `cancel_maker`: Cancel maker order
	// - `cancel_taker`: Cancel taker order
	// - `reduce`: Reduce order quantity
	//
	// Any of "cancel_maker", "cancel_taker", "reduce".
	SmpType SmpTypeEnum `json:"smp_type,omitzero"`
	// Time in force condition.
	//
	// - `fok`: Fill or Kill
	// - `ioc`: Immediate or Cancel
	// - `post_only`: Post only (maker only)
	//
	// Any of "fok", "ioc", "post_only".
	TimeInForce TimeInForceEnum `json:"time_in_force,omitzero"`
	paramObj
}

func (r OrderTestNewParams) MarshalJSON() (data []byte, err error) {
	type shadow OrderTestNewParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *OrderTestNewParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}
