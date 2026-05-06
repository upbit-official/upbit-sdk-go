package constant

import (
	shimjson "github.com/upbit-official/upbit-sdk-go/internal/encoding/json"
)

type Constant[T any] interface {
	Default() T
}

// ValueOf gives the default value of a constant from its type. It's helpful when
// constructing constants as variants in a one-of. Note that empty structs are
// marshalled by default. Usage: constant.ValueOf[constant.Foo]()
func ValueOf[T Constant[T]]() T {
	var t T
	return t.Default()
}

type MyAsset string   // Always "myAsset"
type MyOrder string   // Always "myOrder"
type Orderbook string // Always "orderbook"
type Ticker string    // Always "ticker"
type Trade string     // Always "trade"

func (c MyAsset) Default() MyAsset     { return "myAsset" }
func (c MyOrder) Default() MyOrder     { return "myOrder" }
func (c Orderbook) Default() Orderbook { return "orderbook" }
func (c Ticker) Default() Ticker       { return "ticker" }
func (c Trade) Default() Trade         { return "trade" }

func (c MyAsset) MarshalJSON() ([]byte, error)   { return marshalString(c) }
func (c MyOrder) MarshalJSON() ([]byte, error)   { return marshalString(c) }
func (c Orderbook) MarshalJSON() ([]byte, error) { return marshalString(c) }
func (c Ticker) MarshalJSON() ([]byte, error)    { return marshalString(c) }
func (c Trade) MarshalJSON() ([]byte, error)     { return marshalString(c) }

type constant[T any] interface {
	Constant[T]
	*T
}

func marshalString[T ~string, PT constant[T]](v T) ([]byte, error) {
	var zero T
	if v == zero {
		v = PT(&v).Default()
	}
	return shimjson.Marshal(string(v))
}
