package main

// in, out chan interface{} - this shoud be the input
// dunno what is the output
// job() used as datatype cast at tests
//
//	func SingleHash(data string) string {
//		return DataSignerCrc32(data) + "~" + DataSignerCrc32(DataSignerMd5(data))
//	}
func SingleHash(in, out chan interface{}) {

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
	for _, j := range jobs {
		in = executor(j, in)

	}

	close(in)
}

func executor(j job, in chan interface{}) chan interface{} {
	out := make(chan interface{})

	go j(in, out)

	return out
}
