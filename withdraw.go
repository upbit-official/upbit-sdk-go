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

// WithdrawService contains methods and other services that help with interacting
// with the upbit API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewWithdrawService] method instead.
type WithdrawService struct {
	options []option.RequestOption
}

// NewWithdrawService generates a new service that applies the given options to
// each request. These options are applied after the parent client's options (if
// there is one), and before any request-specific options.
func NewWithdrawService(opts ...option.RequestOption) (r WithdrawService) {
	r = WithdrawService{}
	r.options = opts
	return
}

// Get Withdrawal
//
// GET /v1/withdraw
//
// - docs(kr): https://docs.upbit.com/kr/reference/get-withdrawal
// - docs(sg/th): https://global-docs.upbit.com/reference/get-withdrawal
// - docs(id): https://docs.upbit.co.id/reference/get-withdrawal
//
// Retrieves the latest single withdrawal information. To query a specific
// withdrawal, you can specify the withdrawal’s UUID or transaction ID (TXID).
func (r *WithdrawService) Get(ctx context.Context, query WithdrawGetParams, opts ...option.RequestOption) (res *Withdraw, err error) {
	opts = slices.Concat(r.options, opts)
	path := "v1/withdraw"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, query, &res, opts...)
	return res, err
}

// List Withdrawals
//
// GET /v1/withdraws
//
// - docs(kr): https://docs.upbit.com/kr/reference/list-withdrawals
// - docs(sg/th): https://global-docs.upbit.com/reference/list-withdrawals
// - docs(id): https://docs.upbit.co.id/reference/list-withdrawals
//
// Retrieves latest withdrawal list.
func (r *WithdrawService) List(ctx context.Context, query WithdrawListParams, opts ...option.RequestOption) (res *pagination.UuidCursorIDPage[Withdraw], err error) {
	var raw *http.Response
	opts = slices.Concat(r.options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	path := "v1/withdraws"
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

// List Withdrawals
//
// GET /v1/withdraws
//
// - docs(kr): https://docs.upbit.com/kr/reference/list-withdrawals
// - docs(sg/th): https://global-docs.upbit.com/reference/list-withdrawals
// - docs(id): https://docs.upbit.co.id/reference/list-withdrawals
//
// Retrieves latest withdrawal list.
func (r *WithdrawService) ListAutoPaging(ctx context.Context, query WithdrawListParams, opts ...option.RequestOption) *pagination.UuidCursorIDPageAutoPager[Withdraw] {
	return pagination.NewUuidCursorIDPageAutoPager(r.List(ctx, query, opts...))
}

// Cancel Withdrawal
//
// DELETE /v1/withdraws/coin
//
// - docs(kr): https://docs.upbit.com/kr/reference/cancel-withdrawal
// - docs(sg/th): https://global-docs.upbit.com/reference/cancel-withdrawal
// - docs(id): https://docs.upbit.co.id/reference/cancel-withdrawal
//
// Cancel a completed digital asset withdrawal request by specifying the withdrawal
// UUID.
func (r *WithdrawService) CancelWithdrawal(ctx context.Context, body WithdrawCancelWithdrawalParams, opts ...option.RequestOption) (res *Withdraw, err error) {
	opts = slices.Concat(r.options, opts)
	path := "v1/withdraws/coin"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, body, &res, opts...)
	return res, err
}

// Withdraw KRW
//
// POST /v1/withdraws/krw
//
// - docs(kr): https://docs.upbit.com/kr/reference/withdraw-krw
// - Supported only in Korea.
//
// 원화(KRW) 출금을 요청합니다.
func (r *WithdrawService) NewKrwWithdrawal(ctx context.Context, body WithdrawNewKrwWithdrawalParams, opts ...option.RequestOption) (res *WithdrawNewKrwWithdrawalResponse, err error) {
	opts = slices.Concat(r.options, opts)
	path := "v1/withdraws/krw"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// Withdraw Digital Asset
//
// POST /v1/withdraws/coin
//
// - docs(kr): https://docs.upbit.com/kr/reference/withdraw
// - docs(sg/th): https://global-docs.upbit.com/reference/withdraw
// - docs(id): https://docs.upbit.co.id/reference/withdraw
//
// Requests a withdrawal of digital asset.
func (r *WithdrawService) NewWithdrawal(ctx context.Context, body WithdrawNewWithdrawalParams, opts ...option.RequestOption) (res *Withdraw, err error) {
	opts = slices.Concat(r.options, opts)
	path := "v1/withdraws/coin"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// List Withdrawal Allowed Addresses
//
// GET /v1/withdraws/coin_addresses
//
// - docs(kr): https://docs.upbit.com/kr/reference/list-withdrawal-addresses
// - docs(sg/th): https://global-docs.upbit.com/reference/list-withdrawal-addresses
// - docs(id): https://docs.upbit.co.id/reference/list-withdrawal-addresses
//
// Retrieves the list of withdrawal whitelist addresses registered to the account.
func (r *WithdrawService) ListCoinAddresses(ctx context.Context, opts ...option.RequestOption) (res *[]WithdrawListCoinAddressesResponse, err error) {
	opts = slices.Concat(r.options, opts)
	path := "v1/withdraws/coin_addresses"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// Get Available Withdrawal Information
//
// GET /v1/withdraws/chance
//
//   - docs(kr): https://docs.upbit.com/kr/reference/available-withdrawal-information
//   - docs(sg/th):
//     https://global-docs.upbit.com/reference/available-withdrawal-information
//   - docs(id): https://docs.upbit.co.id/reference/available-withdrawal-information
//
// Retrieves withdrawal availability information for the specified currency,
// including its withdrawal policy and user balance.
func (r *WithdrawService) GetChance(ctx context.Context, query WithdrawGetChanceParams, opts ...option.RequestOption) (res *WithdrawGetChanceResponse, err error) {
	opts = slices.Concat(r.options, opts)
	path := "v1/withdraws/chance"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, query, &res, opts...)
	return res, err
}

type Withdraw struct {
	// Amount of the asset to withdraw.
	Amount string `json:"amount" api:"required"`
	// Timestamp when the withdrawal was created.
	CreatedAt string `json:"created_at" api:"required"`
	// Currency code to be queried.
	Currency string `json:"currency" api:"required"`
	// Withdrawal fee.
	Fee string `json:"fee" api:"required"`
	// Whether the withdrawal can be cancelled.
	IsCancelable bool `json:"is_cancelable" api:"required"`
	// Withdrawal network type.
	//
	// Blockchain network identifier used by Upbit.
	//
	// [Example] "ETH", "TRX", "SOL"
	NetType string `json:"net_type" api:"required"`
	// Withdrawal processing state.
	//
	// - `WAITING`: Waiting
	// - `PROCESSING`: Processing
	// - `DONE`: Completed
	// - `FAILED`: Failed
	// - `CANCELLED`: Cancelled
	// - `REJECTED`: Rejected
	//
	// Any of "WAITING", "PROCESSING", "DONE", "FAILED", "CANCELLED", "REJECTED".
	State WithdrawState `json:"state" api:"required"`
	// Withdrawal transaction type.
	//
	// - `default`: Standard withdrawal
	// - `internal`: Internal (instant) withdrawal
	//
	// Any of "default", "internal".
	TransactionType WithdrawTransactionType `json:"transaction_type" api:"required"`
	// Transaction ID.
	Txid string `json:"txid" api:"required"`
	// Transaction type.
	Type string `json:"type" api:"required"`
	// Unique identifier of the withdrawal.
	Uuid string `json:"uuid" api:"required"`
	// Timestamp when the withdrawal was completed. Returns null if not yet completed.
	DoneAt string `json:"done_at" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Amount          respjson.Field
		CreatedAt       respjson.Field
		Currency        respjson.Field
		Fee             respjson.Field
		IsCancelable    respjson.Field
		NetType         respjson.Field
		State           respjson.Field
		TransactionType respjson.Field
		Txid            respjson.Field
		Type            respjson.Field
		Uuid            respjson.Field
		DoneAt          respjson.Field
		ExtraFields     map[string]respjson.Field
		raw             string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r Withdraw) RawJSON() string { return r.JSON.raw }
func (r *Withdraw) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Withdrawal processing state.
//
// - `WAITING`: Waiting
// - `PROCESSING`: Processing
// - `DONE`: Completed
// - `FAILED`: Failed
// - `CANCELLED`: Cancelled
// - `REJECTED`: Rejected
type WithdrawState string

const (
	WithdrawStateWaiting    WithdrawState = "WAITING"
	WithdrawStateProcessing WithdrawState = "PROCESSING"
	WithdrawStateDone       WithdrawState = "DONE"
	WithdrawStateFailed     WithdrawState = "FAILED"
	WithdrawStateCancelled  WithdrawState = "CANCELLED"
	WithdrawStateRejected   WithdrawState = "REJECTED"
)

// Withdrawal transaction type.
//
// - `default`: Standard withdrawal
// - `internal`: Internal (instant) withdrawal
type WithdrawTransactionType string

const (
	WithdrawTransactionTypeDefault  WithdrawTransactionType = "default"
	WithdrawTransactionTypeInternal WithdrawTransactionType = "internal"
)

type WithdrawNewKrwWithdrawalResponse struct {
	// Amount of the asset to withdraw.
	Amount string `json:"amount" api:"required"`
	// Timestamp when the withdrawal was created.
	CreatedAt string `json:"created_at" api:"required"`
	// Currency code to be queried.
	Currency string `json:"currency" api:"required"`
	// Withdrawal fee.
	Fee string `json:"fee" api:"required"`
	// Whether the withdrawal can be cancelled.
	IsCancelable bool `json:"is_cancelable" api:"required"`
	// Withdrawal network type.
	//
	// Blockchain network identifier used by Upbit.
	//
	// [Example] "ETH", "TRX", "SOL"
	NetType string `json:"net_type" api:"required"`
	// Withdrawal processing state.
	//
	// - `waiting`: Waiting
	// - `processing`: Processing
	// - `done`: Completed
	// - `failed`: Failed
	// - `cancelled`: Cancelled
	// - `rejected`: Rejected
	//
	// Any of "waiting", "processing", "done", "failed", "cancelled", "rejected".
	State WithdrawNewKrwWithdrawalResponseState `json:"state" api:"required"`
	// Withdrawal transaction type.
	//
	// - `default`: Standard withdrawal
	// - `internal`: Internal (instant) withdrawal
	//
	// Any of "default", "internal".
	TransactionType WithdrawNewKrwWithdrawalResponseTransactionType `json:"transaction_type" api:"required"`
	// Transaction ID.
	Txid string `json:"txid" api:"required"`
	// Transaction type.
	Type string `json:"type" api:"required"`
	// Unique identifier of the withdrawal.
	Uuid string `json:"uuid" api:"required"`
	// Timestamp when the withdrawal was completed. Returns null if not yet completed.
	DoneAt string `json:"done_at" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Amount          respjson.Field
		CreatedAt       respjson.Field
		Currency        respjson.Field
		Fee             respjson.Field
		IsCancelable    respjson.Field
		NetType         respjson.Field
		State           respjson.Field
		TransactionType respjson.Field
		Txid            respjson.Field
		Type            respjson.Field
		Uuid            respjson.Field
		DoneAt          respjson.Field
		ExtraFields     map[string]respjson.Field
		raw             string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r WithdrawNewKrwWithdrawalResponse) RawJSON() string { return r.JSON.raw }
func (r *WithdrawNewKrwWithdrawalResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Withdrawal processing state.
//
// - `waiting`: Waiting
// - `processing`: Processing
// - `done`: Completed
// - `failed`: Failed
// - `cancelled`: Cancelled
// - `rejected`: Rejected
type WithdrawNewKrwWithdrawalResponseState string

const (
	WithdrawNewKrwWithdrawalResponseStateWaiting    WithdrawNewKrwWithdrawalResponseState = "waiting"
	WithdrawNewKrwWithdrawalResponseStateProcessing WithdrawNewKrwWithdrawalResponseState = "processing"
	WithdrawNewKrwWithdrawalResponseStateDone       WithdrawNewKrwWithdrawalResponseState = "done"
	WithdrawNewKrwWithdrawalResponseStateFailed     WithdrawNewKrwWithdrawalResponseState = "failed"
	WithdrawNewKrwWithdrawalResponseStateCancelled  WithdrawNewKrwWithdrawalResponseState = "cancelled"
	WithdrawNewKrwWithdrawalResponseStateRejected   WithdrawNewKrwWithdrawalResponseState = "rejected"
)

// Withdrawal transaction type.
//
// - `default`: Standard withdrawal
// - `internal`: Internal (instant) withdrawal
type WithdrawNewKrwWithdrawalResponseTransactionType string

const (
	WithdrawNewKrwWithdrawalResponseTransactionTypeDefault  WithdrawNewKrwWithdrawalResponseTransactionType = "default"
	WithdrawNewKrwWithdrawalResponseTransactionTypeInternal WithdrawNewKrwWithdrawalResponseTransactionType = "internal"
)

type WithdrawListCoinAddressesResponse struct {
	// Currency code of the digital asset to withdraw.
	Currency string `json:"currency" api:"required"`
	// Withdrawal network type.
	//
	// Blockchain network identifier used by Upbit. The `net_type` parameter in
	// withdrawal requests must match this value.
	//
	// [Example] "ETH", "TRX", "SOL"
	NetType string `json:"net_type" api:"required"`
	// Withdrawal network name displayed to users.
	//
	// [Example] "Ethereum", "Bitcoin", "Tron", "Solana"
	NetworkName string `json:"network_name" api:"required"`
	// Recipient address for digital asset withdrawal.
	//
	// Only addresses registered in the allowed withdrawal address list can be used.
	WithdrawAddress string `json:"withdraw_address" api:"required"`
	// Company name of the entity receiving the withdrawn assets.
	BeneficiaryCompanyName string `json:"beneficiary_company_name" api:"nullable"`
	// Name of the wallet owner.
	//
	// Name of the individual or corporate representative receiving the withdrawn
	// assets. Returned only for exchange wallets; null for personal wallets.
	BeneficiaryName string `json:"beneficiary_name"`
	// Account holder type.
	//
	// - `individual`: Personal wallet
	// - `corporate`: Corporate wallet
	BeneficiaryType string `json:"beneficiary_type" api:"nullable"`
	// Name of the exchange where the withdrawal address is registered.
	//
	// [Example] "Binance", "Bybit"
	ExchangeName string `json:"exchange_name" api:"nullable"`
	// Secondary withdrawal address (Destination Tag, Memo, or Message). Some digital
	// assets require a secondary address along with the primary address.
	SecondaryAddress string `json:"secondary_address" api:"nullable"`
	// Personal wallet type.
	//
	// [Example] "MetaMask"
	WalletType string `json:"wallet_type" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Currency               respjson.Field
		NetType                respjson.Field
		NetworkName            respjson.Field
		WithdrawAddress        respjson.Field
		BeneficiaryCompanyName respjson.Field
		BeneficiaryName        respjson.Field
		BeneficiaryType        respjson.Field
		ExchangeName           respjson.Field
		SecondaryAddress       respjson.Field
		WalletType             respjson.Field
		ExtraFields            map[string]respjson.Field
		raw                    string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r WithdrawListCoinAddressesResponse) RawJSON() string { return r.JSON.raw }
func (r *WithdrawListCoinAddressesResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type WithdrawGetChanceResponse struct {
	// Asset balance information.
	Account WithdrawGetChanceResponseAccount `json:"account" api:"required"`
	// Currency information.
	Currency WithdrawGetChanceResponseCurrency `json:"currency" api:"required"`
	// User security level information.
	MemberLevel WithdrawGetChanceResponseMemberLevel `json:"member_level" api:"required"`
	// Withdrawal limit information.
	WithdrawLimit WithdrawGetChanceResponseWithdrawLimit `json:"withdraw_limit" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Account       respjson.Field
		Currency      respjson.Field
		MemberLevel   respjson.Field
		WithdrawLimit respjson.Field
		ExtraFields   map[string]respjson.Field
		raw           string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r WithdrawGetChanceResponse) RawJSON() string { return r.JSON.raw }
func (r *WithdrawGetChanceResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Asset balance information.
type WithdrawGetChanceResponseAccount struct {
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
func (r WithdrawGetChanceResponseAccount) RawJSON() string { return r.JSON.raw }
func (r *WithdrawGetChanceResponseAccount) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Currency information.
type WithdrawGetChanceResponseCurrency struct {
	// Currency code.
	Code string `json:"code" api:"required"`
	// Whether this is a digital asset.
	IsCoin bool `json:"is_coin" api:"required"`
	// Whether deposit/withdrawal has been supported for this asset.
	//
	// - `working`: Has been supported
	// - `unsupported`: Never supported
	//
	// Any of "working", "unsupported".
	WalletState string `json:"wallet_state" api:"required"`
	// Current deposit/withdrawal availability. 'deposit' if deposits are available,
	// 'withdraw' if withdrawals are available. Empty list means neither is available.
	WalletSupport []string `json:"wallet_support" api:"required"`
	// Withdrawal fee.
	WithdrawFee string `json:"withdraw_fee" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Code          respjson.Field
		IsCoin        respjson.Field
		WalletState   respjson.Field
		WalletSupport respjson.Field
		WithdrawFee   respjson.Field
		ExtraFields   map[string]respjson.Field
		raw           string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r WithdrawGetChanceResponseCurrency) RawJSON() string { return r.JSON.raw }
func (r *WithdrawGetChanceResponseCurrency) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// User security level information.
type WithdrawGetChanceResponseMemberLevel struct {
	// Whether bank account is verified.
	BankAccountVerified bool `json:"bank_account_verified" api:"required"`
	// Whether email is verified.
	EmailVerified bool `json:"email_verified" api:"required"`
	// Fee level.
	FeeLevel int64 `json:"fee_level" api:"required"`
	// Whether identity is verified.
	IdentityAuthVerified bool `json:"identity_auth_verified" api:"required"`
	// Whether the account is locked.
	Locked bool `json:"locked" api:"required"`
	// Security level.
	SecurityLevel int64 `json:"security_level" api:"required"`
	// Whether two-factor authentication is enabled.
	TwoFactorAuthVerified bool `json:"two_factor_auth_verified" api:"required"`
	// Whether the wallet (withdrawal) is locked.
	WalletLocked bool `json:"wallet_locked" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		BankAccountVerified   respjson.Field
		EmailVerified         respjson.Field
		FeeLevel              respjson.Field
		IdentityAuthVerified  respjson.Field
		Locked                respjson.Field
		SecurityLevel         respjson.Field
		TwoFactorAuthVerified respjson.Field
		WalletLocked          respjson.Field
		ExtraFields           map[string]respjson.Field
		raw                   string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r WithdrawGetChanceResponseMemberLevel) RawJSON() string { return r.JSON.raw }
func (r *WithdrawGetChanceResponseMemberLevel) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Withdrawal limit information.
type WithdrawGetChanceResponseWithdrawLimit struct {
	// Currency code to be queried.
	Currency string `json:"currency" api:"required"`
	// Reference fiat currency.
	FiatCurrency string `json:"fiat_currency" api:"required"`
	// Integrated remaining daily withdrawal limit (in fiat currency).
	RemainingDailyFiat string `json:"remaining_daily_fiat" api:"required"`
	// Whether withdrawal is currently supported.
	//
	//   - To check current availability, refer to currency.wallet_support for the
	//     'withdraw' value.
	CanWithdraw bool `json:"can_withdraw"`
	// Decimal precision for withdrawal amount/volume.
	Fixed int64 `json:"fixed"`
	// Minimum withdrawal amount or volume.
	Minimum string `json:"minimum"`
	// Amount restricted due to the withdrawal delay policy.
	WithdrawDelayedFiat string `json:"withdraw_delayed_fiat"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Currency            respjson.Field
		FiatCurrency        respjson.Field
		RemainingDailyFiat  respjson.Field
		CanWithdraw         respjson.Field
		Fixed               respjson.Field
		Minimum             respjson.Field
		WithdrawDelayedFiat respjson.Field
		ExtraFields         map[string]respjson.Field
		raw                 string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r WithdrawGetChanceResponseWithdrawLimit) RawJSON() string { return r.JSON.raw }
func (r *WithdrawGetChanceResponseWithdrawLimit) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type WithdrawGetParams struct {
	// Currency code filter. If not specified, the latest withdrawal is returned.
	Currency param.Opt[string] `query:"currency,omitzero" json:"-"`
	// Transaction ID of the withdrawal to query.
	//
	// If neither uuid nor txid is provided, the latest withdrawal is returned.
	Txid param.Opt[string] `query:"txid,omitzero" json:"-"`
	// UUID of the withdrawal to query.
	//
	// If neither uuid nor txid is provided, the latest withdrawal is returned.
	Uuid param.Opt[string] `query:"uuid,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [WithdrawGetParams]'s query parameters as `url.Values`.
func (r WithdrawGetParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type WithdrawListParams struct {
	// Currency code filter. If not specified, all currencies are returned.
	Currency param.Opt[string] `query:"currency,omitzero" json:"-"`
	// Cursor for pagination. Enter a "uuid" from the response to retrieve "limit"
	// withdrawals after that timestamp.
	From param.Opt[string] `query:"from,omitzero" json:"-"`
	// Number of results per request (default: 100, max: 100).
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// Page number for pagination. Default is 1.
	Page param.Opt[int64] `query:"page,omitzero" json:"-"`
	// Cursor for pagination. Enter a "uuid" from the response to retrieve "limit"
	// withdrawals before that timestamp.
	To param.Opt[string] `query:"to,omitzero" json:"-"`
	// Sort order. "desc" (newest first) or "asc" (oldest first). Default is "desc".
	//
	// Any of "asc", "desc".
	OrderBy WithdrawListParamsOrderBy `query:"order_by,omitzero" json:"-"`
	// Withdrawal state filter.
	//
	// - `WAITING`: Waiting
	// - `PROCESSING`: Processing
	// - `DONE`: Completed
	// - `FAILED`: Failed
	// - `CANCELLED`: Cancelled
	// - `REJECTED`: Rejected
	//
	// Any of "WAITING", "PROCESSING", "DONE", "FAILED", "CANCELLED", "REJECTED".
	State WithdrawListParamsState `query:"state,omitzero" json:"-"`
	// List of transaction IDs to query. Maximum 100. Cannot be used together with
	// uuids.
	//
	// [Example] txids[]=txid1&txids[]=txid2
	Txids []string `query:"txids,omitzero" json:"-"`
	// List of UUIDs to query. Maximum 100. Cannot be used together with txids.
	//
	// [Example] uuids[]=uuid1&uuids[]=uuid2
	Uuids []string `query:"uuids,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [WithdrawListParams]'s query parameters as `url.Values`.
func (r WithdrawListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

// Sort order. "desc" (newest first) or "asc" (oldest first). Default is "desc".
type WithdrawListParamsOrderBy string

const (
	WithdrawListParamsOrderByAsc  WithdrawListParamsOrderBy = "asc"
	WithdrawListParamsOrderByDesc WithdrawListParamsOrderBy = "desc"
)

// Withdrawal state filter.
//
// - `WAITING`: Waiting
// - `PROCESSING`: Processing
// - `DONE`: Completed
// - `FAILED`: Failed
// - `CANCELLED`: Cancelled
// - `REJECTED`: Rejected
type WithdrawListParamsState string

const (
	WithdrawListParamsStateWaiting    WithdrawListParamsState = "WAITING"
	WithdrawListParamsStateProcessing WithdrawListParamsState = "PROCESSING"
	WithdrawListParamsStateDone       WithdrawListParamsState = "DONE"
	WithdrawListParamsStateFailed     WithdrawListParamsState = "FAILED"
	WithdrawListParamsStateCancelled  WithdrawListParamsState = "CANCELLED"
	WithdrawListParamsStateRejected   WithdrawListParamsState = "REJECTED"
)

type WithdrawCancelWithdrawalParams struct {
	// UUID of the withdrawal to cancel.
	Uuid string `query:"uuid" api:"required" json:"-"`
	paramObj
}

// URLQuery serializes [WithdrawCancelWithdrawalParams]'s query parameters as
// `url.Values`.
func (r WithdrawCancelWithdrawalParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type WithdrawNewKrwWithdrawalParams struct {
	// KRW amount to withdraw.
	//
	// Must be a numeric string.
	Amount string `json:"amount" api:"required"`
	// Two-factor authentication method for KRW transactions.
	//
	// - `kakao`: Kakao authentication
	// - `naver`: Naver authentication
	// - `hana`: Hana certificate authentication
	//
	// Any of "kakao", "naver", "hana".
	TwoFactorType WithdrawNewKrwWithdrawalParamsTwoFactorType `json:"two_factor_type,omitzero" api:"required"`
	paramObj
}

func (r WithdrawNewKrwWithdrawalParams) MarshalJSON() (data []byte, err error) {
	type shadow WithdrawNewKrwWithdrawalParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *WithdrawNewKrwWithdrawalParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Two-factor authentication method for KRW transactions.
//
// - `kakao`: Kakao authentication
// - `naver`: Naver authentication
// - `hana`: Hana certificate authentication
type WithdrawNewKrwWithdrawalParamsTwoFactorType string

const (
	WithdrawNewKrwWithdrawalParamsTwoFactorTypeKakao WithdrawNewKrwWithdrawalParamsTwoFactorType = "kakao"
	WithdrawNewKrwWithdrawalParamsTwoFactorTypeNaver WithdrawNewKrwWithdrawalParamsTwoFactorType = "naver"
	WithdrawNewKrwWithdrawalParamsTwoFactorTypeHana  WithdrawNewKrwWithdrawalParamsTwoFactorType = "hana"
)

type WithdrawNewWithdrawalParams struct {
	// Recipient address for digital asset withdrawal.
	//
	// Only addresses registered in the allowed withdrawal address list can be used.
	Address string `json:"address" api:"required"`
	// Amount of the asset to withdraw.
	//
	// Must be a numeric string.
	Amount string `json:"amount" api:"required"`
	// Currency code of the digital asset to withdraw.
	Currency string `json:"currency" api:"required"`
	// Check the allowed withdrawal addresses API response to find the available
	// "net_type" value for each address.
	NetType string `json:"net_type" api:"required"`
	// Secondary withdrawal address (Destination Tag, Memo, or Message). If the
	// recipient address includes a secondary address, this field must be included.
	SecondaryAddress param.Opt[string] `json:"secondary_address,omitzero"`
	// Withdrawal transaction type.
	//
	// - `default`: Standard withdrawal
	// - `internal`: Internal (instant) withdrawal
	//
	// Any of "default", "internal".
	TransactionType WithdrawNewWithdrawalParamsTransactionType `json:"transaction_type,omitzero"`
	paramObj
}

func (r WithdrawNewWithdrawalParams) MarshalJSON() (data []byte, err error) {
	type shadow WithdrawNewWithdrawalParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *WithdrawNewWithdrawalParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Withdrawal transaction type.
//
// - `default`: Standard withdrawal
// - `internal`: Internal (instant) withdrawal
type WithdrawNewWithdrawalParamsTransactionType string

const (
	WithdrawNewWithdrawalParamsTransactionTypeDefault  WithdrawNewWithdrawalParamsTransactionType = "default"
	WithdrawNewWithdrawalParamsTransactionTypeInternal WithdrawNewWithdrawalParamsTransactionType = "internal"
)

type WithdrawGetChanceParams struct {
	// Currency code to query withdrawal availability for.
	Currency string `query:"currency" api:"required" json:"-"`
	// Blockchain network identifier for digital asset withdrawal. Required for digital
	// assets.
	NetType string `query:"net_type" api:"required" json:"-"`
	paramObj
}

// URLQuery serializes [WithdrawGetChanceParams]'s query parameters as
// `url.Values`.
func (r WithdrawGetChanceParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}
