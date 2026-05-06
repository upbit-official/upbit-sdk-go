[English](./README.md) | 한국어

# Upbit Go API 라이브러리

<!-- x-release-please-start-version -->

<a href="https://pkg.go.dev/github.com/upbit-official/upbit-sdk-go"><img src="https://pkg.go.dev/badge/github.com/upbit-official/upbit-sdk-go.svg" alt="Go Reference"></a>

<!-- x-release-please-end -->

Upbit Go 라이브러리는 Go 애플리케이션에서 [Upbit REST API](https://docs.upbit.com)를 보다 편리하게 사용할 수 있도록 설계되었습니다. 모든 요청 파라미터와 응답 필드에 대한 타입 정의를 제공합니다.

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
go get -u 'github.com/upbit-official/upbit-sdk-go@v0.9.0'
```

<!-- x-release-please-end -->

## 요구사항

Go 1.25 이상이 필요합니다.

## 사용법

라이브러리의 전체 API는 [api.md](api.md)에서 확인할 수 있습니다.

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

### 요청 필드

Upbit 라이브러리는 `encoding/json`의 [`omitzero`](https://tip.golang.org/doc/go1.24#encodingjsonpkgencodingjson) 동작 방식을 요청 필드에 적용합니다.

필수 기본형 필드(`int64`, `string` 등)에는 <code>`api:"required"`</code> 태그가 지정됩니다. 이 필드들은 값이 0이더라도 항상 직렬화됩니다.

선택적 기본형 필드는 `param.Opt[T]`로 감쌉니다. 값은 `upbit.String(string)`, `upbit.Int(int64)` 등 제공되는 생성자를 사용해 설정할 수 있습니다.

`param.Opt[T]`, 맵, 슬라이스, 구조체, 문자열 열거형에는 <code>`json:"...,omitzero"`</code> 태그가 사용됩니다. 이러한 타입의 제로값은 생략된 것으로 간주됩니다.

`param.IsOmitted(any)` 함수로 `omitzero` 필드가 실제로 생략되었는지 확인할 수 있습니다.

```go
p := upbit.ExampleParams{
	ID:   "id_xxx",            // 필수 필드
	Name: upbit.String("..."), // 선택 필드

	Point: upbit.Point{
		X: 0,            // 필수 필드 — 0이어도 직렬화됨
		Y: upbit.Int(1), // 선택 필드 — 1로 직렬화됨
		// ... 선택하지 않은 필드는 직렬화되지 않음
	},

	Origin: upbit.Origin{}, // [Origin]의 제로값은 생략된 것으로 간주됨
}
```

`param.Opt[T]` 대신 `null`을 보내려면 `param.Null[T]()`를 사용하세요.
구조체 `T` 대신 `null`을 보내려면 `param.NullStruct[T]()`를 사용하세요.

```go
p.Name = param.Null[string]()       // string 대신 'null'
p.Point = param.NullStruct[Point]() // struct 대신 'null'

param.IsNull(p.Name)  // true
param.IsNull(p.Point) // true
```

요청 구조체에는 `.SetExtraFields(map[string]any)` 메서드가 있어 비표준 필드를 요청 본문에 추가할 수 있습니다. 추가 필드는 동일한 키를 가진 구조체 필드를 덮어씁니다. 보안상 `SetExtraFields`는 신뢰할 수 있는 데이터에만 사용하세요.

구조체 대신 사용자 정의 값을 보내려면 `param.Override[T](value)`를 사용하세요.

```go
// API가 특정 타입을 요구하지만 다른 값을 보내야 하는 경우 [SetExtraFields] 사용:
p.SetExtraFields(map[string]any{
	"x": 0.01, // "x"를 int 대신 float으로 전송
})

// 객체 대신 숫자 전송
custom := param.Override[upbit.FooParams](12)
```

### 요청 유니온(Union)

유니온(Union)은 각 변형마다 `Of` 접두사가 붙은 필드를 포함하는 구조체로 표현됩니다. 한 번에 하나의 필드만 non-zero 값을 가질 수 있으며, 해당 non-zero 필드만 직렬화됩니다.

유니온(Union) 구조체의 메서드를 통해 하위 속성에 접근할 수 있습니다. 이 메서드는 해당 데이터가 있으면 수정 가능한 포인터를 반환합니다.

```go
// 하나의 필드만 non-zero 가능, param.IsOmitted()로 설정 여부 확인
type AnimalUnionParam struct {
	OfCat *Cat `json:",omitzero,inline"`
	OfDog *Dog `json:",omitzero,inline"`
}

animal := AnimalUnionParam{
	OfCat: &Cat{
		Name: "Whiskers",
		Owner: PersonParam{
			Address: AddressParam{Street: "3333 Coyote Hill Rd", Zip: 0},
		},
	},
}

// 필드 변경
if address := animal.GetOwner().GetAddress(); address != nil {
	address.ZipCode = 94304
}
```

### 응답 객체

응답 구조체의 모든 필드는 일반 값 타입입니다(포인터나 래퍼 타입이 아님). 또한 각 속성의 메타데이터를 담는 특수 `JSON` 필드가 포함됩니다.

```go
type Animal struct {
	Name   string `json:"name,nullable"`
	Owners int    `json:"owners"`
	Age    int    `json:"age"`
	JSON   struct {
		Name        respjson.Field
		Owner       respjson.Field
		Age         respjson.Field
		ExtraFields map[string]respjson.Field
	} `json:"-"`
}
```

선택적 데이터를 처리하려면 JSON 필드의 `.Valid()` 메서드를 사용하세요. `.Valid()`는 필드가 `null`이 아니고, 누락되지 않았으며, 직렬화에도 실패하지 않았을 때 `true`를 반환합니다.

`.Valid()`가 `false`이면 해당 필드는 제로값으로 처리됩니다.

```go
raw := `{"owners": 1, "name": null}`

var res Animal
json.Unmarshal([]byte(raw), &res)

// 일반 필드 접근

res.Owners // 1
res.Name   // ""
res.Age    // 0

// 선택적 필드 확인

res.JSON.Owners.Valid() // true
res.JSON.Name.Valid()   // false
res.JSON.Age.Valid()    // false

// Raw JSON 값

res.JSON.Owners.Raw()                  // "1"
res.JSON.Name.Raw() == "null"          // true
res.JSON.Name.Raw() == respjson.Null   // true
res.JSON.Age.Raw() == ""               // true
res.JSON.Age.Raw() == respjson.Omitted // true
```

`.JSON` 구조체에는 응답 JSON에는 있지만 구조체에 정의되지 않은 속성을 담는 `ExtraFields` 맵도 포함됩니다. SDK에 아직 반영되지 않은 API 기능을 다룰 때 유용합니다.

```go
body := res.JSON.ExtraFields["my_unexpected_field"].Raw()
```

### 응답 유니온(Union)

응답에서 유니온(Union)은 각 객체 변형에 등장할 수 있는 모든 필드를 포함한 평탄화된 구조체로 표현됩니다. `.AsFooVariant()` 메서드나 `.AsAny()` 메서드(제공되는 경우)를 사용해 특정 변형으로 변환할 수 있습니다.

응답 값 유니온(Union)에 기본형 값이 포함된 경우, 해당 기본형 필드는 `Of` 접두사와 `json:"...,inline"` 태그를 가진 상태로 다른 속성들 사이에 배치됩니다.

```go
type AnimalUnion struct {
	// [Dog], [Cat] 변형에서 공통
	Owner Person `json:"owner"`
	// [Dog] 변형에서 사용
	DogBreed string `json:"dog_breed"`
	// [Cat] 변형에서 사용
	CatBreed string `json:"cat_breed"`
	// ...

	JSON struct {
		Owner respjson.Field
		// ...
	} `json:"-"`
}

// animal 변형 확인
if animal.Owner.Address.ZipCode == "" {
	panic("missing zip code")
}

// 변형에 따라 분기
switch variant := animal.AsAny().(type) {
case Dog:
case Cat:
default:
	panic("unexpected type")
}
```

### RequestOptions

이 라이브러리는 함수형 옵션 패턴을 사용합니다. `option` 패키지에 정의된 함수들은 `RequestConfig`를 변경하는 클로저인 `RequestOption`을 반환합니다. 이 옵션들은 클라이언트 전체 또는 개별 요청에 적용할 수 있습니다.

```go
client := upbit.NewClient(
	// 클라이언트의 모든 요청에 헤더 추가
	option.WithHeader("X-Some-Header", "custom_header_info"),
)

client.Accounts.List(context.TODO(), ...,
	// 헤더 덮어쓰기
	option.WithHeader("X-Some-Header", "some_other_custom_header_info"),
	// sjson 구문으로 요청 본문에 문서화되지 않은 필드 추가
	option.WithJSONSet("some.json.path", map[string]string{"my": "object"}),
)
```

디버깅할 때는 `option.WithDebugLog(nil)` 옵션이 유용합니다.

[전체 요청 옵션 목록](https://pkg.go.dev/github.com/upbit-official/upbit-sdk-go/option)을 참고하세요.

### 페이지네이션

이 라이브러리는 페이지네이션이 있는 목록 엔드포인트를 다루기 위한 편의 기능을 제공합니다.

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

또는 `.List()` 메서드로 단일 페이지를 가져온 뒤 `.GetNextPage()` 등의 헬퍼 메서드를 사용할 수 있습니다.

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

### 오류 처리

API가 성공이 아닌 상태 코드를 반환하면 `*upbit.Error` 타입의 오류가 반환됩니다. 여기에는 `StatusCode`, `*http.Request`, `*http.Response`, 그리고 오류 본문 JSON이 포함됩니다.

오류를 처리할 때는 `errors.As` 패턴을 사용하는 것을 권장합니다.

```go
_, err := client.Accounts.List(context.TODO())
if err != nil {
	var apierr *upbit.Error
	if errors.As(err, &apierr) {
		println(string(apierr.DumpRequest(true)))  // 직렬화된 HTTP 요청 출력
		println(string(apierr.DumpResponse(true))) // 직렬화된 HTTP 응답 출력
	}
	panic(err.Error()) // GET "/v1/accounts": 400 Bad Request { ... }
}
```

그 밖의 오류는 별도로 래핑되지 않고 그대로 반환됩니다. 예를 들어 HTTP 전송 실패 시에는 `*net.OpError`를 감싼 `*url.Error`가 반환될 수 있습니다.

### 타임아웃

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

### 파일 업로드

멀티파트 요청에서 파일 업로드에 해당하는 요청 파라미터는 `io.Reader` 타입입니다. `io.Reader`의 내용은 기본적으로 파일명 `anonymous_file`, content-type `application/octet-stream`으로 전송됩니다.

런타임 타입이 `Name() string` 또는 `ContentType() string`을 구현하면 파일명과 content-type을 사용자 정의할 수 있습니다. `os.File`은 `Name() string`을 구현하므로 `os.Open`으로 연 파일은 디스크의 파일명을 그대로 사용해 전송됩니다.

`upbit.File(reader io.Reader, filename string, contentType string)` 헬퍼를 사용하면 임의의 `io.Reader`에 파일명과 content-type을 지정할 수 있습니다.

### 재시도

일부 오류에 대해서는 짧은 지수 백오프를 적용해 기본적으로 최대 2회까지 자동 재시도합니다. 연결 오류, 408 Request Timeout, 409 Conflict, 429 Rate Limit, 500 이상 내부 오류가 재시도 대상입니다.

`WithMaxRetries` 옵션으로 재시도 설정을 변경하거나 비활성화할 수 있습니다.

```go
// 모든 요청의 기본값 설정:
client := upbit.NewClient(
	option.WithMaxRetries(0), // 기본값은 2
)

// 요청별 오버라이드:
client.Accounts.List(context.TODO(), option.WithMaxRetries(5))
```

### 원시 응답 데이터 접근 (예: 응답 헤더)

`option.WithResponseInto()` 요청 옵션을 사용하면 원시 HTTP 응답 데이터에 접근할 수 있습니다. 응답 헤더나 상태 코드 등을 확인할 때 유용합니다.

```go
// HTTP 응답을 저장할 변수 생성
var response *http.Response
accounts, err := client.Accounts.List(context.TODO(), option.WithResponseInto(&response))
if err != nil {
	// 오류 처리
}
fmt.Printf("%+v\n", accounts)

fmt.Printf("Status Code: %d\n", response.StatusCode)
fmt.Printf("Headers: %+#v\n", response.Header)
```

### 커스텀 요청 / 문서화되지 않은 엔드포인트 사용

이 라이브러리는 문서화된 API를 보다 편리하게 사용할 수 있도록 타입을 정의해 둔 SDK입니다. 문서화되지 않은 엔드포인트, 파라미터, 응답 속성도 라이브러리를 통해 사용할 수 있습니다.

#### 문서화되지 않은 엔드포인트

`client.Get`, `client.Post` 같은 HTTP 메서드를 직접 사용해 문서화되지 않은 엔드포인트에 요청할 수 있습니다. 이때도 재시도 등 클라이언트에 설정된 `RequestOptions`가 동일하게 적용됩니다.

```go
var (
    // params는 io.Reader, []byte, encoding/json 직렬화 가능 객체,
    // 또는 이 라이브러리에 정의된 "…Params" 구조체일 수 있습니다.
    params map[string]any

    // result는 []byte, *http.Response, encoding/json 역직렬화 가능 객체,
    // 또는 이 라이브러리에 정의된 모델일 수 있습니다.
    result *http.Response
)
err := client.Post(context.Background(), "/unspecified", params, &result)
if err != nil {
    …
}
```

#### 문서화되지 않은 요청 파라미터

`option.WithQuery()` 또는 `option.WithJSONSet()` 메서드를 사용하면 문서화되지 않은 파라미터를 포함한 요청을 보낼 수 있습니다.

```go
params := FooNewParams{
    ID:   "id_xxxx",
    Data: FooNewParamsData{
        FirstName: upbit.String("John"),
    },
}
client.Foo.New(context.Background(), params, option.WithJSONSet("data.last_name", "Doe"))
```

#### 문서화되지 않은 응답 속성

`result.JSON.RawJSON()`으로 응답 전체의 raw JSON 문자열에 접근하거나, `result.JSON.Foo.Raw()`로 특정 필드의 raw JSON 값을 가져올 수 있습니다.

응답 구조체에 정의되지 않은 필드도 보존되며, `result.JSON.ExtraFields`를 통해 `map[string]Field` 형태로 접근할 수 있습니다.

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

`option.WithHTTPClient(client)`를 사용하면 기본 `http.Client`를 교체할 수도 있습니다. HTTP 클라이언트는 하나만 지정할 수 있으며(마지막에 지정한 값이 적용됨), 미들웨어가 적용된 뒤 요청을 전달받습니다.

## 버전 관리

이 패키지는 일반적으로 [SemVer](https://semver.org/spec/v2.0.0.html) 규약을 따르지만, 일부 하위 호환성에 영향을 줄 수 있는 변경이 마이너 버전에 포함될 수 있습니다.

1. 기술적으로는 공개되어 있지만 외부 사용을 의도하지 않았거나 문서화되지 않은 라이브러리 내부 변경
2. 대부분의 사용자에게 실질적인 영향이 없을 것으로 판단되는 변경

하위 호환성은 중요하게 고려하고 있으며, 원활한 업그레이드 경험을 제공하기 위해 노력하고 있습니다.

피드백을 환영합니다. 질문, 버그 제보, 개선 제안이 있다면 open-api@upbit.com으로 연락해 주세요.

## 기여

현재 업비트 SDK는 초기 출시 단계로, 외부 Issue 등록 및 PR 기여는 아직 운영하고 있지 않습니다.
버그 제보나 개선 의견은 개발자 지원 채널(open-api@upbit.com)을 통해 전달해 주시기 바랍니다.
외부 기여 채널은 향후 운영 안정화 상황에 맞춰 순차적으로 오픈하는 방안을 검토하고 있습니다.

© 2026 Dunamu Inc. All rights reserved.
