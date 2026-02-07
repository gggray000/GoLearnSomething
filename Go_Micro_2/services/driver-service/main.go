package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"ride-sharing/shared/env"
	"ride-sharing/shared/messaging"

	grpc_server "google.golang.org/grpc"
)

var GrpcAddr = ":9092"

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	rabbitMqURI := env.GetString("RABBITMQ_URI", "amqp://guest:guest@rabbitmq:5672/")

	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
		<-sigCh
		cancel()
	}()

	lis, err := net.Listen("tcp", GrpcAddr)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	rabbitmq, err := messaging.NewRabbitMQConnection(rabbitMqURI)
	if err != nil {
		log.Fatal(err)
	}
	defer rabbitmq.Close()

	service := NewService()

	grpcServer := grpc_server.NewServer()
	NewGRPCHandler(grpcServer, service)

	consumer := NewTripConsumer(rabbitmq, service)
	go func(){
		if err := consumer.Listen(); err != nil{
			log.Fatalf("Failed to listen to the message: %v", err)
		}
	}()

	log.Printf("Starting gRPC server Driver service on port %s", lis.Addr().String())

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("gRPC server failed to serve: %v", err)
			cancel()
		}
	}()

	<-ctx.Done()
	log.Printf("Shutting down gRPC server...")
	grpcServer.GracefulStop()
}
