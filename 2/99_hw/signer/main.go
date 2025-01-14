package main

import "fmt"

func main() {
	res := SingleHash("0")

	fmt.Println("res SingleHash", res)

	res = MultiHash(res)
	fmt.Println("res MultiHash", res)

	res = CombineResults("0")
	fmt.Println("CombineResults", res)
}
