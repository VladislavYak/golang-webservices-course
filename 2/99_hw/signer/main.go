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

	inputData := []int{0, 1, 1, 2, 3, 5, 8, 10, 15, 20, 22, 30, 35, 40, 45, 50, 55}

	hashSignJobs := []job{
		job(func(in, out chan interface{}) {
			for _, fibNum := range inputData {
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

	ExecutePipeline(hashSignJobs...)

	end := time.Since(start)
	fmt.Println("end", end)
}
