package main

import (
	"log"
)

type Context struct {
	name string
	// and other properties
}

type FilterChain struct {
	idx      int
	filters  []Filter
	endpoint func(*Context)
}

func NewFilterChain(filters []Filter, endpoint func(*Context)) *FilterChain {
	return &FilterChain{-1, filters, endpoint}
}

func (chain *FilterChain) DoFilter(ctx *Context) {
	chain.idx++
	if chain.idx >= len(chain.filters) {
		chain.endpoint(ctx)
		return
	}

	chain.filters[chain.idx].DoFilter(ctx, chain)
}

type Filter interface {
	DoFilter(*Context, *FilterChain)
}

type LogFilter struct{ seq int }

func (f *LogFilter) DoFilter(ctx *Context, chain *FilterChain) {
	log.Println(" inbound filter", f.seq)
	// TODO dofilter request
	if f.seq == 33 { // for test
		log.Println("+++++++++++++ breaking return ", f.seq)
	} else {
		chain.DoFilter(ctx)
	}

	// TODO dofilter respone
	log.Println("outbound filter", f.seq)
}

func main() {
	var filters []Filter
	filters = append(filters, &LogFilter{1})
	filters = append(filters, &LogFilter{2})
	filters = append(filters, &LogFilter{3})
	// filters = append(filters, &LogFilter{33}) // try this
	filters = append(filters, &LogFilter{4})
	filters = append(filters, &LogFilter{5})

	chain := NewFilterChain(filters, func(ctx *Context) {
		log.Println("----simulate servlet --------", ctx.name)
	})

	chain.DoFilter(&Context{"this a virtual ctx"})

}
