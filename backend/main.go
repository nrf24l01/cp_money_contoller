package main

import (
	"log"
	"net"
	"os"

	"github.com/joho/godotenv"
	"github.com/nrf24l01/cp_money_contoller/core"
	grpcHandler "github.com/nrf24l01/cp_money_contoller/grpcServer"
	"github.com/nrf24l01/cp_money_contoller/models"
	pb "github.com/nrf24l01/cp_money_contoller/pb"
	"github.com/nrf24l01/cp_money_contoller/redis"
	"google.golang.org/grpc"
)



func main() {
	// Load environment variables from .env file in non-production environments
	if os.Getenv("PRODUCTION_ENV") != "true" {
		err := godotenv.Load(".env")
		if err != nil {
			log.Fatalf("failed to load .env: %v", err)
		}
	}

	// Load configuration from environment variables
	config, err := core.BuildConfigFromEnv()
	if err != nil {
		log.Fatalf("failed to build config: %v", err)
	}

	// Initialize Redis
	redis := redis.NewRedisClient(config)
	defer func() {
		if err := redis.Close(); err != nil {
			log.Printf("failed to close redis: %v", err)
		}
	}()

	// Initialize database
	db := models.RegisterPostgres(config)
	
	// Start gRPC server
    lis, err := net.Listen("tcp", ":50051")
    if err != nil {
        log.Fatalf("failed to listen: %v", err)
    }

	// Create gRPC server and register services
    grpcServer := grpc.NewServer()
    pb.RegisterWorkerServiceServer(grpcServer, &grpcHandler.GrpcServe{Cfg: config, Redis: &redis, DB: db})

	// Start serving
    log.Println("WorkerService server listening on :50051")
    go func() {
        if err := grpcServer.Serve(lis); err != nil {
            log.Fatalf("failed to serve: %v", err)
        }
    }()
    select {}
}
