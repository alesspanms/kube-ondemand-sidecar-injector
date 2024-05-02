package main

import (
	"net/http"
	"os"
	"strings"
	"time"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/alesspanms/kube-ondemand-sidecar-injector/internal/controllers/injector"
	_ "github.com/alesspanms/kube-ondemand-sidecar-injector/internal/docs"
	"github.com/alesspanms/kube-ondemand-sidecar-injector/internal/kube"
	"github.com/alesspanms/kube-ondemand-sidecar-injector/internal/logging"
)

// for generating swagger docs execute the following command at the src/go folder
// go install github.com/swaggo/swag/cmd/swag@latest
// and then
// export PATH=$PATH:$(go env GOPATH)/bin
// and finally the following command
// swag init --dir ./cmd/kube-ondemand-sidecar-injector/,./internal --output ./internal/docs/
//go:generate swag init

//	@title		Kubernetes OnDemand Sidecar Injector API
//	@version	1.0

//	@license.name	MIT
//	@license.url	https://mit-license.org/

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name X-API-KEY

// @externalDocs.description	OpenAPI
// @externalDocs.url			https://swagger.io/resources/open-api/
func main() {
	r := gin.New()

	//initialize required modules and types
	logger := logging.New()
	kubeClient := kube.New(logger, os.Getenv("SIDECAR_NAME_PREFIX"))
	injectorController := injector.New(logger, kubeClient)

	// Attach Zap logger middleware from logger-module
	r.Use(ginzap.Ginzap(logger.Log(), time.RFC3339, true))
	r.Use(AuthMiddleware(os.Getenv("SECRET_API_KEY")))

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Kubernetes OnDemand Sidecar Injector API up & running"})
	})

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Routes Group, Routes and Handlers...
	api := r.Group("/api")
	{
		injectorApi := api.Group("/injector")
		{
			injectorApi.POST("/GetDeployments", injectorController.GetDeployments)
			injectorApi.POST("/GetSingleDeployment", injectorController.GetSingleDeployment)
			injectorApi.POST("/SetSidecar", injectorController.SetSidecar)
			injectorApi.POST("/ClearSidecar", injectorController.ClearSidecar)
		}
	}

	// Run the server
	r.Run(":8080")
}

func AuthMiddleware(apiKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/api/") {
			apiKeyHeader := c.GetHeader("X-API-KEY")
			if apiKeyHeader != apiKey {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
				c.Abort()
				return
			}
		}

		c.Next()
	}
}
