// DCA (recurring buy) automated market-buy scenario.
//
// Demonstrates periodic market-buy of a fixed amount using the Upbit Exchange API.
//   - Fetch current price and verify buy settings
//   - Repeat market-buy for a fixed number of rounds
//   - Buy summary (total spent, total quantity, average price)
//
// * Server note: a DCA strategy requires the bot to run continuously,
//
//	so the server (or scheduler) must be set up by the user.
//	This example assumes such an environment is already in place.
//
// * Order placement note: running this code does not guarantee order acceptance.
//
//	Orders may be rejected due to maintenance, network errors, etc.,
//	so separate handling for order placement errors is required.
//
// Dry run (default):
//
//	Only fetches the current price; does not place orders.
//	Safe to run as-is to verify behavior.
//
// Real execution (DRY_RUN=false):
//
//	Checks balance, then executes real market-buy orders.
//	Requires sufficient KRW balance.
//
// Usage:
//
//	# Dry run (default) — monitoring only, no orders placed
//	UPBIT_ACCESS_KEY=<key> UPBIT_SECRET_KEY=<secret> go run examples/dca.go
//
//	# Real execution — real market-buy
//	DRY_RUN=false UPBIT_ACCESS_KEY=<key> UPBIT_SECRET_KEY=<secret> go run examples/dca.go
package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"
	"time"

	upbit "github.com/upbit-official/upbit-sdk-go"
	"github.com/upbit-official/upbit-sdk-go/option"
	"github.com/upbit-official/upbit-sdk-go/packages/param"
)

const (
	market     = "KRW-BTC"
	buyAmount  = "5000"
	totalRound = 3
	interval   = 5 * time.Second
)

var dryRun = os.Getenv("DRY_RUN") != "false"

type buyResult struct {
	round   int
	spent   *big.Float
	volume  *big.Float
	success bool
}

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

func formatKRW(v *big.Float) string {
	i, _ := v.Int64()
	return formatInt64WithComma(i)
}

func currentPrice(ctx context.Context, client upbit.Client) *big.Float {
	tickers, err := client.Tickers.ListByTradingPairs(ctx, upbit.TickerListByTradingPairsParams{
		Markets: market,
	})
	if err != nil {
		log.Fatalf("Failed to fetch ticker: %v", err)
	}
	price, _, _ := new(big.Float).Parse((*tickers)[0].TradePrice, 10)
	return price
}

// ---------------------------------------------------------------------------
// 1. Pre-flight checks + DCA settings
// ---------------------------------------------------------------------------

func setup(ctx context.Context, client upbit.Client) {
	modeLabel := "DRY RUN"
	if !dryRun {
		modeLabel = "LIVE"
	}
	divider(fmt.Sprintf("1. DCA Settings [%s]", modeLabel))

	current := currentPrice(ctx, client)
	buyAmt, _, _ := new(big.Float).Parse(buyAmount, 10)
	totalNeeded := new(big.Float).Mul(
		new(big.Float).SetInt64(int64(totalRound)),
		buyAmt,
	)

	if dryRun {
		fmt.Printf("  Price:         %14s KRW\n", formatKRW(current))
	} else {
		fiat := strings.Split(market, "-")[0]
		accounts, err := client.Accounts.List(ctx)
		if err != nil {
			log.Fatalf("[Error] Failed to fetch accounts (check API keys): %v", err)
		}

		var available *big.Float
		for _, a := range *accounts {
			if a.Currency == fiat {
				available, _, _ = new(big.Float).Parse(a.Balance, 10)
				break
			}
		}
		if available == nil {
			available = new(big.Float)
		}

		if available.Cmp(totalNeeded) < 0 {
			log.Fatalf("[Error] %s balance (%s) is less than total required (%s).",
				fiat, formatKRW(available), formatKRW(totalNeeded))
		}

		fmt.Printf("  Price:         %14s KRW\n", formatKRW(current))
		fmt.Printf("  %s balance:  %14s KRW\n", fiat, formatKRW(available))
	}

	fmt.Printf("  Market:        %s\n", market)
	fmt.Printf("  Buy per round: %14s KRW\n", formatKRW(buyAmt))
	fmt.Printf("  Round count:   %d\n", totalRound)
	fmt.Printf("  Interval:      %s\n", interval)
	fmt.Printf("  Total needed:  %14s KRW\n", formatKRW(totalNeeded))
	if dryRun {
		fmt.Println("\n  ※ Dry run mode: orders will not be executed.")
	}
}

// ---------------------------------------------------------------------------
// 2. Execute DCA
// ---------------------------------------------------------------------------

func executeDCA(ctx context.Context, client upbit.Client) []buyResult {
	divider("2. Execute DCA")

	if dryRun {
		fmt.Println("  [Skipped] Dry run mode — actual orders are omitted.")
		fmt.Println()
		buyAmt, _, _ := new(big.Float).Parse(buyAmount, 10)
		for r := 1; r <= totalRound; r++ {
			current := currentPrice(ctx, client)
			estVolume := new(big.Float).Quo(buyAmt, current)
			fmt.Printf("  [%d/%d] Price: %14s KRW  (est. volume: %s)\n",
				r, totalRound, formatKRW(current), estVolume.Text('f', 8))
			if r < totalRound {
				time.Sleep(interval)
			}
		}
		return nil
	}

	results := make([]buyResult, 0, totalRound)
	for r := 1; r <= totalRound; r++ {
		fmt.Printf("\n  [%d/%d] Market buy (%s KRW) ...\n", r, totalRound, buyAmount)

		result := executeSingleBuy(ctx, client, r)
		results = append(results, result)

		if r < totalRound {
			if result.success {
				fmt.Printf("  Next buy in %s...\n", interval)
			}
			time.Sleep(interval)
		}
	}

	return results
}

func executeSingleBuy(ctx context.Context, client upbit.Client, r int) buyResult {
	zero := new(big.Float)

	order, err := client.Orders.New(ctx, upbit.OrderNewParams{
		Market:  market,
		Side:    upbit.SideEnumBid,
		OrdType: upbit.OrdTypeEnumPrice,
		Price:   param.NewOpt(buyAmount),
	})
	if err != nil {
		fmt.Printf("  [Error] Order processing error: %v\n", err)
		return buyResult{round: r, spent: zero, volume: zero, success: false}
	}
	fmt.Printf("  UUID:   %s\n", order.Uuid)
	fmt.Printf("  State:  %s\n", order.State)

	// Waiting for order fill
	var info *upbit.Order
	for range 20 {
		info, err = client.Orders.Get(ctx, upbit.OrderGetParams{Uuid: param.NewOpt(order.Uuid)})
		if err != nil {
			fmt.Printf("  [Error] Order retrieval error: %v\n", err)
			return buyResult{round: r, spent: zero, volume: zero, success: false}
		}
		if info.State == "done" || info.State == "cancel" {
			break
		}
		time.Sleep(500 * time.Millisecond)
	}

	if info != nil && string(info.State) == "done" && len(info.Trades) > 0 {
		spent := new(big.Float)
		volume := new(big.Float)
		for _, t := range info.Trades {
			if t.Funds != "" {
				f, _, _ := new(big.Float).Parse(t.Funds, 10)
				spent.Add(spent, f)
			}
			if t.Volume != "" {
				v, _, _ := new(big.Float).Parse(t.Volume, 10)
				volume.Add(volume, v)
			}
		}
		avgPrice := new(big.Float)
		if volume.Cmp(zero) > 0 {
			avgPrice.Quo(spent, volume)
		}
		fmt.Printf("  Filled amount: %14s KRW\n", formatKRW(spent))
		fmt.Printf("  Filled volume: %s\n", volume.Text('f', 8))
		fmt.Printf("  Filled price:  %14s KRW\n", formatKRW(avgPrice))
		return buyResult{round: r, spent: spent, volume: volume, success: true}
	}

	if info != nil {
		fmt.Printf("  Fill state: %s\n", info.State)
	}
	return buyResult{round: r, spent: zero, volume: zero, success: false}
}

// ---------------------------------------------------------------------------
// 3. Buy summary
// ---------------------------------------------------------------------------

func showSummary(results []buyResult) {
	divider("3. Buy Summary")

	if len(results) == 0 {
		fmt.Println("  No results (dry run mode)")
		return
	}

	zero := new(big.Float)
	totalSpent := new(big.Float)
	totalVolume := new(big.Float)
	successes := 0
	failures := 0

	for _, r := range results {
		if r.success {
			successes++
			totalSpent.Add(totalSpent, r.spent)
			totalVolume.Add(totalVolume, r.volume)
		} else {
			failures++
		}
	}

	if successes > 0 {
		avgPrice := new(big.Float)
		if totalVolume.Cmp(zero) > 0 {
			avgPrice.Quo(totalSpent, totalVolume)
		}
		fmt.Printf("  Successes:       %d/%d\n", successes, len(results))
		fmt.Printf("  Total spent:     %14s KRW\n", formatKRW(totalSpent))
		fmt.Printf("  Total quantity:  %s\n", totalVolume.Text('f', 8))
		fmt.Printf("  Average price:   %14s KRW\n", formatKRW(avgPrice))
	} else {
		fmt.Println("  No successful buys.")
	}

	if failures > 0 {
		fmt.Printf("  Failures:        %d\n", failures)
	}
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

	setup(ctx, client)
	results := executeDCA(ctx, client)
	showSummary(results)

	fmt.Printf("\n%s\n", strings.Repeat("=", 60))
	fmt.Println("  DCA scenario complete")
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println()
	fmt.Println("  ※ This example is for educational purposes.")
	fmt.Println("    We are not responsible for any losses incurred when used for actual investment.")
}
