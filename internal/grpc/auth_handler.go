package grpc

import (
	"auth-services/internal/usecase"
	pb "auth-services/proto/auth"
	"context"
)

type AuthGrpcHandler struct {
	pb.UnimplementedAuthServiceServer
	userUsecase usecase.UserUsecase
}

func NewAuthGrpcHandler(userUsecase usecase.UserUsecase) *AuthGrpcHandler {
	return &AuthGrpcHandler{userUsecase: userUsecase}
}

func (h *AuthGrpcHandler) VerifyToken(ctx context.Context, req *pb.VerifyTokenRequest) (*pb.VerifyTokenResponse, error) {
	userID, err := h.userUsecase.VerifyToken(ctx, req.Token)
	if err != nil {
		return &pb.VerifyTokenResponse{
			Valid:  false,
			Error:  err.Error(),
			UserId: "",
		}, nil
	}
	return &pb.VerifyTokenResponse{
		Valid:  true,
		UserId: userID,
	}, nil
}
