// 출금 관리 시나리오.
//
// 업비트 Exchange API를 사용하여 출금 관리 기능을 확인하는 예제입니다.
//   - 출금 허용 주소 목록 조회
//   - 입출금 서비스 상태 확인
//   - 출금 가능 정보 조회
//   - 출금 내역 조회 (목록/개별)
//   - 디지털 자산 출금
//   - 원화(KRW) 출금
//
// Dry run (기본):
//
//	출금 관련 정보 조회만 수행합니다.
//	동작 확인용으로 바로 실행할 수 있습니다.
//
// 실제 실행 (DRY_RUN=false):
//
//	출금 가능 정보를 출력하지만, 실제 출금은 실행하지 않습니다.
//	실제 출금을 원하면 코드 내 주소/금액을 직접 수정한 뒤 실행하세요.
//
//	주의: 출금은 비가역적입니다. 잘못된 주소로 출금하면 자산을 잃을 수 있습니다.
//
// 실행:
//
//	# Dry run (기본) — 조회만
//	UPBIT_ACCESS_KEY=<key> UPBIT_SECRET_KEY=<secret> go run examples/withdrawals.go
//
//	# 실제 실행 — 출금 가능 정보 표시 (실제 출금은 코드 수정 필요)
//	DRY_RUN=false UPBIT_ACCESS_KEY=<key> UPBIT_SECRET_KEY=<secret> go run examples/withdrawals.go
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

var dryRun = os.Getenv("DRY_RUN") != "false"

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

// ---------------------------------------------------------------------------
// 1. 출금 허용 주소 목록 조회
// ---------------------------------------------------------------------------

func listWithdrawalAddresses(ctx context.Context, client upbit.Client) {
	divider("1. 출금 허용 주소 목록")

	addresses, err := client.Withdraws.ListCoinAddresses(ctx)
	if err != nil {
		log.Fatalf("출금 주소 목록 조회 실패: %v", err)
	}
	if len(*addresses) == 0 {
		fmt.Println("  등록된 출금 주소가 없습니다.")
		fmt.Println("  업비트 웹에서 출금 허용 주소를 먼저 등록해주세요.")
		return
	}

	fmt.Printf("  %-8s %-8s %-12s %s\n", "통화", "네트워크", "거래소", "주소")
	fmt.Printf("  %s\n", strings.Repeat("-", 60))
	for i, addr := range *addresses {
		if i >= 10 {
			break
		}
		exchange := addr.ExchangeName
		if exchange == "" {
			exchange = "-"
		}
		fmt.Printf("  %-8s %-8s %-12s %s\n", addr.Currency, addr.NetType, exchange, addr.WithdrawAddress)
	}
}

// ---------------------------------------------------------------------------
// 2. 입출금 서비스 상태 확인
// ---------------------------------------------------------------------------

func checkWalletStatus(ctx context.Context, client upbit.Client) {
	divider("2. 입출금 서비스 상태")

	wallets, err := client.WalletStatus.List(ctx)
	if err != nil {
		log.Fatalf("지갑 상태 조회 실패: %v", err)
	}

	targetCurrencies := map[string]bool{"BTC": true, "ETH": true, "USDT": true}

	fmt.Printf("  전체 지갑 수: %d\n", len(*wallets))
	fmt.Printf("\n  %-8s %-10s %-12s %s\n", "통화", "네트워크", "상태", "블록 상태")
	fmt.Printf("  %s\n", strings.Repeat("-", 48))
	for _, w := range *wallets {
		if targetCurrencies[w.Currency] {
			fmt.Printf("  %-8s %-10s %-12s %s\n",
				w.Currency, w.NetType, string(w.WalletState), string(w.BlockState))
		}
	}
}

// ---------------------------------------------------------------------------
// 3. 출금 가능 정보 조회
// ---------------------------------------------------------------------------

func checkWithdrawalChance(ctx context.Context, client upbit.Client, currency, netType string) {
	divider(fmt.Sprintf("3. %s 출금 가능 정보", currency))

	chance, err := client.Withdraws.GetChance(ctx, upbit.WithdrawGetChanceParams{
		Currency: currency,
		NetType:  netType,
	})
	if err != nil {
		log.Fatalf("출금 가능 정보 조회 실패: %v", err)
	}

	fmt.Printf("  보유 수량:     %s\n", chance.Account.Balance)
	fmt.Printf("  출금 수수료:   %s\n", chance.Currency.WithdrawFee)
	fmt.Printf("  최소 출금금액: %s\n", chance.WithdrawLimit.Minimum)
	fmt.Printf("  1일 잔여한도:  %s %s\n", chance.WithdrawLimit.RemainingDailyFiat, chance.WithdrawLimit.FiatCurrency)
}

// ---------------------------------------------------------------------------
// 4. 출금 내역 조회
// ---------------------------------------------------------------------------

func listWithdrawals(ctx context.Context, client upbit.Client) {
	divider("4. 출금 내역 조회")

	page, err := client.Withdraws.List(ctx, upbit.WithdrawListParams{
		Limit: param.NewOpt[int64](5),
	})
	if err != nil {
		log.Fatalf("출금 내역 조회 실패: %v", err)
	}
	withdrawals := page.Items
	if len(withdrawals) == 0 {
		fmt.Println("  출금 내역이 없습니다.")
		return
	}

	fmt.Printf("  %10s %-6s %14s %-12s\n", "UUID", "통화", "금액", "상태")
	fmt.Printf("  %s\n", strings.Repeat("-", 50))
	for _, w := range withdrawals {
		uuid := w.Uuid
		if len(uuid) > 8 {
			uuid = uuid[:8] + "..."
		}
		fmt.Printf("  %s %-6s %14s %-12s\n", uuid, w.Currency, w.Amount, string(w.State))
	}
}

// ---------------------------------------------------------------------------
// 5. 개별 출금 조회
// ---------------------------------------------------------------------------

func retrieveWithdrawal(ctx context.Context, client upbit.Client) {
	divider("5. 개별 출금 조회")

	withdrawal, err := client.Withdraws.Get(ctx, upbit.WithdrawGetParams{})
	if err != nil {
		fmt.Printf("  출금 조회 실패 (출금 내역이 없을 수 있음): %v\n", err)
		return
	}
	fmt.Printf("  UUID:   %s\n", withdrawal.Uuid)
	fmt.Printf("  통화:   %s\n", withdrawal.Currency)
	fmt.Printf("  금액:   %s\n", withdrawal.Amount)
	fmt.Printf("  상태:   %s\n", withdrawal.State)
	if withdrawal.Txid != "" {
		fmt.Printf("  TxID:   %s\n", withdrawal.Txid)
	}
}

// ---------------------------------------------------------------------------
// 6. 디지털 자산 출금
// ---------------------------------------------------------------------------

func withdrawCoin(ctx context.Context, client upbit.Client, currency, netType, amount, address string) {
	divider("6. 디지털 자산 출금")

	fmt.Printf("  통화:     %s\n", currency)
	fmt.Printf("  네트워크: %s\n", netType)
	fmt.Printf("  금액:     %s\n", amount)
	fmt.Printf("  주소:     %s\n", address)

	if dryRun {
		fmt.Println("\n  [건너뜀] Dry run 모드이므로 실제 출금을 생략합니다.")
		return
	}

	if strings.HasPrefix(address, "<") {
		fmt.Println("\n  [건너뜀] 출금 주소가 설정되지 않았습니다.")
		fmt.Println("  코드에서 address 값을 실제 출금 주소로 변경한 뒤 실행하세요.")
		return
	}

	result, err := client.Withdraws.NewWithdrawal(ctx, upbit.WithdrawNewWithdrawalParams{
		Currency:        currency,
		NetType:         netType,
		Amount:          amount,
		Address:         address,
		TransactionType: upbit.WithdrawNewWithdrawalParamsTransactionTypeDefault,
	})
	if err != nil {
		log.Fatalf("디지털 자산 출금 실패: %v", err)
	}
	fmt.Printf("\n  출금 UUID: %s\n", result.Uuid)
	fmt.Printf("  상태:      %s\n", result.State)
}

// ---------------------------------------------------------------------------
// 7. 원화(KRW) 출금
// ---------------------------------------------------------------------------

func withdrawKRW(ctx context.Context, client upbit.Client, amount string) {
	divider("7. 원화(KRW) 출금")

	amt, _, _ := new(big.Float).Parse(amount, 10)
	amtInt, _ := amt.Int64()
	fmt.Printf("  출금 금액: %s KRW\n", formatInt64WithComma(amtInt))

	if dryRun {
		fmt.Println("\n  [건너뜀] Dry run 모드이므로 실제 출금을 생략합니다.")
		return
	}

	result, err := client.Withdraws.NewKrwWithdrawal(ctx, upbit.WithdrawNewKrwWithdrawalParams{
		Amount:        amount,
		TwoFactorType: upbit.WithdrawNewKrwWithdrawalParamsTwoFactorTypeKakao,
	})
	if err != nil {
		log.Fatalf("원화 출금 실패: %v", err)
	}
	fmt.Printf("\n  출금 UUID: %s\n", result.Uuid)
	fmt.Printf("  상태:      %s\n", result.State)
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

	listWithdrawalAddresses(ctx, client)
	checkWalletStatus(ctx, client)
	checkWithdrawalChance(ctx, client, "BTC", "BTC")
	listWithdrawals(ctx, client)
	retrieveWithdrawal(ctx, client)
	withdrawCoin(ctx, client, "USDT", "TRX", "13.241", "<your_withdrawal_address>")
	withdrawKRW(ctx, client, "10000")

	fmt.Printf("\n%s\n", strings.Repeat("=", 60))
	fmt.Println("  출금 관리 시나리오 완료")
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println()
	fmt.Println("  ※ 이 예제는 교육 목적으로 작성되었습니다.")
	fmt.Println("    실제 투자에 사용할 경우 발생하는 손실에 대해 책임지지 않습니다.")
}
