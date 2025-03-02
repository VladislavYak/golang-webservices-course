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
				out <- fibNum
			}
		}),
		job(SingleHash),
		job(MultiHash),
		// job(CombineResults),
		job(func(in, out chan interface{}) {
			for val := range in {
				fmt.Println("val", val)
			}
		}),
	}

	start := time.Now()

	ExecutePipeline(hashSignJobs...)

	end := time.Since(start)
	fmt.Println("end", end)

	time.Sleep(time.Second * 10)
}
