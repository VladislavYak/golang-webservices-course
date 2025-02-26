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

// 4108050209~502633748
// 2212294583~709660146
func SingleHash(in, out chan interface{}) {

	quotaCh := make(chan struct{}, 1)

	for val := range in {
		crc32OutCh := make(chan interface{})
		md5OutCh := make(chan interface{})
		fmt.Println("VAL", val)
		go func(val interface{}) {
			wg := sync.WaitGroup{}

			wg.Add(3)
			go crc32Wrapper2(val, crc32OutCh, &wg)
			go md5Wrapper(val, quotaCh, md5OutCh, &wg)

			fromCrc32val := <-crc32OutCh
			md5outVal := <-md5OutCh

			// fmt.Println("fromCrc32val, anotherCrc32", fromCrc32val, md5outVal)

			anotherCrc32 := make(chan interface{})

			go crc32Wrapper2(md5outVal, anotherCrc32, &wg)

			finalCrc32 := <-anotherCrc32

			fmt.Println("fromCrc32val, anotherCrc32, finalCrc32", fromCrc32val, md5outVal, finalCrc32)

			fmt.Println("RESULT", fromCrc32val.(string)+"~"+finalCrc32.(string))
			wg.Wait()

		}(val)
		// crc32Wrapper2(anotherCrc32, anotherCrc32)

		// finalCrc32val := <-anotherCrc32

		// fmt.Println(finalCrc32val.(string) + "~" + fromCrc32val)

	}

}

// md5 wrapper
func md5Wrapper(val interface{}, quota chan struct{}, out chan interface{}, wg *sync.WaitGroup) {
	defer wg.Done()
	quota <- struct{}{}

	tmp := val.(int)
	val2 := strconv.Itoa(int(tmp))

	res := DataSignerMd5(val2)

	fmt.Println("res val md5Wrapper", res, val)

	out <- res

	<-quota
}

func crc32Wrapper2(val interface{}, out chan interface{}, wg *sync.WaitGroup) {
	defer wg.Done()
	switch val.(type) {
	case int:
		tmp := val.(int)
		val2 := strconv.Itoa(int(tmp))
		res := DataSignerCrc32(val2)
		fmt.Println("res val crc32Wrapper2", res, val)
		out <- res
	case string:
		tmp := val.(string)
		res := DataSignerCrc32(tmp)
		fmt.Println("res val crc32Wrapper2", res, val)
		out <- res
	default:
		fmt.Println("val is of unknown type")
	}
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
