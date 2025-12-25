package main

// import (
// 	"context"
// 	"fmt"
// 	"runtime"
// 	"time"
// )

// func wait(amount int) {
// 	time.Sleep(time.Duration(amount) * 10 * time.Millisecond)
// }

// func dumpGoroutines(prefix string) {
// 	fmt.Printf("\n=== %s | Goroutines: %d ===\n", prefix, runtime.NumGoroutine())
// 	buf := make([]byte, 1<<20)
// 	n := runtime.Stack(buf, true)
// 	fmt.Print(string(buf[:n]))
// }

// func main() {
// 	println("usage: go test -v")

// 	ctx, finish := context.WithCancel(context.Background())

// 	StartMyMicroservice(ctx, "127.0.0.1:8082", "")

// 	time.Sleep(time.Second * 2)

// 	finish()

// 	fmt.Println("after finish")
// 	fmt.Println("runtime.NumGoroutine()", runtime.NumGoroutine())
// 	dumpGoroutines("main_")
// 	time.Sleep(time.Second * 10)

// 	// wait(1000000)

// 	// fmt.Println("im here1")
// 	// wait(10)
// 	// finish() // при вызове этой функции ваш сервер должен остановиться и освободить порт
// 	// wait(1)

// 	// ctx, finish = context.WithCancel(context.Background())
// 	// StartMyMicroservice(ctx, "127.0.0.1:8082", "")

// 	// fmt.Println("before waiting")
// 	// wait(1)
// 	// finish()
// 	// wait(1)

// 	// fmt.Println("im here2")

// }
