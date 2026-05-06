// Quotation scenario.
//
// Demonstrates how to query market data using the Upbit Quotation API.
//   - List trading pairs
//   - Get ticker (current price) — single/multiple/by quote
//   - Get candles (OHLCV) — minutes/days/weeks/months, to parameter
//   - Get recent trades (including days_ago)
//   - Get orderbook
//
// The Quotation API does not require authentication.
//
// Usage:
//
//	go run examples/quotation.go
package main

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	upbit "github.com/upbit-official/upbit-sdk-go"
	"github.com/upbit-official/upbit-sdk-go/option"
	"github.com/upbit-official/upbit-sdk-go/packages/param"
)

func divider(title string) {
	fmt.Printf("\n%s\n  %s\n%s\n", strings.Repeat("=", 60), title, strings.Repeat("=", 60))
}

// formatKRW formats a big.Float as a comma-separated integer string.
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

// ---------------------------------------------------------------------------
// 1. List trading pairs
// ---------------------------------------------------------------------------

func listMarkets(ctx context.Context, client upbit.Client) {
	divider("1. List Trading Pairs")

	markets, err := client.TradingPairs.List(ctx, upbit.TradingPairListParams{})
	if err != nil {
		fmt.Printf("  Failed to list trading pairs: %v\n", err)
		return
	}
	fmt.Printf("  Total trading pairs: %d\n\n", len(*markets))

	fmt.Printf("  %-12s %-12s %s\n", "Market", "Korean Name", "English Name")
	fmt.Printf("  %s\n", strings.Repeat("-", 50))
	for i, m := range *markets {
		if i >= 5 {
			break
		}
		fmt.Printf("  %-12s %-12s %s\n", m.Market, m.KoreanName, m.EnglishName)
	}

	// Filter caution-flagged markets (MarketWarning == "CAUTION")
	marketsDetail, err := client.TradingPairs.List(ctx, upbit.TradingPairListParams{
		IsDetails: param.NewOpt(true),
	})
	if err != nil {
		fmt.Printf("  Failed to list caution markets: %v\n", err)
		return
	}
	warnings := 0
	for _, m := range *marketsDetail {
		if m.MarketWarning == "CAUTION" {
			warnings++
		}
	}
	fmt.Printf("\n  Caution markets: %d\n", warnings)
}

// ---------------------------------------------------------------------------
// 2. Get ticker (current price)
// ---------------------------------------------------------------------------

func getTicker(ctx context.Context, client upbit.Client) {
	divider("2. Get Ticker (Current Price)")

	// Single trading pair
	tickers, err := client.Tickers.ListByTradingPairs(ctx, upbit.TickerListByTradingPairsParams{
		Markets: "KRW-BTC",
	})
	if err != nil {
		fmt.Printf("  Failed to fetch ticker: %v\n", err)
		return
	}
	t := (*tickers)[0]
	price, _, _ := new(big.Float).Parse(t.TradePrice, 10)
	accTrade, _, _ := new(big.Float).Parse(t.AccTradePrice24h, 10)
	changeRate, _, _ := new(big.Float).Parse(t.SignedChangeRate, 10)
	pct, _ := new(big.Float).Mul(changeRate, new(big.Float).SetInt64(100)).Float64()

	fmt.Printf("  Market:       %s\n", t.Market)
	fmt.Printf("  Price:        %14s KRW\n", formatKRW(price))
	fmt.Printf("  24h Volume:   %14s KRW\n", formatKRW(accTrade))
	fmt.Printf("  24h Change:   %+.2f%%\n", pct)

	// Multiple pairs at once
	fmt.Println()
	tickers, err = client.Tickers.ListByTradingPairs(ctx, upbit.TickerListByTradingPairsParams{
		Markets: "KRW-BTC,KRW-ETH",
	})
	if err != nil {
		fmt.Printf("  Failed to fetch multiple tickers: %v\n", err)
		return
	}
	for _, t := range *tickers {
		p, _, _ := new(big.Float).Parse(t.TradePrice, 10)
		fmt.Printf("  %-12s Price: %14s\n", t.Market, formatKRW(p))
	}

	// All tickers by quote currency (KRW)
	fmt.Println()
	krwTickers, err := client.Tickers.ListByQuoteCurrencies(ctx, upbit.TickerListByQuoteCurrenciesParams{
		QuoteCurrencies: "KRW",
	})
	if err != nil {
		fmt.Printf("  Failed to fetch KRW market tickers: %v\n", err)
		return
	}
	fmt.Printf("  KRW market tickers: %d\n", len(*krwTickers))
}

// ---------------------------------------------------------------------------
// 3. Get candles
// ---------------------------------------------------------------------------

func getCandles(ctx context.Context, client upbit.Client) {
	divider("3. Get Candles")

	// 5-minute candles
	candles, err := client.Candles.ListMinutes(ctx, 5, upbit.CandleListMinutesParams{
		Market: "KRW-BTC",
		Count:  param.NewOpt[int64](3),
	})
	if err != nil {
		fmt.Printf("  Failed to fetch 5-minute candles: %v\n", err)
		return
	}
	fmt.Println("  [5-min candles]")
	fmt.Printf("  %-22s %14s %14s %14s\n", "Time (KST)", "Open", "Close", "Volume")
	fmt.Printf("  %s\n", strings.Repeat("-", 68))
	for _, c := range *candles {
		open, _, _ := new(big.Float).Parse(c.OpeningPrice, 10)
		close, _, _ := new(big.Float).Parse(c.TradePrice, 10)
		vol, _, _ := new(big.Float).Parse(c.CandleAccTradeVolume, 10)
		fmt.Printf("  %-22s %14s %14s %14s\n",
			c.CandleDateTimeKst, formatKRW(open), formatKRW(close), vol.Text('f', 4))
	}

	// Using 'to' parameter for a specific time range
	fmt.Println()
	candlesTo, err := client.Candles.ListMinutes(ctx, 1, upbit.CandleListMinutesParams{
		Market: "KRW-BTC",
		To:     param.NewOpt("2025-01-01T00:00:00Z"),
		Count:  param.NewOpt[int64](1),
	})
	if err != nil {
		fmt.Printf("  Failed to fetch with 'to' parameter: %v\n", err)
		return
	}
	fmt.Println("  [to parameter — 1-min candle before 2025-01-01]")
	if len(*candlesTo) > 0 {
		c := (*candlesTo)[0]
		close, _, _ := new(big.Float).Parse(c.TradePrice, 10)
		fmt.Printf("  %s  Close: %14s\n", c.CandleDateTimeKst, formatKRW(close))
	}

	// Daily candles
	fmt.Println()
	days, err := client.Candles.ListDays(ctx, upbit.CandleListDaysParams{
		Market: "KRW-BTC",
		Count:  param.NewOpt[int64](3),
	})
	if err != nil {
		fmt.Printf("  Failed to fetch daily candles: %v\n", err)
		return
	}
	fmt.Println("  [Daily candles]")
	for _, c := range *days {
		close, _, _ := new(big.Float).Parse(c.TradePrice, 10)
		fmt.Printf("  %s  Close: %14s\n", c.CandleDateTimeKst[:10], formatKRW(close))
	}

	// Weekly candles
	fmt.Println()
	weeks, err := client.Candles.ListWeeks(ctx, upbit.CandleListWeeksParams{
		Market: "KRW-BTC",
		Count:  param.NewOpt[int64](3),
	})
	if err != nil {
		fmt.Printf("  Failed to fetch weekly candles: %v\n", err)
		return
	}
	fmt.Println("  [Weekly candles]")
	for _, c := range *weeks {
		close, _, _ := new(big.Float).Parse(c.TradePrice, 10)
		fmt.Printf("  %s  Close: %14s\n", c.CandleDateTimeKst[:10], formatKRW(close))
	}

	// Monthly candles
	fmt.Println()
	months, err := client.Candles.ListMonths(ctx, upbit.CandleListMonthsParams{
		Market: "KRW-BTC",
		Count:  param.NewOpt[int64](3),
	})
	if err != nil {
		fmt.Printf("  Failed to fetch monthly candles: %v\n", err)
		return
	}
	fmt.Println("  [Monthly candles]")
	for _, c := range *months {
		close, _, _ := new(big.Float).Parse(c.TradePrice, 10)
		fmt.Printf("  %s  Close: %14s\n", c.CandleDateTimeKst[:10], formatKRW(close))
	}
}

// ---------------------------------------------------------------------------
// 4. Recent trade history
// ---------------------------------------------------------------------------

func getTrades(ctx context.Context, client upbit.Client) {
	divider("4. Recent Trades")

	trades, err := client.Trades.List(ctx, upbit.TradeListParams{
		Market: "KRW-BTC",
		Count:  param.NewOpt[int64](5),
	})
	if err != nil {
		fmt.Printf("  Failed to fetch trade history: %v\n", err)
		return
	}
	fmt.Printf("  %-20s %4s %14s %14s\n", "Time (UTC)", "Side", "Price", "Volume")
	fmt.Printf("  %s\n", strings.Repeat("-", 56))
	for _, t := range *trades {
		side := "SELL"
		if t.AskBid == "BID" {
			side = "BUY"
		}
		p, _, _ := new(big.Float).Parse(t.TradePrice, 10)
		v, _, _ := new(big.Float).Parse(t.TradeVolume, 10)
		fmt.Printf("  %s %-9s %4s %14s %14s\n",
			t.TradeDateUtc, t.TradeTimeUtc, side, formatKRW(p), v.Text('f', 8))
	}

	// days_ago parameter
	fmt.Println()
	tradesAgo, err := client.Trades.List(ctx, upbit.TradeListParams{
		Market:  "KRW-BTC",
		DaysAgo: param.NewOpt[int64](1),
		Count:   param.NewOpt[int64](3),
	})
	if err != nil {
		fmt.Printf("  Failed to fetch past trades: %v\n", err)
		return
	}
	fmt.Println("  [Trades from 1 day ago]")
	for _, t := range *tradesAgo {
		side := "SELL"
		if t.AskBid == "BID" {
			side = "BUY"
		}
		p, _, _ := new(big.Float).Parse(t.TradePrice, 10)
		fmt.Printf("  %s %-9s %4s %14s\n", t.TradeDateUtc, t.TradeTimeUtc, side, formatKRW(p))
	}
}

// ---------------------------------------------------------------------------
// 5. Orderbook
// ---------------------------------------------------------------------------

func getOrderbook(ctx context.Context, client upbit.Client) {
	divider("5. Orderbook")

	orderbooks, err := client.Orderbooks.List(ctx, upbit.OrderbookListParams{
		Markets: "KRW-BTC",
	})
	if err != nil {
		fmt.Printf("  Failed to fetch orderbook: %v\n", err)
		return
	}
	ob := (*orderbooks)[0]

	totalAsk, _, _ := new(big.Float).Parse(ob.TotalAskSize, 10)
	totalBid, _, _ := new(big.Float).Parse(ob.TotalBidSize, 10)
	fmt.Printf("  Total ask size: %s\n", totalAsk.Text('f', 8))
	fmt.Printf("  Total bid size: %s\n", totalBid.Text('f', 8))
	fmt.Println()

	fmt.Printf("  %14s %14s   %14s %14s\n", "Ask Price", "Ask Size", "Bid Price", "Bid Size")
	fmt.Printf("  %s\n", strings.Repeat("-", 62))
	for i, u := range ob.OrderbookUnits {
		if i >= 5 {
			break
		}
		askP, _, _ := new(big.Float).Parse(u.AskPrice, 10)
		askS, _, _ := new(big.Float).Parse(u.AskSize, 10)
		bidP, _, _ := new(big.Float).Parse(u.BidPrice, 10)
		bidS, _, _ := new(big.Float).Parse(u.BidSize, 10)
		fmt.Printf("  %14s %14s   %14s %14s\n",
			formatKRW(askP), askS.Text('f', 8),
			formatKRW(bidP), bidS.Text('f', 8))
	}
}

// ---------------------------------------------------------------------------
// Run
// ---------------------------------------------------------------------------

func main() {
	ctx := context.Background()
	client := upbit.NewClient(option.WithEnvironmentKr())

	listMarkets(ctx, client)
	getTicker(ctx, client)
	getCandles(ctx, client)
	getTrades(ctx, client)
	getOrderbook(ctx, client)

	fmt.Printf("\n%s\n", strings.Repeat("=", 60))
	fmt.Println("  Quotation scenario complete")
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println()
	fmt.Println("  ※ This example is for educational purposes.")
	fmt.Println("    We are not responsible for any losses incurred when used for actual investment.")
}
