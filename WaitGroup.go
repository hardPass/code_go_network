package main

import (
  "fmt"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

var flag int64 = 0

func ChnOut(wg *sync.WaitGroup, c <-chan string) {
	defer wg.Done()

	for {
		if atomic.LoadInt64(&flag) == 1 {
			return
		}

		fmt.Println("--", <-c)

	}

}

func ChnIn(wg *sync.WaitGroup, c chan string) {
	defer wg.Done()

	for i := 0; i < 20; i++ {
		if atomic.LoadInt64(&flag) == 1 {
			return
		}

		data := "d--" + strconv.Itoa(i)
		c <- data

		time.Sleep(time.Millisecond * 300)
	}

}

func main() {
	fmt.Println("Started...")
	wg := new(sync.WaitGroup)

	c := make(chan string, 10)

	wg.Add(1)
	go ChnOut(wg, c)

	wg.Add(1)
	go ChnIn(wg, c)

	wg.Add(1)
	go ChnIn(wg, c)

	time.Sleep(time.Millisecond * 100)

	fmt.Println("Prepare to done..")

	atomic.CompareAndSwapInt64(&flag, 0, 1)
	wg.Wait()
	close(c)

	fmt.Println("Done!")

	select {}
}
