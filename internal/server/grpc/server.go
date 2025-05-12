package grpc

import (
	"log"
	"net"
	"os"

	"auth-services/config"
	authgrpc "auth-services/internal/grpc"
	"auth-services/internal/repository"
	"auth-services/internal/usecase"
	authpb "auth-services/proto/auth"

	"google.golang.org/grpc"
)

func StartGRPCServer() error {
	lis, err := net.Listen("tcp", ":"+os.Getenv("PORT"))
	if err != nil {
		return err
	}

	db := config.InitDB()
	userRepo := repository.NewUserRepository(db)
	userUsecase := usecase.NewUserUsecase(userRepo)
	handler := authgrpc.NewAuthGrpcHandler(userUsecase)

	s := grpc.NewServer()
	authpb.RegisterAuthServiceServer(s, handler)

	log.Println("[gRPC] AuthService listening on :" + os.Getenv("PORT"))
	return s.Serve(lis)
}
