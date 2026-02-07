package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/gin-gonic/gin"

	productHttp "github.com/tu-usuario/product-crud-hexagonal/internal/adapters/http"
	"github.com/tu-usuario/product-crud-hexagonal/internal/adapters/repository"
	"github.com/tu-usuario/product-crud-hexagonal/internal/core/services"
	appConfig "github.com/tu-usuario/product-crud-hexagonal/internal/platform/config"
	"github.com/tu-usuario/product-crud-hexagonal/internal/platform/logger"
)

func main() {
	// Load configuration
	cfg := appConfig.LoadConfig()

	// Initialize logger
	appLogger := logger.NewLogger(cfg)
	appLogger.Info("Starting product service", "port", cfg.Port)

	// AWS SDK Configuration
	awsCfg, err := awsconfig.LoadDefaultConfig(context.TODO(), awsconfig.WithRegion(cfg.AWSRegion))
	if err != nil {
		appLogger.Error("unable to load SDK config", "error", err)
		os.Exit(1)
	}

	dbClient := dynamodb.NewFromConfig(awsCfg)

	// Dependency Injection
	productRepo := repository.NewDynamoDBRepository(dbClient, cfg.DynamoDBTable)
	productService := services.NewProductService(productRepo, appLogger)
	productHandler := productHttp.NewProductHandler(productService, appLogger)

	// Router Setup
	if cfg.LogLevel == "debug" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	// Middleware
	router.Use(gin.Recovery())

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "UP",
			"timestamp": time.Now().UTC(),
		})
	})

	// API routes
	v1 := router.Group("/api/v1")
	{
		products := v1.Group("/products")
		{
			products.POST("", productHandler.Create)
			products.GET("", productHandler.List)
			products.GET("/:id", productHandler.Get)
			products.PUT("/:id", productHandler.Update)
			products.DELETE("/:id", productHandler.Delete)
		}
	}

	// Graceful Shutdown
	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	go func() {
		appLogger.Info("Server starting", "port", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			appLogger.Error("listen error", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	appLogger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		appLogger.Error("Server forced to shutdown", "error", err)
		os.Exit(1)
	}

	appLogger.Info("Server exiting")
}
