package main

import (
	"fmt"
	"time"
)

// test hashing - produces correct result (no channels, no pipeline)
// func main() {
// 	res := SingleHash("0")

// 	fmt.Println("res SingleHash", res)

// 	res = MultiHash(res)
// 	fmt.Println("res MultiHash", res)

// 	res = CombineResults("0")
// 	fmt.Println("CombineResults", res)
// }

// test execute pipeline
func main() {
	jobs := []job{
		job(func(in, out chan interface{}) {
			out <- 1
			out <- 2
			out <- 3

		}),
		job(func(in, out chan interface{}) {
			for val := range in {
				fmt.Println("val", val)

				out <- val
			}
		}),
		job(func(in, out chan interface{}) {
			for val := range in {
				fmt.Println("quadratic", val, val.(int)*val.(int))
				// out <- val.(int) * val.(int)
			}
		}),
	}

	ExecutePipeline(jobs...)

	time.Sleep(time.Second * 5)
}

// func gen(out chan int, values ...int) {
// 	defer close(out)
// 	for _, v := range values {
// 		out <- v
// 	}
// }

// func main() {
// 	out := make(chan int)
// 	go gen(out, 1, 2, 3, 4)

// 	for v := range out {
// 		fmt.Println("v", v)
// 	}
// 	time.Sleep(time.Second * 5)
// }
