// Deposits management scenario.
//
// Demonstrates deposit management features using the Upbit Exchange API.
//   - List deposit addresses
//   - Retrieve / create deposit address
//   - Retrieve deposit chance
//   - List / retrieve individual deposits
//   - List Travel Rule supported VASPs
//   - Verify Travel Rule (by deposit UUID)
//
// Dry run (default):
//
//	Only performs deposit-related queries.
//	Safe to run as-is to verify behavior.
//
// Real execution (DRY_RUN=false):
//
//	Creates a deposit address and verifies Travel Rule.
//	(Travel Rule verification requires a real deposit UUID and VASP UUID.)
//
// Usage:
//
//	# Dry run (default) — read-only
//	UPBIT_ACCESS_KEY=<key> UPBIT_SECRET_KEY=<secret> go run examples/deposits.go
//
//	# Real execution — includes deposit address creation and Travel Rule verification
//	DRY_RUN=false UPBIT_ACCESS_KEY=<key> UPBIT_SECRET_KEY=<secret> go run examples/deposits.go
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	upbit "github.com/upbit-official/upbit-sdk-go"
	"github.com/upbit-official/upbit-sdk-go/option"
	"github.com/upbit-official/upbit-sdk-go/packages/param"
)

var dryRun = os.Getenv("DRY_RUN") != "false"

func divider(title string) {
	fmt.Printf("\n%s\n  %s\n%s\n", strings.Repeat("=", 60), title, strings.Repeat("=", 60))
}

// ---------------------------------------------------------------------------
// 1. List deposit addresses
// ---------------------------------------------------------------------------

func listDepositAddresses(ctx context.Context, client upbit.Client) {
	divider("1. List Deposit Addresses")

	addresses, err := client.Deposits.ListCoinAddresses(ctx)
	if err != nil {
		log.Fatalf("Failed to list deposit addresses: %v", err)
	}
	if len(*addresses) == 0 {
		fmt.Println("  No deposit addresses registered.")
		return
	}

	fmt.Printf("  %-8s %-8s %s\n", "Currency", "Network", "Address")
	fmt.Printf("  %s\n", strings.Repeat("-", 50))
	for i, addr := range *addresses {
		if i >= 10 {
			break
		}
		depositAddr := addr.DepositAddress
		if depositAddr == "" {
			depositAddr = "(not generated)"
		}
		fmt.Printf("  %-8s %-8s %s\n", addr.Currency, addr.NetType, depositAddr)
	}
}

// ---------------------------------------------------------------------------
// 2. Retrieve deposit address for a specific coin
// ---------------------------------------------------------------------------

func retrieveDepositAddress(ctx context.Context, client upbit.Client, currency, netType string) {
	divider(fmt.Sprintf("2. Retrieve %s Deposit Address", currency))

	address, err := client.Deposits.GetCoinAddress(ctx, upbit.DepositGetCoinAddressParams{
		Currency: currency,
		NetType:  netType,
	})
	if err != nil {
		fmt.Printf("  Failed to retrieve deposit address: %v\n", err)
		return
	}
	if address.DepositAddress != "" {
		fmt.Printf("  Currency: %s\n", address.Currency)
		fmt.Printf("  Network:  %s\n", address.NetType)
		fmt.Printf("  Address:  %s\n", address.DepositAddress)
		if address.SecondaryAddress != "" {
			fmt.Printf("  Secondary address: %s\n", address.SecondaryAddress)
		}
	} else {
		fmt.Printf("  %s deposit address has not been generated yet.\n", currency)
	}
}

// ---------------------------------------------------------------------------
// 3. Create deposit address
// ---------------------------------------------------------------------------

func createDepositAddress(ctx context.Context, client upbit.Client, currency, netType string) {
	divider(fmt.Sprintf("3. Create %s(%s) Deposit Address", currency, netType))

	if dryRun {
		fmt.Println("  [Skipped] Dry run mode — actual address creation is omitted.")
		fmt.Printf("  Target: %s (%s)\n", currency, netType)
		return
	}

	result, err := client.Deposits.NewCoinAddress(ctx, upbit.DepositNewCoinAddressParams{
		Currency: currency,
		NetType:  netType,
	})
	if err != nil {
		fmt.Printf("  Failed to create deposit address: %v\n", err)
		return
	}

	if result.DepositAddress != "" {
		fmt.Printf("  Currency: %s\n", result.Currency)
		fmt.Printf("  Network:  %s\n", result.NetType)
		fmt.Printf("  Address:  %s\n", result.DepositAddress)
		if result.SecondaryAddress != "" {
			fmt.Printf("  Secondary address: %s\n", result.SecondaryAddress)
		}
	} else {
		fmt.Printf("  Request status: %v\n", result.Success)
		fmt.Printf("  Message:        %s\n", result.Message)
	}
}

// ---------------------------------------------------------------------------
// 4. Retrieve deposit chance
// ---------------------------------------------------------------------------

func retrieveDepositChance(ctx context.Context, client upbit.Client, currency, netType string) {
	divider(fmt.Sprintf("4. %s Deposit Chance", currency))

	chance, err := client.Deposits.GetChance(ctx, upbit.DepositGetChanceParams{
		Currency: currency,
		NetType:  netType,
	})
	if err != nil {
		log.Fatalf("Failed to fetch deposit chance: %v", err)
	}

	depositable := "no"
	if chance.IsDepositPossible {
		depositable = "yes"
	}
	fmt.Printf("  Deposit possible:  %s\n", depositable)
	fmt.Printf("  Minimum amount:    %s\n", chance.MinimumDepositAmount)
	fmt.Printf("  Required confirms: %d\n", chance.MinimumDepositConfirmations)
	if !chance.IsDepositPossible {
		fmt.Printf("  Reason:            %s\n", chance.DepositImpossibleReason)
	}
}

// ---------------------------------------------------------------------------
// 5. List deposits
// ---------------------------------------------------------------------------

func listDeposits(ctx context.Context, client upbit.Client) {
	divider("5. List Deposits")

	page, err := client.Deposits.List(ctx, upbit.DepositListParams{
		Limit: param.NewOpt[int64](5),
	})
	if err != nil {
		log.Fatalf("Failed to list deposits: %v", err)
	}
	deposits := page.Items
	if len(deposits) == 0 {
		fmt.Println("  No deposit history.")
		return
	}

	fmt.Printf("  %10s %-6s %14s %-24s\n", "UUID", "Currency", "Amount", "State")
	fmt.Printf("  %s\n", strings.Repeat("-", 60))
	for _, d := range deposits {
		uuid := d.Uuid
		if len(uuid) > 8 {
			uuid = uuid[:8] + "..."
		}
		fmt.Printf("  %s %-6s %14s %-24s\n", uuid, d.Currency, d.Amount, string(d.State))
	}
}

// ---------------------------------------------------------------------------
// 6. Retrieve individual deposit (UUID / TxID)
// ---------------------------------------------------------------------------

func retrieveDeposit(ctx context.Context, client upbit.Client) {
	divider("6. Retrieve Individual Deposit")

	deposit, err := client.Deposits.Get(ctx, upbit.DepositGetParams{})
	if err != nil {
		fmt.Printf("  Failed to retrieve deposit: %v\n", err)
		return
	}
	fmt.Printf("  UUID:     %s\n", deposit.Uuid)
	fmt.Printf("  Currency: %s\n", deposit.Currency)
	fmt.Printf("  Amount:   %s\n", deposit.Amount)
	fmt.Printf("  State:    %s\n", deposit.State)
	if deposit.Txid != "" {
		fmt.Printf("  TxID:     %s\n", deposit.Txid)
	}

	// Re-fetch by UUID
	fmt.Println()
	same, err := client.Deposits.Get(ctx, upbit.DepositGetParams{Uuid: param.NewOpt(deposit.Uuid)})
	if err != nil {
		fmt.Printf("  Failed to re-fetch by UUID: %v\n", err)
	} else {
		fmt.Printf("  State (by UUID): %s\n", same.State)
	}

	// Re-fetch by TxID
	if deposit.Txid != "" {
		sameByTxid, err := client.Deposits.Get(ctx, upbit.DepositGetParams{
			Txid:     param.NewOpt(deposit.Txid),
			Currency: param.NewOpt(deposit.Currency),
		})
		if err != nil {
			fmt.Printf("  Failed to re-fetch by TxID: %v\n", err)
		} else {
			fmt.Printf("  State (by TxID): %s\n", sameByTxid.State)
		}
	}
}

// ---------------------------------------------------------------------------
// 7. List Travel Rule supported VASPs
// ---------------------------------------------------------------------------

func listTravelRuleVASPs(ctx context.Context, client upbit.Client) {
	divider("7. Travel Rule Supported VASPs")

	vasps, err := client.TravelRule.ListVasps(ctx)
	if err != nil {
		log.Fatalf("Failed to list Travel Rule VASPs: %v", err)
	}
	fmt.Printf("  Supported VASPs: %d\n\n", len(*vasps))

	fmt.Printf("  %-20s %s\n", "VASP Name", "UUID")
	fmt.Printf("  %s\n", strings.Repeat("-", 50))
	for i, v := range *vasps {
		if i >= 10 {
			break
		}
		fmt.Printf("  %-20s %s\n", v.VaspName, v.VaspUuid)
	}
	if len(*vasps) > 10 {
		fmt.Printf("  ... and %d more\n", len(*vasps)-10)
	}
}

// ---------------------------------------------------------------------------
// 8. Travel Rule verification
// ---------------------------------------------------------------------------

func verifyTravelRule(ctx context.Context, client upbit.Client, depositUUID, vaspUUID string) {
	divider("8. Travel Rule Verification")

	fmt.Printf("  Deposit UUID: %s\n", depositUUID)
	fmt.Printf("  VASP UUID:    %s\n", vaspUUID)

	if dryRun {
		fmt.Println("\n  [Skipped] Dry run mode — Travel Rule verification is omitted.")
		return
	}

	if strings.HasPrefix(depositUUID, "<") {
		fmt.Println("\n  [Skipped] Deposit UUID is not set.")
		fmt.Println("  Update the depositUUID value in the code with a real deposit UUID and run again.")
		return
	}

	if strings.HasPrefix(vaspUUID, "<") {
		fmt.Println("\n  [Skipped] VASP UUID is not set.")
		fmt.Println("  Update the vaspUUID value in the code with a real VASP UUID and run again.")
		return
	}

	result, err := client.TravelRule.VerifyDepositByUuid(ctx, upbit.TravelRuleVerifyDepositByUuidParams{
		DepositUuid: depositUUID,
		VaspUuid:    vaspUUID,
	})
	if err != nil {
		fmt.Printf("  Travel Rule verification failed: %v\n", err)
		return
	}
	fmt.Printf("  Deposit UUID:         %s\n", result.DepositUuid)
	fmt.Printf("  Deposit state:        %s\n", result.DepositState)
	fmt.Printf("  Verification result:  %s\n", result.VerificationResult)
}

// ---------------------------------------------------------------------------
// Run
// ---------------------------------------------------------------------------

func main() {
	ctx := context.Background()
	client := upbit.NewClient(
		option.WithAccessKey(os.Getenv("UPBIT_ACCESS_KEY")),
		option.WithSecretKey(os.Getenv("UPBIT_SECRET_KEY")),
	)

	listDepositAddresses(ctx, client)
	retrieveDepositAddress(ctx, client, "BTC", "BTC")
	createDepositAddress(ctx, client, "USDT", "TRX")
	retrieveDepositChance(ctx, client, "BTC", "BTC")
	listDeposits(ctx, client)
	retrieveDeposit(ctx, client)
	listTravelRuleVASPs(ctx, client)
	verifyTravelRule(ctx, client, "<your_deposit_uuid>", "<deposit_vasp_uuid>")

	fmt.Printf("\n%s\n", strings.Repeat("=", 60))
	fmt.Println("  Deposits scenario complete")
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println()
	fmt.Println("  ※ This example is for educational purposes.")
	fmt.Println("    We are not responsible for any losses incurred when used for actual investment.")
}
