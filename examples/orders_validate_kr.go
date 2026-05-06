// 주문 생성 테스트 시나리오.
//
// 주문 생성 테스트 API를 사용하여 실제 주문 없이 다양한 주문 유형을 검증하는 예제입니다.
//   - 주문 가능 정보 확인
//   - 지정가 매수/매도 테스트
//   - 시장가 매수/매도 테스트
//   - 최유리 매수/매도 테스트 (IOC)
//   - 잘못된 주문 요청 검증
//
// 주문 생성 테스트 API는 실제 주문과 동일한 검증 과정을 거치지만,
// 주문이 실제로 생성되지 않으므로 수수료 없이 안전하게 테스트할 수 있습니다.
// 응답으로 반환된 UUID는 주문 조회나 취소에 사용할 수 없습니다.
//
// 실행:
//
//	UPBIT_ACCESS_KEY=<key> UPBIT_SECRET_KEY=<secret> go run examples/orders_test.go
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
// 2. 지정가 매수/매도 테스트
// ---------------------------------------------------------------------------

func testLimitOrders(ctx context.Context, client upbit.Client) {
	divider("2. 지정가 매수/매도 테스트")

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

	// --- 지정가 매수 (3% 할인) ---
	bidPrice := roundPrice(
		new(big.Float).Mul(bestBid, new(big.Float).SetFloat64(0.97)),
		tickSize,
	)
	bp, _, _ := new(big.Float).Parse(bidPrice, 10)
	bpInt, _ := bp.Int64()
	fmt.Printf("\n  [지정가 매수]\n")
	fmt.Printf("  주문 가격:     %s KRW (3%% 할인)\n", formatInt64WithComma(bpInt))
	fmt.Printf("  주문 수량:     %s\n", testVolume)

	bidOrder, err := client.Orders.TestNew(ctx, upbit.OrderTestNewParams{
		Market:  market,
		Side:    upbit.SideEnumBid,
		OrdType: upbit.OrdTypeEnumLimit,
		Price:   param.NewOpt(bidPrice),
		Volume:  param.NewOpt(testVolume),
	})
	if err != nil {
		log.Fatalf("지정가 매수 테스트 실패: %v", err)
	}
	fmt.Printf("  → UUID:  %s\n", bidOrder.Uuid)
	fmt.Printf("  → 상태:  %s\n", bidOrder.State)
	fmt.Printf("  → 유형:  %s\n", bidOrder.OrdType)

	// --- 지정가 매도 (3% 할증) ---
	askPrice := roundPrice(
		new(big.Float).Mul(bestBid, new(big.Float).SetFloat64(1.03)),
		tickSize,
	)
	ap, _, _ := new(big.Float).Parse(askPrice, 10)
	apInt, _ := ap.Int64()
	fmt.Printf("\n  [지정가 매도]\n")
	fmt.Printf("  주문 가격:     %s KRW (3%% 할증)\n", formatInt64WithComma(apInt))
	fmt.Printf("  주문 수량:     %s\n", testVolume)

	askOrder, err := client.Orders.TestNew(ctx, upbit.OrderTestNewParams{
		Market:  market,
		Side:    upbit.SideEnumAsk,
		OrdType: upbit.OrdTypeEnumLimit,
		Price:   param.NewOpt(askPrice),
		Volume:  param.NewOpt(testVolume),
	})
	if err != nil {
		log.Fatalf("지정가 매도 테스트 실패: %v", err)
	}
	fmt.Printf("  → UUID:  %s\n", askOrder.Uuid)
	fmt.Printf("  → 상태:  %s\n", askOrder.State)
	fmt.Printf("  → 유형:  %s\n", askOrder.OrdType)
}

// ---------------------------------------------------------------------------
// 3. 시장가 매수/매도 테스트
// ---------------------------------------------------------------------------

func testMarketOrders(ctx context.Context, client upbit.Client) {
	divider("3. 시장가 매수/매도 테스트")

	// 최소 주문금액 조회
	chance, err := client.Orders.GetChance(ctx, upbit.OrderGetChanceParams{Market: market})
	if err != nil {
		log.Fatalf("주문 가능 정보 조회 실패: %v", err)
	}
	minTotal := chance.Market.Bid.MinTotal
	mt, _, _ := new(big.Float).Parse(minTotal, 10)
	mtInt, _ := mt.Int64()

	// --- 시장가 매수 (총액 지정) ---
	fmt.Printf("  [시장가 매수]\n")
	fmt.Printf("  주문 금액:     %s KRW (최소 주문금액)\n", formatInt64WithComma(mtInt))

	buyOrder, err := client.Orders.TestNew(ctx, upbit.OrderTestNewParams{
		Market:  market,
		Side:    upbit.SideEnumBid,
		OrdType: upbit.OrdTypeEnumPrice,
		Price:   param.NewOpt(minTotal),
	})
	if err != nil {
		log.Fatalf("시장가 매수 테스트 실패: %v", err)
	}
	fmt.Printf("  → UUID:  %s\n", buyOrder.Uuid)
	fmt.Printf("  → 상태:  %s\n", buyOrder.State)
	fmt.Printf("  → 유형:  %s\n", buyOrder.OrdType)

	// --- 시장가 매도 (수량 지정) ---
	fmt.Printf("\n  [시장가 매도]\n")
	fmt.Printf("  주문 수량:     %s\n", testVolume)

	sellOrder, err := client.Orders.TestNew(ctx, upbit.OrderTestNewParams{
		Market:  market,
		Side:    upbit.SideEnumAsk,
		OrdType: upbit.OrdTypeEnumMarket,
		Volume:  param.NewOpt(testVolume),
	})
	if err != nil {
		log.Fatalf("시장가 매도 테스트 실패: %v", err)
	}
	fmt.Printf("  → UUID:  %s\n", sellOrder.Uuid)
	fmt.Printf("  → 상태:  %s\n", sellOrder.State)
	fmt.Printf("  → 유형:  %s\n", sellOrder.OrdType)
}

// ---------------------------------------------------------------------------
// 4. 최유리 매수/매도 테스트 (IOC)
// ---------------------------------------------------------------------------

func testBestOrders(ctx context.Context, client upbit.Client) {
	divider("4. 최유리 매수/매도 테스트 (IOC)")

	// 최소 주문금액 조회
	chance, err := client.Orders.GetChance(ctx, upbit.OrderGetChanceParams{Market: market})
	if err != nil {
		log.Fatalf("주문 가능 정보 조회 실패: %v", err)
	}
	minTotal := chance.Market.Bid.MinTotal
	mt, _, _ := new(big.Float).Parse(minTotal, 10)
	mtInt, _ := mt.Int64()

	// --- 최유리 매수 (총액 지정, IOC) ---
	fmt.Printf("  [최유리 매수 — IOC]\n")
	fmt.Printf("  주문 금액:     %s KRW (최소 주문금액)\n", formatInt64WithComma(mtInt))

	bidOrder, err := client.Orders.TestNew(ctx, upbit.OrderTestNewParams{
		Market:      market,
		Side:        upbit.SideEnumBid,
		OrdType:     upbit.OrdTypeEnumBest,
		Price:       param.NewOpt(minTotal),
		TimeInForce: upbit.TimeInForceEnumIoc,
	})
	if err != nil {
		log.Fatalf("최유리 매수 테스트 실패: %v", err)
	}
	fmt.Printf("  → UUID:  %s\n", bidOrder.Uuid)
	fmt.Printf("  → 상태:  %s\n", bidOrder.State)
	fmt.Printf("  → 유형:  %s\n", bidOrder.OrdType)

	// --- 최유리 매도 (수량 지정, IOC) ---
	fmt.Printf("\n  [최유리 매도 — IOC]\n")
	fmt.Printf("  주문 수량:     %s\n", testVolume)

	askOrder, err := client.Orders.TestNew(ctx, upbit.OrderTestNewParams{
		Market:      market,
		Side:        upbit.SideEnumAsk,
		OrdType:     upbit.OrdTypeEnumBest,
		Volume:      param.NewOpt(testVolume),
		TimeInForce: upbit.TimeInForceEnumIoc,
	})
	if err != nil {
		log.Fatalf("최유리 매도 테스트 실패: %v", err)
	}
	fmt.Printf("  → UUID:  %s\n", askOrder.Uuid)
	fmt.Printf("  → 상태:  %s\n", askOrder.State)
	fmt.Printf("  → 유형:  %s\n", askOrder.OrdType)
}

// ---------------------------------------------------------------------------
// 5. 잘못된 주문 요청 검증
// ---------------------------------------------------------------------------

func testValidation(ctx context.Context, client upbit.Client) {
	divider("5. 잘못된 주문 요청 검증")

	// --- 존재하지 않는 마켓 ---
	fmt.Println("  [검증] 존재하지 않는 마켓 (KRW-INVALID)")
	_, err := client.Orders.TestNew(ctx, upbit.OrderTestNewParams{
		Market:  "KRW-INVALID",
		Side:    upbit.SideEnumBid,
		OrdType: upbit.OrdTypeEnumLimit,
		Price:   param.NewOpt("10000"),
		Volume:  param.NewOpt("1"),
	})
	if err != nil {
		fmt.Printf("  → 에러 발생 (정상): %v\n", err)
	} else {
		fmt.Println("  → 예상과 다르게 성공했습니다.")
	}

	// --- 지정가 주문에 가격 누락 ---
	fmt.Println("\n  [검증] 지정가 주문 가격 누락")
	_, err = client.Orders.TestNew(ctx, upbit.OrderTestNewParams{
		Market:  market,
		Side:    upbit.SideEnumBid,
		OrdType: upbit.OrdTypeEnumLimit,
		Volume:  param.NewOpt(testVolume),
	})
	if err != nil {
		fmt.Printf("  → 에러 발생 (정상): %v\n", err)
	} else {
		fmt.Println("  → 예상과 다르게 성공했습니다.")
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
	testLimitOrders(ctx, client)
	testMarketOrders(ctx, client)
	testBestOrders(ctx, client)
	testValidation(ctx, client)

	fmt.Printf("\n%s\n", strings.Repeat("=", 60))
	fmt.Println("  주문 테스트 시나리오 완료")
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println()
	fmt.Println("  ※ 주문 생성 테스트 API는 실제 주문을 생성하지 않습니다.")
	fmt.Println("    반환된 UUID는 주문 조회/취소에 사용할 수 없습니다.")
}
