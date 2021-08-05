package grpc

import (
	"context"

	"github.com/go-playground/validator/v10"

	"github.com/nortondesenv/Go-Microservice/internal/product"
	"github.com/nortondesenv/Go-Microservice/pkg/logger"
	productsServicePB "github.com/nortondesenv/Go-Microservice/proto/product"
)

// productsService gRPC Service
type productService struct {
	log       logger.Logger
	productUC product.UseCase
	validate  *validator.Validate
	productsServicePB.UnimplementedProductsServiceServer
}

// NewProductService productService constructor
func NewProductService(log logger.Logger, productUC product.UseCase, validate *validator.Validate) *productService {
	return &productService{log: log, productUC: productUC, validate: validate}
}

func (p *productService) Create(ctx context.Context, req *productsServicePB.UpdateReq) (*productsServicePB.UpdateRes, error) {
	panic("implement me")
}

func (p *productService) Update(ctx context.Context, req *productsServicePB.UpdateReq) (*productsServicePB.UpdateRes, error) {
	panic("implement me")
}

func (p *productService) GetByID(ctx context.Context, req *productsServicePB.GetByIDReq) (*productsServicePB.GetByIDRes, error) {
	panic("implement me")
}

func (p *productService) Search(ctx context.Context, req *productsServicePB.SearchReq) (*productsServicePB.SearchRes, error) {
	panic("implement me")
}
