//go:build wireinject
// +build wireinject

package main

import (
	"github.com/ariashabry/boilerplate-go/helpers/env"
	"github.com/ariashabry/boilerplate-go/helpers/log"
	"github.com/ariashabry/boilerplate-go/infras"
	"github.com/ariashabry/boilerplate-go/internal/domains/product/repository"
	"github.com/ariashabry/boilerplate-go/internal/domains/product/service"
	"github.com/ariashabry/boilerplate-go/internal/handler/product"
	"github.com/ariashabry/boilerplate-go/internal/migration"
	transport "github.com/ariashabry/boilerplate-go/transport/http"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

// Wiring for configurations
var configurations = wire.NewSet(
	env.Get,
)

// Wiring for persistences (infrastructure)
var persistences = wire.NewSet(
	infras.ProvideRedis,
	infras.ProvidePostgresConn,
)

// Wiring for product domain
var productDomain = wire.NewSet(
	service.ProvideProductServiceImpl,
	wire.Bind(new(service.ProductService), new(*service.ProductServiceImpl)),
	repository.ProvideProductRepositoryPostgresImpl,
	wire.Bind(new(repository.ProductRepositoryPostgres), new(*repository.ProductRepositoryPostgresImpl)),
)

// Wiring for all domains
var domains = wire.NewSet(
	productDomain,
)

// Wiring for HTTP routing
var routing = wire.NewSet(
	wire.Struct(new(transport.DomainHandlers), "*"),
	product.ProvideProductHandler,
	provideGinEngine,
	transport.ProvideRouter,
)

// Wiring for migrations
var migrations = wire.NewSet(
	migration.ProvideMigrationService,
)

// provideGinEngine creates a new Gin engine
func provideGinEngine() *gin.Engine {
	return gin.New()
}

// Wiring for everything
func InitializeService(l *log.AppLog) *transport.HTTP {
	wire.Build(
		// configurations
		configurations,
		// persistences
		persistences,
		// domains
		domains,
		// routing
		routing,
		// HTTP transport
		transport.ProvideHTTP,
	)
	return &transport.HTTP{}
}

// InitializeMigrations wires up the migration service
func InitializeMigrations(l *log.AppLog) *migration.MigrationServiceImpl {
	wire.Build(
		// configurations
		configurations,
		// persistences
		persistences,
		// migrations
		migrations,
	)
	return nil
}
