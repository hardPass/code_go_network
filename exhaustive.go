package main

import (
	"fmt"
)

// 打印 A B C D 所有排列组合

func main() {
	count := 0
	a := [4]string{}
	for i := 0; i < 4; i++ {
		a[i] = "A"
		for j := 0; j < 4; j++ {
			if i == j {
				continue
			}
			a[j] = "B"
			for k := 0; k < 4; k++ {

				if i == j || j == k || i == k {
					continue
				}
				a[k] = "C"

				for m := 0; m < 4; m++ {
					if i == j || j == k || i == k || j == m || m == i || m == k {
						continue
					}
					a[m] = "D"
					count++
					fmt.Println(count, a)
				}
			}
		}
	}
	fmt.Println("Done.")
}
