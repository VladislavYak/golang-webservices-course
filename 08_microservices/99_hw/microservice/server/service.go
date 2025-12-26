package main

import (
	context "context"
	"encoding/json"
	"fmt"
	"net"
	"strings"
	"sync"

	"gitlab.com/vk-golang/lectures/08_microservices/99_hw/microservice/service"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

// тут вы пишете код
// обращаю ваше внимание - в этом задании запрещены глобальные переменные
// если хочется, то для красоты можно разнести логику по разным файликам

var _ service.BizServer = (*Biz)(nil)
var _ service.AdminServer = (*Admin)(nil)

func StartMyMicroservice(ctx context.Context, addr string, acl string) error {
	var aclProcessed map[string][]string
	if err := json.Unmarshal([]byte(acl), &aclProcessed); err != nil {
		return fmt.Errorf("failed to parse ACL: %w", err)
	}

	loggingChannel := make(chan service.Event)
	subs := map[service.Admin_LoggingServer]bool{}
	addSubCh := make(chan struct{})
	mu := &sync.Mutex{}
	// create channel for statistics

	broadcaster := NewBroadcaster(loggingChannel, addSubCh, &subs, mu)
	go broadcaster.StartBroadcast()

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("cant listen %s: %w", addr, err)
	}

	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			UnaryACLInterceptor(aclProcessed),
			UnaryLoggingInterceptor(loggingChannel, addSubCh),
		),
		grpc.ChainStreamInterceptor(
			StreamACLInterceptor(aclProcessed),
			StreamLoggingInterceptor(loggingChannel, addSubCh),
		),
	)
	service.RegisterBizServer(server, NewBiz())
	service.RegisterAdminServer(server, NewAdmin(&subs, mu))

	fmt.Printf("starting server at %s\n", addr)

	wg := &sync.WaitGroup{}

	wg.Go(func() {
		server.Serve(lis)
	})

	go func() {
		<-ctx.Done()
		close(loggingChannel)
		close(addSubCh)
		server.GracefulStop()
	}()

	return nil
}

type Broadcaster struct {
	loggingChan chan service.Event
	addSubCh    chan struct{}
	mu          *sync.Mutex
	subs        *map[service.Admin_LoggingServer]bool
}

func NewBroadcaster(loggingChan chan service.Event, addSubCh chan struct{}, subs *map[service.Admin_LoggingServer]bool, mu *sync.Mutex) *Broadcaster {

	return &Broadcaster{loggingChan: loggingChan, addSubCh: addSubCh, mu: mu, subs: subs}
}

func (b *Broadcaster) StartBroadcast() {
	for v := range b.loggingChan {
		b.mu.Lock()
		for stream := range *b.subs {

			// yakovlev: add error handling
			stream.Send(&v)
		}

		b.mu.Unlock()
		b.addSubCh <- struct{}{}

	}
}

func getConsumerFromContext(ctx context.Context) string {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if values := md.Get("consumer"); len(values) > 0 {
			return values[0]
		}
	}
	return ""
}

// isAllowed checks if consumer is allowed to call fullMethod based on ACL rules (supports wildcard "/*").
func isAllowed(acl map[string][]string, consumer, fullMethod string) bool {
	allowedMethods, ok := acl[consumer]
	if !ok {
		return false
	}

	for _, rule := range allowedMethods {
		// Exact match
		if rule == fullMethod {
			return true
		}

		// Wildcard support: /main.Biz/*
		if strings.HasSuffix(rule, "/*") {
			prefix := rule[:len(rule)-1] // e.g., /main.Biz/
			if strings.HasPrefix(fullMethod, prefix) {
				return true
			}
		}
	}
	return false
}

func StreamLoggingInterceptor(inChan chan service.Event, addSubCh chan struct{}) grpc.StreamServerInterceptor {
	return func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		ctx := ss.Context()

		// 1. HOST из peer
		host := ""
		if p, ok := peer.FromContext(ctx); ok {
			host = p.Addr.String() // "127.0.0.1:59292"
		}

		// 2. CONSUMER из metadata
		consumer := ""
		if md, ok := metadata.FromIncomingContext(ctx); ok {
			if values := md.Get("consumer"); len(values) > 0 {
				consumer = values[0] // "logger", "biz_user"
			}
		}
		inChan <- service.Event{Consumer: consumer,
			Method: info.FullMethod,
			Host:   host}

		<-addSubCh // Ждём завершения broadcasting'а перед вызовом handler

		return handler(srv, ss)
	}
}

// StreamACLInterceptor возвращает stream интерсептор, проверяющий доступ по ACL
func StreamACLInterceptor(acl map[string][]string) grpc.StreamServerInterceptor {
	return func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		consumer := getConsumerFromContext(ss.Context())
		if consumer == "" || !isAllowed(acl, consumer, info.FullMethod) {
			return status.Errorf(codes.Unauthenticated, "access denied for consumer '%s' to method '%s'", consumer, info.FullMethod)
		}
		return handler(srv, ss)
	}
}

func UnaryLoggingInterceptor(inChan chan service.Event, addSubCh chan struct{}) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {

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
		inChan <- service.Event{Method: info.FullMethod, Host: host, Consumer: consumer}

		<-addSubCh // Ждём завершения broadcasting'а перед вызовом handler

		return handler(ctx, req)
	}
}

// UnaryACLInterceptor возвращает unary интерсептор, проверяющий доступ по ACL
func UnaryACLInterceptor(acl map[string][]string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		consumer := getConsumerFromContext(ctx)
		if consumer == "" || !isAllowed(acl, consumer, info.FullMethod) {
			return nil, status.Errorf(codes.Unauthenticated, "access denied for consumer '%s' to method '%s'", consumer, info.FullMethod)
		}
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
	a.mu.Lock()
	(*a.subs)[myStream] = true // ← stream уникален!
	a.mu.Unlock()

	// yakovlev: have no idea what is is, but lets just add that
	<-myStream.Context().Done()

	a.mu.Lock()
	delete(*a.subs, myStream) // ← stream уникален!
	a.mu.Unlock()
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
