// 지표 산출 시나리오.
//
// 업비트 Quotation API 응답을 활용하여 투자 지표를 산출하는 예제입니다.
//   - 24시간 누적 거래대금 상위 5개 페어
//   - RSI(상대강도지수) 산출
//
// Quotation API는 인증이 필요하지 않습니다.
//
// 실행:
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

// formatKor formats a big.Float value in Korean units (조/억).
func formatKor(v *big.Float) string {
	trillion := new(big.Float).SetInt64(1_0000_0000_0000)
	billion := new(big.Float).SetInt64(1_0000_0000)

	if v.Cmp(trillion) >= 0 {
		q := new(big.Float).Quo(v, trillion)
		f, _ := q.Float64()
		return fmt.Sprintf("약 %.1f조", f)
	}
	if v.Cmp(billion) >= 0 {
		q := new(big.Float).Quo(v, billion)
		f, _ := q.Float64()
		return fmt.Sprintf("약 %.0f억", f)
	}
	return formatKRW(v)
}

// ---------------------------------------------------------------------------
// 1. 24시간 누적 거래대금 상위 5개 페어
// ---------------------------------------------------------------------------

func top5ByTradeVolume(ctx context.Context, client upbit.Client) {
	divider("1. 24시간 누적 거래대금 상위 5개 (KRW 마켓)")

	tickers, err := client.Tickers.ListByQuoteCurrencies(ctx, upbit.TickerListByQuoteCurrenciesParams{
		QuoteCurrencies: "KRW",
	})
	if err != nil {
		fmt.Printf("  티커 조회 실패: %v\n", err)
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

	fmt.Printf("  %-12s %14s %18s  %s\n", "마켓", "현재가", "24h 거래대금", "단위")
	fmt.Printf("  %s\n", strings.Repeat("-", 60))
	for i := 0; i < 5 && i < len(items); i++ {
		t := items[i]
		fmt.Printf("  %-12s %14s %18s  (%s)\n",
			t.Market, formatKRW(t.Price), formatKRW(t.AccTrade), formatKor(t.AccTrade))
	}
}

// ---------------------------------------------------------------------------
// 2. RSI(상대강도지수) 산출
// ---------------------------------------------------------------------------

func calculateRSI(ctx context.Context, client upbit.Client) {
	market := "KRW-BTC"
	period := 14

	divider(fmt.Sprintf("2. RSI 지표 산출 (%s, %d일)", market, period))

	// 일봉 데이터 200개 조회 (최신순 → 오름차순 정렬)
	days, err := client.Candles.ListDays(ctx, upbit.CandleListDaysParams{
		Market: market,
		Count:  param.NewOpt[int64](200),
	})
	if err != nil {
		fmt.Printf("  일봉 조회 실패: %v\n", err)
		return
	}

	candles := *days
	if len(candles) < period {
		fmt.Printf("  오류: 최소 %d개의 캔들이 필요합니다 (현재 %d개)\n", period, len(candles))
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

	// 초기 평균 상승/하락폭 (SMA)
	p := new(big.Float).SetInt64(int64(period))
	au := new(big.Float)
	ad := new(big.Float)
	for i := 0; i < period; i++ {
		au.Add(au, gains[i])
		ad.Add(ad, losses[i])
	}
	au.Quo(au, p)
	ad.Quo(ad, p)

	// Wilder 평활법
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

	fmt.Printf("  데이터 기간: %s ~ %s\n", firstDate, lastDate)
	fmt.Printf("  사용 캔들 수: %d개\n", len(candles))
	fmt.Println()
	fmt.Printf("  AU (평균 상승폭): %s\n", au.Text('f', 4))
	fmt.Printf("  AD (평균 하락폭): %s\n", ad.Text('f', 4))
	fmt.Printf("  RS (상대강도):    %s\n", rsStr)
	fmt.Printf("  RSI:              %s\n", rsiVal.Text('f', 4))
	fmt.Println()

	rsi64, _ := rsiVal.Float64()
	if rsi64 >= 70 {
		fmt.Println("  해석: 과매수 구간 (RSI >= 70)")
	} else if rsi64 <= 30 {
		fmt.Println("  해석: 과매도 구간 (RSI <= 30)")
	} else {
		fmt.Println("  해석: 중립 구간 (30 < RSI < 70)")
	}

	fmt.Println()
	fmt.Println("  * RSI는 보조 지표이며 매매 판단의 절대 근거로 사용할 수 없습니다.")
}

// ---------------------------------------------------------------------------
// 실행
// ---------------------------------------------------------------------------

func main() {
	ctx := context.Background()
	client := upbit.NewClient(option.WithEnvironmentKr())

	top5ByTradeVolume(ctx, client)
	calculateRSI(ctx, client)

	fmt.Printf("\n%s\n", strings.Repeat("=", 60))
	fmt.Println("  지표 산출 시나리오 완료")
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println()
	fmt.Println("  ※ 이 예제는 교육 목적으로 작성되었습니다.")
	fmt.Println("    실제 투자에 사용할 경우 발생하는 손실에 대해 책임지지 않습니다.")
}
