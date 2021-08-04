package grpc

import (
	"github.com/nortondesenv/Go-Microservice/internal/product"
	"github.com/nortondesenv/Go-Microservice/pkg/logger"
)

// productGRPCService gRPC Service
type productGRPCService struct {
	log       logger.Logger
	productUC product.UseCase
}
