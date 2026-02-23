package http

import (
	"fmt"
	"time"

	"github.com/ariashabry/boilerplate-go/docs"
	"github.com/ariashabry/boilerplate-go/helpers/env"
	applog "github.com/ariashabry/boilerplate-go/helpers/log"
	"github.com/ariashabry/boilerplate-go/infras"
	"github.com/gin-contrib/cors"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// HTTP is the HTTP server struct
type HTTP struct {
	Config *env.Config
	DB     *infras.PostgresConn
	Router Router
	Log    *applog.AppLog
}

// ProvideHTTP is the provider function for Wire DI
func ProvideHTTP(db *infras.PostgresConn, config *env.Config, router Router, log *applog.AppLog) *HTTP {
	return &HTTP{
		DB:     db,
		Config: config,
		Router: router,
		Log:    log,
	}
}

// SetupCORS configures CORS middleware
func (h *HTTP) SetupCORS() {
	if !h.Config.AppCorsEnable {
		return
	}

	c := cors.Config{
		AllowWildcard:    h.Config.AppCorsAllowWildcard,
		AllowCredentials: h.Config.AppCorsAllowCredentials,
		AllowOrigins:     h.Config.AppCorsAllowedOrigins,
		AllowHeaders:     h.Config.AppCorsAllowedHeaders,
		AllowMethods:     h.Config.AppCorsAllowedMethods,
		MaxAge:           time.Duration(h.Config.AppCorsMaxAgeSeconds) * time.Second,
	}

	h.Router.Gin.Use(cors.New(c))
}

func (h *HTTP) setupSwaggerDocs() {
	if h.Config.AppEnv == "development" {
		docs.SwaggerInfo.Title = h.Config.AppName
		docs.SwaggerInfo.Version = h.Config.AppRevision
		swaggerURL := fmt.Sprintf("%s/swagger/doc.json", h.Config.AppURL)
		h.Router.Gin.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL(swaggerURL)))
		h.Log.Info("Swagger documentation enabled.", "url", swaggerURL)
	}
}

// SetupAndServe sets up the server and starts serving
func (h *HTTP) SetupAndServe() error {
	h.Router.Gin.SetTrustedProxies([]string{"localhost"})
	h.SetupCORS()
	h.setupSwaggerDocs()
	h.Router.SetupRoutes("")

	address := fmt.Sprintf("%s:%d", h.Config.AppHost, h.Config.AppPort)

	h.Log.WithField("address", address).Info("Starting HTTP server")

	return h.Router.Gin.Run(address)
}
