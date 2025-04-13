package main

import (
	"context"
	pb "github.com/ether-echo/protos/messageProcessor"
	"github.com/ether-echo/telegram-api-service/adapter/grpc_server"
	"github.com/ether-echo/telegram-api-service/adapter/producer"
	"github.com/ether-echo/telegram-api-service/internal/handler"
	"github.com/ether-echo/telegram-api-service/internal/repository"
	"github.com/ether-echo/telegram-api-service/internal/service"
	"github.com/ether-echo/telegram-api-service/pkg/config"
	"github.com/ether-echo/telegram-api-service/pkg/logger"
	"google.golang.org/grpc"
	"net"
)

var (
	log = logger.Logger().Named("main").Sugar()
)

func main() {

	conf, err := config.ReadConfig()
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	prod := producer.NewKafkaProducer([]string{"kafka:9092"})

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	repo := repository.NewRepository(conf, prod)

	grpcServer := grpc.NewServer()
	pb.RegisterMessageServiceServer(grpcServer, &grpc_server.MessageServer{
		IMessage: repo,
	})

	log.Info("Notification gRPC server started on port 50051")

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	serv := service.NewService(repo)

	hand := handler.NewHandler(serv)

	hand.StartBot(ctx, repo)

}
