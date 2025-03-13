package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
)

func SingleHash(in chan interface{}, out chan interface{}) {

	quotaCh := make(chan struct{}, 1)

	mainWgSH := sync.WaitGroup{}

	for val := range in {
		mainWgSH.Add(1)
		crc32OutCh := make(chan interface{})
		md5OutCh := make(chan interface{})
		anotherCrc32OutCh := make(chan interface{})
		go func(val interface{}) {
			defer mainWgSH.Done()

			go crc32Wrapper(val, crc32OutCh)

			go md5Wrapper(val, quotaCh, md5OutCh)

			md5outVal := <-md5OutCh

			go crc32Wrapper(md5outVal, anotherCrc32OutCh)

			finalCrc32 := <-anotherCrc32OutCh
			fromCrc32val := <-crc32OutCh

			out <- fromCrc32val.(string) + "~" + finalCrc32.(string)
		}(val)

	}

	mainWgSH.Wait()
}

func md5Wrapper(val interface{}, quota chan struct{}, out chan interface{}) {
	quota <- struct{}{}

	switch val.(type) {
	case string:
		out <- DataSignerMd5(val.(string))
	case int:
		tmp := strconv.Itoa(val.(int))
		out <- DataSignerMd5(tmp)
	default:
		// yakovlev neet to raise error somehow
	}

	<-quota
}

func crc32Wrapper(val interface{}, out chan interface{}) {
	switch val.(type) {
	case int:
		val2 := strconv.Itoa(val.(int))
		out <- DataSignerCrc32(val2)
	case string:
		out <- DataSignerCrc32(val.(string))
	default:
		// yakovlev: need to raise error somehow
		fmt.Println("val is of unknown type")
	}
}

func ctc32WrapperMultiHash(val interface{}, i int, out chan interface{}, wg *sync.WaitGroup) {
	defer wg.Done()

	switch val.(type) {
	case int:
		tmp := val.(int)
		val2 := strconv.Itoa(tmp)
		res := DataSignerCrc32(val2)

		res_2 := strconv.Itoa(i) + "_" + res
		out <- res_2
	case string:
		tmp := val.(string)
		res := DataSignerCrc32(tmp)

		res_2 := strconv.Itoa(i) + "_" + res
		out <- res_2
	default:
		// yakovlev: need to raise error somehow
		fmt.Println("val is of unknown type")
	}
}

func MultiHash(in, out chan interface{}) {
	mainWgMh := sync.WaitGroup{}
	for val := range in {
		mainWgMh.Add(1)

		hashingOutCh := make(chan interface{})

		go func(val interface{}) {
			defer mainWgMh.Done()
			wg := sync.WaitGroup{}
			wg.Add(6)

			for i := 0; i < 6; i++ {
				go func(i int, val interface{}) {

					switch val.(type) {
					case int:
						convertedVal := strconv.Itoa(i) + strconv.Itoa(val.(int))

						ctc32WrapperMultiHash(convertedVal, i, hashingOutCh, &wg)
					case string:
						convertedVal := strconv.Itoa(i) + val.(string)

						ctc32WrapperMultiHash(convertedVal, i, hashingOutCh, &wg)
					default:
						// yakovlev: need to raise error somehow
						fmt.Println("val is of unknown type")
					}

				}(i, val)
			}

			go func() {
				wg.Wait()
				close(hashingOutCh)
			}()

			myArr := []string{}
			for h_val := range hashingOutCh {
				myArr = append(myArr, h_val.(string))
			}

			out <- FormatMultiHashResult(myArr)

		}(val)

	}
	mainWgMh.Wait()

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

func CombineResults(in, out chan interface{}) {
	res := []string{}
	for val := range in {
		res = append(res, val.(string))
	}
	sort.Strings(res)

	out <- strings.Join(res, "_")

}

func ExecutePipeline(jobs ...job) {

	in := make(chan interface{})

	mainWg := sync.WaitGroup{}
	mainWg.Add(len(jobs))

	for _, myJob := range jobs {
		out := make(chan interface{})
		wg := sync.WaitGroup{}
		wg.Add(1)

		go func(in, out chan interface{}, j job) {
			defer wg.Done()
			defer mainWg.Done()

			j(in, out)

		}(in, out, myJob)

		in = out

		go func(out chan interface{}, j job) {
			wg.Wait()
			close(out)
		}(out, myJob)

	}

	mainWg.Wait()
}
