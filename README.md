# Upbit SDK for Go

English | [한국어](./README_KR.md)

<!-- x-release-please-start-version -->

<a href="https://pkg.go.dev/github.com/upbit-official/upbit-sdk-go"><img src="https://pkg.go.dev/badge/github.com/upbit-official/upbit-sdk-go.svg" alt="Go Reference"></a>

<!-- x-release-please-end -->

The Upbit SDK for Go is designed to make it easy to use the Upbit REST API in Go 1.25+.
It includes type definitions for all request parameters and response fields. For the complete SDK API reference, see [api.md](api.md).

## Installation

<!-- x-release-please-start-version -->

```go
import (
	"github.com/upbit-official/upbit-sdk-go" // imported as upbit
)
```

<!-- x-release-please-end -->

Or to pin the version:

<!-- x-release-please-start-version -->

```sh
go get -u 'github.com/upbit-official/upbit-sdk-go@v1.0.0'
```

<!-- x-release-please-end -->

## Usage

The full API of this SDK can be found in [api.md](api.md).

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
		// option.WithEnvironmentSg(), // or WithEnvironmentKr() | WithEnvironmentId() | WithEnvironmentTh(); defaults to WithEnvironmentKr()
	)
	accounts, err := client.Accounts.List(context.TODO())
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("%+v\n", accounts)
}
```

For more runnable examples, see the scripts in [`examples/`](examples/).

## Request & Response types

This SDK includes Go types for all documented request params and response fields. You may import and use them like so:

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

Required request fields are ordinary value types and are tagged with <code>`api:"required"`</code>. Optional primitive fields are wrapped in `param.Opt[T]` and can be set with constructors such as `upbit.String(string)`, `upbit.Int(int64)`, `upbit.Bool(bool)`, and `upbit.Float(float64)`. The SDK uses the [`omitzero`](https://tip.golang.org/doc/go1.24#encodingjsonpkgencodingjson) semantics from `encoding/json` for omitted request fields.

Request body params use the same pattern:

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

All fields in response structs are ordinary value types. Each response struct also includes a `JSON` field with metadata for each property. Use `.Valid()` to distinguish an explicit value from `null`, an omitted field, or a value that could not be decoded.

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

Some responses are unions. For example, a `DepositNewCoinAddressResponseUnion` value returned from `client.Deposits.NewCoinAddress` can be converted to one of its documented variants with the generated `As...` methods.

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

Documentation for each method, request param, and response field is available in Go doc comments and on [pkg.go.dev](https://pkg.go.dev/github.com/upbit-official/upbit-sdk-go).

## Handling errors

When the API returns a non-success status code, the SDK returns an error of type `*upbit.Error`. This contains the `StatusCode`, `*http.Request`, and `*http.Response` values of the request, as well as the JSON of the error body.

To handle errors, use the `errors.As` pattern:

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

When other errors occur, they are returned unwrapped. For example, if HTTP transport fails, you might receive `*url.Error` wrapping `*net.OpError`.

## Retries

Certain errors will be automatically retried 2 times by default, with a short exponential backoff.
We retry by default all connection errors, 408 Request Timeout, 409 Conflict, 429 Rate Limit, and >=500 Internal errors.

You can use the `WithMaxRetries` option to configure or disable this:

```go
// Configure the default for all requests:
client := upbit.NewClient(
	option.WithMaxRetries(0), // default is 2
)

// Override per-request:
client.Accounts.List(context.TODO(), option.WithMaxRetries(5))
```

## Timeouts

Requests do not time out by default. Use context to configure a timeout for a request lifecycle.

Note that if a request is [retried](#retries), the context timeout does not start over.
To set a per-retry timeout, use `option.WithRequestTimeout()`.

```go
// This sets the timeout for the request, including all the retries.
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
defer cancel()
client.Accounts.List(
	ctx,
	// This sets the per-retry timeout.
	option.WithRequestTimeout(20*time.Second),
)
```

## Pagination

This SDK provides conveniences for working with paginated list endpoints.

You can use `.ListAutoPaging()` methods to iterate through items across all pages:

```go
iter := client.Orders.ListOpenAutoPaging(context.TODO(), upbit.OrderListOpenParams{})
// Automatically fetches more pages as needed.
for iter.Next() {
	order := iter.Current()
	fmt.Printf("%+v\n", order)
}
if err := iter.Err(); err != nil {
	panic(err.Error())
}
```

Or you can use simple `.List()` methods to fetch a single page and receive a standard response object with helper methods like `.GetNextPage()`:

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

## Advanced Usage

### RequestOptions

This SDK uses the functional options pattern. Functions defined in the `option` package return a `RequestOption`, which is a closure that mutates a `RequestConfig`. These options can be supplied to the client or to individual requests.

```go
client := upbit.NewClient(
	// Adds a header to every request made by the client.
	option.WithHeader("X-Some-Header", "custom_header_info"),
)

_, err := client.Accounts.List(
	context.TODO(),
	// Overrides the header for this request.
	option.WithHeader("X-Some-Header", "some_other_custom_header_info"),
)
```

The request option `option.WithDebugLog(nil)` may be helpful while debugging.

See the [full list of request options](https://pkg.go.dev/github.com/upbit-official/upbit-sdk-go/option).

### Accessing raw response data

You can access raw HTTP response data by using the `option.WithResponseInto()` request option. This is useful when you need to examine response headers, status codes, or other details.

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

### Making custom/undocumented requests

This SDK is typed for convenient access to the documented API. If you need to access undocumented endpoints, params, or response properties, the SDK can still be used.

#### Undocumented endpoints

To make requests to undocumented endpoints, you can use `client.Get`, `client.Post`, and other HTTP verbs. `RequestOptions` on the client, such as retries, will be respected when making these requests.

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

#### Undocumented request params

To send undocumented query parameters, use `option.WithQuery()` on a real SDK method call:

```go
candles, err := client.Candles.ListDays(
	context.Background(),
	upbit.CandleListDaysParams{
		Market: "KRW-BTC",
	},
	option.WithQuery("include_experimental", "true"),
)
```

To send undocumented JSON body parameters, use `option.WithJSONSet()`:

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

#### Undocumented response properties

To access undocumented response properties, use `RawJSON()` for the full raw JSON payload or `JSON.ExtraFields` for fields not represented on the response struct.

```go
if len(*accounts) > 0 {
	account := (*accounts)[0]
	fmt.Println(account.RawJSON())
	fmt.Println(account.JSON.ExtraFields["new_field"].Raw())
}
```

### Middleware

We provide `option.WithMiddleware`, which applies the given middleware to requests.

```go
func Logger(req *http.Request, next option.MiddlewareNext) (res *http.Response, err error) {
	// Before the request.
	start := time.Now()
	LogReq(req)

	// Forward the request to the next handler.
	res, err = next(req)

	// Handle work after the request.
	end := time.Now()
	LogRes(res, err, end.Sub(start))

	return res, err
}

client := upbit.NewClient(
	option.WithMiddleware(Logger),
)
```

When multiple middlewares are provided as variadic arguments, the middlewares are applied left to right. If `option.WithMiddleware` is given multiple times, for example first in the client then the method, the middleware in the client will run first and the middleware given in the method will run next.

You may also replace the default `http.Client` with `option.WithHTTPClient(client)`. Only one HTTP client is accepted, overwriting any previous client, and receives requests after any middleware has been applied.

## Versioning

This package generally follows [SemVer](https://semver.org/spec/v2.0.0.html) conventions, though certain backwards-incompatible changes may be released as minor versions:

1. Changes to SDK internals which are technically public but not intended or documented for external use.
2. Changes that we do not expect to impact the vast majority of users in practice.

We take backwards-compatibility seriously and work hard to ensure you can rely on a smooth upgrade experience.

We are keen for your feedback; please contact us at open-api@upbit.com with questions, bugs, or suggestions.

## Requirements

This SDK requires Go 1.25+.

## Contribution

The Upbit SDK is in its initial release phase, and public contributions (Issues/PRs) are currently closed.
For bug reports and feedback, please email open-api@upbit.com.
We are considering opening external contribution channels in phases as the SDK becomes more stable.

© 2026 Dunamu Inc. All rights reserved.
