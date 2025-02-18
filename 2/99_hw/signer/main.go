package main

import (
	"fmt"
	"strconv"
	"time"
)

func main() {
	in := make(chan interface{})
	out := make(chan interface{})

	go func(out chan interface{}) {
		defer close(out)
		out <- uint32(5)
		out <- uint32(1)
		time.Sleep(time.Second * 3)
		out <- uint32(7)

	}(out)

	in = out

	out2 := make(chan interface{})
	ParallelWorker(in, out2, DataSignerCrc32)

	time.Sleep(time.Second * 10)
}

func ParallelWorker(in, out chan interface{}, baseFunc func(data string) string) {

	for val := range in {
		go func(val interface{}, out chan interface{}) {
			defer close(out)
			fmt.Println("val", val)

			tmp := val.(uint32)
			val2 := strconv.Itoa(int(tmp))

			res := baseFunc(val2)

			fmt.Println("res for val", res, val)
			out <- res

		}(val, out)
	}
}

// OneInATimeWorker executes baseFunc sequentially
func OneInATimeWorker(in, out chan interface{}, baseFunc func(data string) string) {
	//
}
