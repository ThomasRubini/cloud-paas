package entrypoint

import "fmt"

func businessLogic(a, b int) int {
	return a + b
}

func Entrypoint() {
	fmt.Println("Hello, World!")
	fmt.Println("1+2=", businessLogic(1, 2))
}