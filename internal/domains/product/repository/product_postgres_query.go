package repository

import (
	"context"

	"github.com/ariashabry/boilerplate-go/internal/domains/product/model/dto"
)

type ProductPostgresRepo interface {
	GetAllProducts(ctx context.Context) ([]dto.Product, error)
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
