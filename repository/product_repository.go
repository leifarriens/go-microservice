package repository

import (
	"context"
	"log"

	"github.com/leifarriens/go-microservice/model"
	"gorm.io/gorm"
)

type ProductRepository interface {
	Create(ctx context.Context, product *model.Product) (*uint, error)
	FindAll(ctx context.Context, limit int, offset int) ([]*model.Product, error)
	FindByID(ctx context.Context, id string) (*model.Product, error)
}

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	err := db.AutoMigrate(&model.Product{})

	if err != nil {
		log.Fatalln(err)
	}

	return &productRepository{db: db}
}

func (r *productRepository) Create(ctx context.Context, product *model.Product) (*uint, error) {
	result := r.db.Create(&product)

	err := result.Error
	if err != nil {
		return nil, err
	}

	return &product.ID, nil
}

func (r *productRepository) FindAll(ctx context.Context, limit int, offset int) ([]*model.Product, error) {
	var products []*model.Product

	result := r.db.Limit(limit).Offset(offset).Find(&products)
	if result.Error != nil {
		return nil, result.Error
	}

	return products, nil
}

func (r *productRepository) FindByID(ctx context.Context, id string) (*model.Product, error) {
	var products model.Product

	result := r.db.First(&products, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return &products, nil
}
