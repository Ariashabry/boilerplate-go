package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/ariashabry/boilerplate-go/helpers/log"
	"github.com/ariashabry/boilerplate-go/infras"
	"github.com/ariashabry/boilerplate-go/internal/domains/product/model/dto"
	"github.com/ariashabry/boilerplate-go/internal/domains/product/repository"
	"github.com/redis/go-redis/v9"
)

type ProductService interface {
	GetProduct(ctx context.Context) (res []dto.Product, err error)
	GetList(ctx context.Context) (res []dto.Product, err error)
	GetProductById(ctx context.Context, id int) (res dto.Product, err error)
}

// info: Dependency Injection Container
type ProductServiceImpl struct {
	Repo        repository.ProductRepositoryPostgres
	Log         *log.AppLog
	RedisClient *infras.Redis
}

// info: provider func / constructor
// wire akan membaca func ini untuk melakukan dependency injection

func ProvideProductServiceImpl(repo repository.ProductRepositoryPostgres, l *log.AppLog, redis *infras.Redis) *ProductServiceImpl {
	return &ProductServiceImpl{Repo: repo, Log: l, RedisClient: redis}
}

func (s *ProductServiceImpl) GetProduct(ctx context.Context) (res []dto.Product, err error) {
	// do something
	// _ = s.Repo.GetFoo(ctx)
	s.Log.Info("[GetProduct] Success")
	return
}

func (s *ProductServiceImpl) GetList(ctx context.Context) (res []dto.Product, err error) {

	// info: get data from database
	s.Log.Info("Get from database")
	productData, err := s.Repo.GetAllProducts(ctx)
	if err != nil {
		s.Log.Error("Failed to get data from database")
		return nil, err
	}
	s.Log.Info("[GetList] Success")
	return productData, nil

}

func (s *ProductServiceImpl) GetProductById(ctx context.Context, id int) (res dto.Product, err error) {
	// info: get data from redis
	redisClient := s.RedisClient.Client
	cacheKey := fmt.Sprintf("product_%d", id)
	val, err := redisClient.Get(ctx, cacheKey).Result()

	// if data not found
	if err == redis.Nil {
		// get data from database
		s.Log.Info("Get from database")
		productData, err := s.Repo.GetProductByID(ctx, id)
		if err != nil {
			s.Log.Error("Failed to get data from database")
			return res, err
		}

		// set productData to cache
		cacheData, err := json.Marshal(productData)
		if err != nil {
			s.Log.Error("Failed to marshal product data")
			return res, err
		}
		if redisClient != nil {
			err = redisClient.Set(ctx, cacheKey, cacheData, 20*time.Minute).Err()
			if err != nil {
				s.Log.Error("Failed to set product cache")
			}
		} else {
			s.Log.Warn("Redis client is nil, skipping cache set")
		}

		return productData, nil
	} else if err != nil {
		s.Log.Error("Failed to get data from cache")
		return res, err
	} else {
		// Jika data ditemukan di Redis, unmarshall data
		s.Log.Info("Cache hit, using cached data")
		err = json.Unmarshal([]byte(val), &res)
		if err != nil {
			s.Log.Error("Failed to unmarshal cached data")
			return res, err
		}
		s.Log.Info("[GetProductByID] Success")
		return res, nil
	}
}
