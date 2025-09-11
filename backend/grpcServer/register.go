package grpcServer

import (
	"context"

	pb "github.com/nrf24l01/cp_money_contoller/pb"
)

func (s *GrpcServe) RegisterWorker(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
    uuid, key, err := s.Redis.NewWorker(req.Id)
    if err != nil {
        return nil, err
    }
    return &pb.RegisterResponse{
        Uuid: uuid,
        Key:  key,
    }, nil
}

func (s *GrpcServe) GetBuilding(ctx context.Context, req *pb.GetBuildingRequest) (*pb.GetBuildingResponse, error) {
    // Пример: возвращаем фиктивное здание
    return &pb.GetBuildingResponse{
        Uuid:     "building-1234",
        Task:     "collect_coins",
        Payload:  "some_payload_data",
        UnixTime: 1234567890,
    }, nil
}

func (s *GrpcServe) CompleteTask(ctx context.Context, req *pb.CompleteTaskRequest) (*pb.CompleteTaskResponse, error) {
    // Пример: всегда возвращаем успех
    return &pb.CompleteTaskResponse{
        Ok: true,
    }, nil
}