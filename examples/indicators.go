// Indicators scenario.
//
// Demonstrates how to compute investment indicators from Upbit Quotation API responses.
//   - Top 5 pairs by 24h cumulative trading value
//   - RSI (Relative Strength Index) calculation
//
// The Quotation API does not require authentication.
//
// Usage:
//
//	go run examples/indicators.go
package main

import (
	"context"
	"fmt"
	"math/big"
	"sort"
	"strings"

	upbit "github.com/upbit-official/upbit-sdk-go"
	"github.com/upbit-official/upbit-sdk-go/option"
	"github.com/upbit-official/upbit-sdk-go/packages/param"
)

func divider(title string) {
	fmt.Printf("\n%s\n  %s\n%s\n", strings.Repeat("=", 60), title, strings.Repeat("=", 60))
}

func formatKRW(v *big.Float) string {
	i, _ := v.Int64()
	return formatInt64WithComma(i)
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

// formatShort formats a big.Float value with short scale (T/B/M).
func formatShort(v *big.Float) string {
	trillion := new(big.Float).SetInt64(1_000_000_000_000)
	billion := new(big.Float).SetInt64(1_000_000_000)
	million := new(big.Float).SetInt64(1_000_000)

	if v.Cmp(trillion) >= 0 {
		q := new(big.Float).Quo(v, trillion)
		f, _ := q.Float64()
		return fmt.Sprintf("~%.1fT", f)
	}
	if v.Cmp(billion) >= 0 {
		q := new(big.Float).Quo(v, billion)
		f, _ := q.Float64()
		return fmt.Sprintf("~%.1fB", f)
	}
	if v.Cmp(million) >= 0 {
		q := new(big.Float).Quo(v, million)
		f, _ := q.Float64()
		return fmt.Sprintf("~%.0fM", f)
	}
	return formatKRW(v)
}

// ---------------------------------------------------------------------------
// 1. Top 5 pairs by 24h cumulative trading value
// ---------------------------------------------------------------------------

func top5ByTradeVolume(ctx context.Context, client upbit.Client) {
	divider("1. Top 5 by 24h Cumulative Trading Value (KRW market)")

	tickers, err := client.Tickers.ListByQuoteCurrencies(ctx, upbit.TickerListByQuoteCurrenciesParams{
		QuoteCurrencies: "KRW",
	})
	if err != nil {
		fmt.Printf("  Failed to fetch tickers: %v\n", err)
		return
	}

	type tickerItem struct {
		Market   string
		Price    *big.Float
		AccTrade *big.Float
	}

	items := make([]tickerItem, len(*tickers))
	for i, t := range *tickers {
		p, _, _ := new(big.Float).Parse(t.TradePrice, 10)
		a, _, _ := new(big.Float).Parse(t.AccTradePrice24h, 10)
		items[i] = tickerItem{Market: t.Market, Price: p, AccTrade: a}
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].AccTrade.Cmp(items[j].AccTrade) > 0
	})

	fmt.Printf("  %-12s %14s %18s  %s\n", "Market", "Price", "24h Volume", "Short")
	fmt.Printf("  %s\n", strings.Repeat("-", 60))
	for i := 0; i < 5 && i < len(items); i++ {
		t := items[i]
		fmt.Printf("  %-12s %14s %18s  (%s)\n",
			t.Market, formatKRW(t.Price), formatKRW(t.AccTrade), formatShort(t.AccTrade))
	}
}

// ---------------------------------------------------------------------------
// 2. RSI (Relative Strength Index) calculation
// ---------------------------------------------------------------------------

func calculateRSI(ctx context.Context, client upbit.Client) {
	market := "KRW-BTC"
	period := 14

	divider(fmt.Sprintf("2. RSI Calculation (%s, %d-day)", market, period))

	// Fetch 200 daily candles (newest first -> reverse to ascending)
	days, err := client.Candles.ListDays(ctx, upbit.CandleListDaysParams{
		Market: market,
		Count:  param.NewOpt[int64](200),
	})
	if err != nil {
		fmt.Printf("  Failed to fetch daily candles: %v\n", err)
		return
	}

	candles := *days
	if len(candles) < period {
		fmt.Printf("  Error: at least %d candles required (got %d)\n", period, len(candles))
		return
	}

	// Reverse to ascending order
	for i, j := 0, len(candles)-1; i < j; i, j = i+1, j-1 {
		candles[i], candles[j] = candles[j], candles[i]
	}

	zero := new(big.Float)
	gains := make([]*big.Float, len(candles))
	losses := make([]*big.Float, len(candles))

	for i, c := range candles {
		change, _, _ := new(big.Float).Parse(c.ChangePrice, 10)
		if change == nil {
			change = new(big.Float)
		}
		if change.Cmp(zero) > 0 {
			gains[i] = new(big.Float).Copy(change)
			losses[i] = new(big.Float)
		} else if change.Cmp(zero) < 0 {
			gains[i] = new(big.Float)
			losses[i] = new(big.Float).Abs(change)
		} else {
			gains[i] = new(big.Float)
			losses[i] = new(big.Float)
		}
	}

	// Initial average gain/loss (SMA)
	p := new(big.Float).SetInt64(int64(period))
	au := new(big.Float)
	ad := new(big.Float)
	for i := 0; i < period; i++ {
		au.Add(au, gains[i])
		ad.Add(ad, losses[i])
	}
	au.Quo(au, p)
	ad.Quo(ad, p)

	// Wilder smoothing
	pMinus1 := new(big.Float).SetInt64(int64(period - 1))
	for i := period; i < len(gains); i++ {
		au.Mul(au, pMinus1)
		au.Add(au, gains[i])
		au.Quo(au, p)

		ad.Mul(ad, pMinus1)
		ad.Add(ad, losses[i])
		ad.Quo(ad, p)
	}

	// RS & RSI
	hundred := new(big.Float).SetInt64(100)
	one := new(big.Float).SetInt64(1)

	var rsStr string
	var rsiVal *big.Float

	if ad.Cmp(zero) == 0 {
		rsStr = "INF"
		rsiVal = new(big.Float).Copy(hundred)
	} else {
		rs := new(big.Float).Quo(au, ad)
		denom := new(big.Float).Add(one, rs)
		rsiVal = new(big.Float).Sub(hundred, new(big.Float).Quo(hundred, denom))
		rsStr = rs.Text('f', 4)
	}

	firstDate := candles[0].CandleDateTimeKst
	lastDate := candles[len(candles)-1].CandleDateTimeKst
	if len(firstDate) >= 10 {
		firstDate = firstDate[:10]
	}
	if len(lastDate) >= 10 {
		lastDate = lastDate[:10]
	}

	fmt.Printf("  Data range: %s ~ %s\n", firstDate, lastDate)
	fmt.Printf("  Candles used: %d\n", len(candles))
	fmt.Println()
	fmt.Printf("  AU (Avg Gain):    %s\n", au.Text('f', 4))
	fmt.Printf("  AD (Avg Loss):    %s\n", ad.Text('f', 4))
	fmt.Printf("  RS:               %s\n", rsStr)
	fmt.Printf("  RSI:              %s\n", rsiVal.Text('f', 4))
	fmt.Println()

	rsi64, _ := rsiVal.Float64()
	if rsi64 >= 70 {
		fmt.Println("  Interpretation: Overbought zone (RSI >= 70)")
	} else if rsi64 <= 30 {
		fmt.Println("  Interpretation: Oversold zone (RSI <= 30)")
	} else {
		fmt.Println("  Interpretation: Neutral zone (30 < RSI < 70)")
	}

	fmt.Println()
	fmt.Println("  * RSI is an auxiliary indicator and should not be the sole basis for trading decisions.")
}

// ---------------------------------------------------------------------------
// Run
// ---------------------------------------------------------------------------

func main() {
	ctx := context.Background()
	client := upbit.NewClient(option.WithEnvironmentKr())

	top5ByTradeVolume(ctx, client)
	calculateRSI(ctx, client)

	fmt.Printf("\n%s\n", strings.Repeat("=", 60))
	fmt.Println("  Indicators scenario complete")
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println()
	fmt.Println("  ※ This example is for educational purposes.")
	fmt.Println("    We are not responsible for any losses incurred when used for actual investment.")
}
