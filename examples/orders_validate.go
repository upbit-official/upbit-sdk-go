// Order validation (test new) scenario.
//
// Demonstrates how to validate various order types using the order test API
// without placing real orders.
//   - Check order availability
//   - Test limit buy/sell
//   - Test market buy/sell
//   - Test best-price buy/sell (IOC)
//   - Validate invalid order requests
//
// The order test API performs the same validation as real orders, but
// since no real order is created, you can safely test without incurring fees.
// The UUID returned in the response cannot be used to retrieve or cancel orders.
//
// Usage:
//
//	UPBIT_ACCESS_KEY=<key> UPBIT_SECRET_KEY=<secret> go run examples/orders_validate.go
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

const (
	market     = "KRW-BTC"
	testVolume = "0.0001"
)

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
// 2. Test limit buy/sell
// ---------------------------------------------------------------------------

func testLimitOrders(ctx context.Context, client upbit.Client) {
	divider("2. Test Limit Buy/Sell")

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

	// --- Limit buy (3% discount) ---
	bidPrice := roundPrice(
		new(big.Float).Mul(bestBid, new(big.Float).SetFloat64(0.97)),
		tickSize,
	)
	bp, _, _ := new(big.Float).Parse(bidPrice, 10)
	bpInt, _ := bp.Int64()
	fmt.Printf("\n  [Limit Buy]\n")
	fmt.Printf("  Order price:     %s KRW (3%% discount)\n", formatInt64WithComma(bpInt))
	fmt.Printf("  Order volume:    %s\n", testVolume)

	bidOrder, err := client.Orders.TestNew(ctx, upbit.OrderTestNewParams{
		Market:  market,
		Side:    upbit.SideEnumBid,
		OrdType: upbit.OrdTypeEnumLimit,
		Price:   param.NewOpt(bidPrice),
		Volume:  param.NewOpt(testVolume),
	})
	if err != nil {
		log.Fatalf("Limit buy test failed: %v", err)
	}
	fmt.Printf("  → UUID:  %s\n", bidOrder.Uuid)
	fmt.Printf("  → State: %s\n", bidOrder.State)
	fmt.Printf("  → Type:  %s\n", bidOrder.OrdType)

	// --- Limit sell (3% premium) ---
	askPrice := roundPrice(
		new(big.Float).Mul(bestBid, new(big.Float).SetFloat64(1.03)),
		tickSize,
	)
	ap, _, _ := new(big.Float).Parse(askPrice, 10)
	apInt, _ := ap.Int64()
	fmt.Printf("\n  [Limit Sell]\n")
	fmt.Printf("  Order price:     %s KRW (3%% premium)\n", formatInt64WithComma(apInt))
	fmt.Printf("  Order volume:    %s\n", testVolume)

	askOrder, err := client.Orders.TestNew(ctx, upbit.OrderTestNewParams{
		Market:  market,
		Side:    upbit.SideEnumAsk,
		OrdType: upbit.OrdTypeEnumLimit,
		Price:   param.NewOpt(askPrice),
		Volume:  param.NewOpt(testVolume),
	})
	if err != nil {
		log.Fatalf("Limit sell test failed: %v", err)
	}
	fmt.Printf("  → UUID:  %s\n", askOrder.Uuid)
	fmt.Printf("  → State: %s\n", askOrder.State)
	fmt.Printf("  → Type:  %s\n", askOrder.OrdType)
}

// ---------------------------------------------------------------------------
// 3. Test market buy/sell
// ---------------------------------------------------------------------------

func testMarketOrders(ctx context.Context, client upbit.Client) {
	divider("3. Test Market Buy/Sell")

	// Fetch minimum order amount
	chance, err := client.Orders.GetChance(ctx, upbit.OrderGetChanceParams{Market: market})
	if err != nil {
		log.Fatalf("Failed to fetch order chance: %v", err)
	}
	minTotal := chance.Market.Bid.MinTotal
	mt, _, _ := new(big.Float).Parse(minTotal, 10)
	mtInt, _ := mt.Int64()

	// --- Market buy (specify total amount) ---
	fmt.Printf("  [Market Buy]\n")
	fmt.Printf("  Order amount:    %s KRW (minimum order amount)\n", formatInt64WithComma(mtInt))

	buyOrder, err := client.Orders.TestNew(ctx, upbit.OrderTestNewParams{
		Market:  market,
		Side:    upbit.SideEnumBid,
		OrdType: upbit.OrdTypeEnumPrice,
		Price:   param.NewOpt(minTotal),
	})
	if err != nil {
		log.Fatalf("Market buy test failed: %v", err)
	}
	fmt.Printf("  → UUID:  %s\n", buyOrder.Uuid)
	fmt.Printf("  → State: %s\n", buyOrder.State)
	fmt.Printf("  → Type:  %s\n", buyOrder.OrdType)

	// --- Market sell (specify volume) ---
	fmt.Printf("\n  [Market Sell]\n")
	fmt.Printf("  Order volume:    %s\n", testVolume)

	sellOrder, err := client.Orders.TestNew(ctx, upbit.OrderTestNewParams{
		Market:  market,
		Side:    upbit.SideEnumAsk,
		OrdType: upbit.OrdTypeEnumMarket,
		Volume:  param.NewOpt(testVolume),
	})
	if err != nil {
		log.Fatalf("Market sell test failed: %v", err)
	}
	fmt.Printf("  → UUID:  %s\n", sellOrder.Uuid)
	fmt.Printf("  → State: %s\n", sellOrder.State)
	fmt.Printf("  → Type:  %s\n", sellOrder.OrdType)
}

// ---------------------------------------------------------------------------
// 4. Test best-price buy/sell (IOC)
// ---------------------------------------------------------------------------

func testBestOrders(ctx context.Context, client upbit.Client) {
	divider("4. Test Best-Price Buy/Sell (IOC)")

	// Fetch minimum order amount
	chance, err := client.Orders.GetChance(ctx, upbit.OrderGetChanceParams{Market: market})
	if err != nil {
		log.Fatalf("Failed to fetch order chance: %v", err)
	}
	minTotal := chance.Market.Bid.MinTotal
	mt, _, _ := new(big.Float).Parse(minTotal, 10)
	mtInt, _ := mt.Int64()

	// --- Best-price buy (specify total amount, IOC) ---
	fmt.Printf("  [Best-Price Buy — IOC]\n")
	fmt.Printf("  Order amount:    %s KRW (minimum order amount)\n", formatInt64WithComma(mtInt))

	bidOrder, err := client.Orders.TestNew(ctx, upbit.OrderTestNewParams{
		Market:      market,
		Side:        upbit.SideEnumBid,
		OrdType:     upbit.OrdTypeEnumBest,
		Price:       param.NewOpt(minTotal),
		TimeInForce: upbit.TimeInForceEnumIoc,
	})
	if err != nil {
		log.Fatalf("Best-price buy test failed: %v", err)
	}
	fmt.Printf("  → UUID:  %s\n", bidOrder.Uuid)
	fmt.Printf("  → State: %s\n", bidOrder.State)
	fmt.Printf("  → Type:  %s\n", bidOrder.OrdType)

	// --- Best-price sell (specify volume, IOC) ---
	fmt.Printf("\n  [Best-Price Sell — IOC]\n")
	fmt.Printf("  Order volume:    %s\n", testVolume)

	askOrder, err := client.Orders.TestNew(ctx, upbit.OrderTestNewParams{
		Market:      market,
		Side:        upbit.SideEnumAsk,
		OrdType:     upbit.OrdTypeEnumBest,
		Volume:      param.NewOpt(testVolume),
		TimeInForce: upbit.TimeInForceEnumIoc,
	})
	if err != nil {
		log.Fatalf("Best-price sell test failed: %v", err)
	}
	fmt.Printf("  → UUID:  %s\n", askOrder.Uuid)
	fmt.Printf("  → State: %s\n", askOrder.State)
	fmt.Printf("  → Type:  %s\n", askOrder.OrdType)
}

// ---------------------------------------------------------------------------
// 5. Validate invalid order requests
// ---------------------------------------------------------------------------

func testValidation(ctx context.Context, client upbit.Client) {
	divider("5. Validate Invalid Order Requests")

	// --- Non-existent market ---
	fmt.Println("  [Validation] Non-existent market (KRW-INVALID)")
	_, err := client.Orders.TestNew(ctx, upbit.OrderTestNewParams{
		Market:  "KRW-INVALID",
		Side:    upbit.SideEnumBid,
		OrdType: upbit.OrdTypeEnumLimit,
		Price:   param.NewOpt("10000"),
		Volume:  param.NewOpt("1"),
	})
	if err != nil {
		fmt.Printf("  → Error raised (expected): %v\n", err)
	} else {
		fmt.Println("  → Unexpectedly succeeded.")
	}

	// --- Missing price on limit order ---
	fmt.Println("\n  [Validation] Limit order missing price")
	_, err = client.Orders.TestNew(ctx, upbit.OrderTestNewParams{
		Market:  market,
		Side:    upbit.SideEnumBid,
		OrdType: upbit.OrdTypeEnumLimit,
		Volume:  param.NewOpt(testVolume),
	})
	if err != nil {
		fmt.Printf("  → Error raised (expected): %v\n", err)
	} else {
		fmt.Println("  → Unexpectedly succeeded.")
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
	testLimitOrders(ctx, client)
	testMarketOrders(ctx, client)
	testBestOrders(ctx, client)
	testValidation(ctx, client)

	fmt.Printf("\n%s\n", strings.Repeat("=", 60))
	fmt.Println("  Order validation scenario complete")
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println()
	fmt.Println("  ※ The order test API does not create real orders.")
	fmt.Println("    The returned UUID cannot be used to retrieve or cancel orders.")
}
