package main

import (
  "fmt"
	"runtime"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

var (
	flag int64 = 0
	rw   sync.RWMutex

	c = make(chan string)
)

func Recv(n int) {
	for {
		select {
		case data, ok := <-c:
			if !ok {
				fmt.Println("Receiver ", n, " left.")
				return
			}

			fmt.Println("Receiver ", n, " received:", data)
		}
		//runtime.Gosched()
	}

}

func Send(n int) {
	for i := 0; ; i++ {
		if atomic.LoadInt64(&flag) == 1 {
			fmt.Println("Sender ", n, " left.")
			return
		}
		time.Sleep(time.Millisecond * 100)
		MustSend("<data " + strconv.Itoa(i) + "> from Sender" + strconv.Itoa(n))
		// select {
		// case <-exit:
		// 	fmt.Println("Sender ", n, " left.")
		// 	return
		// default:
		// 	data := "<data " + strconv.Itoa(i) + "> from Sender" + strconv.Itoa(n)
		// 	c <- data

		// 	time.Sleep(time.Millisecond * 100)
		// }
	}

}

func MustSend(data string) {
	rw.RLock()
	defer rw.RUnlock()

	if atomic.LoadInt64(&flag) == 0 {
		c <- data
	}

}

func Close() {
	rw.Lock()
	defer rw.Unlock()

	atomic.CompareAndSwapInt64(&flag, 0, 1)
	close(c)
}

func main() {
	runtime.GOMAXPROCS(4)
	fmt.Println("Started...")

	go Recv(1)
	go Recv(2)
	go Recv(3)

	go Send(1)
	go Send(2)
	go Send(3)

	time.Sleep(time.Millisecond * 1000)

	fmt.Println("Prepare to done..")
	Close()

	fmt.Println("Done!")
	time.Sleep(1e10 * 2)

}
