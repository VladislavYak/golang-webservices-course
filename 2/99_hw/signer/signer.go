package main

import (
	"strconv"
	"strings"
)

// in, out chan interface{} - this shoud be the input
// dunno what is the output
// job() used as datatype cast at tests
func SingleHash(data string) string {
	return DataSignerCrc32(data) + "~" + DataSignerCrc32(DataSignerMd5(data))
}

// in, out chan interface{} - this shoud be the input
// dunno what is the output
// job() used as datatype cast at tests
func MultiHash(data string) string {
	values := []int{0, 1, 2, 3, 4, 5}
	res := ""
	for _, val := range values {
		res += DataSignerCrc32(strconv.Itoa(val) + data)
	}
	return res
}

func CombineResults(data string) string {
	// probably this func should take several values as input
	sh := SingleHash(data)
	res := []string{sh, MultiHash(sh)}
	// sorting should be done somehow

	return strings.Join(res, "_")

	// concat above
}

func ExecutePipeline() {
	// who the fuck are you

}
