package usecase

import (
	"context"
	"errors"
	"time"

	"auth-services/config"
	"auth-services/internal/dto"
	"auth-services/internal/entity"
	"auth-services/internal/repository"
	"auth-services/pkg/utils"

	"encoding/json"

	"golang.org/x/crypto/bcrypt"
)

type UserUsecase interface {
	Register(dto.RegisterRequest) error
	Login(dto.LoginRequest) (string, error)
	GetProfile(uint) (*dto.UserResponse, error)
	VerifyToken(context.Context, string) (string, error)
}

type userUsecase struct {
	repo repository.UserRepository
}

func NewUserUsecase(r repository.UserRepository) UserUsecase {
	return &userUsecase{r}
}

func (u *userUsecase) Register(req dto.RegisterRequest) error {
	hashed, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	user := &entity.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashed),
	}
	return u.repo.Create(user)
}

func (u *userUsecase) Login(req dto.LoginRequest) (string, error) {
	user, err := u.repo.GetByEmail(req.Email)
	if err != nil || bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)) != nil {
		return "", errors.New("invalid credentials")
	}
	return utils.GenerateJWT(user.ID)
}

func (u *userUsecase) GetProfile(id uint) (*dto.UserResponse, error) {
	ctx := context.Background()
	cacheKey := "user_profile:" + string(rune(id))

	val, err := config.Redis.Get(ctx, cacheKey).Result()
	if err == nil {
		var user dto.UserResponse
		json.Unmarshal([]byte(val), &user)
		return &user, nil
	}

	user, err := u.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	res := &dto.UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}

	b, _ := json.Marshal(res)
	config.Redis.Set(ctx, cacheKey, b, 10*time.Minute)

	return res, nil
}

func (u *userUsecase) VerifyToken(ctx context.Context, token string) (string, error) {
	claims, err := utils.ValidateJWT(token)
	if err != nil {
		return "", err
	}
	return claims.UserID, nil
}
