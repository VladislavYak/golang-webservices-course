package main

import (
	"fmt"
	"time"
)

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
