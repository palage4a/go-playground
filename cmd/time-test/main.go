package main

import (
	"fmt"
	"time"
)

func main() {
	a := time.Now()
	for i := 0; i < 1000000; i++ {
		a = time.Now()
	}
	fmt.Println(a)
}
