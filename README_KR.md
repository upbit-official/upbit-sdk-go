# Upbit SDK for Go

[English](./README.md) | 한국어

<!-- x-release-please-start-version -->

<a href="https://pkg.go.dev/github.com/upbit-official/upbit-sdk-go"><img src="https://pkg.go.dev/badge/github.com/upbit-official/upbit-sdk-go.svg" alt="Go Reference"></a>

<!-- x-release-please-end -->

Upbit SDK for Go는 Go 1.25 이상 환경에서 Upbit REST API를 편리하게 사용할 수 있도록 설계되었습니다.
모든 요청 파라미터와 응답 필드에 대한 타입 정의를 포함하며, SDK의 전체 API 레퍼런스는 [api.md](api.md)에서 확인할 수 있습니다.

## 설치

<!-- x-release-please-start-version -->

```go
import (
	"github.com/upbit-official/upbit-sdk-go" // imported as upbit
)
```

<!-- x-release-please-end -->

특정 버전으로 고정하려면 다음과 같이 설치하세요.

<!-- x-release-please-start-version -->

```sh
go get -u 'github.com/upbit-official/upbit-sdk-go@v1.0.1'
```

<!-- x-release-please-end -->

## 사용법

### 인도네시아 API 주소

v1.0.1에서 `option.WithEnvironmentId()`는 `https://api.upbit.co.id/`를 사용한다. 구버전에서는 `option.WithBaseURL("https://api.upbit.co.id/")`로 주소를 지정한다. 직접 설정한 주소는 SDK 업데이트 후에도 유지되므로 별도로 갱신한다.

국가별 API 문서: [한국](https://docs.upbit.com/kr), [인도네시아](https://docs.upbit.co.id), [싱가포르·태국](https://global-docs.upbit.com).

SDK의 전체 API는 [api.md](api.md)에서 확인할 수 있습니다.

```go
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/upbit-official/upbit-sdk-go"
	"github.com/upbit-official/upbit-sdk-go/option"
)

func main() {
	client := upbit.NewClient(
		option.WithAccessKey(os.Getenv("UPBIT_ACCESS_KEY")),
		option.WithSecretKey(os.Getenv("UPBIT_SECRET_KEY")),
		// option.WithEnvironmentSg(), // 또는 WithEnvironmentKr() | WithEnvironmentId() | WithEnvironmentTh(); 기본값은 WithEnvironmentKr()
	)
	accounts, err := client.Accounts.List(context.TODO())
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("%+v\n", accounts)
}
```

실행 가능한 예제는 [`examples/`](examples/) 디렉터리에서 확인할 수 있습니다.

## 요청 및 응답 타입

이 SDK는 문서화된 모든 요청 파라미터와 응답 필드에 대한 Go 타입을 제공합니다. 다음과 같이 가져와 사용할 수 있습니다.

```go
params := upbit.CandleListDaysParams{
	Market: "KRW-BTC",
	Count:  upbit.Int(10),
	To:     upbit.String("2024-01-01T00:00:00"),
}

candles, err := client.Candles.ListDays(context.TODO(), params)
if err != nil {
	panic(err.Error())
}
fmt.Printf("%+v\n", candles)
```

필수 요청 필드는 일반 값 타입이며 <code>`api:"required"`</code> 태그가 지정됩니다. 선택적 기본형 필드는 `param.Opt[T]`로 감싸며 `upbit.String(string)`, `upbit.Int(int64)`, `upbit.Bool(bool)`, `upbit.Float(float64)` 같은 생성자로 설정할 수 있습니다. SDK는 생략된 요청 필드에 대해 `encoding/json`의 [`omitzero`](https://tip.golang.org/doc/go1.24#encodingjsonpkgencodingjson) 동작 방식을 사용합니다.

요청 본문 파라미터도 같은 패턴을 사용합니다.

```go
params := upbit.OrderNewParams{
	Market:  "KRW-BTC",
	Side:    upbit.SideEnumBid,
	OrdType: upbit.OrdTypeEnumLimit,
	Price:   upbit.String("14000000"),
	Volume:  upbit.String("1"),
}

order, err := client.Orders.New(context.TODO(), params)
if err != nil {
	panic(err.Error())
}
fmt.Printf("%+v\n", order)
```

응답 구조체의 모든 필드는 일반 값 타입입니다. 각 응답 구조체에는 속성별 메타데이터를 담는 `JSON` 필드도 포함됩니다. 명시적인 값과 `null`, 누락된 필드, 디코딩할 수 없는 값을 구분하려면 `.Valid()`를 사용하세요.

```go
accounts, err := client.Accounts.List(context.TODO())
if err != nil {
	panic(err.Error())
}

var account upbit.AccountListResponse
if len(*accounts) > 0 {
	account = (*accounts)[0]
}

fmt.Println(account.Currency)

if account.JSON.Balance.Valid() {
	fmt.Println(account.Balance)
}

raw := account.RawJSON()
extra := account.JSON.ExtraFields["new_field"].Raw()
fmt.Println(raw, extra)
```

일부 응답은 유니온 타입입니다. 예를 들어 `client.Deposits.NewCoinAddress`에서 반환된 `DepositNewCoinAddressResponseUnion` 값은 생성된 `As...` 메서드로 문서화된 변형 중 하나로 변환할 수 있습니다.

```go
result, err := client.Deposits.NewCoinAddress(context.TODO(), params)
if err != nil {
	panic(err.Error())
}

if result.JSON.DepositAddress.Valid() {
	address := result.AsDepositAddress()
	fmt.Println(address.DepositAddress)
} else if result.JSON.Message.Valid() {
	status := result.AsDepositNewCoinAddressResponseGenerateCoinAddressResponse()
	fmt.Println(status.Message)
}
```

각 메서드, 요청 파라미터, 응답 필드에 대한 문서는 Go doc 주석과 [pkg.go.dev](https://pkg.go.dev/github.com/upbit-official/upbit-sdk-go)에서 확인할 수 있습니다.

## 오류 처리

API가 성공이 아닌 상태 코드를 반환하면 `*upbit.Error` 타입의 오류가 반환됩니다. 여기에는 `StatusCode`, `*http.Request`, `*http.Response`, 그리고 오류 본문 JSON이 포함됩니다.

오류를 처리할 때는 `errors.As` 패턴을 사용하세요.

```go
_, err := client.Accounts.List(context.TODO())
if err != nil {
	var apierr *upbit.Error
	if errors.As(err, &apierr) {
		fmt.Println(apierr.StatusCode)
		fmt.Println(apierr.Name)
		fmt.Println(apierr.Message)
		fmt.Println(string(apierr.DumpRequest(true)))
		fmt.Println(string(apierr.DumpResponse(true)))
	}
	panic(err.Error())
}
```

그 밖의 오류는 별도로 래핑되지 않고 그대로 반환됩니다. 예를 들어 HTTP 전송 실패 시에는 `*net.OpError`를 감싼 `*url.Error`가 반환될 수 있습니다.

## 재시도

일부 오류에 대해서는 짧은 지수 백오프를 적용해 기본적으로 최대 2회까지 자동 재시도합니다.
연결 오류, 408 Request Timeout, 409 Conflict, 429 Rate Limit, 500 이상 내부 오류가 기본 재시도 대상입니다.

`WithMaxRetries` 옵션으로 재시도 설정을 변경하거나 비활성화할 수 있습니다.

```go
// 모든 요청의 기본값 설정:
client := upbit.NewClient(
	option.WithMaxRetries(0), // 기본값은 2
)

// 요청별 오버라이드:
client.Accounts.List(context.TODO(), option.WithMaxRetries(5))
```

## 타임아웃

요청에는 기본 타임아웃이 설정되어 있지 않습니다. 요청의 전체 생명주기에 대한 타임아웃은 context를 사용해 지정하세요.

[재시도](#재시도) 중에는 context 타임아웃이 초기화되지 않습니다. 재시도마다 별도 타임아웃을 적용하려면 `option.WithRequestTimeout()`을 사용하세요.

```go
// 모든 재시도를 포함한 전체 요청 타임아웃 설정
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
defer cancel()
client.Accounts.List(
	ctx,
	// 재시도별 타임아웃 설정
	option.WithRequestTimeout(20*time.Second),
)
```

## 페이지네이션

이 SDK는 페이지네이션이 있는 목록 엔드포인트를 다루기 위한 편의 기능을 제공합니다.

`.ListAutoPaging()` 메서드를 사용하면 모든 페이지의 항목을 순회할 수 있습니다.

```go
iter := client.Orders.ListOpenAutoPaging(context.TODO(), upbit.OrderListOpenParams{})
// 필요할 때 자동으로 다음 페이지를 가져옵니다.
for iter.Next() {
	order := iter.Current()
	fmt.Printf("%+v\n", order)
}
if err := iter.Err(); err != nil {
	panic(err.Error())
}
```

또는 `.List()` 메서드로 단일 페이지를 가져온 뒤 `.GetNextPage()` 같은 헬퍼 메서드를 사용할 수 있습니다.

```go
page, err := client.Orders.ListOpen(context.TODO(), upbit.OrderListOpenParams{})
for page != nil {
	for _, order := range page.Items {
		fmt.Printf("%+v\n", order)
	}
	page, err = page.GetNextPage()
}
if err != nil {
	panic(err.Error())
}
```

## 고급 사용법

### RequestOptions

이 SDK는 함수형 옵션 패턴을 사용합니다. `option` 패키지에 정의된 함수들은 `RequestConfig`를 변경하는 클로저인 `RequestOption`을 반환합니다. 이 옵션들은 클라이언트 전체 또는 개별 요청에 적용할 수 있습니다.

```go
client := upbit.NewClient(
	// 클라이언트의 모든 요청에 헤더 추가
	option.WithHeader("X-Some-Header", "custom_header_info"),
)

_, err := client.Accounts.List(
	context.TODO(),
	// 이 요청에서 헤더 덮어쓰기
	option.WithHeader("X-Some-Header", "some_other_custom_header_info"),
)
```

디버깅할 때는 `option.WithDebugLog(nil)` 옵션이 유용합니다.

[전체 요청 옵션 목록](https://pkg.go.dev/github.com/upbit-official/upbit-sdk-go/option)을 참고하세요.

### 원시 응답 데이터 접근

`option.WithResponseInto()` 요청 옵션을 사용하면 원시 HTTP 응답 데이터에 접근할 수 있습니다. 응답 헤더나 상태 코드 등을 확인할 때 유용합니다.

```go
var response *http.Response
accounts, err := client.Accounts.List(context.TODO(), option.WithResponseInto(&response))
if err != nil {
	panic(err.Error())
}
fmt.Printf("%+v\n", accounts)

fmt.Printf("Status Code: %d\n", response.StatusCode)
fmt.Printf("Headers: %+#v\n", response.Header)
```

### 커스텀 요청 / 문서화되지 않은 엔드포인트 사용

이 SDK는 문서화된 API를 편리하게 사용할 수 있도록 타입을 정의한 SDK입니다. 문서화되지 않은 엔드포인트, 파라미터, 응답 속성도 SDK를 통해 사용할 수 있습니다.

#### 문서화되지 않은 엔드포인트

`client.Get`, `client.Post` 같은 HTTP 메서드를 직접 사용해 문서화되지 않은 엔드포인트에 요청할 수 있습니다. 이때도 재시도 등 클라이언트에 설정된 `RequestOptions`가 동일하게 적용됩니다.

```go
var result map[string]any

err := client.Get(
	context.Background(),
	"/v1/undocumented",
	nil,
	&result,
	option.WithQuery("market", "KRW-BTC"),
)
if err != nil {
	panic(err.Error())
}
```

#### 문서화되지 않은 요청 파라미터

문서화되지 않은 쿼리 파라미터를 보내려면 실제 SDK 메서드 호출에 `option.WithQuery()`를 사용하세요.

```go
candles, err := client.Candles.ListDays(
	context.Background(),
	upbit.CandleListDaysParams{
		Market: "KRW-BTC",
	},
	option.WithQuery("include_experimental", "true"),
)
```

문서화되지 않은 JSON 본문 파라미터를 보내려면 `option.WithJSONSet()`을 사용하세요.

```go
params := upbit.OrderNewParams{
	Market:  "KRW-BTC",
	Side:    upbit.SideEnumBid,
	OrdType: upbit.OrdTypeEnumLimit,
	Price:   upbit.String("14000000"),
	Volume:  upbit.String("1"),
}

order, err := client.Orders.New(
	context.Background(),
	params,
	option.WithJSONSet("metadata.strategy_id", "manual-1"),
)
```

#### 문서화되지 않은 응답 속성

문서화되지 않은 응답 속성은 전체 raw JSON을 반환하는 `RawJSON()` 또는 응답 구조체에 정의되지 않은 필드를 보존하는 `JSON.ExtraFields`로 접근할 수 있습니다.

```go
if len(*accounts) > 0 {
	account := (*accounts)[0]
	fmt.Println(account.RawJSON())
	fmt.Println(account.JSON.ExtraFields["new_field"].Raw())
}
```

### 미들웨어

`option.WithMiddleware`를 사용하면 요청에 미들웨어를 적용할 수 있습니다.

```go
func Logger(req *http.Request, next option.MiddlewareNext) (res *http.Response, err error) {
	// 요청 전
	start := time.Now()
	LogReq(req)

	// 다음 핸들러로 전달
	res, err = next(req)

	// 요청 후 처리
	end := time.Now()
	LogRes(res, err, end.Sub(start))

	return res, err
}

client := upbit.NewClient(
	option.WithMiddleware(Logger),
)
```

여러 미들웨어를 가변 인수로 전달하면 왼쪽에서 오른쪽 순서로 적용됩니다. `option.WithMiddleware`가 클라이언트와 메서드 양쪽에 지정된 경우에는 클라이언트에 지정한 미들웨어가 먼저 실행됩니다.

`option.WithHTTPClient(client)`를 사용하면 기본 `http.Client`를 교체할 수도 있습니다. HTTP 클라이언트는 하나만 지정할 수 있으며 마지막에 지정한 값이 적용됩니다. 요청은 미들웨어가 적용된 뒤 HTTP 클라이언트로 전달됩니다.

## 버전 관리

이 패키지는 일반적으로 [SemVer](https://semver.org/spec/v2.0.0.html) 규약을 따르지만, 일부 하위 호환성에 영향을 줄 수 있는 변경이 마이너 버전에 포함될 수 있습니다.

1. 기술적으로는 공개되어 있지만 외부 사용을 의도하지 않았거나 문서화되지 않은 SDK 내부 변경
2. 대부분의 사용자에게 실질적인 영향이 없을 것으로 판단되는 변경

하위 호환성은 중요하게 고려하고 있으며, 원활한 업그레이드 경험을 제공하기 위해 노력하고 있습니다.

피드백을 환영합니다. 질문, 버그 제보, 개선 제안이 있다면 open-api@upbit.com으로 연락해 주세요.

## 요구사항

Go 1.25 이상이 필요합니다.

## 기여

현재 업비트 SDK는 초기 출시 단계로, 외부 Issue 등록 및 PR 기여는 아직 운영하고 있지 않습니다.
버그 제보나 개선 의견은 개발자 지원 채널(open-api@upbit.com)을 통해 전달해 주시기 바랍니다.
외부 기여 채널은 향후 운영 안정화 상황에 맞춰 순차적으로 오픈하는 방안을 검토하고 있습니다.

© 2026 Dunamu Inc. All rights reserved.
