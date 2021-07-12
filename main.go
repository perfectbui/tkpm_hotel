package main

import (
	"fmt"
	"time"
)

func main() {
	defer CalculateExecutionTime("alo")()
	fmt.Println("b")
}

// CalculateExecutionTime ...
func CalculateExecutionTime(fnName string) func() {
	start := time.Now()
	fmt.Println("a")
	return func() {
		fmt.Println("c")
		fmt.Printf("%s took %vs\n", fnName, time.Since(start).Seconds())
	}
}
