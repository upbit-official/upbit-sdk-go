// DCA (정기 매수) 자동 시장가 매수 시나리오.
//
// 업비트 Exchange API를 사용하여 일정 금액을 주기적으로 시장가 매수하는 예제입니다.
//   - 현재가 조회 및 매수 설정 확인
//   - 지정 횟수만큼 시장가 매수 반복
//   - 매수 결과 요약 (총 매수금액, 총 매수수량, 평균 단가)
//
// * 서버 안내: DCA 전략은 봇을 지속적으로 실행해야 하므로
//
//	서버(또는 스케줄러)는 사용자가 직접 구성해야 합니다.
//	이 예제는 서버 환경이 갖추어져 있음을 전제로 작성되었습니다.
//
// * 주문 접수 안내: 이 코드의 실행이 주문의 접수를 보장하지 않습니다.
//
//	점검, 네트워크 오류 등으로 주문이 거절될 수 있으며,
//	주문 접수 오류에 대한 별도의 처리가 필요합니다.
//
// Dry run (기본):
//
//	현재가 조회만 수행하며, 주문을 실행하지 않습니다.
//	동작 확인용으로 바로 실행할 수 있습니다.
//
// 실제 실행 (DRY_RUN=false):
//
//	잔고 확인 후 실제 시장가 매수를 실행합니다.
//	충분한 KRW 잔고가 있어야 합니다.
//
// 실행:
//
//	# Dry run (기본) — 모니터링만, 주문 없음
//	UPBIT_ACCESS_KEY=<key> UPBIT_SECRET_KEY=<secret> go run examples/dca.go
//
//	# 실제 실행 — 실제 시장가 매수
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
		log.Fatalf("현재가 조회 실패: %v", err)
	}
	price, _, _ := new(big.Float).Parse((*tickers)[0].TradePrice, 10)
	return price
}

// ---------------------------------------------------------------------------
// 1. 사전 점검 + DCA 설정
// ---------------------------------------------------------------------------

func setup(ctx context.Context, client upbit.Client) {
	modeLabel := "DRY RUN"
	if !dryRun {
		modeLabel = "실제 실행"
	}
	divider(fmt.Sprintf("1. DCA 설정 [%s]", modeLabel))

	current := currentPrice(ctx, client)
	buyAmt, _, _ := new(big.Float).Parse(buyAmount, 10)
	totalNeeded := new(big.Float).Mul(
		new(big.Float).SetInt64(int64(totalRound)),
		buyAmt,
	)

	if dryRun {
		fmt.Printf("  현재가:       %14s KRW\n", formatKRW(current))
	} else {
		fiat := strings.Split(market, "-")[0]
		accounts, err := client.Accounts.List(ctx)
		if err != nil {
			log.Fatalf("[오류] 계좌 조회 실패 (API 키를 확인하세요): %v", err)
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
			log.Fatalf("[오류] %s 잔고(%s)가 총 필요 금액(%s)보다 적습니다.",
				fiat, formatKRW(available), formatKRW(totalNeeded))
		}

		fmt.Printf("  현재가:       %14s KRW\n", formatKRW(current))
		fmt.Printf("  %s 잔고:  %14s KRW\n", fiat, formatKRW(available))
	}

	fmt.Printf("  마켓:         %s\n", market)
	fmt.Printf("  1회 매수금액: %14s KRW\n", formatKRW(buyAmt))
	fmt.Printf("  매수 횟수:    %d회\n", totalRound)
	fmt.Printf("  매수 간격:    %s\n", interval)
	fmt.Printf("  총 필요금액:  %14s KRW\n", formatKRW(totalNeeded))
	if dryRun {
		fmt.Println("\n  ※ Dry run 모드: 주문이 실행되지 않습니다.")
	}
}

// ---------------------------------------------------------------------------
// 2. DCA 매수 실행
// ---------------------------------------------------------------------------

func executeDCA(ctx context.Context, client upbit.Client) []buyResult {
	divider("2. DCA 매수 실행")

	if dryRun {
		fmt.Println("  [건너뜀] Dry run 모드이므로 실제 주문을 생략합니다.")
		fmt.Println()
		buyAmt, _, _ := new(big.Float).Parse(buyAmount, 10)
		for r := 1; r <= totalRound; r++ {
			current := currentPrice(ctx, client)
			estVolume := new(big.Float).Quo(buyAmt, current)
			fmt.Printf("  [%d/%d] 현재가: %14s KRW  (예상 수량: %s)\n",
				r, totalRound, formatKRW(current), estVolume.Text('f', 8))
			if r < totalRound {
				time.Sleep(interval)
			}
		}
		return nil
	}

	results := make([]buyResult, 0, totalRound)
	for r := 1; r <= totalRound; r++ {
		fmt.Printf("\n  [%d/%d] 시장가 매수 (%s KRW) ...\n", r, totalRound, buyAmount)

		result := executeSingleBuy(ctx, client, r)
		results = append(results, result)

		if r < totalRound {
			if result.success {
				fmt.Printf("  %s 후 다음 매수...\n", interval)
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
		fmt.Printf("  [오류] 주문 처리 중 에러: %v\n", err)
		return buyResult{round: r, spent: zero, volume: zero, success: false}
	}
	fmt.Printf("  UUID:   %s\n", order.Uuid)
	fmt.Printf("  상태:   %s\n", order.State)

	// 체결 대기
	var info *upbit.Order
	for range 20 {
		info, err = client.Orders.Get(ctx, upbit.OrderGetParams{Uuid: param.NewOpt(order.Uuid)})
		if err != nil {
			fmt.Printf("  [오류] 주문 조회 에러: %v\n", err)
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
		fmt.Printf("  체결 금액: %14s KRW\n", formatKRW(spent))
		fmt.Printf("  체결 수량: %s\n", volume.Text('f', 8))
		fmt.Printf("  체결 단가: %14s KRW\n", formatKRW(avgPrice))
		return buyResult{round: r, spent: spent, volume: volume, success: true}
	}

	if info != nil {
		fmt.Printf("  체결 상태: %s\n", info.State)
	}
	return buyResult{round: r, spent: zero, volume: zero, success: false}
}

// ---------------------------------------------------------------------------
// 3. 매수 결과 요약
// ---------------------------------------------------------------------------

func showSummary(results []buyResult) {
	divider("3. 매수 결과 요약")

	if len(results) == 0 {
		fmt.Println("  결과 없음 (Dry run 모드)")
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
		fmt.Printf("  성공 횟수:     %d/%d\n", successes, len(results))
		fmt.Printf("  총 매수금액:   %14s KRW\n", formatKRW(totalSpent))
		fmt.Printf("  총 매수수량:   %s\n", totalVolume.Text('f', 8))
		fmt.Printf("  평균 매수단가: %14s KRW\n", formatKRW(avgPrice))
	} else {
		fmt.Println("  성공한 매수가 없습니다.")
	}

	if failures > 0 {
		fmt.Printf("  실패 횟수:     %d\n", failures)
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

	setup(ctx, client)
	results := executeDCA(ctx, client)
	showSummary(results)

	fmt.Printf("\n%s\n", strings.Repeat("=", 60))
	fmt.Println("  DCA 시나리오 완료")
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println()
	fmt.Println("  ※ 이 예제는 교육 목적으로 작성되었습니다.")
	fmt.Println("    실제 투자에 사용할 경우 발생하는 손실에 대해 책임지지 않습니다.")
}
