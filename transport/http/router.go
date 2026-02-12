package http

import (
	"net/http"

	"github.com/ariashabry/boilerplate-go/internal/handler/product"
	"github.com/gin-gonic/gin"
)

// DomainHandlers contains all domain-specific handlers
type DomainHandlers struct {
	ProductHandler product.ProductHandler
}

// Router is the router struct containing handlers and gin engine
type Router struct {
	DomainHandlers DomainHandlers
	Gin            *gin.Engine
}

// ProvideRouter is the provider function for Wire DI
func ProvideRouter(domainHandlers DomainHandlers, gin *gin.Engine) Router {
	return Router{
		DomainHandlers: domainHandlers,
		Gin:            gin,
	}
}

// SetupRoutes sets up all routing for this server
func (r *Router) SetupRoutes(group string) {
	public := r.Gin.Group(group)

	// Health check endpoint
	public.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "ok",
		})
	})

	// Add product routes from product handler
	r.DomainHandlers.ProductHandler.Router(public.Group(""))
}
