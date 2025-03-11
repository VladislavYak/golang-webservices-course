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

	sh_wg := sync.WaitGroup{}

	for val := range in {
		sh_wg.Add(1)
		crc32OutCh := make(chan interface{})
		md5OutCh := make(chan interface{})
		anotherCrc32 := make(chan interface{})
		wg := sync.WaitGroup{}
		wg.Add(3)
		go func(val interface{}) {
			defer sh_wg.Done()

			go crc32Wrapper2(val, crc32OutCh, &wg)
			go md5Wrapper(val, quotaCh, md5OutCh, &wg)

			fromCrc32val := <-crc32OutCh
			md5outVal := <-md5OutCh

			go crc32Wrapper2(md5outVal, anotherCrc32, &wg)

			finalCrc32 := <-anotherCrc32

			res := fromCrc32val.(string) + "~" + finalCrc32.(string)
			wg.Wait()
			out <- res

		}(val)

	}

	sh_wg.Wait()
}

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
	m_h_wg := sync.WaitGroup{}
	for val := range in {
		m_h_wg.Add(1)

		hashingOut := make(chan interface{})

		go func(val interface{}) {
			defer m_h_wg.Done()
			wg := sync.WaitGroup{}
			wg.Add(6)

			for i := 0; i < 6; i++ {
				go func(i int, val interface{}) {

					switch val.(type) {
					case int:
						toF := strconv.Itoa(i) + strconv.Itoa(val.(int))

						ctc32WrapperMultiHash(toF, i, hashingOut, &wg)
					case string:

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
				myArr = append(myArr, h_val.(string))
			}

			result := FormatMultiHashResult(myArr)
			out <- result

		}(val)

	}
	m_h_wg.Wait()

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
