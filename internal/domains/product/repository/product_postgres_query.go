package repository

import (
	"context"

	"github.com/ariashabry/boilerplate-go/internal/domains/product/model/dto"
)

type ProductPostgresRepo interface {
	GetAllProducts(ctx context.Context) ([]dto.Product, error)
	GetProductByID(ctx context.Context, id int) (dto.Product, error)
}

func (repo *ProductRepositoryPostgresImpl) GetAllProducts(ctx context.Context) ([]dto.Product, error) {
	query := `SELECT id, name, price, description, category, status FROM product`
	var products []dto.Product
	// Use Read connection for read operations
	err := repo.DB.Read.Raw(query).Scan(&products).Error
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (repo *ProductRepositoryPostgresImpl) GetProductByID(ctx context.Context, id int) (dto.Product, error) {
	var product dto.Product
	err := repo.DB.Read.Where("id = ?", id).First(&product).Error
	if err != nil {
		return dto.Product{}, err
	}
	return product, nil
}
