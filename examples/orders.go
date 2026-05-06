// Order creation and management scenario.
//
// Demonstrates how to create and manage orders using the Upbit Exchange API.
//   - Check order availability
//   - Create limit buy order -> retrieve -> cancel
//   - Check market sell availability
//   - List closed orders
//
// Dry run (default):
//
//	Only queries order-related information; does not create real orders.
//
// Real execution (DRY_RUN=false):
//
//	Executes limit buy order creation -> retrieval -> cancellation.
//
// Usage:
//
//	# Dry run (default)
//	UPBIT_ACCESS_KEY=<key> UPBIT_SECRET_KEY=<secret> go run examples/orders.go
//
//	# Real execution
//	DRY_RUN=false UPBIT_ACCESS_KEY=<key> UPBIT_SECRET_KEY=<secret> go run examples/orders.go
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

const market = "KRW-BTC"

var dryRun = os.Getenv("DRY_RUN") != "false"

func divider(title string) {
	fmt.Printf("\n%s\n  %s\n%s\n", strings.Repeat("=", 60), title, strings.Repeat("=", 60))
}

// roundPrice rounds price down to the nearest tick size.
func roundPrice(price, tickSize *big.Float) string {
	pInt, _ := price.Int(nil)
	tInt, _ := tickSize.Int(nil)

	if tInt.Sign() > 0 && tInt.Int64() >= 1 {
		pInt.Div(pInt, tInt)
		pInt.Mul(pInt, tInt)
		return pInt.String()
	}

	q := new(big.Float).Quo(price, tickSize)
	qInt, _ := q.Int(nil)
	result := new(big.Float).Mul(new(big.Float).SetInt(qInt), tickSize)
	return result.Text('f', 8)
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
// 1. Check order availability
// ---------------------------------------------------------------------------

func checkOrderChance(ctx context.Context, client upbit.Client) {
	divider("1. Check Order Availability")

	chance, err := client.Orders.GetChance(ctx, upbit.OrderGetChanceParams{
		Market: market,
	})
	if err != nil {
		log.Fatalf("[Error] Failed to fetch order chance: %v", err)
	}

	fmt.Printf("  Market:        %s\n", market)
	fmt.Printf("  Bid fee:       %s\n", chance.BidFee)
	fmt.Printf("  Ask fee:       %s\n", chance.AskFee)
	fmt.Printf("  Bid types:     %v\n", chance.Market.BidTypes)
	fmt.Printf("  Ask types:     %v\n", chance.Market.AskTypes)
	fmt.Printf("  Min bid total: %s\n", chance.Market.Bid.MinTotal)
	fmt.Printf("  Min ask total: %s\n", chance.Market.Ask.MinTotal)
}

// ---------------------------------------------------------------------------
// 2. Limit buy -> retrieve -> cancel
// ---------------------------------------------------------------------------

func limitBidAndCancel(ctx context.Context, client upbit.Client) {
	divider("2. Limit Buy -> Retrieve -> Cancel")

	// Fetch tick size
	// WARNING: this example queries the tick size via API before every order.
	// In production, querying tick size per order is not recommended.
	// See the Upbit official documentation for tick size policy.
	instruments, err := client.Orderbooks.ListInstruments(ctx, upbit.OrderbookListInstrumentsParams{
		Markets: market,
	})
	if err != nil {
		log.Fatalf("Failed to fetch tick size: %v", err)
	}
	tickSize, _, _ := new(big.Float).Parse((*instruments)[0].TickSize, 10)
	fmt.Printf("  Tick size:       %s\n", (*instruments)[0].TickSize)

	// Fetch best bid price
	orderbooks, err := client.Orderbooks.List(ctx, upbit.OrderbookListParams{Markets: market})
	if err != nil {
		log.Fatalf("Failed to fetch orderbook: %v", err)
	}
	bestBid, _, _ := new(big.Float).Parse((*orderbooks)[0].OrderbookUnits[0].BidPrice, 10)
	bestBidInt, _ := bestBid.Int64()
	fmt.Printf("  Best bid price:  %s KRW\n", formatInt64WithComma(bestBidInt))

	// Target price 3% lower, aligned to tick size
	targetPrice := roundPrice(
		new(big.Float).Mul(bestBid, new(big.Float).SetFloat64(0.97)),
		tickSize,
	)
	volume := "0.0001"
	tp, _, _ := new(big.Float).Parse(targetPrice, 10)
	tpInt, _ := tp.Int64()
	fmt.Printf("  Order price:     %s KRW (3%% discount)\n", formatInt64WithComma(tpInt))
	fmt.Printf("  Order volume:    %s\n", volume)

	if dryRun {
		fmt.Println("\n  [Skipped] Dry run mode — actual order creation is omitted.")
		return
	}

	// Create order
	order, err := client.Orders.New(ctx, upbit.OrderNewParams{
		Market:  market,
		Side:    upbit.SideEnumBid,
		OrdType: upbit.OrdTypeEnumLimit,
		Price:   param.NewOpt(targetPrice),
		Volume:  param.NewOpt(volume),
	})
	if err != nil {
		log.Fatalf("Failed to create order: %v", err)
	}
	fmt.Println("\n  Order created")
	fmt.Printf("  UUID:   %s\n", order.Uuid)
	fmt.Printf("  State:  %s\n", order.State)

	// Retrieve order
	info, err := client.Orders.Get(ctx, upbit.OrderGetParams{Uuid: param.NewOpt(order.Uuid)})
	if err != nil {
		log.Fatalf("Failed to retrieve order: %v", err)
	}
	fmt.Printf("  Retrieve: state=%s, executed_volume=%s\n", info.State, info.ExecutedVolume)

	// Cancel order
	_, err = client.Orders.Cancel(ctx, upbit.OrderCancelParams{Uuid: param.NewOpt(order.Uuid)})
	if err != nil {
		log.Fatalf("Failed to cancel order: %v", err)
	}
	var cancelled *upbit.Order
	for range 10 {
		cancelled, err = client.Orders.Get(ctx, upbit.OrderGetParams{Uuid: param.NewOpt(order.Uuid)})
		if err != nil {
			log.Fatalf("Failed to retrieve order: %v", err)
		}
		if cancelled.State == "done" || cancelled.State == "cancel" {
			break
		}
		time.Sleep(500 * time.Millisecond)
	}
	if cancelled != nil {
		fmt.Printf("  Cancel:   state=%s\n", cancelled.State)
	}
}

// ---------------------------------------------------------------------------
// 3. Check holdings + market sell availability
// ---------------------------------------------------------------------------

func checkMarketSell(ctx context.Context, client upbit.Client) {
	divider("3. Check Holdings + Market Sell Availability")

	coin := strings.Split(market, "-")[1]

	accounts, err := client.Accounts.List(ctx)
	if err != nil {
		log.Fatalf("[Error] Failed to fetch accounts: %v", err)
	}

	found := false
	for _, a := range *accounts {
		if a.Currency == coin {
			balance, _, _ := new(big.Float).Parse(a.Balance, 10)
			zero := new(big.Float)
			if balance.Cmp(zero) == 0 {
				break
			}
			avgBuy, _, _ := new(big.Float).Parse(a.AvgBuyPrice, 10)
			avgInt, _ := avgBuy.Int64()
			fmt.Printf("  Asset:             %s\n", coin)
			fmt.Printf("  Balance:           %s\n", balance.Text('f', 8))
			fmt.Printf("  Average buy price: %s KRW\n", formatInt64WithComma(avgInt))
			found = true
			break
		}
	}
	if !found {
		fmt.Printf("  No %s holdings.\n", coin)
	}

	// Check market sell availability
	chance, err := client.Orders.GetChance(ctx, upbit.OrderGetChanceParams{Market: market})
	if err != nil {
		log.Fatalf("Failed to fetch order chance: %v", err)
	}
	hasMarket := false
	for _, t := range chance.Market.AskTypes {
		if string(t) == "market" {
			hasMarket = true
			break
		}
	}
	if hasMarket {
		fmt.Println("  Market sell: available")
	} else {
		fmt.Printf("  Market sell: unavailable (ask_types=%v)\n", chance.Market.AskTypes)
	}
}

// ---------------------------------------------------------------------------
// 4. List closed orders
// ---------------------------------------------------------------------------

func listClosedOrders(ctx context.Context, client upbit.Client) {
	divider("4. List Closed Orders")

	closed, err := client.Orders.ListClosed(ctx, upbit.OrderListClosedParams{
		Market: param.NewOpt(market),
		Limit:  param.NewOpt[int64](5),
	})
	if err != nil {
		log.Fatalf("Failed to list closed orders: %v", err)
	}
	if len(*closed) == 0 {
		fmt.Println("  No closed orders.")
		return
	}

	fmt.Printf("  %10s %4s %8s %6s %14s\n", "UUID", "Side", "Type", "State", "Executed Vol")
	fmt.Printf("  %s\n", strings.Repeat("-", 50))
	for _, o := range *closed {
		side := "SELL"
		if o.Side == "bid" {
			side = "BUY"
		}
		uuid := o.Uuid
		if len(uuid) > 8 {
			uuid = uuid[:8] + "..."
		}
		fmt.Printf("  %s %4s %8s %6s %14s\n",
			uuid, side, string(o.OrdType), string(o.State), o.ExecutedVolume)
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

	checkOrderChance(ctx, client)
	limitBidAndCancel(ctx, client)
	checkMarketSell(ctx, client)
	listClosedOrders(ctx, client)

	fmt.Printf("\n%s\n", strings.Repeat("=", 60))
	fmt.Println("  Orders scenario complete")
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println()
	fmt.Println("  ※ This example is for educational purposes.")
	fmt.Println("    We are not responsible for any losses incurred when used for actual investment.")
}
