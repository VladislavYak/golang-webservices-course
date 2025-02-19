package main

import (
	"fmt"
	"strconv"
	"sync"
)

// in, out chan interface{} - this shoud be the input
// dunno what is the output
// job() used as datatype cast at tests
//
//	func SingleHash(data string) string {
//		return DataSignerCrc32(data) + "~" + DataSignerCrc32(DataSignerMd5(data))
//	}
func SingleHash(in, out chan interface{}) {

	// parallel computations
	// for val := range in {
	// 	go func(val interface{}, out chan interface{}) {
	// 		// defer close(out)

	// 		tmp := val.(int)
	// 		val2 := strconv.Itoa(int(tmp))

	// 		res := DataSignerCrc32(val2)

	// 		fmt.Println("res for val", res, val)
	// 		out <- res

	// 	}(val, out)
	// }

	// sequential execution
	quotaCh := make(chan struct{}, 1)
	go func() {
		quotaCh <- struct{}{}

		for val := range in {
			tmp := val.(int)
			val2 := strconv.Itoa(int(tmp))

			res := DataSignerMd5(val2)

			fmt.Println("res for val", res, val)

			out <- res
		}

		<-quotaCh
	}()
}

// in, out chan interface{} - this shoud be the input
// dunno what is the output
// job() used as datatype cast at tests
//
//	func MultiHash(data string) string {
//		values := []int{0, 1, 2, 3, 4, 5}
//		res := ""
//		for _, val := range values {
//			res += DataSignerCrc32(strconv.Itoa(val) + data)
//		}
//		return res
//	}
func MultiHash(in, out chan interface{}) {

}

// it gets arbitrary number of values as input and should process them concurrently
// func CombineResults(data string) string {
// 	// probably this func should take several values as input
// 	sh := SingleHash(data)
// 	res := []string{sh, MultiHash(sh)}

// 	sort.Strings(res)
// 	return strings.Join(res, "_")

// }
func CombineResults(in, out chan interface{}) {

}

func ExecutePipeline(jobs ...job) {

	in := make(chan interface{})
	wg := sync.WaitGroup{}
	wg.Add(len(jobs))
	for _, j := range jobs {
		out := make(chan interface{})

		go func(j job, in, out chan interface{}) {
			defer wg.Done()
			j(in, out)
			// close(out)
		}(j, in, out)
		in = out

	}

	wg.Wait()
}
