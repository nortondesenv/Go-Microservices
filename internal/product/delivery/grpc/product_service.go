package grpc

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/opentracing/opentracing-go"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/nortondesenv/Go-Microservice/internal/models"
	"github.com/nortondesenv/Go-Microservice/internal/product"
	grpcErrors "github.com/nortondesenv/Go-Microservice/pkg/grpc_errors"
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

func (p *productService) Create(ctx context.Context, req *productsServicePB.CreateReq) (*productsServicePB.CreateRes, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "productService.Create")
	defer span.Finish()

	catID, err := primitive.ObjectIDFromHex(req.GetCategoryID())
	if err != nil {
		p.log.Errorf("primitive.ObjectIDFromHex: %v", err)
		return nil, grpcErrors.ErrorResponse(err, err.Error())
	}

	prod := &models.Product{
		CategoryID:  catID,
		Name:        req.GetName(),
		Description: req.GetDescription(),
		Price:       req.GetPrice(),
		ImageURL:    req.GetImageURL(),
		Photos:      req.GetPhotos(),
		Quantity:    req.GetQuantity(),
		Rating:      int(req.GetRating()),
	}

	if err := p.validate.StructCtx(ctx, prod); err != nil {
		p.log.Errorf("validate.StructCtx: %v", err)
		return nil, grpcErrors.ErrorResponse(err, err.Error())
	}

	created, err := p.productUC.Create(ctx, prod)
	if err != nil {
		p.log.Errorf("productUC.Create: %v", err)
		return nil, grpcErrors.ErrorResponse(err, err.Error())
	}

	return &productsServicePB.CreateRes{Product: created.ToProto()}, nil
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
