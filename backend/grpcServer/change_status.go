package grpcServer

import (
	"context"
	"encoding/json"

	"github.com/nrf24l01/cp_money_contoller/models"
	"github.com/nrf24l01/cp_money_contoller/pb"
)

func (s *GrpcServe) ChangeStatus(ctx context.Context, req *pb.ChangeStatusRequest) (*pb.ChangeStatusResponse, error) {
	var taskStatus models.TaskStatus
	if err := s.DB.Where("id = ?", req.Uuid).First(&taskStatus).Error; err != nil {
		return nil, err
	}

	taskStatus.Status = req.Status
	var jsonData interface{}
	if err := json.Unmarshal([]byte(req.Payload), &jsonData); err != nil {
		return nil, err
	}
	jsonBytes, err := json.Marshal(jsonData)
	if err != nil {
		return nil, err
	}
	jsonString := string(jsonBytes)
	taskStatus.Result = &jsonString
	if err := s.DB.Save(&taskStatus).Error; err != nil {
		return nil, err
	}

	return &pb.ChangeStatusResponse{
		Ok: true,
	}, nil
}