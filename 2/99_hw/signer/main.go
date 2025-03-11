package main

import (
	"fmt"
	"time"
)

// func main() {
// 	res := FormatMultiHashResult([]string{"3_3407918797", "0_2956866606", "4_2730963093", "1_803518384", "2_1425683795"})

// 	fmt.Println("res", res)
// }

func main() {

	inputData := []int{0, 1}

	hashSignJobs := []job{
		job(func(in, out chan interface{}) {
			for _, fibNum := range inputData {
				// time.Sleep(time.Second * 2)
				fmt.Println("fibNum", fibNum)
				out <- fibNum
				fmt.Println("gen after insertion")
			}

			fmt.Println("below loop generator")
		}),
		job(SingleHash),
		job(MultiHash),
		job(CombineResults),
		job(func(in, out chan interface{}) {
			fmt.Println("fin read")
			for val := range in {
				fmt.Println("val final 2", val)
			}
			fmt.Println("after fin")
		}),
	}

	start := time.Now()

	ExecutePipeline2(hashSignJobs...)

	end := time.Since(start)
	fmt.Println("end", end)

	time.Sleep(time.Second * 20)
}
