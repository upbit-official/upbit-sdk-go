package pagination

import (
	"fmt"
	"iter"
	"net/http"
	"reflect"
	"strconv"

	"github.com/upbit-official/upbit-sdk-go/internal/apijson"
	"github.com/upbit-official/upbit-sdk-go/internal/requestconfig"
	"github.com/upbit-official/upbit-sdk-go/option"
	"github.com/upbit-official/upbit-sdk-go/packages/param"
	"github.com/upbit-official/upbit-sdk-go/packages/respjson"
)

// aliased to make [param.APIUnion] private when embedding
type paramUnion = param.APIUnion

// aliased to make [param.APIObject] private when embedding
type paramObj = param.APIObject

type PageNumberPage[T any] struct {
	Items []T `json:",inline"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Items       respjson.Field
		ExtraFields map[string]respjson.Field
		Raw         string
	} `json:"-"`
	cfg *requestconfig.RequestConfig
	res *http.Response
}

// RawJSON returns the unmodified JSON received from the API.
func (r PageNumberPage[T]) RawJSON() string { return r.JSON.Raw }
func (r *PageNumberPage[T]) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// GetNextPage returns the next page as defined by this pagination style. When
// there is no next page, this function will return a 'nil' for the page value, but
// will not return an error
func (r *PageNumberPage[T]) GetNextPage() (res *PageNumberPage[T], err error) {
	if len(r.Items) == 0 {
		return nil, nil
	}
	u := r.cfg.Request.URL
	currentPage, err := strconv.ParseInt(u.Query().Get("page"), 10, 64)
	if err != nil {
		currentPage = 1
	}
	cfg := r.cfg.Clone(r.cfg.Context)
	query := cfg.Request.URL.Query()
	query.Set("page", fmt.Sprintf("%d", currentPage+1))
	cfg.Request.URL.RawQuery = query.Encode()
	var raw *http.Response
	cfg.ResponseInto = &raw
	cfg.ResponseBodyInto = &res
	err = cfg.Execute()
	if err != nil {
		return nil, err
	}
	res.SetPageConfig(cfg, raw)
	return res, nil
}

func (r *PageNumberPage[T]) SetPageConfig(cfg *requestconfig.RequestConfig, res *http.Response) {
	if r == nil {
		r = &PageNumberPage[T]{}
	}
	r.cfg = cfg
	r.res = res
}

type PageNumberPageAutoPager[T any] struct {
	page *PageNumberPage[T]
	cur  T
	idx  int
	run  int
	err  error
	paramObj
}

func NewPageNumberPageAutoPager[T any](page *PageNumberPage[T], err error) *PageNumberPageAutoPager[T] {
	return &PageNumberPageAutoPager[T]{
		page: page,
		err:  err,
	}
}

func (r *PageNumberPageAutoPager[T]) Next() bool {
	if r.page == nil || len(r.page.Items) == 0 {
		return false
	}
	if r.idx >= len(r.page.Items) {
		r.idx = 0
		r.page, r.err = r.page.GetNextPage()
		if r.err != nil || r.page == nil || len(r.page.Items) == 0 {
			return false
		}
	}
	r.cur = r.page.Items[r.idx]
	r.run += 1
	r.idx += 1
	return true
}

func (r *PageNumberPageAutoPager[T]) Current() T {
	return r.cur
}

func (r *PageNumberPageAutoPager[T]) Err() error {
	return r.err
}

func (r *PageNumberPageAutoPager[T]) Index() int {
	return r.run
}

func (r *PageNumberPageAutoPager[T]) All() iter.Seq[T] {
	return func(yield func(T) bool) {
		for r.Next() {
			if !yield(r.Current()) {
				break
			}
		}
	}
}

type UuidCursorIDPage[T any] struct {
	Items []T `json:",inline"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Items       respjson.Field
		ExtraFields map[string]respjson.Field
		Raw         string
	} `json:"-"`
	cfg *requestconfig.RequestConfig
	res *http.Response
}

// RawJSON returns the unmodified JSON received from the API.
func (r UuidCursorIDPage[T]) RawJSON() string { return r.JSON.Raw }
func (r *UuidCursorIDPage[T]) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// GetNextPage returns the next page as defined by this pagination style. When
// there is no next page, this function will return a 'nil' for the page value, but
// will not return an error
func (r *UuidCursorIDPage[T]) GetNextPage() (res *UuidCursorIDPage[T], err error) {
	if len(r.Items) == 0 {
		return nil, nil
	}
	items := r.Items
	cfg := r.cfg.Clone(r.cfg.Context)
	value := reflect.ValueOf(items[len(items)-1])
	field := value.FieldByName("Uuid")
	err = cfg.Apply(option.WithQuery("to", field.Interface().(string)))
	if err != nil {
		return nil, err
	}
	var raw *http.Response
	cfg.ResponseInto = &raw
	cfg.ResponseBodyInto = &res
	err = cfg.Execute()
	if err != nil {
		return nil, err
	}
	res.SetPageConfig(cfg, raw)
	return res, nil
}

func (r *UuidCursorIDPage[T]) SetPageConfig(cfg *requestconfig.RequestConfig, res *http.Response) {
	if r == nil {
		r = &UuidCursorIDPage[T]{}
	}
	r.cfg = cfg
	r.res = res
}

type UuidCursorIDPageAutoPager[T any] struct {
	page *UuidCursorIDPage[T]
	cur  T
	idx  int
	run  int
	err  error
	paramObj
}

func NewUuidCursorIDPageAutoPager[T any](page *UuidCursorIDPage[T], err error) *UuidCursorIDPageAutoPager[T] {
	return &UuidCursorIDPageAutoPager[T]{
		page: page,
		err:  err,
	}
}

func (r *UuidCursorIDPageAutoPager[T]) Next() bool {
	if r.page == nil || len(r.page.Items) == 0 {
		return false
	}
	if r.idx >= len(r.page.Items) {
		r.idx = 0
		r.page, r.err = r.page.GetNextPage()
		if r.err != nil || r.page == nil || len(r.page.Items) == 0 {
			return false
		}
	}
	r.cur = r.page.Items[r.idx]
	r.run += 1
	r.idx += 1
	return true
}

func (r *UuidCursorIDPageAutoPager[T]) Current() T {
	return r.cur
}

func (r *UuidCursorIDPageAutoPager[T]) Err() error {
	return r.err
}

func (r *UuidCursorIDPageAutoPager[T]) Index() int {
	return r.run
}

func (r *UuidCursorIDPageAutoPager[T]) All() iter.Seq[T] {
	return func(yield func(T) bool) {
		for r.Next() {
			if !yield(r.Current()) {
				break
			}
		}
	}
}
