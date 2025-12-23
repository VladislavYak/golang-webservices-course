package main

// import (
// 	"context"
// 	"fmt"
// 	"time"
// )

// func wait(amount int) {
// 	time.Sleep(time.Duration(amount) * 10 * time.Millisecond)
// }

// func main() {
// 	println("usage: go test -v")

// 	ctx, finish := context.WithCancel(context.Background())

// 	StartMyMicroservice(ctx, "127.0.0.1:8082", "")

// 	fmt.Println("im here1")
// 	wait(10)
// 	finish() // при вызове этой функции ваш сервер должен остановиться и освободить порт
// 	wait(1)

// 	ctx, finish = context.WithCancel(context.Background())
// 	StartMyMicroservice(ctx, "127.0.0.1:8082", "")

// 	fmt.Println("before waiting")
// 	wait(1)
// 	finish()
// 	wait(1)

// 	fmt.Println("im here2")

// }
