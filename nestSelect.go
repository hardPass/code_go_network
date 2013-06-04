package main

import (
  "fmt"
	"runtime"
	"time"
)

func tryPlain() {
	runtime.GOMAXPROCS(3)

	highTick := time.Tick(1e8 * 5)
	lowTick := time.Tick(1e9 * 1) // for remove intervals
	for {
		select {
		case <-highTick:
			fmt.Println("hello, high precedence")
		case <-lowTick:
			fmt.Println("hello, low precedence")
		default:
		}
	}
}

func tryNest() {
	runtime.GOMAXPROCS(3)

	highTick := time.Tick(1e8 * 5)
	lowTick := time.Tick(1e9 * 1) // for remove intervals
	for {
		select {
		case <-highTick:
			fmt.Println("hello, high precedence")
		default:
			//fmt.Println("low block-------------------------")
			select {
			case <-lowTick:
				fmt.Println("hello, low precedence")
			default:
			}
		}
	}
}

func main() {
	//tryPlain()
	tryNest()
}
