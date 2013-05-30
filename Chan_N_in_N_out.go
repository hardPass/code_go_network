package main

import (
	"fmt"
	"strconv"
	"sync/atomic"
	"time"
)

var flag int64 = 0

var exitChan chan int = make(chan int)

func ChnOut(n int, c <-chan string) {
	for {
		data, ok := <-c

		if !ok {
			return
		}

		fmt.Println(n, "--", data)

	}

}

func ChnIn(n int, c chan<- string) {
	for i := 0; i < 20; i++ {
		if atomic.LoadInt64(&flag) == 1 {
			close(c)
			return
		}

		data := strconv.Itoa(n) + "-" + strconv.Itoa(i)
		c <- data

		time.Sleep(time.Millisecond * 300)
	}

}

func main() {
	fmt.Println("Started...")

	c := make(chan string, 10)

	go ChnOut(1, c)
	go ChnOut(2, c)
	go ChnOut(3, c)

	go ChnIn(1, c)
	go ChnIn(2, c)
	go ChnIn(3, c)

	time.Sleep(time.Millisecond * 2000)

	fmt.Println("Prepare to done..")

	atomic.CompareAndSwapInt64(&flag, 0, 1)

	fmt.Println("Done!")

}
