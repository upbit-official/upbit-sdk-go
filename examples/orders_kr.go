// 주문 생성 및 관리 시나리오.
//
// 업비트 Exchange API를 사용하여 주문을 생성하고 관리하는 예제입니다.
//   - 주문 가능 정보 확인
//   - 지정가 매수 주문 생성 → 조회 → 취소
//   - 시장가 매도 가능 여부 확인
//   - 완료 주문 목록 조회
//
// Dry run (기본):
//
//	주문 관련 정보 조회만 수행하며, 실제 주문을 생성하지 않습니다.
//
// 실제 실행 (DRY_RUN=false):
//
//	지정가 매수 주문 생성 → 조회 → 취소를 실행합니다.
//
// 실행:
//
//	# Dry run (기본)
//	UPBIT_ACCESS_KEY=<key> UPBIT_SECRET_KEY=<secret> go run examples/orders.go
//
//	# 실제 실행
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
// 1. 주문 가능 정보 확인
// ---------------------------------------------------------------------------

func checkOrderChance(ctx context.Context, client upbit.Client) {
	divider("1. 주문 가능 정보 확인")

	chance, err := client.Orders.GetChance(ctx, upbit.OrderGetChanceParams{
		Market: market,
	})
	if err != nil {
		log.Fatalf("[오류] 주문 가능 정보 조회 실패: %v", err)
	}

	fmt.Printf("  마켓:          %s\n", market)
	fmt.Printf("  매수 수수료:   %s\n", chance.BidFee)
	fmt.Printf("  매도 수수료:   %s\n", chance.AskFee)
	fmt.Printf("  매수 주문유형: %v\n", chance.Market.BidTypes)
	fmt.Printf("  매도 주문유형: %v\n", chance.Market.AskTypes)
	fmt.Printf("  최소 매수금액: %s\n", chance.Market.Bid.MinTotal)
	fmt.Printf("  최소 매도금액: %s\n", chance.Market.Ask.MinTotal)
}

// ---------------------------------------------------------------------------
// 2. 지정가 매수 → 조회 → 취소
// ---------------------------------------------------------------------------

func limitBidAndCancel(ctx context.Context, client upbit.Client) {
	divider("2. 지정가 매수 → 조회 → 취소")

	// 호가 단위(tick_size) 조회
	// WARNING: 이 예제는 주문 시마다 호가 단위를 API로 조회합니다.
	// 실 서비스에서는 매 주문마다 호가 단위를 조회하는 것은 권장하지 않습니다.
	// 정확한 호가 단위 정책은 업비트 공식 문서를 참고하세요.
	instruments, err := client.Orderbooks.ListInstruments(ctx, upbit.OrderbookListInstrumentsParams{
		Markets: market,
	})
	if err != nil {
		log.Fatalf("호가 단위 조회 실패: %v", err)
	}
	tickSize, _, _ := new(big.Float).Parse((*instruments)[0].TickSize, 10)
	fmt.Printf("  호가 단위:     %s\n", (*instruments)[0].TickSize)

	// 최고 매수호가 조회
	orderbooks, err := client.Orderbooks.List(ctx, upbit.OrderbookListParams{Markets: market})
	if err != nil {
		log.Fatalf("호가 조회 실패: %v", err)
	}
	bestBid, _, _ := new(big.Float).Parse((*orderbooks)[0].OrderbookUnits[0].BidPrice, 10)
	bestBidInt, _ := bestBid.Int64()
	fmt.Printf("  최고 매수호가: %s KRW\n", formatInt64WithComma(bestBidInt))

	// 3% 낮은 가격으로 호가 단위 맞춤
	targetPrice := roundPrice(
		new(big.Float).Mul(bestBid, new(big.Float).SetFloat64(0.97)),
		tickSize,
	)
	volume := "0.0001"
	tp, _, _ := new(big.Float).Parse(targetPrice, 10)
	tpInt, _ := tp.Int64()
	fmt.Printf("  주문 가격:     %s KRW (3%% 할인)\n", formatInt64WithComma(tpInt))
	fmt.Printf("  주문 수량:     %s\n", volume)

	if dryRun {
		fmt.Println("\n  [건너뜀] Dry run 모드이므로 실제 주문을 생략합니다.")
		return
	}

	// 주문 생성
	order, err := client.Orders.New(ctx, upbit.OrderNewParams{
		Market:  market,
		Side:    upbit.SideEnumBid,
		OrdType: upbit.OrdTypeEnumLimit,
		Price:   param.NewOpt(targetPrice),
		Volume:  param.NewOpt(volume),
	})
	if err != nil {
		log.Fatalf("주문 생성 실패: %v", err)
	}
	fmt.Println("\n  주문 생성 완료")
	fmt.Printf("  UUID:   %s\n", order.Uuid)
	fmt.Printf("  상태:   %s\n", order.State)

	// 주문 조회
	info, err := client.Orders.Get(ctx, upbit.OrderGetParams{Uuid: param.NewOpt(order.Uuid)})
	if err != nil {
		log.Fatalf("주문 조회 실패: %v", err)
	}
	fmt.Printf("  조회:   상태=%s, 체결수량=%s\n", info.State, info.ExecutedVolume)

	// 주문 취소
	_, err = client.Orders.Cancel(ctx, upbit.OrderCancelParams{Uuid: param.NewOpt(order.Uuid)})
	if err != nil {
		log.Fatalf("주문 취소 실패: %v", err)
	}
	var cancelled *upbit.Order
	for range 10 {
		cancelled, err = client.Orders.Get(ctx, upbit.OrderGetParams{Uuid: param.NewOpt(order.Uuid)})
		if err != nil {
			log.Fatalf("주문 조회 실패: %v", err)
		}
		if cancelled.State == "done" || cancelled.State == "cancel" {
			break
		}
		time.Sleep(500 * time.Millisecond)
	}
	if cancelled != nil {
		fmt.Printf("  취소:   상태=%s\n", cancelled.State)
	}
}

// ---------------------------------------------------------------------------
// 3. 보유 자산 확인 + 시장가 매도 가능 여부
// ---------------------------------------------------------------------------

func checkMarketSell(ctx context.Context, client upbit.Client) {
	divider("3. 보유 자산 확인 + 시장가 매도 가능 여부")

	coin := strings.Split(market, "-")[1]

	accounts, err := client.Accounts.List(ctx)
	if err != nil {
		log.Fatalf("[오류] 계좌 조회 실패: %v", err)
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
			fmt.Printf("  보유 자산:   %s\n", coin)
			fmt.Printf("  보유 수량:   %s\n", balance.Text('f', 8))
			fmt.Printf("  평균 매수가: %s KRW\n", formatInt64WithComma(avgInt))
			found = true
			break
		}
	}
	if !found {
		fmt.Printf("  %s 보유량이 없습니다.\n", coin)
	}

	// 시장가 매도 가능 확인
	chance, err := client.Orders.GetChance(ctx, upbit.OrderGetChanceParams{Market: market})
	if err != nil {
		log.Fatalf("주문 가능 정보 조회 실패: %v", err)
	}
	hasMarket := false
	for _, t := range chance.Market.AskTypes {
		if string(t) == "market" {
			hasMarket = true
			break
		}
	}
	if hasMarket {
		fmt.Println("  시장가 매도: 가능")
	} else {
		fmt.Printf("  시장가 매도: 불가능 (ask_types=%v)\n", chance.Market.AskTypes)
	}
}

// ---------------------------------------------------------------------------
// 4. 완료 주문 목록 조회
// ---------------------------------------------------------------------------

func listClosedOrders(ctx context.Context, client upbit.Client) {
	divider("4. 완료 주문 목록 조회")

	closed, err := client.Orders.ListClosed(ctx, upbit.OrderListClosedParams{
		Market: param.NewOpt(market),
		Limit:  param.NewOpt[int64](5),
	})
	if err != nil {
		log.Fatalf("완료 주문 조회 실패: %v", err)
	}
	if len(*closed) == 0 {
		fmt.Println("  완료 주문이 없습니다.")
		return
	}

	fmt.Printf("  %10s %4s %8s %6s %14s\n", "UUID", "구분", "유형", "상태", "체결수량")
	fmt.Printf("  %s\n", strings.Repeat("-", 50))
	for _, o := range *closed {
		side := "매도"
		if o.Side == "bid" {
			side = "매수"
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
// 실행
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
	fmt.Println("  주문 시나리오 완료")
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println()
	fmt.Println("  ※ 이 예제는 교육 목적으로 작성되었습니다.")
	fmt.Println("    실제 투자에 사용할 경우 발생하는 손실에 대해 책임지지 않습니다.")
}
