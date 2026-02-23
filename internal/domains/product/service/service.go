package service

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/ariashabry/boilerplate-go/helpers/log"
	"github.com/ariashabry/boilerplate-go/infras"
	"github.com/ariashabry/boilerplate-go/internal/domains/product/model/dto"
	"github.com/ariashabry/boilerplate-go/internal/domains/product/repository"
	"github.com/redis/go-redis/v9"
)

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

type ProductService interface {
	GetProduct(ctx context.Context) (res []dto.Product, err error)
	GetList(ctx context.Context) (res []dto.Product, err error)
	GetProductById(ctx context.Context, id int) (res dto.Product, err error)
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
	redisClient := s.RedisClient.Client
	cacheKey := "product:" + strconv.Itoa(id)

	val, err := redisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		s.Log.Info("Cache hit, using cached data")
		err = json.Unmarshal([]byte(val), &res)
		if err == nil {
			return res, nil
		}
		s.Log.WithError(err).Warn("Failed to unmarshal cache, fallback to DB")
	}
	if err != redis.Nil {
		s.Log.WithError(err).Warn("Redis error, fallback to DB")
	}
	s.Log.Info("Get from database")
	productData, err := s.Repo.GetProductByID(ctx, id)
	if err != nil {
		s.Log.WithError(err).Error("Failed to get data from database")
		return res, err
	}
	cacheData, _ := json.Marshal(productData)
	if redisClient != nil {
		if err := redisClient.Set(ctx, cacheKey, cacheData, 20*time.Minute).Err(); err != nil {
			s.Log.WithError(err).Warn("Failed to set cache, skipping")
		}
	}
	return productData, nil
}
