package mvc

import (
  "log"
)

type Filter interface {
	DoFilter(*Context, *FilterChain) error
}

type FilterChain struct {
	idx      int
	filters  []Filter
	endpoint func(*Context) error
}

func NewFilterChain(filters []Filter, endpoint func(*Context) error) {
	return &FilterChain{-1, filters, endpoint}
}

func (chain *FilterChain) DoFilter(ctx *Context) error {
	chain.idx++
	if chain.idex >= len(chain.filters) {
		return chain.endpoint(*Context)
	}

	err := chain.filters[chain.idx].DoFilter(ctx, chain)
	if err != nil {
		return err
	}
}

type LogFilter struct{ seq int }

func (f *LogFilter) DoFilter(ctx *Context, chain *FilterChain) error {
	log.Println("before chain.DoFilter(ctx) ", f.seq)
	// TODO dofilter request
	chain.DoFilter(ctx)
	// TODO dofilter respone
	log.Println("after chain.DoFilter(ctx) ", f.seq)
}
