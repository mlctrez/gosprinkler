package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("starting")
	for {
		// fmt.Println(time.Now())
		time.Sleep(60 * time.Second)
	}
}
