// Withdrawals management scenario.
//
// Demonstrates withdrawal management features using the Upbit Exchange API.
//   - List allowed withdrawal addresses
//   - Check deposit/withdrawal service status
//   - Retrieve withdrawal chance
//   - List / retrieve individual withdrawals
//   - Withdraw digital asset
//   - Withdraw KRW (fiat)
//
// Dry run (default):
//
//	Only performs withdrawal-related queries.
//	Safe to run as-is to verify behavior.
//
// Real execution (DRY_RUN=false):
//
//	Requests a KRW withdrawal using the default 10,000 KRW amount and kakao authentication.
//	Digital asset withdrawals require an actual withdrawal address in the code.
//
//	WARNING: withdrawals are irreversible. Sending to a wrong address may lose your funds.
//
// Usage:
//
//	# Dry run (default) — read-only
//	UPBIT_ACCESS_KEY=<key> UPBIT_SECRET_KEY=<secret> go run examples/withdrawals.go
//
//	# Real execution — includes a KRW withdrawal request (default: 10,000 KRW, kakao)
//	DRY_RUN=false UPBIT_ACCESS_KEY=<key> UPBIT_SECRET_KEY=<secret> go run examples/withdrawals.go
package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
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

func formatInt64WithComma(n int64) string {
	sign := ""
	if n < 0 {
		sign = "-"
		n = -n
	}
	s := fmt.Sprintf("%d", n)
	if len(s) <= 3 {
		return sign + s
	}
	var result []byte
	for i, c := range s {
		if i > 0 && (len(s)-i)%3 == 0 {
			result = append(result, ',')
		}
		result = append(result, byte(c))
	}
	return sign + string(result)
}

// ---------------------------------------------------------------------------
// 1. List allowed withdrawal addresses
// ---------------------------------------------------------------------------

func listWithdrawalAddresses(ctx context.Context, client upbit.Client) {
	divider("1. Allowed Withdrawal Addresses")

	addresses, err := client.Withdraws.ListCoinAddresses(ctx)
	if err != nil {
		log.Fatalf("Failed to list withdrawal addresses: %v", err)
	}
	if len(*addresses) == 0 {
		fmt.Println("  No allowed withdrawal addresses registered.")
		fmt.Println("  Please register allowed withdrawal addresses on the Upbit website first.")
		return
	}

	fmt.Printf("  %-8s %-8s %-12s %s\n", "Currency", "Network", "Exchange", "Address")
	fmt.Printf("  %s\n", strings.Repeat("-", 60))
	for i, addr := range *addresses {
		if i >= 10 {
			break
		}
		exchange := addr.ExchangeName
		if exchange == "" {
			exchange = "-"
		}
		fmt.Printf("  %-8s %-8s %-12s %s\n", addr.Currency, addr.NetType, exchange, addr.WithdrawAddress)
	}
}

// ---------------------------------------------------------------------------
// 2. Check deposit/withdrawal service status
// ---------------------------------------------------------------------------

func checkWalletStatus(ctx context.Context, client upbit.Client) {
	divider("2. Deposit/Withdrawal Service Status")

	wallets, err := client.WalletStatus.List(ctx)
	if err != nil {
		log.Fatalf("Failed to fetch wallet status: %v", err)
	}

	targetCurrencies := map[string]bool{"BTC": true, "ETH": true, "USDT": true}

	fmt.Printf("  Total wallets: %d\n", len(*wallets))
	fmt.Printf("\n  %-8s %-10s %-12s %s\n", "Currency", "Network", "Wallet", "Block")
	fmt.Printf("  %s\n", strings.Repeat("-", 48))
	for _, w := range *wallets {
		if targetCurrencies[w.Currency] {
			fmt.Printf("  %-8s %-10s %-12s %s\n",
				w.Currency, w.NetType, string(w.WalletState), string(w.BlockState))
		}
	}
}

// ---------------------------------------------------------------------------
// 3. Retrieve withdrawal chance
// ---------------------------------------------------------------------------

func checkWithdrawalChance(ctx context.Context, client upbit.Client, currency, netType string) {
	divider(fmt.Sprintf("3. %s Withdrawal Chance", currency))

	chance, err := client.Withdraws.GetChance(ctx, upbit.WithdrawGetChanceParams{
		Currency: currency,
		NetType:  netType,
	})
	if err != nil {
		log.Fatalf("Failed to fetch withdrawal chance: %v", err)
	}

	fmt.Printf("  Balance:           %s\n", chance.Account.Balance)
	fmt.Printf("  Withdrawal fee:    %s\n", chance.Currency.WithdrawFee)
	fmt.Printf("  Minimum amount:    %s\n", chance.WithdrawLimit.Minimum)
	fmt.Printf("  Remaining daily:   %s %s\n", chance.WithdrawLimit.RemainingDailyFiat, chance.WithdrawLimit.FiatCurrency)
}

// ---------------------------------------------------------------------------
// 4. List withdrawals
// ---------------------------------------------------------------------------

func listWithdrawals(ctx context.Context, client upbit.Client) {
	divider("4. List Withdrawals")

	page, err := client.Withdraws.List(ctx, upbit.WithdrawListParams{
		Limit: param.NewOpt[int64](5),
	})
	if err != nil {
		log.Fatalf("Failed to list withdrawals: %v", err)
	}
	withdrawals := page.Items
	if len(withdrawals) == 0 {
		fmt.Println("  No withdrawal history.")
		return
	}

	fmt.Printf("  %10s %-6s %14s %-12s\n", "UUID", "Currency", "Amount", "State")
	fmt.Printf("  %s\n", strings.Repeat("-", 50))
	for _, w := range withdrawals {
		uuid := w.Uuid
		if len(uuid) > 8 {
			uuid = uuid[:8] + "..."
		}
		fmt.Printf("  %s %-6s %14s %-12s\n", uuid, w.Currency, w.Amount, string(w.State))
	}
}

// ---------------------------------------------------------------------------
// 5. Retrieve individual withdrawal
// ---------------------------------------------------------------------------

func retrieveWithdrawal(ctx context.Context, client upbit.Client) {
	divider("5. Retrieve Individual Withdrawal")

	withdrawal, err := client.Withdraws.Get(ctx, upbit.WithdrawGetParams{})
	if err != nil {
		fmt.Printf("  Failed to retrieve withdrawal (no withdrawal history?): %v\n", err)
		return
	}
	fmt.Printf("  UUID:     %s\n", withdrawal.Uuid)
	fmt.Printf("  Currency: %s\n", withdrawal.Currency)
	fmt.Printf("  Amount:   %s\n", withdrawal.Amount)
	fmt.Printf("  State:    %s\n", withdrawal.State)
	if withdrawal.Txid != "" {
		fmt.Printf("  TxID:     %s\n", withdrawal.Txid)
	}
}

// ---------------------------------------------------------------------------
// 6. Withdraw digital asset
// ---------------------------------------------------------------------------

func withdrawCoin(ctx context.Context, client upbit.Client, currency, netType, amount, address string) {
	divider("6. Withdraw Digital Asset")

	fmt.Printf("  Currency: %s\n", currency)
	fmt.Printf("  Network:  %s\n", netType)
	fmt.Printf("  Amount:   %s\n", amount)
	fmt.Printf("  Address:  %s\n", address)

	if dryRun {
		fmt.Println("\n  [Skipped] Dry run mode — actual withdrawal is omitted.")
		return
	}

	if strings.HasPrefix(address, "<") {
		fmt.Println("\n  [Skipped] Withdrawal address is not set.")
		fmt.Println("  Update the address value in the code with a real withdrawal address and run again.")
		return
	}

	result, err := client.Withdraws.NewWithdrawal(ctx, upbit.WithdrawNewWithdrawalParams{
		Currency:        currency,
		NetType:         netType,
		Amount:          amount,
		Address:         address,
		TransactionType: upbit.WithdrawNewWithdrawalParamsTransactionTypeDefault,
	})
	if err != nil {
		log.Fatalf("Digital asset withdrawal failed: %v", err)
	}
	fmt.Printf("\n  Withdrawal UUID: %s\n", result.Uuid)
	fmt.Printf("  State:           %s\n", result.State)
}

// ---------------------------------------------------------------------------
// 7. Withdraw KRW (fiat)
// ---------------------------------------------------------------------------

func withdrawKRW(ctx context.Context, client upbit.Client, amount string) {
	divider("7. Withdraw KRW")

	amt, _, _ := new(big.Float).Parse(amount, 10)
	amtInt, _ := amt.Int64()
	fmt.Printf("  Withdrawal amount: %s KRW\n", formatInt64WithComma(amtInt))

	if dryRun {
		fmt.Println("\n  [Skipped] Dry run mode — actual withdrawal is omitted.")
		return
	}

	result, err := client.Withdraws.NewKrwWithdrawal(ctx, upbit.WithdrawNewKrwWithdrawalParams{
		Amount:        amount,
		TwoFactorType: upbit.WithdrawNewKrwWithdrawalParamsTwoFactorTypeKakao,
	})
	if err != nil {
		log.Fatalf("KRW withdrawal failed: %v", err)
	}
	fmt.Printf("\n  Withdrawal UUID: %s\n", result.Uuid)
	fmt.Printf("  State:           %s\n", result.State)
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

	listWithdrawalAddresses(ctx, client)
	checkWalletStatus(ctx, client)
	checkWithdrawalChance(ctx, client, "BTC", "BTC")
	listWithdrawals(ctx, client)
	retrieveWithdrawal(ctx, client)
	withdrawCoin(ctx, client, "USDT", "TRX", "13.241", "<your_withdrawal_address>")
	withdrawKRW(ctx, client, "10000")

	fmt.Printf("\n%s\n", strings.Repeat("=", 60))
	fmt.Println("  Withdrawals scenario complete")
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println()
	fmt.Println("  ※ This example is for educational purposes.")
	fmt.Println("    We are not responsible for any losses incurred when used for actual investment.")
}
