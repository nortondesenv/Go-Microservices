package usecase

import (
	"context"

	"github.com/nortondesenv/Go-Microservice/internal/models"
	"github.com/nortondesenv/Go-Microservice/internal/product"
	"github.com/nortondesenv/Go-Microservice/pkg/logger"
)

// productUC
type productUC struct {
	productRepo product.MongoRepository
	log         logger.Logger
}

func (p *productUC) Create(ctx context.Context, product *models.Product) (*models.Product, error) {
	panic("implement me")
}

func (p *productUC) Update(ctx context.Context, product *models.Product) (*models.Product, error) {
	panic("implement me")
}

func (p *productUC) GetByID(ctx context.Context, productID string) (*models.Product, error) {
	panic("implement me")
}

func (p *productUC) Search(ctx context.Context, search string, page, size int64) ([]*models.Product, error) {
	panic("implement me")
}
