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
