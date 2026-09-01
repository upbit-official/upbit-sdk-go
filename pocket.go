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

// PocketService contains methods and other services that help with interacting
// with the upbit API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewPocketService] method instead.
type PocketService struct {
	options []option.RequestOption
}

// NewPocketService generates a new service that applies the given options to each
// request. These options are applied after the parent client's options (if there
// is one), and before any request-specific options.
func NewPocketService(opts ...option.RequestOption) (r PocketService) {
	r = PocketService{}
	r.options = opts
	return
}

// List Pockets
//
// GET /v1/pockets
//
// - docs(kr): https://docs.upbit.com/kr/reference/list-pockets
// - Supported only in Korea.
//
// Retrieve the UUIDs and names of the main pocket and user-created pockets.
func (r *PocketService) List(ctx context.Context, opts ...option.RequestOption) (res *[]Pocket, err error) {
	opts = slices.Concat(r.options, opts)
	path := "v1/pockets"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// List Pocket API Keys
//
// GET /v1/pockets/api_keys
//
// - docs(kr): https://docs.upbit.com/kr/reference/list-pocket-api-keys
// - Supported only in Korea.
//
// List API keys and related information for each created pocket using main-pocket
// API key permissions.
func (r *PocketService) ListAPIKeys(ctx context.Context, query PocketListAPIKeysParams, opts ...option.RequestOption) (res *[]PocketAPIKeyGroup, err error) {
	opts = slices.Concat(r.options, opts)
	path := "v1/pockets/api_keys"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, query, &res, opts...)
	return res, err
}

// List Transfers
//
// GET /v1/pockets/transfers
//
// - docs(kr): https://docs.upbit.com/kr/reference/list-transfers
// - Supported only in Korea.
//
// List asset transfer history between pockets that occurred from the sub-pocket
// using sub-pocket API key permissions.
func (r *PocketService) ListTransfers(ctx context.Context, query PocketListTransfersParams, opts ...option.RequestOption) (res *[]PocketTransfer, err error) {
	opts = slices.Concat(r.options, opts)
	path := "v1/pockets/transfers"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, query, &res, opts...)
	return res, err
}

// List Universal Transfers
//
// GET /v1/pockets/universal_transfers
//
// - docs(kr): https://docs.upbit.com/kr/reference/list-universal-transfers
// - Supported only in Korea.
//
// List asset transfer history between pockets, such as main pocket to sub-pocket,
// using main-pocket API key permissions.
func (r *PocketService) ListUniversalTransfers(ctx context.Context, query PocketListUniversalTransfersParams, opts ...option.RequestOption) (res *[]PocketTransfer, err error) {
	opts = slices.Concat(r.options, opts)
	path := "v1/pockets/universal_transfers"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, query, &res, opts...)
	return res, err
}

// Get Sub-Pocket Balance
//
// GET /v1/pockets/assets
//
// - docs(kr): https://docs.upbit.com/kr/reference/get-sub-pocket-balance
// - Supported only in Korea.
//
// Retrieve balances for a specific sub-pocket using main-pocket API key
// permissions.
func (r *PocketService) GetBalance(ctx context.Context, query PocketGetBalanceParams, opts ...option.RequestOption) (res *[]PocketBalance, err error) {
	opts = slices.Concat(r.options, opts)
	path := "v1/pockets/assets"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, query, &res, opts...)
	return res, err
}

// Transfer
//
// POST /v1/pockets/transfers
//
// - docs(kr): https://docs.upbit.com/kr/reference/transfer
// - Supported only in Korea.
//
// Request an asset transfer between pockets, such as sub-pocket to sub-pocket,
// using sub-pocket API key permissions.
func (r *PocketService) Transfer(ctx context.Context, body PocketTransferParams, opts ...option.RequestOption) (res *PocketTransfer, err error) {
	opts = slices.Concat(r.options, opts)
	path := "v1/pockets/transfers"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// Universal Transfer
//
// POST /v1/pockets/universal_transfers
//
// - docs(kr): https://docs.upbit.com/kr/reference/universal-transfer
// - Supported only in Korea.
//
// Request an asset transfer between pockets, such as main pocket to sub-pocket,
// using main-pocket API key permissions.
func (r *PocketService) UniversalTransfer(ctx context.Context, body PocketUniversalTransferParams, opts ...option.RequestOption) (res *PocketTransfer, err error) {
	opts = slices.Concat(r.options, opts)
	path := "v1/pockets/universal_transfers"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

type Pocket struct {
	// Pocket name.
	Name string `json:"name" api:"required"`
	// Pocket type.
	//
	// - `main`: Main pocket
	// - `user_spot_trading`: User-created pocket
	//
	// Any of "main", "user_spot_trading".
	Type PocketType `json:"type" api:"required"`
	// Pocket UUID.
	Uuid string `json:"uuid" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Name        respjson.Field
		Type        respjson.Field
		Uuid        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r Pocket) RawJSON() string { return r.JSON.raw }
func (r *Pocket) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type PocketAPIKey struct {
	// Access key of the API key.
	AccessKey string `json:"access_key" api:"required"`
	// List of IP addresses allowed to make API requests.
	AllowedIPs []string `json:"allowed_ips" api:"required"`
	// Timestamp when the API key was created.
	//
	// [Format] yyyy-MM-dd'T'HH:mm:ss+09:00
	CreatedAt string `json:"created_at" api:"required"`
	// Timestamp when the API key expires.
	//
	// [Format] yyyy-MM-dd'T'HH:mm:ss+09:00
	ExpiredAt string `json:"expired_at" api:"required"`
	// List of access permissions granted to the API key.
	Permissions []PocketAPIKeyPermission `json:"permissions" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		AccessKey   respjson.Field
		AllowedIPs  respjson.Field
		CreatedAt   respjson.Field
		ExpiredAt   respjson.Field
		Permissions respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r PocketAPIKey) RawJSON() string { return r.JSON.raw }
func (r *PocketAPIKey) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type PocketAPIKeyGroup struct {
	// List of API keys issued for the pocket.
	Keys []PocketAPIKey `json:"keys" api:"required"`
	// Pocket UUID.
	Uuid string `json:"uuid" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Keys        respjson.Field
		Uuid        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r PocketAPIKeyGroup) RawJSON() string { return r.JSON.raw }
func (r *PocketAPIKeyGroup) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Access permission granted to the API key.
//
// - `view_account`: View assets/accounts
// - `view_orders`: View order history
// - `make_orders`: Create and execute orders
// - `view_withdrawals`: View withdrawal history
// - `withdraw`: Withdraw digital assets/KRW
// - `view_deposits`: View deposit history
// - `deposit`: Create deposit addresses and deposits
// - `transfer`: Transfer assets between pockets
// - `manage_pockets`: Create and manage pockets
type PocketAPIKeyPermission string

const (
	PocketAPIKeyPermissionViewAccount     PocketAPIKeyPermission = "view_account"
	PocketAPIKeyPermissionViewOrders      PocketAPIKeyPermission = "view_orders"
	PocketAPIKeyPermissionMakeOrders      PocketAPIKeyPermission = "make_orders"
	PocketAPIKeyPermissionViewWithdrawals PocketAPIKeyPermission = "view_withdrawals"
	PocketAPIKeyPermissionWithdraw        PocketAPIKeyPermission = "withdraw"
	PocketAPIKeyPermissionViewDeposits    PocketAPIKeyPermission = "view_deposits"
	PocketAPIKeyPermissionDeposit         PocketAPIKeyPermission = "deposit"
	PocketAPIKeyPermissionTransfer        PocketAPIKeyPermission = "transfer"
	PocketAPIKeyPermissionManagePockets   PocketAPIKeyPermission = "manage_pockets"
)

type PocketBalance struct {
	// Average buy price of the asset.
	AvgBuyPrice string `json:"avg_buy_price" api:"required"`
	// Indicates whether the average buy price has been modified.
	AvgBuyPriceModified bool `json:"avg_buy_price_modified" api:"required"`
	// Available quantity or amount for orders.
	//
	// - Digital assets: quantity
	// - Fiat currency (KRW): amount
	Balance string `json:"balance" api:"required"`
	// Currency code of the held asset.
	//
	// [Example] KRW, BTC, ETH
	Currency string `json:"currency" api:"required"`
	// Balance locked by withdrawals, pending order execution, or similar operations.
	Locked string `json:"locked" api:"required"`
	// Currency unit used as the basis for `avg_buy_price`.
	//
	// [Example] KRW, BTC, USDT
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
func (r PocketBalance) RawJSON() string { return r.JSON.raw }
func (r *PocketBalance) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type PocketTransfer struct {
	// Transferred asset quantity.
	Amount string `json:"amount" api:"required"`
	// Timestamp when the request was created.
	CreatedAt string `json:"created_at" api:"required"`
	// Transferred asset code.
	//
	// [Example] KRW, BTC, ETH
	Currency string `json:"currency" api:"required"`
	// UUID of the pocket that sent the asset.
	From string `json:"from" api:"required"`
	// Client-provided identifier for the asset transfer request. Returns null if not
	// specified.
	Identifier string `json:"identifier" api:"required"`
	// Asset transfer request state.
	//
	// - `submitted`: Submitted
	// - `processing`: Processing
	// - `done`: Completed
	// - `failed`: Failed
	//
	// Any of "submitted", "processing", "done", "failed".
	State PocketTransferState `json:"state" api:"required"`
	// UUID of the pocket that received the asset.
	To string `json:"to" api:"required"`
	// UUID identifying the asset transfer request.
	Uuid string `json:"uuid" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Amount      respjson.Field
		CreatedAt   respjson.Field
		Currency    respjson.Field
		From        respjson.Field
		Identifier  respjson.Field
		State       respjson.Field
		To          respjson.Field
		Uuid        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r PocketTransfer) RawJSON() string { return r.JSON.raw }
func (r *PocketTransfer) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Transfer direction.
//
// - `in`: Transfers into the sub-pocket
// - `out`: Transfers out of the sub-pocket
// - `all`: Transfers in all directions
type PocketTransferDirection string

const (
	PocketTransferDirectionIn  PocketTransferDirection = "in"
	PocketTransferDirectionOut PocketTransferDirection = "out"
	PocketTransferDirectionAll PocketTransferDirection = "all"
)

// Asset transfer request state.
//
// - `submitted`: Submitted
// - `processing`: Processing
// - `done`: Completed
// - `failed`: Failed
type PocketTransferState string

const (
	PocketTransferStateSubmitted  PocketTransferState = "submitted"
	PocketTransferStateProcessing PocketTransferState = "processing"
	PocketTransferStateDone       PocketTransferState = "done"
	PocketTransferStateFailed     PocketTransferState = "failed"
)

// Pocket type.
//
// - `main`: Main pocket
// - `user_spot_trading`: User-created pocket
type PocketType string

const (
	PocketTypeMain            PocketType = "main"
	PocketTypeUserSpotTrading PocketType = "user_spot_trading"
)

type PocketListAPIKeysParams struct {
	// Whether to include expired API keys in the results.
	IncludeExpired param.Opt[bool] `query:"include_expired,omitzero" json:"-"`
	// Array of pocket UUIDs to query.
	//
	// [Example] uuids[]=uuid1&uuids[]=uuid2
	Uuids []string `query:"uuids,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [PocketListAPIKeysParams]'s query parameters as
// `url.Values`.
func (r PocketListAPIKeysParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type PocketListTransfersParams struct {
	// Asset code to query.
	Currency param.Opt[string] `query:"currency,omitzero" json:"-"`
	// End time of the query period. Retrieves records created from `start_time` to the
	// specified time.
	//
	// Default: Current time. Limit: Maximum query range is 7 days. Allowed format: ISO
	// 8601 with timezone or millisecond timestamp.
	EndTime param.Opt[string] `query:"end_time,omitzero" json:"-"`
	// Number of results to request (default 20, max 100).
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// Start time of the query period. Retrieves records created from the specified
	// time to `end_time`.
	//
	// Allowed format: ISO 8601 with timezone or millisecond timestamp.
	StartTime param.Opt[string] `query:"start_time,omitzero" json:"-"`
	// Transfer direction.
	//
	// - `in`: Transfers into the sub-pocket
	// - `out`: Transfers out of the sub-pocket
	// - `all`: Transfers in all directions
	//
	// Any of "in", "out", "all".
	Direction PocketTransferDirection `query:"direction,omitzero" json:"-"`
	// List of client-provided identifiers specified when requesting transfers.
	//
	// [Example] identifiers[]=id1&identifiers[]=id2
	//
	// Limit: Up to 20 items.
	Identifiers []string `query:"identifiers,omitzero" json:"-"`
	// Sort order. `asc` for oldest first, `desc` for newest first.
	//
	// Any of "asc", "desc".
	OrderBy PocketListTransfersParamsOrderBy `query:"order_by,omitzero" json:"-"`
	// Asset transfer request states. Returns only transfers in the specified states.
	// If omitted, returns transfers in all states.
	//
	// [Example] states[]=done&states[]=failed
	States []PocketTransferState `query:"states,omitzero" json:"-"`
	// List of UUIDs identifying asset transfer requests. Returns only transfers
	// matching the specified UUIDs.
	//
	// [Example] uuids[]=uuid1&uuids[]=uuid2
	//
	// Limit: Up to 20 items.
	Uuids []string `query:"uuids,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [PocketListTransfersParams]'s query parameters as
// `url.Values`.
func (r PocketListTransfersParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

// Sort order. `asc` for oldest first, `desc` for newest first.
type PocketListTransfersParamsOrderBy string

const (
	PocketListTransfersParamsOrderByAsc  PocketListTransfersParamsOrderBy = "asc"
	PocketListTransfersParamsOrderByDesc PocketListTransfersParamsOrderBy = "desc"
)

type PocketListUniversalTransfersParams struct {
	// Asset code to query.
	Currency param.Opt[string] `query:"currency,omitzero" json:"-"`
	// End time of the query period. Retrieves records created from `start_time` to the
	// specified time.
	//
	// Default: Current time. Limit: Maximum query range is 7 days. Allowed format: ISO
	// 8601 with timezone or millisecond timestamp.
	EndTime param.Opt[string] `query:"end_time,omitzero" json:"-"`
	// UUID of the pocket that sent the asset.
	From param.Opt[string] `query:"from,omitzero" json:"-"`
	// Number of results to request (default 20, max 100).
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// Start time of the query period. Retrieves records created from the specified
	// time to `end_time`.
	//
	// Allowed format: ISO 8601 with timezone or millisecond timestamp.
	StartTime param.Opt[string] `query:"start_time,omitzero" json:"-"`
	// UUID of the pocket that received the asset.
	To param.Opt[string] `query:"to,omitzero" json:"-"`
	// List of client-provided identifiers specified when requesting transfers.
	//
	// [Example] identifiers[]=id1&identifiers[]=id2
	//
	// Limit: Up to 20 items.
	Identifiers []string `query:"identifiers,omitzero" json:"-"`
	// Sort order. `asc` for oldest first, `desc` for newest first.
	//
	// Any of "asc", "desc".
	OrderBy PocketListUniversalTransfersParamsOrderBy `query:"order_by,omitzero" json:"-"`
	// Asset transfer request states. Available values: `submitted`, `processing`,
	// `done`, `failed`. Returns only transfers in the specified states. If omitted,
	// returns transfers in all states.
	//
	// [Example] states[]=done&states[]=failed
	States []PocketTransferState `query:"states,omitzero" json:"-"`
	// List of UUIDs identifying asset transfer requests. Returns only transfers
	// matching the specified UUIDs.
	//
	// [Example] uuids[]=uuid1&uuids[]=uuid2
	//
	// Limit: Up to 20 items.
	Uuids []string `query:"uuids,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [PocketListUniversalTransfersParams]'s query parameters as
// `url.Values`.
func (r PocketListUniversalTransfersParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

// Sort order. `asc` for oldest first, `desc` for newest first.
type PocketListUniversalTransfersParamsOrderBy string

const (
	PocketListUniversalTransfersParamsOrderByAsc  PocketListUniversalTransfersParamsOrderBy = "asc"
	PocketListUniversalTransfersParamsOrderByDesc PocketListUniversalTransfersParamsOrderBy = "desc"
)

type PocketGetBalanceParams struct {
	// UUID of the sub-pocket to query.
	Uuid string `query:"uuid" api:"required" json:"-"`
	paramObj
}

// URLQuery serializes [PocketGetBalanceParams]'s query parameters as `url.Values`.
func (r PocketGetBalanceParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type PocketTransferParams struct {
	// Asset quantity to transfer.
	Amount string `json:"amount" api:"required"`
	// Asset code to transfer.
	//
	// [Example] KRW, BTC, ETH
	Currency string `json:"currency" api:"required"`
	// UUID of the receiving pocket.
	To string `json:"to" api:"required"`
	// Client-provided identifier for the asset transfer request. Returns null if not
	// specified.
	//
	// Allowed format: 1 to 64 letters, numbers, underscores (\_), periods (.), or
	// hyphens (-).
	Identifier param.Opt[string] `json:"identifier,omitzero"`
	paramObj
}

func (r PocketTransferParams) MarshalJSON() (data []byte, err error) {
	type shadow PocketTransferParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *PocketTransferParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type PocketUniversalTransferParams struct {
	// Asset quantity to transfer.
	Amount string `json:"amount" api:"required"`
	// Asset code to transfer.
	//
	// [Example] KRW, BTC, ETH
	Currency string `json:"currency" api:"required"`
	// UUID of the receiving pocket.
	To string `json:"to" api:"required"`
	// UUID of the sending pocket. If omitted, defaults to the pocket associated with
	// the API key used for the request.
	From param.Opt[string] `json:"from,omitzero"`
	// Client-provided identifier for the asset transfer request. Returns null if not
	// specified.
	//
	// Allowed format: 1 to 64 letters, numbers, underscores (\_), periods (.), or
	// hyphens (-).
	Identifier param.Opt[string] `json:"identifier,omitzero"`
	paramObj
}

func (r PocketUniversalTransferParams) MarshalJSON() (data []byte, err error) {
	type shadow PocketUniversalTransferParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *PocketUniversalTransferParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}
