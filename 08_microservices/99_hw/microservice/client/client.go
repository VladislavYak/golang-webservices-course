package main

import (
	"context"
	"fmt"
	"log"

	"gitlab.com/vk-golang/lectures/08_microservices/99_hw/microservice/service"
	"google.golang.org/grpc"
)

func main() {
	grcpConn, err := grpc.Dial(
		"127.0.0.1:8082",
		// grpc.WithUnaryInterceptor(timingInterceptor),
		// grpc.WithPerRPCCredentials(&tokenAuth{"100500"}),
		grpc.WithInsecure(),
	)
	if err != nil {
		log.Fatalf("cant connect to grpc")
	}
	defer grcpConn.Close()

	bizClient := service.NewBizClient(grcpConn)
	adminClient := service.NewAdminClient(grcpConn)

	_, err = bizClient.Add(context.Background(), &service.Nothing{})
	mySream, err := adminClient.Logging(context.Background(), &service.Nothing{})

	for {
		evnt, err := mySream.Recv()
		if err != nil {
			panic(err)
		}
		fmt.Println("evnt", evnt.Method)
	}
	if err != nil {
		panic(err)
	}

}
