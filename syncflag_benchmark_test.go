package syncflag

import (
	"container/list"
	// "log"
	"sync"
	"sync/atomic"
	"testing"
)

var (
	nLen = 100000
)

type Client struct {
	validChan chan bool
	valid     int32
}

func BenchmarkFlagByChanInit(b *testing.B) {
	b.StartTimer()
	clients := list.New()

	for i := 0; i < nLen; i++ {
		c := &Client{validChan: make(chan bool)}
		clients.PushBack(c)
		close(c.validChan)
	}
	// log.Println("clients.Len():", clients.Len())
}

func BenchmarkFlagByChanCheck(b *testing.B) {
	b.StopTimer()

	clients := list.New()

	for i := 0; i < nLen; i++ {
		c := &Client{validChan: make(chan bool)}
		clients.PushBack(c)
		close(c.validChan)
	}

	b.StartTimer()
	wg := &sync.WaitGroup{}
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			count := 0

			e := clients.Front()
			for e != nil {
				c := e.Value.(*Client)
				_, ok := <-c.validChan
				if !ok {
					count++
				}
				e = e.Next()
			}
			wg.Done()
		}()
	}

	wg.Wait()
}

func BenchmarkFlagByAtomicInit(b *testing.B) {
	b.StartTimer()

	clients := list.New()

	for i := 0; i < nLen; i++ {
		c := &Client{valid: 0}
		clients.PushBack(c)
		atomic.CompareAndSwapInt32(&c.valid, 0, 1)
	}
	// log.Println("clients.Len():", clients.Len())
}

func BenchmarkFlagByAtomicCheck(b *testing.B) {
	b.StopTimer()

	clients := list.New()

	for i := 0; i < nLen; i++ {
		c := &Client{valid: 0}
		clients.PushBack(c)
		atomic.CompareAndSwapInt32(&c.valid, 0, 1)
	}

	b.StartTimer()
	wg := &sync.WaitGroup{}

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			count := 0

			e := clients.Front()
			for e != nil {
				c := e.Value.(*Client)

				if atomic.LoadInt32(&c.valid) == 1 {
					count++
				}
				e = e.Next()
			}
			wg.Done()
		}()

	}

	wg.Wait()
}
