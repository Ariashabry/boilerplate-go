package repository

import (
	"github.com/ariashabry/boilerplate-go/infras"
)

type ProductRepositoryPostgres interface {
	ProductPostgresRepo
}

type ProductRepositoryPostgresImpl struct {
	DB *infras.PostgresConn
}

func ProvideProductRepositoryPostgresImpl(db *infras.PostgresConn) *ProductRepositoryPostgresImpl {
	return &ProductRepositoryPostgresImpl{DB: db}
}
