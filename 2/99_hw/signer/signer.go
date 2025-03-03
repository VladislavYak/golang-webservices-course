package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
)

// 4108050209~502633748
// 2212294583~709660146
func SingleHash(in, out chan interface{}) {

	quotaCh := make(chan struct{}, 1)

	// closeChWg := sync.WaitGroup{}

	for val := range in {
		crc32OutCh := make(chan interface{})
		md5OutCh := make(chan interface{})
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

			res := fromCrc32val.(string) + "~" + finalCrc32.(string)
			fmt.Println("RESULT", res)
			out <- res
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
		out <- res
	case string:
		tmp := val.(string)
		res := DataSignerCrc32(tmp)
		out <- res
	default:
		fmt.Println("val is of unknown type")
	}
}

func ctc32WrapperMultiHash(val interface{}, i int, out chan interface{}, wg *sync.WaitGroup) {
	defer wg.Done()

	switch val.(type) {
	case int:
		tmp := val.(int)
		val2 := strconv.Itoa(int(tmp))
		res := DataSignerCrc32(val2)

		res_2 := strconv.Itoa(i) + "_" + res
		out <- res_2
	case string:
		tmp := val.(string)
		res := DataSignerCrc32(tmp)

		res_2 := strconv.Itoa(i) + "_" + res
		out <- res_2
	default:
		fmt.Println("val is of unknown type")
	}
}

func MultiHash(in, out chan interface{}) {

	for val := range in {

		hashingOut := make(chan interface{})

		go func(val interface{}) {
			wg := sync.WaitGroup{}
			wg.Add(6)

			for i := 0; i < 6; i++ {
				go func(i int, val interface{}) {

					switch val.(type) {
					case int:
						tmp := val.(int)
						val2 := strconv.Itoa(int(tmp))
						res := DataSignerCrc32(val2)
						out <- res

						fmt.Println("MULTIHASH VAL", val)

						toF := strconv.Itoa(i) + strconv.Itoa(val.(int))

						ctc32WrapperMultiHash(toF, i, hashingOut, &wg)
					case string:
						fmt.Println("MULTIHASH VAL", val)

						toF := strconv.Itoa(i) + val.(string)

						ctc32WrapperMultiHash(toF, i, hashingOut, &wg)
					default:
						fmt.Println("val is of unknown type")
					}

				}(i, val)
			}

			go func() {
				wg.Wait()
				close(hashingOut)
			}()

			myArr := []string{}
			for h_val := range hashingOut {
				fmt.Println("MultiHash hashingOut h_val, val", h_val, val)
				myArr = append(myArr, h_val.(string))
			}

			result := FormatMultiHashResult(myArr)
			fmt.Println("multiHash result", result)
			out <- result

		}(val)

	}

}

func FormatMultiHashResult(arr []string) string {
	sort.Strings(arr)

	res := ""

	for _, v := range arr {
		a := strings.SplitAfter(v, "_")

		res += a[len(a)-1]

	}

	return res
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

		// go func() {
		// 	wg.Wait()
		// 	close(out)
		// }()
		in = out

	}

	wg.Wait()
}
