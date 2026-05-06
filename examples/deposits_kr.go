// 입금 관리 시나리오.
//
// 업비트 Exchange API를 사용하여 입금 관리 기능을 확인하는 예제입니다.
//   - 입금 주소 목록 조회
//   - 입금 주소 조회 / 생성
//   - 입금 가능 정보 조회
//   - 입금 내역 조회 (목록/개별)
//   - 트래블룰 지원 거래소 목록 조회
//   - 트래블룰 검증 (입금 UUID 기준)
//
// Dry run (기본):
//
//	입금 관련 정보 조회만 수행합니다.
//	동작 확인용으로 바로 실행할 수 있습니다.
//
// 실제 실행 (DRY_RUN=false):
//
//	입금 주소 생성 및 트래블룰 검증을 실행합니다.
//	(트래블룰 검증은 실제 입금 UUID와 거래소 UUID가 필요합니다.)
//
// 실행:
//
//	# Dry run (기본) — 조회만
//	UPBIT_ACCESS_KEY=<key> UPBIT_SECRET_KEY=<secret> go run examples/deposits.go
//
//	# 실제 실행 — 입금 주소 생성 및 트래블룰 검증 포함
//	DRY_RUN=false UPBIT_ACCESS_KEY=<key> UPBIT_SECRET_KEY=<secret> go run examples/deposits.go
package main

import (
	"context"
	"fmt"
	"log"
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

// ---------------------------------------------------------------------------
// 1. 입금 주소 목록 조회
// ---------------------------------------------------------------------------

func listDepositAddresses(ctx context.Context, client upbit.Client) {
	divider("1. 입금 주소 목록 조회")

	addresses, err := client.Deposits.ListCoinAddresses(ctx)
	if err != nil {
		log.Fatalf("입금 주소 목록 조회 실패: %v", err)
	}
	if len(*addresses) == 0 {
		fmt.Println("  등록된 입금 주소가 없습니다.")
		return
	}

	fmt.Printf("  %-8s %-8s %s\n", "통화", "네트워크", "주소")
	fmt.Printf("  %s\n", strings.Repeat("-", 50))
	for i, addr := range *addresses {
		if i >= 10 {
			break
		}
		depositAddr := addr.DepositAddress
		if depositAddr == "" {
			depositAddr = "(미생성)"
		}
		fmt.Printf("  %-8s %-8s %s\n", addr.Currency, addr.NetType, depositAddr)
	}
}

// ---------------------------------------------------------------------------
// 2. 특정 코인 입금 주소 조회
// ---------------------------------------------------------------------------

func retrieveDepositAddress(ctx context.Context, client upbit.Client, currency, netType string) {
	divider(fmt.Sprintf("2. %s 입금 주소 조회", currency))

	address, err := client.Deposits.GetCoinAddress(ctx, upbit.DepositGetCoinAddressParams{
		Currency: currency,
		NetType:  netType,
	})
	if err != nil {
		fmt.Printf("  입금 주소 조회 실패: %v\n", err)
		return
	}
	if address.DepositAddress != "" {
		fmt.Printf("  통화:     %s\n", address.Currency)
		fmt.Printf("  네트워크: %s\n", address.NetType)
		fmt.Printf("  주소:     %s\n", address.DepositAddress)
		if address.SecondaryAddress != "" {
			fmt.Printf("  2차 주소: %s\n", address.SecondaryAddress)
		}
	} else {
		fmt.Printf("  %s 입금 주소가 아직 생성되지 않았습니다.\n", currency)
	}
}

// ---------------------------------------------------------------------------
// 3. 입금 주소 생성
// ---------------------------------------------------------------------------

func createDepositAddress(ctx context.Context, client upbit.Client, currency, netType string) {
	divider(fmt.Sprintf("3. %s(%s) 입금 주소 생성", currency, netType))

	if dryRun {
		fmt.Println("  [건너뜀] Dry run 모드이므로 실제 주소 생성을 생략합니다.")
		fmt.Printf("  생성 대상: %s (%s)\n", currency, netType)
		return
	}

	result, err := client.Deposits.NewCoinAddress(ctx, upbit.DepositNewCoinAddressParams{
		Currency: currency,
		NetType:  netType,
	})
	if err != nil {
		fmt.Printf("  입금 주소 생성 실패: %v\n", err)
		return
	}

	if result.DepositAddress != "" {
		fmt.Printf("  통화:     %s\n", result.Currency)
		fmt.Printf("  네트워크: %s\n", result.NetType)
		fmt.Printf("  주소:     %s\n", result.DepositAddress)
		if result.SecondaryAddress != "" {
			fmt.Printf("  2차 주소: %s\n", result.SecondaryAddress)
		}
	} else {
		fmt.Printf("  생성 요청 상태: %v\n", result.Success)
		fmt.Printf("  메시지:         %s\n", result.Message)
	}
}

// ---------------------------------------------------------------------------
// 4. 입금 가능 정보 조회
// ---------------------------------------------------------------------------

func retrieveDepositChance(ctx context.Context, client upbit.Client, currency, netType string) {
	divider(fmt.Sprintf("4. %s 입금 가능 정보", currency))

	chance, err := client.Deposits.GetChance(ctx, upbit.DepositGetChanceParams{
		Currency: currency,
		NetType:  netType,
	})
	if err != nil {
		log.Fatalf("입금 가능 정보 조회 실패: %v", err)
	}

	depositable := "아니오"
	if chance.IsDepositPossible {
		depositable = "예"
	}
	fmt.Printf("  입금 가능:     %s\n", depositable)
	fmt.Printf("  최소 입금금액: %s\n", chance.MinimumDepositAmount)
	fmt.Printf("  필요 확인 수:  %d\n", chance.MinimumDepositConfirmations)
	if !chance.IsDepositPossible {
		fmt.Printf("  불가 사유:     %s\n", chance.DepositImpossibleReason)
	}
}

// ---------------------------------------------------------------------------
// 5. 입금 내역 조회
// ---------------------------------------------------------------------------

func listDeposits(ctx context.Context, client upbit.Client) {
	divider("5. 입금 내역 조회")

	page, err := client.Deposits.List(ctx, upbit.DepositListParams{
		Limit: param.NewOpt[int64](5),
	})
	if err != nil {
		log.Fatalf("입금 내역 조회 실패: %v", err)
	}
	deposits := page.Items
	if len(deposits) == 0 {
		fmt.Println("  입금 내역이 없습니다.")
		return
	}

	fmt.Printf("  %10s %-6s %14s %-24s\n", "UUID", "통화", "금액", "상태")
	fmt.Printf("  %s\n", strings.Repeat("-", 60))
	for _, d := range deposits {
		uuid := d.Uuid
		if len(uuid) > 8 {
			uuid = uuid[:8] + "..."
		}
		fmt.Printf("  %s %-6s %14s %-24s\n", uuid, d.Currency, d.Amount, string(d.State))
	}
}

// ---------------------------------------------------------------------------
// 6. 개별 입금 조회 (UUID / TxID)
// ---------------------------------------------------------------------------

func retrieveDeposit(ctx context.Context, client upbit.Client) {
	divider("6. 개별 입금 조회")

	deposit, err := client.Deposits.Get(ctx, upbit.DepositGetParams{})
	if err != nil {
		fmt.Printf("  입금 조회 실패: %v\n", err)
		return
	}
	fmt.Printf("  UUID:   %s\n", deposit.Uuid)
	fmt.Printf("  통화:   %s\n", deposit.Currency)
	fmt.Printf("  금액:   %s\n", deposit.Amount)
	fmt.Printf("  상태:   %s\n", deposit.State)
	if deposit.Txid != "" {
		fmt.Printf("  TxID:   %s\n", deposit.Txid)
	}

	// UUID로 동일 건 재조회
	fmt.Println()
	same, err := client.Deposits.Get(ctx, upbit.DepositGetParams{Uuid: param.NewOpt(deposit.Uuid)})
	if err != nil {
		fmt.Printf("  UUID 재조회 실패: %v\n", err)
	} else {
		fmt.Printf("  UUID 재조회 상태: %s\n", same.State)
	}

	// TxID로 동일 건 재조회
	if deposit.Txid != "" {
		sameByTxid, err := client.Deposits.Get(ctx, upbit.DepositGetParams{
			Txid:     param.NewOpt(deposit.Txid),
			Currency: param.NewOpt(deposit.Currency),
		})
		if err != nil {
			fmt.Printf("  TxID 재조회 실패: %v\n", err)
		} else {
			fmt.Printf("  TxID 재조회 상태: %s\n", sameByTxid.State)
		}
	}
}

// ---------------------------------------------------------------------------
// 7. 트래블룰 지원 거래소 목록 조회
// ---------------------------------------------------------------------------

func listTravelRuleVASPs(ctx context.Context, client upbit.Client) {
	divider("7. 트래블룰 지원 거래소 목록")

	vasps, err := client.TravelRule.ListVasps(ctx)
	if err != nil {
		log.Fatalf("트래블룰 거래소 목록 조회 실패: %v", err)
	}
	fmt.Printf("  지원 거래소 수: %d\n\n", len(*vasps))

	fmt.Printf("  %-20s %s\n", "거래소명", "UUID")
	fmt.Printf("  %s\n", strings.Repeat("-", 50))
	for i, v := range *vasps {
		if i >= 10 {
			break
		}
		fmt.Printf("  %-20s %s\n", v.VaspName, v.VaspUuid)
	}
	if len(*vasps) > 10 {
		fmt.Printf("  ... 외 %d개\n", len(*vasps)-10)
	}
}

// ---------------------------------------------------------------------------
// 8. 트래블룰 검증
// ---------------------------------------------------------------------------

func verifyTravelRule(ctx context.Context, client upbit.Client, depositUUID, vaspUUID string) {
	divider("8. 트래블룰 검증")

	fmt.Printf("  Deposit UUID: %s\n", depositUUID)
	fmt.Printf("  거래소 UUID:  %s\n", vaspUUID)

	if dryRun {
		fmt.Println("\n  [건너뜀] Dry run 모드이므로 트래블룰 검증을 생략합니다.")
		return
	}

	if strings.HasPrefix(depositUUID, "<") {
		fmt.Println("\n  [건너뜀] 입금 UUID가 설정되지 않았습니다.")
		fmt.Println("  코드에서 depositUUID 값을 실제 입금 UUID로 변경한 뒤 실행하세요.")
		return
	}

	if strings.HasPrefix(vaspUUID, "<") {
		fmt.Println("\n  [건너뜀] 거래소 UUID가 설정되지 않았습니다.")
		fmt.Println("  코드에서 vaspUUID 값을 실제 거래소 UUID로 변경한 뒤 실행하세요.")
		return
	}

	result, err := client.TravelRule.VerifyDepositByUuid(ctx, upbit.TravelRuleVerifyDepositByUuidParams{
		DepositUuid: depositUUID,
		VaspUuid:    vaspUUID,
	})
	if err != nil {
		fmt.Printf("  트래블룰 검증 실패: %v\n", err)
		return
	}
	fmt.Printf("  입금 UUID:  %s\n", result.DepositUuid)
	fmt.Printf("  입금 상태:  %s\n", result.DepositState)
	fmt.Printf("  검증 결과:  %s\n", result.VerificationResult)
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

	listDepositAddresses(ctx, client)
	retrieveDepositAddress(ctx, client, "BTC", "BTC")
	createDepositAddress(ctx, client, "USDT", "TRX")
	retrieveDepositChance(ctx, client, "BTC", "BTC")
	listDeposits(ctx, client)
	retrieveDeposit(ctx, client)
	listTravelRuleVASPs(ctx, client)
	verifyTravelRule(ctx, client, "<your_deposit_uuid>", "<deposit_vasp_uuid>")

	fmt.Printf("\n%s\n", strings.Repeat("=", 60))
	fmt.Println("  입금 관리 시나리오 완료")
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println()
	fmt.Println("  ※ 이 예제는 교육 목적으로 작성되었습니다.")
	fmt.Println("    실제 투자에 사용할 경우 발생하는 손실에 대해 책임지지 않습니다.")
}
