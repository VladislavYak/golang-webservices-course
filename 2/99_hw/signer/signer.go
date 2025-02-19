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
	// ебануть разные каналы для парралельно и последовательной обработки?
	// создать канал в который должна записать последовательная который после будет читать паралельной обработкой?

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
	// quotaCh := make(chan struct{}, 1)
	// go func() {
	// 	quotaCh <- struct{}{}

	// 	for val := range in {
	// 		tmp := val.(int)
	// 		val2 := strconv.Itoa(int(tmp))

	// 		res := DataSignerMd5(val2)

	// 		fmt.Println("res for val", res, val)

	// 		out <- res
	// 	}

	// 	<-quotaCh
	// }()

	// attempt to do several computiations in one place
	quotaCh := make(chan struct{}, 1)
	crc32OutCh := make(chan interface{})
	md5OutCh := make(chan interface{})
	for val := range in {
		// go func(val interface{}, out chan interface{}) {
		// 	tmp := val.(int)
		// 	val2 := strconv.Itoa(int(tmp))

		// 	res := DataSignerCrc32(val2)
		// 	fmt.Println("res", res)

		// }(val, out)
		go md5Wrapper(val, quotaCh, md5OutCh)
		go crc32Wrapper2(val, crc32OutCh)

		// go md5Wrapper(val)
		// close(out)
		// go md5Wrapper(val, quotaCh, md5OutCh)

	}

	// for val := range md5OutCh {
	// 	fmt.Println("READ val | md5OutCh", val)
	// }

	// for val := range crc32OutCh {
	// 	fmt.Println("READ val | crc32OutCh", val)
	// }

	for {
		v1 := <-md5OutCh

		out := make(chan interface{})
		go crc32Wrapper2(v1, out)
		res := <-out

		v2 := <-crc32OutCh

		// fmt.Println("READ val | md5OutCh", v1)
		// fmt.Println("READ val | crc32OutCh", v2)

		fmt.Println("V1 & V2 & res: ", v1, v2, res)

		fmt.Println("RESULT", res.(string)+"~"+v2.(string))
		// оно перепуталось
		// должно быть
		// 0 SingleHash result 4108050209~502633748
		// 2212294583~709660146

		// RESULT 709660146~2212294583
		// res val crc32Wrapper2 502633748 cfcd208495d565ef66e7dff9f98764da
		// V1 & V2 & res:  cfcd208495d565ef66e7dff9f98764da 4108050209 502633748
		// RESULT 502633748~4108050209
	}

	// for val := range crc32OutCh {
	// 	fmt.Println("val", val)
	// }

	// for val := range md5OutCh {
	// 	fmt.Println("v md5OutCh", val)
	// }

	// 	for val := range in {
	// 		md5Out := make(chan string)
	// 		wg := sync.WaitGroup{}
	// 		wg.Add(3)
	// 		// wg.Add(2)
	// 		// wg.Add(1)
	// 		go md5Wrapper(val, quotaCh, md5Out, &wg)

	// 		crc32Out := make(chan interface{})
	// 		go crc32Wrapper(val, crc32Out, &wg)

	// 		tmp := <-md5Out
	// 		tmp1 := <-crc32Out

	// 		go crc32Wrapper(tmp, crc32Out, &wg)

	// 		tmp2 := <-crc32Out

	// 		wg.Wait()
	// 		fmt.Println("tmp", tmp)
	// 		fmt.Println("tmp1", tmp1)
	// 		fmt.Println("tmp2", tmp2)

	// 	}
}

// md5 wrapper
func md5Wrapper(val interface{}, quota chan struct{}, out chan interface{},

// wg *sync.WaitGroup
) {
	// defer wg.Done()
	quota <- struct{}{}

	tmp := val.(int)
	val2 := strconv.Itoa(int(tmp))

	res := DataSignerMd5(val2)

	fmt.Println("res val md5Wrapper", res, val)

	out <- res

	<-quota
}

func crc32Wrapper2(val interface{}, out chan interface{}) {
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
