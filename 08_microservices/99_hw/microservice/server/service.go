package main

import (
	context "context"
	"fmt"
	"net"
	"sync"

	"gitlab.com/vk-golang/lectures/08_microservices/99_hw/microservice/service"
	grpc "google.golang.org/grpc"
)

// тут вы пишете код
// обращаю ваше внимание - в этом задании запрещены глобальные переменные
// если хочется, то для красоты можно разнести логику по разным файликам

var _ service.BizServer = (*Biz)(nil)
var _ service.AdminServer = (*Admin)(nil)

func StartMyMicroservice(ctx context.Context, addr string, acl string) error {
	// create channel for statistics
	// create channel for logging
	// need to push data to that channels somehow

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("cant listen %s: %w", addr, err)
	}

	server := grpc.NewServer()
	service.RegisterBizServer(server, NewBiz())
	service.RegisterAdminServer(server, NewAdmin())

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

type Admin struct {
	service.UnimplementedAdminServer
}

func NewAdmin() *Admin {
	return &Admin{}
}

func (a *Admin) Logging(*service.Nothing, grpc.ServerStreamingServer[service.Event]) error {
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
