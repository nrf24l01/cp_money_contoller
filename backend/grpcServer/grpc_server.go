package grpcServer

import (
	"github.com/nrf24l01/cp_money_contoller/core"
	"github.com/nrf24l01/cp_money_contoller/pb"
	"github.com/nrf24l01/cp_money_contoller/redis"
)

type GrpcServe struct {
    pb.UnimplementedWorkerServiceServer
    Cfg *core.Config
    Redis *redis.Redis
}