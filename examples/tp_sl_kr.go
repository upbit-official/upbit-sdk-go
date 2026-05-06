// TP/SL (익절/손절) 자동 매도 시나리오.
//
// 업비트 Exchange API를 사용하여 현재가를 폴링하고,
// 목표가(익절) 또는 손절가에 도달하면 시장가 매도를 실행하는 예제입니다.
//   - 현재가 조회 및 TP/SL 가격 계산
//   - 가격 모니터링 (REST 폴링)
//   - 시장가 매도 실행
//
// 업비트 API는 예약주문(스톱/리밋)을 지원하지 않으므로,
// 클라이언트 사이드에서 폴링 방식으로 구현합니다.
//
// Dry run (기본):
//
//	현재가 기준, 주문 실행 안 함. 동작 확인용으로 바로 실행할 수 있습니다.
//
// 실제 실행 (DRY_RUN=false):
//
//	평균 매수가 기준, 실제 매도 주문 실행. 반드시 보유 자산이 있어야 합니다.
//
// 실행:
//
//	# Dry run (기본) — 모니터링만, 주문 없음
//	UPBIT_ACCESS_KEY=<key> UPBIT_SECRET_KEY=<secret> go run examples/tp_sl.go
//
//	# 실제 실행 — 평균 매수가 기준, 실제 매도
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
	tpPercent = new(big.Float).SetInt64(3) // 익절 기준: +3%
	slPercent = new(big.Float).SetInt64(2) // 손절 기준: -2%
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
		log.Fatalf("현재가 조회 실패: %v", err)
	}
	price, _, _ := new(big.Float).Parse((*tickers)[0].TradePrice, 10)
	return price
}

// ---------------------------------------------------------------------------
// 1. 사전 점검 + TP/SL 설정
// ---------------------------------------------------------------------------

func setupTPSL(ctx context.Context, client upbit.Client) (*big.Float, *big.Float) {
	modeLabel := "DRY RUN"
	if !dryRun {
		modeLabel = "실제 실행"
	}
	divider(fmt.Sprintf("1. TP/SL 설정 [%s]", modeLabel))

	current := fetchCurrentPrice(ctx, client)
	hundred := new(big.Float).SetInt64(100)
	one := new(big.Float).SetInt64(1)

	var basePrice *big.Float

	if dryRun {
		basePrice = current
		fmt.Printf("  기준가:     %14s KRW (현재가)\n", formatKRW(basePrice))
	} else {
		coin := strings.Split(market, "-")[1]
		accounts, err := client.Accounts.List(ctx)
		if err != nil {
			log.Fatalf("[오류] 계좌 조회 실패 (API 키를 확인하세요): %v", err)
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
			log.Fatalf("[오류] %s 보유량이 없습니다.", coin)
		}

		sv, _, _ := new(big.Float).Parse(sellVolume, 10)
		if balance.Cmp(sv) < 0 {
			log.Fatalf("[오류] %s 보유량(%s)이 매도 수량(%s)보다 적습니다.",
				coin, balance.Text('f', 8), sellVolume)
		}

		basePrice = avgBuyPrice
		fmt.Printf("  잔고:       %s %s\n", coin, balance.Text('f', 8))
		fmt.Printf("  기준가:     %14s KRW (평균 매수가)\n", formatKRW(basePrice))
		fmt.Printf("  현재가:     %14s KRW\n", formatKRW(current))
	}

	// tp_price = base_price * (1 + tp_percent / 100)
	tpMul := new(big.Float).Quo(tpPercent, hundred)
	tpMul.Add(tpMul, one)
	tpPrice := new(big.Float).Mul(basePrice, tpMul)

	// sl_price = base_price * (1 - sl_percent / 100)
	slMul := new(big.Float).Quo(slPercent, hundred)
	slMul.Sub(one, slMul)
	slPrice := new(big.Float).Mul(basePrice, slMul)

	fmt.Printf("  마켓:       %s\n", market)
	fmt.Printf("  매도 수량:  %s\n", sellVolume)
	fmt.Printf("  익절 가격:  %14s KRW (+%s%%)\n", formatKRW(tpPrice), tpPercent.Text('f', 0))
	fmt.Printf("  손절 가격:  %14s KRW (-%s%%)\n", formatKRW(slPrice), slPercent.Text('f', 0))
	fmt.Printf("  폴링 간격:  %s × 최대 %d회\n", pollInterval, maxPolls)
	if dryRun {
		fmt.Println("\n  ※ Dry run 모드: 주문이 실행되지 않습니다.")
	}

	return tpPrice, slPrice
}

// ---------------------------------------------------------------------------
// 2. 가격 모니터링
// ---------------------------------------------------------------------------

func monitorPrice(ctx context.Context, client upbit.Client, tpPrice, slPrice *big.Float) string {
	divider("2. 가격 모니터링 시작")

	for poll := 1; poll <= maxPolls; poll++ {
		current := fetchCurrentPrice(ctx, client)

		fmt.Printf("  [%3d/%d] 현재가: %14s KRW", poll, maxPolls, formatKRW(current))

		if current.Cmp(tpPrice) >= 0 {
			fmt.Println("  → 익절 도달!")
			return "tp"
		}
		if current.Cmp(slPrice) <= 0 {
			fmt.Println("  → 손절 도달!")
			return "sl"
		}

		fmt.Println()

		if poll < maxPolls {
			time.Sleep(pollInterval)
		}
	}

	fmt.Printf("\n  최대 폴링 횟수(%d회)에 도달했습니다.\n", maxPolls)
	return "timeout"
}

// ---------------------------------------------------------------------------
// 3. 시장가 매도 실행
// ---------------------------------------------------------------------------

func executeSell(ctx context.Context, client upbit.Client, trigger string) {
	label := "익절"
	if trigger == "sl" {
		label = "손절"
	}
	divider(fmt.Sprintf("3. 매도 실행 (트리거: %s)", label))

	if dryRun {
		fmt.Println("  [건너뜀] Dry run 모드이므로 실제 주문을 생략합니다.")
		return
	}

	order, err := client.Orders.New(ctx, upbit.OrderNewParams{
		Market:  market,
		Side:    upbit.SideEnumAsk,
		OrdType: upbit.OrdTypeEnumMarket,
		Volume:  param.NewOpt(sellVolume),
	})
	if err != nil {
		log.Fatalf("매도 주문 생성 실패: %v", err)
	}
	fmt.Println("  주문 생성 완료")
	fmt.Printf("  UUID:   %s\n", order.Uuid)
	fmt.Printf("  상태:   %s\n", order.State)

	// 체결 대기
	var info *upbit.Order
	for range 20 {
		info, err = client.Orders.Get(ctx, upbit.OrderGetParams{Uuid: param.NewOpt(order.Uuid)})
		if err != nil {
			log.Fatalf("주문 조회 실패: %v", err)
		}
		if info.State == "done" || info.State == "cancel" {
			break
		}
		time.Sleep(500 * time.Millisecond)
	}

	if info != nil {
		fmt.Printf("  체결 상태: %s\n", info.State)
		fmt.Printf("  체결 수량: %s\n", info.ExecutedVolume)
		if len(info.Trades) > 0 {
			total := new(big.Float)
			for _, t := range info.Trades {
				if t.Funds != "" {
					f, _, _ := new(big.Float).Parse(t.Funds, 10)
					total.Add(total, f)
				}
			}
			fmt.Printf("  체결 금액: %14s KRW\n", formatKRW(total))
		}
	}
}

// ---------------------------------------------------------------------------
// 실행
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
		fmt.Println("\n  시간 초과로 매도를 실행하지 않습니다.")
	} else {
		executeSell(ctx, client, trigger)
	}

	fmt.Printf("\n%s\n", strings.Repeat("=", 60))
	fmt.Println("  TP/SL 시나리오 완료")
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println()
	fmt.Println("  ※ 이 예제는 교육 목적으로 작성되었습니다.")
	fmt.Println("    실제 투자에 사용할 경우 발생하는 손실에 대해 책임지지 않습니다.")
}
