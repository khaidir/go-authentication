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

	name := os.Getenv("SERVICE_NAME")

	config.DB = config.InitDB()
	config.Redis = config.InitRedis()
	logger.InitLogger()

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(logger.TraceMiddleware())
	r.Use(middleware.RequestContextMiddleware())
	r.Use(gzip.Gzip(gzip.BestCompression))
	r.Use(otelgin.Middleware(name))
	r.Use(middleware.GinLoggerMiddleware())

	if os.Getenv("APP_ENV") == "local" {
		observability.InitTracerLocal()
	} else {
		observability.InitTracerOTLP(name)
	}

	routes.InitRoutes(r)

	// Jalankan HTTP server dalam goroutine
	go func() {
		log.Println("Starting HTTP server on port " + os.Getenv("PORT"))
		if err := r.Run(":" + os.Getenv("PORT")); err != nil {
			log.Fatalf("failed to start HTTP server: %v", err)
		}
	}()

	// Jalankan gRPC server di thread utama
	log.Println("Starting gRPC Auth Service...")
	if err := grpc.StartGRPCServer(); err != nil {
		log.Fatalf("failed to start gRPC server: %v", err)
	}
}
