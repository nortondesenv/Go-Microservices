package usecase

import (
	"context"

	"github.com/nortondesenv/Go-Microservice/internal/models"
	"github.com/nortondesenv/Go-Microservice/internal/product"
	"github.com/nortondesenv/Go-Microservice/pkg/logger"
	"github.com/nortondesenv/Go-Microservice/pkg/utils"
	"github.com/opentracing/opentracing-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// productUC
type productUC struct {
	productRepo product.MongoRepository
	log         logger.Logger
}

// NewProductUC productUC constructor
func NewProductUC(productRepo product.MongoRepository, log logger.Logger) *productUC {
	return &productUC{productRepo: productRepo, log: log}
}

// Create Create new product
func (p *productUC) Create(ctx context.Context, product *models.Product) (*models.Product, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "productUC.Create")
	defer span.Finish()
	return p.productRepo.Create(ctx, product)
}

// Update single product
func (p *productUC) Update(ctx context.Context, product *models.Product) (*models.Product, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "productUC.Update")
	defer span.Finish()
	return p.productRepo.Update(ctx, product)
}

// GetByID Get single product by id
func (p *productUC) GetByID(ctx context.Context, productID primitive.ObjectID) (*models.Product, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "productUC.GetByID")
	defer span.Finish()
	return p.productRepo.GetByID(ctx, productID)
}

// Search Search products
func (p *productUC) Search(ctx context.Context, search string, pagination *utils.Pagination) (*models.ProductsList, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "productUC.Search")
	defer span.Finish()
	return p.productRepo.Search(ctx, search, pagination)
}
