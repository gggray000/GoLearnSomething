package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"ride-sharing/services/trip-service/internal/infrastructure/events"
	grpc_handler "ride-sharing/services/trip-service/internal/infrastructure/grpc"
	"ride-sharing/services/trip-service/internal/infrastructure/repository"
	"ride-sharing/services/trip-service/internal/service"
	"syscall"

	"ride-sharing/shared/db"
	"ride-sharing/shared/env"
	"ride-sharing/shared/messaging"
	"ride-sharing/shared/tracing"

	grpc_server "google.golang.org/grpc"
)

var GrpcAddr = ":9093"

func main() {
	
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mongoClient, err := db.NewMongoClient(ctx, db.NewMongoDefaultConfig())
	if err != nil {
		log.Fatalf("Failed to initialize MongoDB: %v", err)
	}
	defer mongoClient.Disconnect(ctx)

	mongoDb := db.GetDatabase(mongoClient, db.NewMongoDefaultConfig())

	mongoDBRepo := repository.NewMongoRepository(mongoDb)
	if err := mongoDBRepo.CreateRideFareTTLIndex(); err != nil {
		log.Fatal("failed to create TTL index:", err)
	}

	svc := service.NewService(mongoDBRepo)

	traceCfg := tracing.Config{
		ServiceName: "driver-service",
		Environment: env.GetString("ENVIRONMENT","development"),
		JaegerEndpoint: env.GetString("JAEGER_ENDPOINT", "http://jaeger:14268/api/traces"),
	}

	shutdownTracing, err := tracing.InitTracer(traceCfg)
	if err != nil {
		log.Fatalf("Failed to initialize the tracer: %w", err)
	}

	defer shutdownTracing(ctx)

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

	publisher := events.NewTripEventPublisher(rabbitmq)

	driverConsumer := events.NewDriverConsumer(rabbitmq, svc)
	go driverConsumer.Listen()

	paymentConsumer := events.NewPaymentConsumer(rabbitmq, svc)
	go paymentConsumer.Listen()

	grpcServer := grpc_server.NewServer(tracing.WithTracingInterceptors()...)

	grpc_handler.NewGRPCHandler(grpcServer, svc, publisher)

	log.Printf("Starting gRPC server Trip service on port %s", lis.Addr().String())

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
