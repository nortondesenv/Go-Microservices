package repository

import (
	"context"
	"log"
	"time"

	"github.com/nortondesenv/Go-Microservice/internal/models"
	productErrors "github.com/nortondesenv/Go-Microservice/pkg/product_errors"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	productsDB         = "products"
	productsCollection = "products"
)

// productMongoRepo
type productMongoRepo struct {
	mongoDB *mongo.Client
}

// NewProductMongoRepo productMongoRepo constructor
func NewProductMongoRepo(mongoDB *mongo.Client) *productMongoRepo {
	return &productMongoRepo{mongoDB: mongoDB}
}
func (p *productMongoRepo) Create(ctx context.Context, product *models.Product) (*models.Product, error) {
	collection := p.mongoDB.Database(productsDB).Collection(productsCollection)

	product.CreatedAt = time.Now()
	product.UpdatedAt = time.Now()

	result, err := collection.InsertOne(ctx, product, &options.InsertOneOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "InsertOne")
	}

	objectID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, errors.Wrap(productErrors.ErrObjectIDTypeConversion, "result.InsertedID")
	}

	product.ProductID = objectID

	log.Printf("CREATED PRODUCT: %-v", product)

	return product, nil
}

func (p *productMongoRepo) Update(ctx context.Context, product *models.Product) (*models.Product, error) {
	panic("implement me")
}

func (p *productMongoRepo) GetByID(ctx context.Context, productID primitive.ObjectID) (*models.Product, error) {
	panic("implement me")
}

func (p *productMongoRepo) Search(ctx context.Context, search string, page, size int64) ([]*models.Product, error) {
	panic("implement me")
}