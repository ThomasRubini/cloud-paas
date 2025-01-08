package main

import "fmt"

func businessLogic(a, b int) int {
	return a + b
}

func main() {
	fmt.Println("CLI main")
	fmt.Println("Hello, World!")
	fmt.Println("1+2=", businessLogic(1, 2))
}
