// 시세 조회 시나리오.
//
// 업비트 Quotation API를 사용하여 시세 정보를 조회하는 예제입니다.
//   - 페어 목록 조회
//   - 현재가(Ticker) 조회 (개별/복수/마켓 단위)
//   - 캔들(OHLCV) 조회 (분/일/주/월봉, to 파라미터)
//   - 최근 체결 내역 조회 (days_ago 포함)
//   - 호가(Orderbook) 조회
//
// Quotation API는 인증이 필요하지 않습니다.
//
// 실행:
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
// 1. 페어 목록 조회
// ---------------------------------------------------------------------------

func listMarkets(ctx context.Context, client upbit.Client) {
	divider("1. 페어 목록 조회")

	markets, err := client.TradingPairs.List(ctx, upbit.TradingPairListParams{})
	if err != nil {
		fmt.Printf("  페어 목록 조회 실패: %v\n", err)
		return
	}
	fmt.Printf("  전체 페어 수: %d\n\n", len(*markets))

	fmt.Printf("  %-12s %-12s %s\n", "마켓", "한글명", "영문명")
	fmt.Printf("  %s\n", strings.Repeat("-", 50))
	for i, m := range *markets {
		if i >= 5 {
			break
		}
		fmt.Printf("  %-12s %-12s %s\n", m.Market, m.KoreanName, m.EnglishName)
	}

	// 유의 종목 필터링 (MarketWarning == "CAUTION")
	marketsDetail, err := client.TradingPairs.List(ctx, upbit.TradingPairListParams{
		IsDetails: param.NewOpt(true),
	})
	if err != nil {
		fmt.Printf("  유의 종목 조회 실패: %v\n", err)
		return
	}
	warnings := 0
	for _, m := range *marketsDetail {
		if m.MarketWarning == "CAUTION" {
			warnings++
		}
	}
	fmt.Printf("\n  유의 종목 수: %d\n", warnings)
}

// ---------------------------------------------------------------------------
// 2. 현재가(Ticker) 조회
// ---------------------------------------------------------------------------

func getTicker(ctx context.Context, client upbit.Client) {
	divider("2. 현재가 조회")

	// 단일 페어
	tickers, err := client.Tickers.ListByTradingPairs(ctx, upbit.TickerListByTradingPairsParams{
		Markets: "KRW-BTC",
	})
	if err != nil {
		fmt.Printf("  현재가 조회 실패: %v\n", err)
		return
	}
	t := (*tickers)[0]
	price, _, _ := new(big.Float).Parse(t.TradePrice, 10)
	accTrade, _, _ := new(big.Float).Parse(t.AccTradePrice24h, 10)
	changeRate, _, _ := new(big.Float).Parse(t.SignedChangeRate, 10)
	pct, _ := new(big.Float).Mul(changeRate, new(big.Float).SetInt64(100)).Float64()

	fmt.Printf("  마켓:         %s\n", t.Market)
	fmt.Printf("  현재가:       %14s KRW\n", formatKRW(price))
	fmt.Printf("  24h 거래대금: %14s KRW\n", formatKRW(accTrade))
	fmt.Printf("  전일 대비:    %+.2f%%\n", pct)

	// 여러 페어 동시 조회
	fmt.Println()
	tickers, err = client.Tickers.ListByTradingPairs(ctx, upbit.TickerListByTradingPairsParams{
		Markets: "KRW-BTC,KRW-ETH",
	})
	if err != nil {
		fmt.Printf("  복수 조회 실패: %v\n", err)
		return
	}
	for _, t := range *tickers {
		p, _, _ := new(big.Float).Parse(t.TradePrice, 10)
		fmt.Printf("  %-12s 현재가: %14s\n", t.Market, formatKRW(p))
	}

	// 마켓 단위(KRW) 전체 조회
	fmt.Println()
	krwTickers, err := client.Tickers.ListByQuoteCurrencies(ctx, upbit.TickerListByQuoteCurrenciesParams{
		QuoteCurrencies: "KRW",
	})
	if err != nil {
		fmt.Printf("  KRW 마켓 조회 실패: %v\n", err)
		return
	}
	fmt.Printf("  KRW 마켓 전체 티커 수: %d\n", len(*krwTickers))
}

// ---------------------------------------------------------------------------
// 3. 캔들 조회
// ---------------------------------------------------------------------------

func getCandles(ctx context.Context, client upbit.Client) {
	divider("3. 캔들 조회")

	// 5분봉
	candles, err := client.Candles.ListMinutes(ctx, 5, upbit.CandleListMinutesParams{
		Market: "KRW-BTC",
		Count:  param.NewOpt[int64](3),
	})
	if err != nil {
		fmt.Printf("  5분봉 조회 실패: %v\n", err)
		return
	}
	fmt.Println("  [5분봉]")
	fmt.Printf("  %-22s %14s %14s %14s\n", "시각(KST)", "시가", "종가", "거래량")
	fmt.Printf("  %s\n", strings.Repeat("-", 68))
	for _, c := range *candles {
		open, _, _ := new(big.Float).Parse(c.OpeningPrice, 10)
		close, _, _ := new(big.Float).Parse(c.TradePrice, 10)
		vol, _, _ := new(big.Float).Parse(c.CandleAccTradeVolume, 10)
		fmt.Printf("  %-22s %14s %14s %14s\n",
			c.CandleDateTimeKst, formatKRW(open), formatKRW(close), vol.Text('f', 4))
	}

	// to 파라미터로 특정 시간대 조회
	fmt.Println()
	candlesTo, err := client.Candles.ListMinutes(ctx, 1, upbit.CandleListMinutesParams{
		Market: "KRW-BTC",
		To:     param.NewOpt("2025-01-01T00:00:00Z"),
		Count:  param.NewOpt[int64](1),
	})
	if err != nil {
		fmt.Printf("  to 파라미터 조회 실패: %v\n", err)
		return
	}
	fmt.Println("  [to 파라미터 — 2025-01-01 이전 1분봉]")
	if len(*candlesTo) > 0 {
		c := (*candlesTo)[0]
		close, _, _ := new(big.Float).Parse(c.TradePrice, 10)
		fmt.Printf("  %s  종가: %14s\n", c.CandleDateTimeKst, formatKRW(close))
	}

	// 일봉
	fmt.Println()
	days, err := client.Candles.ListDays(ctx, upbit.CandleListDaysParams{
		Market: "KRW-BTC",
		Count:  param.NewOpt[int64](3),
	})
	if err != nil {
		fmt.Printf("  일봉 조회 실패: %v\n", err)
		return
	}
	fmt.Println("  [일봉]")
	for _, c := range *days {
		close, _, _ := new(big.Float).Parse(c.TradePrice, 10)
		fmt.Printf("  %s  종가: %14s\n", c.CandleDateTimeKst[:10], formatKRW(close))
	}

	// 주봉
	fmt.Println()
	weeks, err := client.Candles.ListWeeks(ctx, upbit.CandleListWeeksParams{
		Market: "KRW-BTC",
		Count:  param.NewOpt[int64](3),
	})
	if err != nil {
		fmt.Printf("  주봉 조회 실패: %v\n", err)
		return
	}
	fmt.Println("  [주봉]")
	for _, c := range *weeks {
		close, _, _ := new(big.Float).Parse(c.TradePrice, 10)
		fmt.Printf("  %s  종가: %14s\n", c.CandleDateTimeKst[:10], formatKRW(close))
	}

	// 월봉
	fmt.Println()
	months, err := client.Candles.ListMonths(ctx, upbit.CandleListMonthsParams{
		Market: "KRW-BTC",
		Count:  param.NewOpt[int64](3),
	})
	if err != nil {
		fmt.Printf("  월봉 조회 실패: %v\n", err)
		return
	}
	fmt.Println("  [월봉]")
	for _, c := range *months {
		close, _, _ := new(big.Float).Parse(c.TradePrice, 10)
		fmt.Printf("  %s  종가: %14s\n", c.CandleDateTimeKst[:10], formatKRW(close))
	}
}

// ---------------------------------------------------------------------------
// 4. 최근 체결 내역 조회
// ---------------------------------------------------------------------------

func getTrades(ctx context.Context, client upbit.Client) {
	divider("4. 최근 체결 내역")

	trades, err := client.Trades.List(ctx, upbit.TradeListParams{
		Market: "KRW-BTC",
		Count:  param.NewOpt[int64](5),
	})
	if err != nil {
		fmt.Printf("  체결 내역 조회 실패: %v\n", err)
		return
	}
	fmt.Printf("  %-20s %4s %14s %14s\n", "시각(UTC)", "구분", "체결가", "체결량")
	fmt.Printf("  %s\n", strings.Repeat("-", 56))
	for _, t := range *trades {
		side := "매도"
		if t.AskBid == "BID" {
			side = "매수"
		}
		p, _, _ := new(big.Float).Parse(t.TradePrice, 10)
		v, _, _ := new(big.Float).Parse(t.TradeVolume, 10)
		fmt.Printf("  %s %-9s %4s %14s %14s\n",
			t.TradeDateUtc, t.TradeTimeUtc, side, formatKRW(p), v.Text('f', 8))
	}

	// days_ago 파라미터
	fmt.Println()
	tradesAgo, err := client.Trades.List(ctx, upbit.TradeListParams{
		Market:  "KRW-BTC",
		DaysAgo: param.NewOpt[int64](1),
		Count:   param.NewOpt[int64](3),
	})
	if err != nil {
		fmt.Printf("  과거 체결 조회 실패: %v\n", err)
		return
	}
	fmt.Println("  [1일 전 체결 내역]")
	for _, t := range *tradesAgo {
		side := "매도"
		if t.AskBid == "BID" {
			side = "매수"
		}
		p, _, _ := new(big.Float).Parse(t.TradePrice, 10)
		fmt.Printf("  %s %-9s %4s %14s\n", t.TradeDateUtc, t.TradeTimeUtc, side, formatKRW(p))
	}
}

// ---------------------------------------------------------------------------
// 5. 호가 조회
// ---------------------------------------------------------------------------

func getOrderbook(ctx context.Context, client upbit.Client) {
	divider("5. 호가 조회")

	orderbooks, err := client.Orderbooks.List(ctx, upbit.OrderbookListParams{
		Markets: "KRW-BTC",
	})
	if err != nil {
		fmt.Printf("  호가 조회 실패: %v\n", err)
		return
	}
	ob := (*orderbooks)[0]

	totalAsk, _, _ := new(big.Float).Parse(ob.TotalAskSize, 10)
	totalBid, _, _ := new(big.Float).Parse(ob.TotalBidSize, 10)
	fmt.Printf("  총 매도 잔량: %s\n", totalAsk.Text('f', 8))
	fmt.Printf("  총 매수 잔량: %s\n", totalBid.Text('f', 8))
	fmt.Println()

	fmt.Printf("  %14s %14s   %14s %14s\n", "매도호가", "매도잔량", "매수호가", "매수잔량")
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
// 실행
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
	fmt.Println("  시세 조회 시나리오 완료")
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println()
	fmt.Println("  ※ 이 예제는 교육 목적으로 작성되었습니다.")
	fmt.Println("    실제 투자에 사용할 경우 발생하는 손실에 대해 책임지지 않습니다.")
}
