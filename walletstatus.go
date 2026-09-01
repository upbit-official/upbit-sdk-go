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

// WalletStatusService contains methods and other services that help with
// interacting with the upbit API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewWalletStatusService] method instead.
type WalletStatusService struct {
	options []option.RequestOption
}

// NewWalletStatusService generates a new service that applies the given options to
// each request. These options are applied after the parent client's options (if
// there is one), and before any request-specific options.
func NewWalletStatusService(opts ...option.RequestOption) (r WalletStatusService) {
	r = WalletStatusService{}
	r.options = opts
	return
}

// Get Deposit/Withdrawal Service Status
//
// GET /v1/status/wallet
//
// - docs(kr): https://docs.upbit.com/kr/reference/get-service-status
// - docs(global): https://global-docs.upbit.com/reference/get-service-status
//
// Retrieves deposit and withdrawal service status for all currencies.
func (r *WalletStatusService) List(ctx context.Context, opts ...option.RequestOption) (res *[]WalletStatusListResponse, err error) {
	opts = slices.Concat(r.options, opts)
	path := "v1/status/wallet"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

type WalletStatusListResponse struct {
	// Minutes elapsed since the last block update. May return null.
	BlockElapsedMinutes int64 `json:"block_elapsed_minutes" api:"required"`
	// Currency code to be queried.
	Currency string `json:"currency" api:"required"`
	// Deposit/withdrawal network type.
	//
	// Blockchain network identifier used by Upbit.
	//
	// [Example] "ETH", "TRX", "SOL"
	NetType string `json:"net_type" api:"required"`
	// Deposit/withdrawal network name displayed to users.
	//
	// [Example] "Ethereum", "Bitcoin", "Tron", "Solana"
	NetworkName string `json:"network_name" api:"required"`
	// Deposit/withdrawal status.
	//
	// - `working`: Both available
	// - `withdraw_only`: Withdrawal only
	// - `deposit_only`: Deposit only
	// - `paused`: Both suspended
	// - `unsupported`: Not supported
	//
	// Any of "working", "withdraw_only", "deposit_only", "paused", "unsupported".
	WalletState WalletStatusListResponseWalletState `json:"wallet_state" api:"required"`
	// Current confirmed block height. May return null.
	BlockHeight int64 `json:"block_height"`
	// Blockchain network state. May return null depending on wallet or exchange
	// status.
	//
	// - `normal`: Normal
	// - `delayed`: Delayed
	// - `inactive`: Inactive
	//
	// Any of "normal", "delayed", "inactive".
	BlockState WalletStatusListResponseBlockState `json:"block_state"`
	// Timestamp when the block height was last updated (UTC). May return null.
	//
	// [Format] yyyy-MM-dd'T'HH:mm:sss+00:00
	BlockUpdatedAt string `json:"block_updated_at"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		BlockElapsedMinutes respjson.Field
		Currency            respjson.Field
		NetType             respjson.Field
		NetworkName         respjson.Field
		WalletState         respjson.Field
		BlockHeight         respjson.Field
		BlockState          respjson.Field
		BlockUpdatedAt      respjson.Field
		ExtraFields         map[string]respjson.Field
		raw                 string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r WalletStatusListResponse) RawJSON() string { return r.JSON.raw }
func (r *WalletStatusListResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Deposit/withdrawal status.
//
// - `working`: Both available
// - `withdraw_only`: Withdrawal only
// - `deposit_only`: Deposit only
// - `paused`: Both suspended
// - `unsupported`: Not supported
type WalletStatusListResponseWalletState string

const (
	WalletStatusListResponseWalletStateWorking      WalletStatusListResponseWalletState = "working"
	WalletStatusListResponseWalletStateWithdrawOnly WalletStatusListResponseWalletState = "withdraw_only"
	WalletStatusListResponseWalletStateDepositOnly  WalletStatusListResponseWalletState = "deposit_only"
	WalletStatusListResponseWalletStatePaused       WalletStatusListResponseWalletState = "paused"
	WalletStatusListResponseWalletStateUnsupported  WalletStatusListResponseWalletState = "unsupported"
)

// Blockchain network state. May return null depending on wallet or exchange
// status.
//
// - `normal`: Normal
// - `delayed`: Delayed
// - `inactive`: Inactive
type WalletStatusListResponseBlockState string

const (
	WalletStatusListResponseBlockStateNormal   WalletStatusListResponseBlockState = "normal"
	WalletStatusListResponseBlockStateDelayed  WalletStatusListResponseBlockState = "delayed"
	WalletStatusListResponseBlockStateInactive WalletStatusListResponseBlockState = "inactive"
)
