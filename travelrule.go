package upbit

import (
	"context"
	"net/http"
	"slices"

	"github.com/upbit-official/upbit-sdk-go/internal/apijson"
	"github.com/upbit-official/upbit-sdk-go/internal/requestconfig"
	"github.com/upbit-official/upbit-sdk-go/option"
	"github.com/upbit-official/upbit-sdk-go/packages/param"
	"github.com/upbit-official/upbit-sdk-go/packages/respjson"
)

// TravelRuleService contains methods and other services that help with interacting
// with the upbit API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewTravelRuleService] method instead.
type TravelRuleService struct {
	options []option.RequestOption
}

// NewTravelRuleService generates a new service that applies the given options to
// each request. These options are applied after the parent client's options (if
// there is one), and before any request-specific options.
func NewTravelRuleService(opts ...option.RequestOption) (r TravelRuleService) {
	r = TravelRuleService{}
	r.options = opts
	return
}

// List Travel Rule Supporting VASPs
//
// GET /v1/travel_rule/vasps
//
// - docs(kr): https://docs.upbit.com/kr/reference/list-travelrule-vasps
// - docs(global): https://global-docs.upbit.com/reference/list-travelrule-vasps
func (r *TravelRuleService) ListVasps(ctx context.Context, opts ...option.RequestOption) (res *[]TravelRuleListVaspsResponse, err error) {
	opts = slices.Concat(r.options, opts)
	path := "v1/travel_rule/vasps"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// Verify Travel Rule by Deposit TxID
//
// POST /v1/travel_rule/deposit/txid
//
//   - docs(kr): https://docs.upbit.com/kr/reference/verify-travelrule-by-txid
//   - docs(global):
//     https://global-docs.upbit.com/reference/verify-travelrule-by-txid
func (r *TravelRuleService) VerifyDepositByTxid(ctx context.Context, body TravelRuleVerifyDepositByTxidParams, opts ...option.RequestOption) (res *VerificationResult, err error) {
	opts = slices.Concat(r.options, opts)
	path := "v1/travel_rule/deposit/txid"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// Verify Travel Rule by Deposit UUID
//
// POST /v1/travel_rule/deposit/uuid
//
//   - docs(kr): https://docs.upbit.com/kr/reference/verify-travelrule-by-uuid
//   - docs(global):
//     https://global-docs.upbit.com/reference/verify-travelrule-by-uuid
func (r *TravelRuleService) VerifyDepositByUuid(ctx context.Context, body TravelRuleVerifyDepositByUuidParams, opts ...option.RequestOption) (res *VerificationResult, err error) {
	opts = slices.Concat(r.options, opts)
	path := "v1/travel_rule/deposit/uuid"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

type VerificationResult struct {
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
	DepositState VerificationResultDepositState `json:"deposit_state" api:"required"`
	// Unique identifier (UUID) for the deposit to verify.
	DepositUuid string `json:"deposit_uuid" api:"required"`
	// Result of the account verification process.
	VerificationResult string `json:"verification_result" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		DepositState       respjson.Field
		DepositUuid        respjson.Field
		VerificationResult respjson.Field
		ExtraFields        map[string]respjson.Field
		Raw                string
	} `json:"-"`
}

// RawJSON returns the unmodified JSON received from the API.
func (r VerificationResult) RawJSON() string { return r.JSON.Raw }
func (r *VerificationResult) UnmarshalJSON(data []byte) error {
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
type VerificationResultDepositState string

const (
	VerificationResultDepositStateProcessing          VerificationResultDepositState = "PROCESSING"
	VerificationResultDepositStateAccepted            VerificationResultDepositState = "ACCEPTED"
	VerificationResultDepositStateCancelled           VerificationResultDepositState = "CANCELLED"
	VerificationResultDepositStateRejected            VerificationResultDepositState = "REJECTED"
	VerificationResultDepositStateTravelRuleSuspected VerificationResultDepositState = "TRAVEL_RULE_SUSPECTED"
	VerificationResultDepositStateRefunding           VerificationResultDepositState = "REFUNDING"
	VerificationResultDepositStateRefunded            VerificationResultDepositState = "REFUNDED"
)

type TravelRuleListVaspsResponse struct {
	// Indicates whether deposits are available from the VASP.
	Depositable bool `json:"depositable" api:"required"`
	// Exchange name. Name of the VASP (Virtual Asset Service Provider).
	VaspName string `json:"vasp_name" api:"required"`
	// Unique identifier (UUID) of the exchange. Uniquely assigned VASP identifier for
	// the transfer.
	VaspUuid string `json:"vasp_uuid" api:"required"`
	// Indicates whether withdrawals are available to the VASP.
	Withdrawable bool `json:"withdrawable" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Depositable  respjson.Field
		VaspName     respjson.Field
		VaspUuid     respjson.Field
		Withdrawable respjson.Field
		ExtraFields  map[string]respjson.Field
		Raw          string
	} `json:"-"`
}

// RawJSON returns the unmodified JSON received from the API.
func (r TravelRuleListVaspsResponse) RawJSON() string { return r.JSON.Raw }
func (r *TravelRuleListVaspsResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type TravelRuleVerifyDepositByTxidParams struct {
	// Currency code to be queried.
	Currency string `json:"currency" api:"required"`
	// Blockchain network identifier for deposit/withdrawal of digital assets. Used as
	// a filter parameter to specify the network.
	NetType string `json:"net_type" api:"required"`
	// Transaction ID of the deposit to be verified.
	Txid string `json:"txid" api:"required"`
	// Unique identifier (UUID) of the counterparty exchange from which the asset was
	// withdrawn.
	VaspUuid string `json:"vasp_uuid" api:"required"`
	paramObj
}

func (r TravelRuleVerifyDepositByTxidParams) MarshalJSON() (data []byte, err error) {
	type shadow TravelRuleVerifyDepositByTxidParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *TravelRuleVerifyDepositByTxidParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type TravelRuleVerifyDepositByUuidParams struct {
	// Unique identifier (UUID) for the deposit to verify.
	DepositUuid string `json:"deposit_uuid" api:"required"`
	// Unique identifier (UUID) of the counterparty exchange from which the asset was
	// withdrawn.
	VaspUuid string `json:"vasp_uuid" api:"required"`
	paramObj
}

func (r TravelRuleVerifyDepositByUuidParams) MarshalJSON() (data []byte, err error) {
	type shadow TravelRuleVerifyDepositByUuidParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *TravelRuleVerifyDepositByUuidParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}
