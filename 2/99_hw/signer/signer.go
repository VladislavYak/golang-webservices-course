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
func SingleHash(in chan interface{}, out chan interface{}) {

	quotaCh := make(chan struct{}, 1)

	// closeChWg := sync.WaitGroup{}

	for val := range in {
		fmt.Println("SingleHash | SingleHash, val", val)
		crc32OutCh := make(chan interface{})
		md5OutCh := make(chan interface{})
		anotherCrc32 := make(chan interface{})
		wg := sync.WaitGroup{}
		wg.Add(3)
		go func(val interface{}) {

			wg.Add(3)
			go crc32Wrapper2(val, crc32OutCh, &wg)
			go md5Wrapper(val, quotaCh, md5OutCh, &wg)

			fromCrc32val := <-crc32OutCh
			md5outVal := <-md5OutCh

			fmt.Println("SingleHash | fromCrc32val, anotherCrc32", fromCrc32val, md5outVal)

			// anotherCrc32 := make(chan interface{})

			go crc32Wrapper2(md5outVal, anotherCrc32, &wg)

			finalCrc32 := <-anotherCrc32

			fmt.Println("SingleHash| finalCrc32", finalCrc32)

			res := fromCrc32val.(string) + "~" + finalCrc32.(string)
			fmt.Println("SingleHash | RESULT", res)
			out <- res
			wg.Wait()

		}(val)

		// crc32Wrapper2(anotherCrc32, anotherCrc32)

		// finalCrc32val := <-anotherCrc32

		// fmt.Println(finalCrc32val.(string) + "~" + fromCrc32val)

	}

	fmt.Println("i can leave single Hash loop!")

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
	fmt.Println("MULTIHASH START")
	for val := range in {

		hashingOut := make(chan interface{})

		go func(val interface{}) {
			wg := sync.WaitGroup{}
			wg.Add(6)

			for i := 0; i < 6; i++ {
				go func(i int, val interface{}) {

					switch val.(type) {
					case int:
						// tmp := val.(int)
						// val2 := strconv.Itoa(int(tmp))
						// res := DataSignerCrc32(val2)
						// out <- res

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
				fmt.Println("closing multihash")
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
	fmt.Println("inside combine results")
	res := []string{}
	for val := range in {
		fmt.Println("CombineResults val", val)
		res = append(res, val.(string))
	}
	fmt.Println("combine results after looping")
	sort.Strings(res)

	finalString := strings.Join(res, "_")

	fmt.Println("COMBINE RESULTS finalString", finalString)
	out <- finalString

	// close(in)
	// close(out)

}

func ExecutePipeline3(jobs ...job) {
	in := make(chan interface{})
	out := make(chan interface{})
	myWg := sync.WaitGroup{}
	myWg.Add(1)

	go func(in, out chan interface{}, wg *sync.WaitGroup, j job) {
		defer wg.Done()
		j(in, out)
	}(in, out, &myWg, jobs[0])

	go func(wg *sync.WaitGroup, myOut chan interface{}) {
		wg.Wait()
		close(myOut)
	}(&myWg, out)

	myWg2 := sync.WaitGroup{}
	myWg2.Add(1)
	in = make(chan interface{})

	go func(in, out chan interface{}, wg *sync.WaitGroup, j job) {
		defer wg.Done()
		j(in, out)
	}(out, in, &myWg2, jobs[1])

	// go func(in, out chan interface{}, wg *sync.WaitGroup, j job) {
	// 	defer wg.Done()
	// 	jobs[0](in, out)
	// }(in, out, &myWg, myJob)

	// go func(in, out chan interface{}) {
	// 	jobs[0](in, out)
	// }(out, in)
}

// for _, myJob := range jobs {
// 	myWg := sync.WaitGroup{}
// 	myWg.Add(1)

// 	go func(in, out chan interface{}, wg *sync.WaitGroup, j job) {
// 		defer wg.Done()
// 		j(in, out)
// 	}(in, out, &myWg, myJob)

// 	go func(wg *sync.WaitGroup, myOut chan interface{}) {
// 		wg.Wait()
// 		close(myOut)
// 	}(&myWg, out)

// 	in = out
// 	out = make(chan interface{})
// }

func ExecutePipeline2(jobs ...job) {

	in := make(chan interface{})

	for _, myJob := range jobs {
		out := make(chan interface{})
		myWg := sync.WaitGroup{}
		myWg.Add(1)

		go func(in, out chan interface{}, wg *sync.WaitGroup, j job) {
			defer wg.Done()
			j(in, out)
		}(in, out, &myWg, myJob)

		go func(wg *sync.WaitGroup, myOut chan interface{}) {
			wg.Wait()
			close(myOut)
		}(&myWg, out)

		in = out
	}
}

func ExecutePipeline(jobs ...job) {

	in := make(chan interface{})
	for _, j := range jobs {
		out := make(chan interface{})

		innerWg := sync.WaitGroup{}
		innerWg.Add(1)

		go func(j job, in, out chan interface{}, wg *sync.WaitGroup) {
			fmt.Println("ExecutePipeline | EXECUTION J ExecutePipeline goroutine here for j", j)
			defer wg.Done()
			// defer close(out)
			// defer wg2.Done()
			j(in, out)

			fmt.Println("ExecutePipeline | after j", j)
			// in = out
			// close(out)
		}(j, in, out, &innerWg)

		go func(wg *sync.WaitGroup, j job) {
			fmt.Println("ExecutePipeline | CLOSE ExecutePipeline goroutine here for j", j)
			wg.Wait()
			fmt.Println("ExecutePipeline | ive closed out", j)
			close(out)
			// in = out
		}(&innerWg, j)

		// close(out)

		fmt.Println("ExecutePipeline | in = out for j", j)
		fmt.Println("ExecutePipeline | ----")
		in = out

	}

	// go func() {
	// 	wg2.Wait()
	// }()
}
