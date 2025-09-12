package grpcServer

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/nrf24l01/cp_money_contoller/models"
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

func (s *GrpcServe) GetTask(ctx context.Context, req *pb.GetTaskRequest) (*pb.GetTaskResponse, error) {
    if !s.Redis.ValidateWorker(req.WorkerUuid, req.WorkerKey) {
        return nil, errors.New("invalid worker credentials")
    }
    var task models.Task
    tx := s.DB.Begin()
    err := tx.Raw(`
        SELECT t.* FROM tasks t
        WHERE NOT EXISTS (SELECT 1 FROM task_statuses ts WHERE ts.task_id = t.id)
        LIMIT 1
        FOR UPDATE SKIP LOCKED
    `).Scan(&task).Error
    if err != nil {
        tx.Rollback()
        return nil, err
    }
    if task.ID.String() == "00000000-0000-0000-0000-000000000000" {
        tx.Rollback()
        return &pb.GetTaskResponse{}, nil
    }
    taskStatus := models.TaskStatus{}
    taskStatus.ID = uuid.New()
    taskStatus.TaskID = task.ID
    taskStatus.Status = "assigned"
    if err := tx.Create(&taskStatus).Error; err != nil {
        tx.Rollback()
        return nil, err
    }
    if err := tx.Commit().Error; err != nil {
        return nil, err
    }
    unixTime := uint64(0)
    if !task.CreatedAt.IsZero() {
        unixTime = uint64(task.CreatedAt.Unix())
    }
    return &pb.GetTaskResponse{
        Uuid:     task.ID.String(),
        Task:     task.Type,
        Payload:  task.Payload,
        UnixTime: unixTime,
    }, nil
}

func (s *GrpcServe) CompleteTask(ctx context.Context, req *pb.CompleteTaskRequest) (*pb.CompleteTaskResponse, error) {
    // Пример: всегда возвращаем успех
    return &pb.CompleteTaskResponse{
        Ok: true,
    }, nil
}