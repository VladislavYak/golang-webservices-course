package main

import (
	context "context"
	"fmt"
	"net"
	"sync"

	"gitlab.com/vk-golang/lectures/08_microservices/99_hw/microservice/service"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

// тут вы пишете код
// обращаю ваше внимание - в этом задании запрещены глобальные переменные
// если хочется, то для красоты можно разнести логику по разным файликам

var _ service.BizServer = (*Biz)(nil)
var _ service.AdminServer = (*Admin)(nil)

func StartMyMicroservice(ctx context.Context, addr string, acl string) error {
	loggingChannel := make(chan service.Event)
	subs := map[service.Admin_LoggingServer]bool{}
	mu := &sync.Mutex{}
	// create channel for statistics
	// create channel for logging
	// need to push data to that channels somehow

	broadcaster := NewBroadcaster(loggingChannel, &subs, mu)
	go broadcaster.StartBroadcast()

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("cant listen %s: %w", addr, err)
	}

	server := grpc.NewServer(
		grpc.StreamInterceptor(StreamLoggingInterceptor(loggingChannel)),
		grpc.UnaryInterceptor(UnaryLoggingInterceptor(loggingChannel)),
	)
	service.RegisterBizServer(server, NewBiz())
	service.RegisterAdminServer(server, NewAdmin(&subs, mu))

	fmt.Printf("starting server at %s\n", addr)

	wg := &sync.WaitGroup{}

	wg.Go(func() {
		server.Serve(lis)
	})

	go func() {
		fmt.Println("im before ctx.done")
		<-ctx.Done()
		fmt.Println("context canceled, stopping server")
		server.Stop()
	}()

	fmt.Println("im afer go func")

	return nil
}

type Broadcaster struct {
	loggingChan chan service.Event
	mu          *sync.Mutex
	subs        *map[service.Admin_LoggingServer]bool
}

func NewBroadcaster(loggingChan chan service.Event, subs *map[service.Admin_LoggingServer]bool, mu *sync.Mutex) *Broadcaster {

	return &Broadcaster{loggingChan: loggingChan, mu: mu, subs: subs}
}

func (b *Broadcaster) StartBroadcast() {
	// yakovlev:
	for v := range b.loggingChan {
		fmt.Println("Broadcaster | v.Method", v.Method)
		// fmt.Println("im here read from channel and will be sending to all subs")
		fmt.Println("len b.subs", len(*b.subs))
		for stream := range *b.subs {

			// fmt.Println("stream", stream)
			fmt.Println("Broadcaster | v.Method", v.Method)

			go func() {
				stream.Send(&v)
			}()
			fmt.Println("Broadcaster | After send")
		}

	}
}

type wrappedServerStream struct {
	grpc.ServerStream
	myChan chan service.Event
	info   *grpc.StreamServerInfo
}

func (w *wrappedServerStream) SendMsg(m any) error {
	// fmt.Printf("Intercepted server Send: %v\n", m)

	// w.myChan <- service.Event{Method: w.info.FullMethod}

	// Call the original Send method
	return w.ServerStream.SendMsg(m)
}

func StreamLoggingInterceptor(inChan chan service.Event) grpc.StreamServerInterceptor {
	return func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {

		wrappedStream := &wrappedServerStream{
			ServerStream: ss,
			info:         info,
			myChan:       inChan,
		}
		// tmp hardcode
		inChan <- service.Event{Method: info.FullMethod, Consumer: "logger", Host: "127.0.0.1:58879"}
		// fmt.Println("FullMethod", info.FullMethod)

		return handler(srv, wrappedStream)
	}
}

func UnaryLoggingInterceptor(inChan chan service.Event) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		fmt.Println("im in unary logginc interceptor, info.FullMethod", info.FullMethod)

		host := ""
		if p, ok := peer.FromContext(ctx); ok {
			host = p.Addr.String()
		}

		// Consumer из metadata
		consumer := ""
		if md, ok := metadata.FromIncomingContext(ctx); ok {
			if values := md.Get("consumer"); len(values) > 0 {
				consumer = values[0]
			}
		}
		fmt.Println("Host: host", host)
		inChan <- service.Event{Method: info.FullMethod, Host: host, Consumer: consumer}

		fmt.Println("after send UnaryLoggingInterceptor")

		return handler(ctx, req)
	}
}

type Admin struct {
	service.UnimplementedAdminServer
	mu   *sync.Mutex
	subs *map[service.Admin_LoggingServer]bool
}

func NewAdmin(subs *map[service.Admin_LoggingServer]bool, mu *sync.Mutex) *Admin {
	return &Admin{subs: subs, mu: mu}
}

func (a *Admin) Logging(myNothing *service.Nothing, myStream grpc.ServerStreamingServer[service.Event]) error {
	fmt.Println("before reading from Loggin channel")

	a.mu.Lock()
	(*a.subs)[myStream] = true // ← stream уникален!
	a.mu.Unlock()

	// yakovlev: have no idea what is is, but lets just add that
	<-myStream.Context().Done()

	fmt.Println("im afater <-myStream.Context().Done() Logging")

	a.mu.Lock()
	delete(*a.subs, myStream) // ← stream уникален!
	a.mu.Unlock()

	// for methodName := range a.myChan {
	// 	// fmt.Println("im inside loop logging")
	// 	// fmt.Println("read", methodName)
	// 	myStream.Send(&service.Event{Method: methodName})

	// }

	// for {
	// 	fmt.Println("im here")
	// 	methodName := <-a.myChan
	// 	fmt.Println("after recieve")

	// 	myStream.Send(&service.Event{Method: methodName})

	// 	// myStream.SendMsg(&service.Event{Method: "test"})

	// 	time.Sleep(time.Second * 5)
	// }

	// for {
	// 	fmt.Println("im here")
	// 	// methodName := <-a.myChan
	// 	fmt.Println("after recieve")

	// 	myStream.Send(&service.Event{Method: "test"})

	// 	// myStream.SendMsg(&service.Event{Method: "test"})

	// 	time.Sleep(time.Second * 5)
	// }
	return nil

}

func (a *Admin) Statistics(*service.StatInterval, grpc.ServerStreamingServer[service.Stat]) error {

	return nil
}

type Biz struct {
	service.UnimplementedBizServer
}

func NewBiz() *Biz {
	return &Biz{}
}

func (b *Biz) Check(ctx context.Context, nothing *service.Nothing) (*service.Nothing, error) {
	return &service.Nothing{}, nil
}

func (b *Biz) Add(ctx context.Context, nothing *service.Nothing) (*service.Nothing, error) {
	return &service.Nothing{}, nil
}

func (b *Biz) Test(ctx context.Context, nothing *service.Nothing) (*service.Nothing, error) {
	return &service.Nothing{}, nil
}
