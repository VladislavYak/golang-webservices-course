package main

import (
	"context"
	"log"

	"gitlab.com/vk-golang/lectures/08_microservices/99_hw/microservice/service"
	"google.golang.org/grpc"
)

func main() {
	grcpConn, err := grpc.Dial(
		"127.0.0.1:8081",
		// grpc.WithUnaryInterceptor(timingInterceptor),
		// grpc.WithPerRPCCredentials(&tokenAuth{"100500"}),
		grpc.WithInsecure(),
	)
	if err != nil {
		log.Fatalf("cant connect to grpc")
	}
	defer grcpConn.Close()

	bizClient := service.NewBizClient(grcpConn)

	_, err = bizClient.Add(context.Background(), &service.Nothing{})
	if err != nil {
		panic(err)
	}

}
