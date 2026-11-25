package upbit

import (
	"context"
	"encoding/json"
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

// DepositService contains methods and other services that help with interacting
// with the upbit API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewDepositService] method instead.
type DepositService struct {
	options []option.RequestOption
}

// NewDepositService generates a new service that applies the given options to each
// request. These options are applied after the parent client's options (if there
// is one), and before any request-specific options.
func NewDepositService(opts ...option.RequestOption) (r DepositService) {
	r = DepositService{}
	r.options = opts
	return
}

// Get Deposit
//
// GET /v1/deposit
//
// - docs(kr): https://docs.upbit.com/kr/reference/get-deposit
// - docs(sg/th): https://global-docs.upbit.com/reference/get-deposit
// - docs(id): https://docs.upbit.co.id/reference/get-deposit
//
// Retrieves the latest deposit history. To query specific deposit information, you
// can specify the deposit’s UUID, transaction ID (TXID), or currency code.
func (r *DepositService) Get(ctx context.Context, query DepositGetParams, opts ...option.RequestOption) (res *Deposit, err error) {
	opts = slices.Concat(r.options, opts)
	path := "v1/deposit"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, query, &res, opts...)
	return res, err
}

// List Deposits
//
// GET /v1/deposits
//
// - docs(kr): https://docs.upbit.com/kr/reference/list-deposits
// - docs(sg/th): https://global-docs.upbit.com/reference/list-deposits
// - docs(id): https://docs.upbit.co.id/reference/list-deposits
//
// Retrieves the latest deposit list.
func (r *DepositService) List(ctx context.Context, query DepositListParams, opts ...option.RequestOption) (res *pagination.UuidCursorIDPage[Deposit], err error) {
	var raw *http.Response
	opts = slices.Concat(r.options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	path := "v1/deposits"
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

// List Deposits
//
// GET /v1/deposits
//
// - docs(kr): https://docs.upbit.com/kr/reference/list-deposits
// - docs(sg/th): https://global-docs.upbit.com/reference/list-deposits
// - docs(id): https://docs.upbit.co.id/reference/list-deposits
//
// Retrieves the latest deposit list.
func (r *DepositService) ListAutoPaging(ctx context.Context, query DepositListParams, opts ...option.RequestOption) *pagination.UuidCursorIDPageAutoPager[Deposit] {
	return pagination.NewUuidCursorIDPageAutoPager(r.List(ctx, query, opts...))
}

// Create Deposit Address
//
// POST /v1/deposits/generate_coin_address
//
// - docs(kr): https://docs.upbit.com/kr/reference/create-deposit-address
// - docs(sg/th): https://global-docs.upbit.com/reference/create-deposit-address
// - docs(id): https://docs.upbit.co.id/reference/create-deposit-address
//
// Requests the creation of a deposit address for depositing assets from a personal
// wallet or another exchange into Upbit.
func (r *DepositService) NewCoinAddress(ctx context.Context, body DepositNewCoinAddressParams, opts ...option.RequestOption) (res *DepositNewCoinAddressResponseUnion, err error) {
	opts = slices.Concat(r.options, opts)
	path := "v1/deposits/generate_coin_address"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// Deposit KRW
//
// POST /v1/deposits/krw
//
// - docs(kr): https://docs.upbit.com/kr/reference/deposit-krw
// - Supported only in Korea.
//
// 입출금 계좌로부터 원화를 지정한 금액만큼 입금합니다.
func (r *DepositService) DepositKrw(ctx context.Context, body DepositDepositKrwParams, opts ...option.RequestOption) (res *Deposit, err error) {
	opts = slices.Concat(r.options, opts)
	path := "v1/deposits/krw"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// List Deposit Addresses
//
// GET /v1/deposits/coin_addresses
//
// - docs(kr): https://docs.upbit.com/kr/reference/list-deposit-addresses
// - docs(sg/th): https://global-docs.upbit.com/reference/list-deposit-addresses
// - docs(id): https://docs.upbit.co.id/reference/list-deposit-addresses
//
// Retrieves deposit address information for the specified currency. Return all
// successfully generated deposit addresses.
func (r *DepositService) ListCoinAddresses(ctx context.Context, opts ...option.RequestOption) (res *[]DepositAddress, err error) {
	opts = slices.Concat(r.options, opts)
	path := "v1/deposits/coin_addresses"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// Get Available Deposit Information
//
// GET /v1/deposits/chance/coin
//
//   - docs(kr): https://docs.upbit.com/kr/reference/available-deposit-information
//   - docs(sg/th):
//     https://global-docs.upbit.com/reference/available-deposit-information
//   - docs(id): https://docs.upbit.co.id/reference/available-deposit-information
//
// Deposits availability information for the specified currency is retrieved.
func (r *DepositService) GetChance(ctx context.Context, query DepositGetChanceParams, opts ...option.RequestOption) (res *DepositGetChanceResponse, err error) {
	opts = slices.Concat(r.options, opts)
	path := "v1/deposits/chance/coin"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, query, &res, opts...)
	return res, err
}

// Get Deposit Address
//
// GET /v1/deposits/coin_address
//
// - docs(kr): https://docs.upbit.com/kr/reference/get-deposit-address
// - docs(sg/th): https://global-docs.upbit.com/reference/get-deposit-address
// - docs(id): https://docs.upbit.co.id/reference/get-deposit-address
//
// Retrieves deposit address information for the specified currency.
func (r *DepositService) GetCoinAddress(ctx context.Context, query DepositGetCoinAddressParams, opts ...option.RequestOption) (res *DepositAddress, err error) {
	opts = slices.Concat(r.options, opts)
	path := "v1/deposits/coin_address"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, query, &res, opts...)
	return res, err
}

type Deposit struct {
	// Deposit amount.
	Amount string `json:"amount" api:"required"`
	// Timestamp when the deposit was requested.
	//
	// [Format] yyyy-MM-ddTHH:mm:ss+00:00
	CreatedAt string `json:"created_at" api:"required"`
	// Currency code of the deposit.
	Currency string `json:"currency" api:"required"`
	// Timestamp when the deposit was completed.
	//
	// [Format] yyyy-MM-ddTHH:mm:ss+00:00
	DoneAt string `json:"done_at" api:"required"`
	// Deposit fee.
	Fee string `json:"fee" api:"required"`
	// Deposit status.
	//
	//   - `PROCESSING`: Deposit in progress (digital asset only)
	//   - `ACCEPTED`: Completed
	//   - `CANCELLED`: Cancelled
	//   - `REJECTED`: Rejected
	//   - `TRAVEL_RULE_SUSPECTED`: Awaiting additional Travel Rule verification (digital
	//     assets only)
	//   - `REFUNDING`: Refund in progress
	//   - `REFUNDED`: Refunded (digital asset only)
	//
	// Any of "PROCESSING", "ACCEPTED", "CANCELLED", "REJECTED",
	// "TRAVEL_RULE_SUSPECTED", "REFUNDING", "REFUNDED".
	State DepositState `json:"state" api:"required"`
	// Deposit transaction type.
	//
	// - `default`: Standard deposit
	// - `internal`: Internal deposit (between Upbit accounts)
	//
	// Any of "default", "internal".
	TransactionType DepositTransactionType `json:"transaction_type" api:"required"`
	// Deposit transaction ID.
	Txid string `json:"txid" api:"required"`
	// Transaction type.
	Type string `json:"type" api:"required"`
	// Unique identifier (UUID) of the deposit.
	Uuid string `json:"uuid" api:"required"`
	// Deposit network type. Returns null for fiat (KRW) deposits.
	//
	// [Example] ETH, TRX, SOL
	NetType string `json:"net_type" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Amount          respjson.Field
		CreatedAt       respjson.Field
		Currency        respjson.Field
		DoneAt          respjson.Field
		Fee             respjson.Field
		State           respjson.Field
		TransactionType respjson.Field
		Txid            respjson.Field
		Type            respjson.Field
		Uuid            respjson.Field
		NetType         respjson.Field
		ExtraFields     map[string]respjson.Field
		raw             string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r Deposit) RawJSON() string { return r.JSON.raw }
func (r *Deposit) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Deposit status.
//
//   - `PROCESSING`: Deposit in progress (digital asset only)
//   - `ACCEPTED`: Completed
//   - `CANCELLED`: Cancelled
//   - `REJECTED`: Rejected
//   - `TRAVEL_RULE_SUSPECTED`: Awaiting additional Travel Rule verification (digital
//     assets only)
//   - `REFUNDING`: Refund in progress
//   - `REFUNDED`: Refunded (digital asset only)
type DepositState string

const (
	DepositStateProcessing          DepositState = "PROCESSING"
	DepositStateAccepted            DepositState = "ACCEPTED"
	DepositStateCancelled           DepositState = "CANCELLED"
	DepositStateRejected            DepositState = "REJECTED"
	DepositStateTravelRuleSuspected DepositState = "TRAVEL_RULE_SUSPECTED"
	DepositStateRefunding           DepositState = "REFUNDING"
	DepositStateRefunded            DepositState = "REFUNDED"
)

// Deposit transaction type.
//
// - `default`: Standard deposit
// - `internal`: Internal deposit (between Upbit accounts)
type DepositTransactionType string

const (
	DepositTransactionTypeDefault  DepositTransactionType = "default"
	DepositTransactionTypeInternal DepositTransactionType = "internal"
)

type DepositAddress struct {
	// Currency code for which the deposit address was created.
	Currency string `json:"currency" api:"required"`
	// Deposit address.
	DepositAddress string `json:"deposit_address" api:"required"`
	// Deposit network type.
	//
	// Blockchain network identifier used by Upbit.
	//
	// [Example] ETH, TRX, SOL
	NetType string `json:"net_type" api:"required"`
	// Secondary deposit address (Destination Tag, Memo, or Message).
	SecondaryAddress string `json:"secondary_address" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Currency         respjson.Field
		DepositAddress   respjson.Field
		NetType          respjson.Field
		SecondaryAddress respjson.Field
		ExtraFields      map[string]respjson.Field
		raw              string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r DepositAddress) RawJSON() string { return r.JSON.raw }
func (r *DepositAddress) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// DepositNewCoinAddressResponseUnion contains all possible properties and values
// from [DepositAddress],
// [DepositNewCoinAddressResponseGenerateCoinAddressResponse].
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type DepositNewCoinAddressResponseUnion struct {
	// This field is from variant [DepositAddress].
	Currency string `json:"currency"`
	// This field is from variant [DepositAddress].
	DepositAddress string `json:"deposit_address"`
	// This field is from variant [DepositAddress].
	NetType string `json:"net_type"`
	// This field is from variant [DepositAddress].
	SecondaryAddress string `json:"secondary_address"`
	// This field is from variant
	// [DepositNewCoinAddressResponseGenerateCoinAddressResponse].
	Message string `json:"message"`
	// This field is from variant
	// [DepositNewCoinAddressResponseGenerateCoinAddressResponse].
	Success bool `json:"success"`
	JSON    struct {
		Currency         respjson.Field
		DepositAddress   respjson.Field
		NetType          respjson.Field
		SecondaryAddress respjson.Field
		Message          respjson.Field
		Success          respjson.Field
		raw              string
	} `json:"-"`
}

func (u DepositNewCoinAddressResponseUnion) AsDepositAddress() (v DepositAddress) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u DepositNewCoinAddressResponseUnion) AsDepositNewCoinAddressResponseGenerateCoinAddressResponse() (v DepositNewCoinAddressResponseGenerateCoinAddressResponse) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u DepositNewCoinAddressResponseUnion) RawJSON() string { return u.JSON.raw }

func (r *DepositNewCoinAddressResponseUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type DepositNewCoinAddressResponseGenerateCoinAddressResponse struct {
	// Message about the deposit address creation request.
	Message string `json:"message" api:"required"`
	// Whether the deposit address creation request was accepted.
	Success bool `json:"success" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Message     respjson.Field
		Success     respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r DepositNewCoinAddressResponseGenerateCoinAddressResponse) RawJSON() string { return r.JSON.raw }
func (r *DepositNewCoinAddressResponseGenerateCoinAddressResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type DepositGetChanceResponse struct {
	// Currency code to be queried.
	Currency string `json:"currency" api:"required"`
	// Decimal precision applied to deposits.
	DecimalPrecision int64 `json:"decimal_precision" api:"required"`
	// Reason why deposits are not possible.
	//
	// ※ Provided when "is_deposit_possible" is "false".
	DepositImpossibleReason string `json:"deposit_impossible_reason" api:"required"`
	// Whether deposits are currently possible.
	IsDepositPossible bool `json:"is_deposit_possible" api:"required"`
	// Minimum deposit amount.
	MinimumDepositAmount string `json:"minimum_deposit_amount" api:"required"`
	// Minimum number of blockchain confirmations required for the deposit to be
	// credited.
	MinimumDepositConfirmations int64 `json:"minimum_deposit_confirmations" api:"required"`
	// Deposit network type.
	//
	// Blockchain network identifier used by Upbit.
	//
	// [Example] ETH, TRX, SOL
	NetType string `json:"net_type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Currency                    respjson.Field
		DecimalPrecision            respjson.Field
		DepositImpossibleReason     respjson.Field
		IsDepositPossible           respjson.Field
		MinimumDepositAmount        respjson.Field
		MinimumDepositConfirmations respjson.Field
		NetType                     respjson.Field
		ExtraFields                 map[string]respjson.Field
		raw                         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r DepositGetChanceResponse) RawJSON() string { return r.JSON.raw }
func (r *DepositGetChanceResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type DepositGetParams struct {
	// Currency code filter.
	Currency param.Opt[string] `query:"currency,omitzero" json:"-"`
	// Transaction ID of the deposit to query.
	//
	// If neither uuid nor txid is provided, the latest deposit is returned.
	Txid param.Opt[string] `query:"txid,omitzero" json:"-"`
	// UUID of the deposit to query.
	//
	// If neither uuid nor txid is provided, the latest deposit is returned.
	Uuid param.Opt[string] `query:"uuid,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [DepositGetParams]'s query parameters as `url.Values`.
func (r DepositGetParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type DepositListParams struct {
	// Currency code filter. If not specified, the latest deposits are returned.
	Currency param.Opt[string] `query:"currency,omitzero" json:"-"`
	// Cursor for pagination. Enter a "uuid" from the response to retrieve "limit"
	// deposits after that timestamp.
	From param.Opt[string] `query:"from,omitzero" json:"-"`
	// Number of results per request (default: 100, max: 100).
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// Page number for pagination. Default is 1.
	Page param.Opt[int64] `query:"page,omitzero" json:"-"`
	// Cursor for pagination. Enter a "uuid" from the response to retrieve "limit"
	// deposits before that timestamp.
	To param.Opt[string] `query:"to,omitzero" json:"-"`
	// Sort order. "desc" (newest first) or "asc" (oldest first). Default is "desc".
	//
	// Any of "asc", "desc".
	OrderBy DepositListParamsOrderBy `query:"order_by,omitzero" json:"-"`
	// Deposit state filter.
	//
	// - `PROCESSING`: Processing
	// - `ACCEPTED`: Completed
	// - `CANCELLED`: Cancelled
	// - `REJECTED`: Rejected
	// - `TRAVEL_RULE_SUSPECTED`: Pending Travel Rule verification
	// - `REFUNDING`: Refund in progress
	// - `REFUNDED`: Refunded
	//
	// Any of "PROCESSING", "ACCEPTED", "CANCELLED", "REJECTED",
	// "TRAVEL_RULE_SUSPECTED", "REFUNDING", "REFUNDED".
	State DepositListParamsState `query:"state,omitzero" json:"-"`
	// List of transaction IDs to query.
	//
	// [Example] txids[]=txid1&txids[]=txid2
	Txids []string `query:"txids,omitzero" json:"-"`
	// List of UUIDs to query.
	//
	// [Example] uuids[]=uuid1&uuids[]=uuid2
	Uuids []string `query:"uuids,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [DepositListParams]'s query parameters as `url.Values`.
func (r DepositListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

// Sort order. "desc" (newest first) or "asc" (oldest first). Default is "desc".
type DepositListParamsOrderBy string

const (
	DepositListParamsOrderByAsc  DepositListParamsOrderBy = "asc"
	DepositListParamsOrderByDesc DepositListParamsOrderBy = "desc"
)

// Deposit state filter.
//
// - `PROCESSING`: Processing
// - `ACCEPTED`: Completed
// - `CANCELLED`: Cancelled
// - `REJECTED`: Rejected
// - `TRAVEL_RULE_SUSPECTED`: Pending Travel Rule verification
// - `REFUNDING`: Refund in progress
// - `REFUNDED`: Refunded
type DepositListParamsState string

const (
	DepositListParamsStateProcessing          DepositListParamsState = "PROCESSING"
	DepositListParamsStateAccepted            DepositListParamsState = "ACCEPTED"
	DepositListParamsStateCancelled           DepositListParamsState = "CANCELLED"
	DepositListParamsStateRejected            DepositListParamsState = "REJECTED"
	DepositListParamsStateTravelRuleSuspected DepositListParamsState = "TRAVEL_RULE_SUSPECTED"
	DepositListParamsStateRefunding           DepositListParamsState = "REFUNDING"
	DepositListParamsStateRefunded            DepositListParamsState = "REFUNDED"
)

type DepositNewCoinAddressParams struct {
	// Currency code for which to create a deposit address.
	Currency string `json:"currency" api:"required"`
	// Network type.
	NetType string `json:"net_type" api:"required"`
	paramObj
}

func (r DepositNewCoinAddressParams) MarshalJSON() (data []byte, err error) {
	type shadow DepositNewCoinAddressParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *DepositNewCoinAddressParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type DepositDepositKrwParams struct {
	// KRW amount to deposit.
	Amount string `json:"amount" api:"required"`
	// Two-factor authentication method for KRW transactions.
	//
	// - `kakao`: Kakao authentication
	// - `naver`: Naver authentication
	// - `hana`: Hana certificate authentication
	//
	// Any of "kakao", "naver", "hana".
	TwoFactorType DepositDepositKrwParamsTwoFactorType `json:"two_factor_type,omitzero" api:"required"`
	paramObj
}

func (r DepositDepositKrwParams) MarshalJSON() (data []byte, err error) {
	type shadow DepositDepositKrwParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *DepositDepositKrwParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Two-factor authentication method for KRW transactions.
//
// - `kakao`: Kakao authentication
// - `naver`: Naver authentication
// - `hana`: Hana certificate authentication
type DepositDepositKrwParamsTwoFactorType string

const (
	DepositDepositKrwParamsTwoFactorTypeKakao DepositDepositKrwParamsTwoFactorType = "kakao"
	DepositDepositKrwParamsTwoFactorTypeNaver DepositDepositKrwParamsTwoFactorType = "naver"
	DepositDepositKrwParamsTwoFactorTypeHana  DepositDepositKrwParamsTwoFactorType = "hana"
)

type DepositGetChanceParams struct {
	// Currency code to query deposit availability for.
	Currency string `query:"currency" api:"required" json:"-"`
	// Blockchain network identifier for digital asset deposit.
	//
	// Used to filter by network.
	NetType string `query:"net_type" api:"required" json:"-"`
	paramObj
}

// URLQuery serializes [DepositGetChanceParams]'s query parameters as `url.Values`.
func (r DepositGetChanceParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type DepositGetCoinAddressParams struct {
	// Currency code to query.
	Currency string `query:"currency" api:"required" json:"-"`
	// Blockchain network identifier.
	//
	// Used to filter by network.
	NetType string `query:"net_type" api:"required" json:"-"`
	paramObj
}

// URLQuery serializes [DepositGetCoinAddressParams]'s query parameters as
// `url.Values`.
func (r DepositGetCoinAddressParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}
