package service

import (
	"context"

	"github.com/ariashabry/boilerplate-go/helpers/log"
	"github.com/ariashabry/boilerplate-go/internal/domains/product/model/dto"
	"github.com/ariashabry/boilerplate-go/internal/domains/product/repository"
)

type ProductService interface {
	GetProduct(ctx context.Context) (res []dto.Product, err error)
}

type ProductServiceImpl struct {
	Repo repository.ProductRepositoryPostgres
	// inject ext service for use in this service hrex
	Log *log.AppLog
}

func ProvideProductServiceImpl(repo repository.ProductRepositoryPostgres, l *log.AppLog) *ProductServiceImpl {
	return &ProductServiceImpl{Repo: repo, Log: l}
}

func (s *ProductServiceImpl) GetProduct(ctx context.Context) (res []dto.Product, err error) {
	// do something
	// _ = s.Repo.GetFoo(ctx)
	s.Log.Info("[GetProduct] Success")
	return
}

func (s *ProductServiceImpl) GetList(ctx context.Context, queries map[string][]string) (res []dto.Product, err error) {

	res, err = s.Repo.GetAllProducts(ctx)
	if err != nil {
		s.Log.Error("[GetList] Failed to get List Product")
	}
	s.Log.Info("[GetList] Success")
	return
}
