// TP/SL (take-profit / stop-loss) automated sell scenario.
//
// Demonstrates polling the current price using the Upbit Exchange API and
// executing a market-sell when the target (take-profit) or stop-loss is reached.
//   - Fetch current price and compute TP/SL prices
//   - Monitor price (REST polling)
//   - Execute market-sell
//
// The Upbit API does not support stop/limit reserved orders, so this is
// implemented client-side with polling.
//
// Dry run (default):
//
//	Based on the current price, no orders are executed.
//	Safe to run as-is to verify behavior.
//
// Real execution (DRY_RUN=false):
//
//	Based on the average buy price, real sell orders are executed.
//	You must hold the asset.
//
// Usage:
//
//	# Dry run (default) — monitoring only, no orders placed
//	UPBIT_ACCESS_KEY=<key> UPBIT_SECRET_KEY=<secret> go run examples/tp_sl.go
//
//	# Real execution — based on average buy price, real sell
//	DRY_RUN=false UPBIT_ACCESS_KEY=<key> UPBIT_SECRET_KEY=<secret> go run examples/tp_sl.go
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
	market       = "KRW-BTC"
	sellVolume   = "0.0001"
	pollInterval = 1 * time.Second
	maxPolls     = 10
)

var (
	dryRun    = os.Getenv("DRY_RUN") != "false"
	tpPercent = new(big.Float).SetInt64(3) // Take-profit threshold: +3%
	slPercent = new(big.Float).SetInt64(2) // Stop-loss threshold: -2%
)

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

func fetchCurrentPrice(ctx context.Context, client upbit.Client) *big.Float {
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
// 1. Pre-flight checks + TP/SL settings
// ---------------------------------------------------------------------------

func setupTPSL(ctx context.Context, client upbit.Client) (*big.Float, *big.Float) {
	modeLabel := "DRY RUN"
	if !dryRun {
		modeLabel = "LIVE"
	}
	divider(fmt.Sprintf("1. TP/SL Settings [%s]", modeLabel))

	current := fetchCurrentPrice(ctx, client)
	hundred := new(big.Float).SetInt64(100)
	one := new(big.Float).SetInt64(1)

	var basePrice *big.Float

	if dryRun {
		basePrice = current
		fmt.Printf("  Base price: %14s KRW (current price)\n", formatKRW(basePrice))
	} else {
		coin := strings.Split(market, "-")[1]
		accounts, err := client.Accounts.List(ctx)
		if err != nil {
			log.Fatalf("[Error] Failed to fetch accounts (check API keys): %v", err)
		}

		var balance *big.Float
		var avgBuyPrice *big.Float
		for _, a := range *accounts {
			if a.Currency == coin {
				balance, _, _ = new(big.Float).Parse(a.Balance, 10)
				avgBuyPrice, _, _ = new(big.Float).Parse(a.AvgBuyPrice, 10)
				break
			}
		}

		zero := new(big.Float)
		if balance == nil || balance.Cmp(zero) == 0 {
			log.Fatalf("[Error] No %s holdings.", coin)
		}

		sv, _, _ := new(big.Float).Parse(sellVolume, 10)
		if balance.Cmp(sv) < 0 {
			log.Fatalf("[Error] %s balance (%s) is less than sell volume (%s).",
				coin, balance.Text('f', 8), sellVolume)
		}

		basePrice = avgBuyPrice
		fmt.Printf("  Balance:    %s %s\n", coin, balance.Text('f', 8))
		fmt.Printf("  Base price: %14s KRW (average buy price)\n", formatKRW(basePrice))
		fmt.Printf("  Price:      %14s KRW\n", formatKRW(current))
	}

	// tp_price = base_price * (1 + tp_percent / 100)
	tpMul := new(big.Float).Quo(tpPercent, hundred)
	tpMul.Add(tpMul, one)
	tpPrice := new(big.Float).Mul(basePrice, tpMul)

	// sl_price = base_price * (1 - sl_percent / 100)
	slMul := new(big.Float).Quo(slPercent, hundred)
	slMul.Sub(one, slMul)
	slPrice := new(big.Float).Mul(basePrice, slMul)

	fmt.Printf("  Market:     %s\n", market)
	fmt.Printf("  Sell volume: %s\n", sellVolume)
	fmt.Printf("  TP price:   %14s KRW (+%s%%)\n", formatKRW(tpPrice), tpPercent.Text('f', 0))
	fmt.Printf("  SL price:   %14s KRW (-%s%%)\n", formatKRW(slPrice), slPercent.Text('f', 0))
	fmt.Printf("  Polling interval: %s × max %d times\n", pollInterval, maxPolls)
	if dryRun {
		fmt.Println("\n  ※ Dry run mode: orders will not be executed.")
	}

	return tpPrice, slPrice
}

// ---------------------------------------------------------------------------
// 2. Price monitoring
// ---------------------------------------------------------------------------

func monitorPrice(ctx context.Context, client upbit.Client, tpPrice, slPrice *big.Float) string {
	divider("2. Start Price Monitoring")

	for poll := 1; poll <= maxPolls; poll++ {
		current := fetchCurrentPrice(ctx, client)

		fmt.Printf("  [%3d/%d] Price: %14s KRW", poll, maxPolls, formatKRW(current))

		if current.Cmp(tpPrice) >= 0 {
			fmt.Println("  → Take-profit hit!")
			return "tp"
		}
		if current.Cmp(slPrice) <= 0 {
			fmt.Println("  → Stop-loss hit!")
			return "sl"
		}

		fmt.Println()

		if poll < maxPolls {
			time.Sleep(pollInterval)
		}
	}

	fmt.Printf("\n  Reached maximum polling count (%d).\n", maxPolls)
	return "timeout"
}

// ---------------------------------------------------------------------------
// 3. Execute market-sell
// ---------------------------------------------------------------------------

func executeSell(ctx context.Context, client upbit.Client, trigger string) {
	label := "take-profit"
	if trigger == "sl" {
		label = "stop-loss"
	}
	divider(fmt.Sprintf("3. Execute Sell (trigger: %s)", label))

	if dryRun {
		fmt.Println("  [Skipped] Dry run mode — actual order is omitted.")
		return
	}

	order, err := client.Orders.New(ctx, upbit.OrderNewParams{
		Market:  market,
		Side:    upbit.SideEnumAsk,
		OrdType: upbit.OrdTypeEnumMarket,
		Volume:  param.NewOpt(sellVolume),
	})
	if err != nil {
		log.Fatalf("Failed to create sell order: %v", err)
	}
	fmt.Println("  Order created")
	fmt.Printf("  UUID:   %s\n", order.Uuid)
	fmt.Printf("  State:  %s\n", order.State)

	// Waiting for order fill
	var info *upbit.Order
	for range 20 {
		info, err = client.Orders.Get(ctx, upbit.OrderGetParams{Uuid: param.NewOpt(order.Uuid)})
		if err != nil {
			log.Fatalf("Failed to retrieve order: %v", err)
		}
		if info.State == "done" || info.State == "cancel" {
			break
		}
		time.Sleep(500 * time.Millisecond)
	}

	if info != nil {
		fmt.Printf("  Fill state:    %s\n", info.State)
		fmt.Printf("  Filled volume: %s\n", info.ExecutedVolume)
		if len(info.Trades) > 0 {
			total := new(big.Float)
			for _, t := range info.Trades {
				if t.Funds != "" {
					f, _, _ := new(big.Float).Parse(t.Funds, 10)
					total.Add(total, f)
				}
			}
			fmt.Printf("  Filled amount: %14s KRW\n", formatKRW(total))
		}
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

	tpPrice, slPrice := setupTPSL(ctx, client)
	trigger := monitorPrice(ctx, client, tpPrice, slPrice)

	if trigger == "timeout" {
		fmt.Println("\n  Timeout — not executing sell.")
	} else {
		executeSell(ctx, client, trigger)
	}

	fmt.Printf("\n%s\n", strings.Repeat("=", 60))
	fmt.Println("  TP/SL scenario complete")
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println()
	fmt.Println("  ※ This example is for educational purposes.")
	fmt.Println("    We are not responsible for any losses incurred when used for actual investment.")
}
