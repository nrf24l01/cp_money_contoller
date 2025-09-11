package main

import (
	"log"
	"net"
	"os"

	"github.com/joho/godotenv"
	"github.com/nrf24l01/cp_money_contoller/core"
	grpcHandler "github.com/nrf24l01/cp_money_contoller/grpcServer"
	pb "github.com/nrf24l01/cp_money_contoller/pb"
	"github.com/nrf24l01/cp_money_contoller/redis"
	"google.golang.org/grpc"
)






func main() {
	if os.Getenv("PRODUCTION_ENV") != "true" {
		err := godotenv.Load(".env")
		if err != nil {
			log.Fatalf("failed to load .env: %v", err)
		}
	}

	config, err := core.BuildConfigFromEnv()
	if err != nil {
		log.Fatalf("failed to build config: %v", err)
	}

    lis, err := net.Listen("tcp", ":50051")
    if err != nil {
        log.Fatalf("failed to listen: %v", err)
    }

	redis := redis.NewRedisClient(config)
	defer func() {
		if err := redis.Close(); err != nil {
			log.Printf("failed to close redis: %v", err)
		}
	}()

    grpcServer := grpc.NewServer()
    pb.RegisterWorkerServiceServer(grpcServer, &grpcHandler.GrpcServe{Cfg: config, Redis: &redis})

    log.Println("WorkerService server listening on :50051")
    go func() {
        if err := grpcServer.Serve(lis); err != nil {
            log.Fatalf("failed to serve: %v", err)
        }
    }()

    // Keep the main goroutine alive
    select {}
}
