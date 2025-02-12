package main

import (
	"fmt"
	"sync/atomic"
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

	var recieved uint32
	jobs := []job{
		job(func(in, out chan interface{}) {
			out <- uint32(1)
			out <- uint32(3)
			out <- uint32(4)
		}),
		job(func(in, out chan interface{}) {
			for val := range in {
				out <- val.(uint32) * 3
				time.Sleep(time.Millisecond * 100)
			}
		}),
		job(func(in, out chan interface{}) {
			for val := range in {
				fmt.Println("collected", val)
				atomic.AddUint32(&recieved, val.(uint32))

				fmt.Println("IN GOROUTINE recieved BECOME", recieved)
			}
		}),
	}

	ExecutePipeline(jobs...)

	fmt.Println("result recieved", recieved)

	time.Sleep(time.Second * 5)
}
