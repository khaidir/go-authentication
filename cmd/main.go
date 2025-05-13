package main

import (
	"log"
	"os"

	"auth-services/config"
	"auth-services/internal/http/routes"
	"auth-services/internal/server/grpc"
	"auth-services/pkg/logger"
	"auth-services/pkg/middleware"
	"auth-services/pkg/observability"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func main() {
	config.DB = config.InitDB()
	config.Redis = config.InitRedis()

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(logger.TraceMiddleware())
	r.Use(gzip.Gzip(gzip.BestCompression))
	r.Use(middleware.RequestContextMiddleware())
	r.Use(otelgin.Middleware("auth-service"))
	r.Use(middleware.GinLoggerMiddleware())

	if os.Getenv("APP_ENV") == "local" {
		observability.InitTracerLocal()
	} else {
		observability.InitTracerOTLP("auth-service")
	}

	routes.InitRoutes(r)
	r.Run(":" + os.Getenv("PORT"))

	log.Println("Starting gRPC Auth Service...")
	if err := grpc.StartGRPCServer(); err != nil {
		log.Fatalf("failed to start gRPC server: %v", err)
	}
}
