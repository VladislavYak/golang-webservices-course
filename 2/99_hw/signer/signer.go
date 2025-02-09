package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// in, out chan interface{} - this shoud be the input
// dunno what is the output
// job() used as datatype cast at tests
func SingleHash(data string) string {
	return DataSignerCrc32(data) + "~" + DataSignerCrc32(DataSignerMd5(data))
}

// in, out chan interface{} - this shoud be the input
// dunno what is the output
// job() used as datatype cast at tests
func MultiHash(data string) string {
	values := []int{0, 1, 2, 3, 4, 5}
	res := ""
	for _, val := range values {
		res += DataSignerCrc32(strconv.Itoa(val) + data)
	}
	return res
}

// it gets arbitrary number of values as input and should process them concurrently
func CombineResults(data string) string {
	// probably this func should take several values as input
	sh := SingleHash(data)
	res := []string{sh, MultiHash(sh)}

	sort.Strings(res)
	return strings.Join(res, "_")

}

// jobs := []job{
// 	job(func(in, out chan interface{}) {
// 		// out <- 1
// 		out <- 2
// 		out <- 3

// 		close(out)
// 	}),
// 	job(func(in, out chan interface{}) {
// 		for val := range in {
// 			fmt.Println("val", val)

// 			out <- val
// 		}
// 	}),
// 	// job(func(in, out chan interface{}) {
// 	// 	for val := range in {
// 	// 		fmt.Println("quadratic", val, val.(int)*val.(int))
// 	// 		out <- val.(int) * val.(int)
// 	// 	}
// 	// }),
// }

func ExecutePipeline(jobs ...job) {

	in := make(chan interface{})
	// out := make(chan interface{})

	j0 := jobs[0]
	j1 := jobs[1]

	o0 := executor(j0, in)
	o1 := executor(j1, o0)

	for v := range o1 {
		fmt.Println("v", v)
	}

	// for _, j := range jobs {
	// 	out := executor(j, in)

	// 	in = out

	// }

}

func executor(j job, in chan interface{}) chan interface{} {
	out := make(chan interface{})

	go j(in, out)

	return out
}
