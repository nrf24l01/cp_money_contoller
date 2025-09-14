package grpcServer

import (
	"context"
	"encoding/json"

	"github.com/nrf24l01/cp_money_contoller/models"
	"github.com/nrf24l01/cp_money_contoller/pb"
)

func (s *GrpcServe) ChangeStatus(ctx context.Context, req *pb.ChangeStatusRequest) (*pb.ChangeStatusResponse, error) {
	// parse payload JSON
	var jsonData interface{}
	if err := json.Unmarshal([]byte(req.Payload), &jsonData); err != nil {
		return nil, err
	}
	jsonBytes, err := json.Marshal(jsonData)
	if err != nil {
		return nil, err
	}
	jsonString := string(jsonBytes)

	// parse logs JSON array
	var logs []string
	if err := json.Unmarshal([]byte(req.Logs), &logs); err != nil {
		return nil, err
	}
	// use raw JSON string for logs field to match JSON column
	logsString := req.Logs

	// update fields without loading full record to avoid Scan error
	if err := s.DB.Model(&models.TaskStatus{}).
		Where("id = ?", req.Uuid).
		Updates(map[string]interface{}{"status": req.Status, "result": jsonString, "logs": logsString}).Error; err != nil {
		return nil, err
	}

	return &pb.ChangeStatusResponse{
		Ok: true,
	}, nil
}