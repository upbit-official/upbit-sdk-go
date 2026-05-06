// Auto-generated example — do not edit manually.
// Create Order (POST /v1/orders)
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/upbit-official/upbit-sdk-go"
	"github.com/upbit-official/upbit-sdk-go/option"
)

// default
func exampleNewOrder() {
	client := upbit.NewClient(
		option.WithAccessKey(os.Getenv("UPBIT_ACCESS_KEY")),
		option.WithSecretKey(os.Getenv("UPBIT_SECRET_KEY")),
	)
	params := upbit.OrderNewParams{
		Market:  "KRW-BTC",
		Side:    upbit.SideEnumBid,
		Volume:  upbit.String("1"),
		Price:   upbit.String("14000000"),
		OrdType: upbit.OrdTypeEnumLimit,
	}
	res, err := client.Orders.New(context.TODO(), params)
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
}

// market_sell
func exampleNewOrderMarketSell() {
	client := upbit.NewClient(
		option.WithAccessKey(os.Getenv("UPBIT_ACCESS_KEY")),
		option.WithSecretKey(os.Getenv("UPBIT_SECRET_KEY")),
	)
	params := upbit.OrderNewParams{
		Market:  "KRW-BTC",
		Side:    upbit.SideEnumAsk,
		Volume:  upbit.String("0.001"),
		OrdType: upbit.OrdTypeEnumMarket,
	}
	res, err := client.Orders.New(context.TODO(), params)
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
}

// market_buy
func exampleNewOrderMarketBuy() {
	client := upbit.NewClient(
		option.WithAccessKey(os.Getenv("UPBIT_ACCESS_KEY")),
		option.WithSecretKey(os.Getenv("UPBIT_SECRET_KEY")),
	)
	params := upbit.OrderNewParams{
		Market:  "KRW-BTC",
		Side:    upbit.SideEnumBid,
		Price:   upbit.String("10000"),
		OrdType: upbit.OrdTypeEnumPrice,
	}
	res, err := client.Orders.New(context.TODO(), params)
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
}

// best_buy
func exampleNewOrderBestBuy() {
	client := upbit.NewClient(
		option.WithAccessKey(os.Getenv("UPBIT_ACCESS_KEY")),
		option.WithSecretKey(os.Getenv("UPBIT_SECRET_KEY")),
	)
	params := upbit.OrderNewParams{
		Market:      "KRW-BTC",
		Side:        upbit.SideEnumBid,
		Price:       upbit.String("10000"),
		OrdType:     upbit.OrdTypeEnumBest,
		TimeInForce: upbit.TimeInForceEnumIoc,
	}
	res, err := client.Orders.New(context.TODO(), params)
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
}

// best_sell
func exampleNewOrderBestSell() {
	client := upbit.NewClient(
		option.WithAccessKey(os.Getenv("UPBIT_ACCESS_KEY")),
		option.WithSecretKey(os.Getenv("UPBIT_SECRET_KEY")),
	)
	params := upbit.OrderNewParams{
		Market:      "KRW-BTC",
		Side:        upbit.SideEnumAsk,
		Volume:      upbit.String("0.001"),
		OrdType:     upbit.OrdTypeEnumBest,
		TimeInForce: upbit.TimeInForceEnumIoc,
	}
	res, err := client.Orders.New(context.TODO(), params)
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
}

// limit_ioc
func exampleNewOrderLimitIoc() {
	client := upbit.NewClient(
		option.WithAccessKey(os.Getenv("UPBIT_ACCESS_KEY")),
		option.WithSecretKey(os.Getenv("UPBIT_SECRET_KEY")),
	)
	params := upbit.OrderNewParams{
		Market:      "KRW-BTC",
		Side:        upbit.SideEnumBid,
		Volume:      upbit.String("0.001"),
		Price:       upbit.String("10000000"),
		OrdType:     upbit.OrdTypeEnumLimit,
		TimeInForce: upbit.TimeInForceEnumIoc,
	}
	res, err := client.Orders.New(context.TODO(), params)
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
}

// limit_smp_reduce
func exampleNewOrderLimitSmpReduce() {
	client := upbit.NewClient(
		option.WithAccessKey(os.Getenv("UPBIT_ACCESS_KEY")),
		option.WithSecretKey(os.Getenv("UPBIT_SECRET_KEY")),
	)
	params := upbit.OrderNewParams{
		Market:  "KRW-BTC",
		Side:    upbit.SideEnumBid,
		Volume:  upbit.String("0.001"),
		Price:   upbit.String("10000000"),
		OrdType: upbit.OrdTypeEnumLimit,
		SmpType: upbit.SmpTypeEnumReduce,
	}
	res, err := client.Orders.New(context.TODO(), params)
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
}

func main() {
	exampleNewOrder()
	exampleNewOrderMarketSell()
	exampleNewOrderMarketBuy()
	exampleNewOrderBestBuy()
	exampleNewOrderBestSell()
	exampleNewOrderLimitIoc()
	exampleNewOrderLimitSmpReduce()
}
