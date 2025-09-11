package main

import (
	"context"
	"log"
	"net"
	"time"

	pb "github.com/nrf24l01/cp_money_contoller/pb"
	"google.golang.org/grpc"
)

type server struct {
    pb.UnimplementedWorkerServiceServer
}

func (s *server) RegisterWorker(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
    // Пример: генерируем фиктивный UUID и ключ
    return &pb.RegisterResponse{
        Uuid: "worker-1234",
        Key:  "secret-key",
    }, nil
}

func (s *server) GetBuilding(ctx context.Context, req *pb.GetBuildingRequest) (*pb.GetBuildingResponse, error) {
    return &pb.GetBuildingResponse{
        Uuid:     "task-5678",
        Task:     "build-house",
        Payload:  "{\"materials\": [\"wood\", \"stone\"]}",
        UnixTime: uint64(time.Now().Unix()),
    }, nil
}

func (s *server) CompleteTask(ctx context.Context, req *pb.CompleteTaskRequest) (*pb.CompleteTaskResponse, error) {
    log.Printf("Task completed: UUID=%s, Payload=%s, Time=%d", req.Uuid, req.Payload, req.UnixTime)
    return &pb.CompleteTaskResponse{Ok: true}, nil
}

func main() {
    lis, err := net.Listen("tcp", ":50051")
    if err != nil {
        log.Fatalf("failed to listen: %v", err)
    }

    grpcServer := grpc.NewServer()
    pb.RegisterWorkerServiceServer(grpcServer, &server{})

    log.Println("WorkerService server listening on :50051")
    if err := grpcServer.Serve(lis); err != nil {
        log.Fatalf("failed to serve: %v", err)
    }
}
