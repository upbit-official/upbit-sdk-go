// API 에러 처리 시나리오.
//
// 400 / 401 / 404 세 가지 에러 케이스를 실제 API 호출로 시연합니다.
//   - 400 BadRequest  : limit 주문에 price 누락
//   - 401 AuthError   : 잘못된 credentials 로 계좌 조회
//   - 404 NotFound    : 존재하지 않는 UUID 로 주문 조회
//
// 실행:
//
//	UPBIT_ACCESS_KEY=<key> UPBIT_SECRET_KEY=<secret> go run examples/error_handling_kr.go
package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	upbit "github.com/upbit-official/upbit-sdk-go"
	"github.com/upbit-official/upbit-sdk-go/option"
)

func printErr(label string, err error) {
	fmt.Printf("\n[%s]\n", label)
	if err == nil {
		fmt.Println("  에러 없음 (예상과 다름)")
		return
	}

	var apiErr *upbit.Error
	if errors.As(err, &apiErr) {
		fmt.Printf("  StatusCode : %d\n", apiErr.StatusCode)
		fmt.Printf("  Name       : %q\n", apiErr.Name)
		fmt.Printf("  Message    : %q\n", apiErr.Message)
	} else {
		fmt.Printf("  %v\n", err)
	}
}

func main() {
	ctx := context.Background()

	client := upbit.NewClient(
		option.WithAccessKey(os.Getenv("UPBIT_ACCESS_KEY")),
		option.WithSecretKey(os.Getenv("UPBIT_SECRET_KEY")),
		option.WithEnvironmentKr(),
	)

	// ── 400 BadRequest ────────────────────────────────────────────────────────
	// limit 주문에 price 를 빠뜨리면 invalid_parameter 에러 발생
	_, err := client.Orders.New(ctx, upbit.OrderNewParams{
		Market:  "KRW-BTC",
		OrdType: upbit.OrdTypeEnumLimit,
		Side:    upbit.SideEnumBid,
		Volume:  upbit.String("0.0001"),
		// Price 의도적으로 누락
	})
	printErr("400 BadRequest", err)
	var apiErr *upbit.Error
	if errors.As(err, &apiErr) && apiErr.StatusCode == 400 {
		fmt.Printf("  → 400 확인: Name=%q\n", apiErr.Name)
	}

	// ── 401 AuthenticationError ───────────────────────────────────────────────
	// 잘못된 키로 클라이언트를 만들면 인증 에러 발생
	badClient := upbit.NewClient(
		option.WithAccessKey("invalid_access_key"),
		option.WithSecretKey("invalid_secret_key"),
		option.WithEnvironmentKr(),
	)
	_, err = badClient.Accounts.List(ctx)
	printErr("401 AuthenticationError", err)
	if errors.As(err, &apiErr) && apiErr.StatusCode == 401 {
		fmt.Printf("  → 401 확인: Name=%q\n", apiErr.Name)
	}

	// ── 404 NotFound ─────────────────────────────────────────────────────────
	// 존재하지 않는 UUID 로 주문 조회
	_, err = client.Orders.Get(ctx, upbit.OrderGetParams{
		Uuid: upbit.String("00000000-0000-0000-0000-000000000000"),
	})
	printErr("404 NotFound", err)
	if errors.As(err, &apiErr) && apiErr.StatusCode == 404 {
		fmt.Printf("  → 404 확인: Name=%q\n", apiErr.Name)
	}

	// ── 429 RateLimitError ────────────────────────────────────────────────────
	// 짧은 시간에 대량 요청을 보내야 발생 — 예제에서는 핸들러 패턴만 시연
	//
	// if errors.As(err, &apiErr) && apiErr.StatusCode == 429 {
	//     fmt.Println("요청 횟수 초과. 잠시 후 재시도.")
	// }
}
