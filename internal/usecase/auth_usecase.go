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

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserUsecase interface {
	Register(dto.RegisterRequest) error
	Login(dto.LoginRequest) (*dto.LoginResponse, error)
	GetProfile(string) (*dto.UserResponse, error)
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
		ID:       uuid.NewString(),
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashed),
		Role:     req.Role,
	}
	return u.repo.Create(user)
}

func (u *userUsecase) Login(req dto.LoginRequest) (*dto.LoginResponse, error) {
	user, err := u.repo.GetByEmail(req.Email)
	if err != nil || bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)) != nil {
		return nil, errors.New("invalid credentials")
	}
	// return utils.GenerateJWT(user.ID)
	token, err := utils.GenerateJWT(user.ID)
	if err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		ID:    user.ID,
		Email: user.Email,
		Token: token,
	}, nil
}

func (u *userUsecase) GetProfile(id string) (*dto.UserResponse, error) {
	ctx := context.Background()
	cacheKey := "user_profile:" + id

	val, err := config.Redis.Get(ctx, cacheKey).Result()
	if err == nil {
		var user dto.UserResponse
		if err := json.Unmarshal([]byte(val), &user); err == nil {
			return &user, nil
		}
	}

	user, err := u.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, entity.ErrUserNotFound) {
			return nil, err
		}
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
	id, ok := claims["id"].(string)
	if !ok {
		return "", errors.New("invalid token id")
	}
	return id, nil
}
