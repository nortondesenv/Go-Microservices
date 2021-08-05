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
		ImageURL:    &req.ImageURL,
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
	span, ctx := opentracing.StartSpanFromContext(ctx, "productService.Update")
	defer span.Finish()

	prodID, err := primitive.ObjectIDFromHex(req.GetProductID())
	if err != nil {
		p.log.Errorf("primitive.ObjectIDFromHex: %v", err)
		return nil, grpcErrors.ErrorResponse(err, err.Error())
	}
	catID, err := primitive.ObjectIDFromHex(req.GetCategoryID())
	if err != nil {
		p.log.Errorf("primitive.ObjectIDFromHex: %v", err)
		return nil, grpcErrors.ErrorResponse(err, err.Error())
	}

	prod := &models.Product{
		ProductID:   prodID,
		CategoryID:  catID,
		Name:        req.GetName(),
		Description: req.GetDescription(),
		Price:       req.GetPrice(),
		ImageURL:    &req.ImageURL,
		Photos:      req.GetPhotos(),
		Quantity:    req.GetQuantity(),
		Rating:      int(req.GetRating()),
	}

	if err := p.validate.StructCtx(ctx, prod); err != nil {
		p.log.Errorf("validate.StructCtx: %v", err)
		return nil, grpcErrors.ErrorResponse(err, err.Error())
	}

	update, err := p.productUC.Update(ctx, prod)
	if err != nil {
		p.log.Errorf("productUC.Update: %v", err)
		return nil, grpcErrors.ErrorResponse(err, err.Error())
	}

	return &productsServicePB.UpdateRes{Product: update.ToProto()}, nil
}

// GetByID Get single product by id
func (p *productService) GetByID(ctx context.Context, req *productsServicePB.GetByIDReq) (*productsServicePB.GetByIDRes, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "productService.GetByID")
	defer span.Finish()

	prodID, err := primitive.ObjectIDFromHex(req.GetProductID())
	if err != nil {
		p.log.Errorf("primitive.ObjectIDFromHex: %v", err)
		return nil, grpcErrors.ErrorResponse(err, err.Error())
	}

	prod, err := p.productUC.GetByID(ctx, prodID)
	if err != nil {
		p.log.Errorf("productUC.GetByID: %v", err)
		return nil, grpcErrors.ErrorResponse(err, err.Error())
	}

	return &productsServicePB.GetByIDRes{Product: prod.ToProto()}, nil
}

// Search Search products
func (p *productService) Search(ctx context.Context, req *productsServicePB.SearchReq) (*productsServicePB.SearchRes, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "productService.Search")
	defer span.Finish()

	products, err := p.productUC.Search(ctx, req.GetSearch(), req.GetPage(), req.GetSize())
	if err != nil {
		p.log.Errorf("productUC.Search: %v", err)
		return nil, grpcErrors.ErrorResponse(err, err.Error())
	}

	p.log.Infof("PRODUCTS: %-v", products)

	return nil, nil
}