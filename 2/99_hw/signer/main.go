package main

import (
	"time"
)

func main() {
	inputData := []int{0, 1}
	in := make(chan interface{})

	go func(out chan interface{}) {

		for _, val := range inputData {
			out <- val
		}
		close(out)

	}(in)

	shOut := make(chan interface{})
	go SingleHash(in, shOut)

	// for val := range shOut {
	// 	fmt.Println("val", val)
	// }

	time.Sleep(time.Second * 10)

}

// inputData := []int{0, 1, 1, 2, 3, 5, 8}
// inputData := []int{0, 1, 1, 2}

// 	hashSignJobs := []job{
// 		job(func(in, out chan interface{}) {
// 			for _, fibNum := range inputData {
// 				out <- fibNum
// 			}
// 		}),
// 		job(SingleHash),
// 		// job(MultiHash),
// 		// job(CombineResults),
// 		job(func(in, out chan interface{}) {
// 			for val := range in {
// 				fmt.Println("val", val)
// 			}
// 		}),
// 	}

// 	start := time.Now()

// 	ExecutePipeline(hashSignJobs...)

// 	end := time.Since(start)
// 	fmt.Println("end", end)

// 	time.Sleep(time.Second * 10)
// }

// func ParallelWorker(in, out chan interface{}, baseFunc func(data string) string) {

// 	for val := range in {
// 		go func(val interface{}, out chan interface{}) {
// 			defer close(out)
// 			fmt.Println("val", val)

// 			tmp := val.(uint32)
// 			val2 := strconv.Itoa(int(tmp))

// 			res := baseFunc(val2)

// 			fmt.Println("res for val", res, val)
// 			out <- res

// 		}(val, out)
// 	}
// }

// OneInATimeWorker executes baseFunc sequentially
// func OneInATimeWorker(in, out chan interface{}, baseFunc func(data string) string) {
// 	//
// }
