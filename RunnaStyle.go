package main

import (
  "fmt"
	"sync/atomic"
	"time"
)

type Runnable func() bool // true: performed,   false: not performed

type RunnaStyle struct {
	RunInterval time.Duration
	maxidle     int32
	sem         int32
}

func (r *RunnaStyle) ShutDown() {
	atomic.StoreInt32(&r.sem, -1)
}

func (r *RunnaStyle) Run(f Runnable) {
	if atomic.CompareAndSwapInt32(&r.sem, 0, 1) {
		go func() {
			ticker := time.NewTicker(r.RunInterval)
			defer ticker.Stop()

			fmt.Println("let's run")
			var idle int32 = 0
			for {
				if atomic.LoadInt32(&r.sem) == -1 {
					return
				}
				<-ticker.C

				if f() {
					idle = 0
				} else {
					idle++
				}

				if idle >= r.maxidle {
					fmt.Println("idle for", idle)
					break
				}
			}

			atomic.CompareAndSwapInt32(&r.sem, 1, 0)
		}()
	} else {
		fmt.Println("already running")
	}
}

// for test -----------------------------------------
type T struct {
	count int32
	msg   string
	r     *RunnaStyle
}

func NewT() *T {
	return &T{
		r: &RunnaStyle{
			RunInterval: 1 * time.Second,
			maxidle:     4,
		},
		msg: "name--",
	}
}

func (t *T) Run() {
	t.r.Run(func() bool {
		t.count++
		fmt.Println("---", t.msg, t.count)
		return true
	})
}

func main() {
	t := NewT()
	time.AfterFunc(1*time.Second, func() {
		t.Run()
	})

	time.AfterFunc(6*time.Second, func() {
		t.Run()
	})

	c := make(chan (int))
	time.AfterFunc(10*time.Second, func() {
		t.r.ShutDown()
		fmt.Println("ShutDown")
		c <- 1
	})

	<-c

	time.Sleep(300 * time.Nanosecond)
	fmt.Println("game over")
}
