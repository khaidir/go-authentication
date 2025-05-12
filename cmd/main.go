package main

import (
	"log"
	"os"

	"auth-services/config"
	"auth-services/internal/http/routes"
	"auth-services/internal/server/grpc"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func main() {
	config.InitDB()
	config.InitRedis()

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(gzip.Gzip(gzip.BestCompression))
	r.Use(otelgin.Middleware("user-service"))

	routes.InitRoutes(r)

	r.Run(":" + os.Getenv("PORT"))

	log.Println("Starting gRPC Auth Service...")
	if err := grpc.StartGRPCServer(); err != nil {
		log.Fatalf("failed to start gRPC server: %v", err)
	}
}
