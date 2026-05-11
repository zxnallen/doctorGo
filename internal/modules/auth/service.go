package auth

import (
	"context"
	"errors"
	"time"

	"doctor-go/internal/config"
	"doctor-go/internal/modules/user"
	"doctor-go/internal/pkg/password"
	"doctor-go/internal/pkg/token"
)

var (
	ErrUserExists     = errors.New("user already exists")
	ErrBadCredentials = errors.New("bad credentials")
)

type Service struct {
	cfg          config.JWTConfig
	repo         *Repository
	tokenService *token.Service
}

func NewService(cfg config.JWTConfig, repo *Repository, tokenService *token.Service) *Service {
	return &Service{cfg: cfg, repo: repo, tokenService: tokenService}
}

func (s *Service) Register(ctx context.Context, req RegisterRequest) (*LoginResponse, error) {
	exists, err := s.repo.FindByUsername(req.Username)
	if err != nil {
		return nil, err
	}
	if exists != nil {
		return nil, ErrUserExists
	}

	hash, err := password.Hash(req.Password)
	if err != nil {
		return nil, err
	}

	u := &user.User{
		Username:     req.Username,
		Phone:        req.Phone,
		Email:        req.Email,
		PasswordHash: hash,
		Status:       1,
	}
	if err := s.repo.Create(u); err != nil {
		return nil, err
	}

	pair, err := s.tokenService.GeneratePair(ctx, u.ID, "user")
	if err != nil {
		return nil, err
	}
	return loginResponse(pair, u), nil
}

func (s *Service) Login(ctx context.Context, req LoginRequest) (*LoginResponse, error) {
	u, err := s.repo.FindByUsername(req.Username)
	if err != nil {
		return nil, err
	}
	if u == nil || !password.Check(u.PasswordHash, req.Password) {
		return nil, ErrBadCredentials
	}

	now := time.Now()
	_ = s.repo.UpdateLastLogin(u.ID, now)
	u.LastLoginAt = &now

	pair, err := s.tokenService.GeneratePair(ctx, u.ID, "user")
	if err != nil {
		return nil, err
	}
	return loginResponse(pair, u), nil
}

func (s *Service) Refresh(ctx context.Context, refreshToken string) (*token.Pair, error) {
	return s.tokenService.Refresh(ctx, refreshToken, "user")
}

func (s *Service) Logout(ctx context.Context, accessToken string, refreshToken string) {
	s.tokenService.Logout(ctx, accessToken, refreshToken)
}

func loginResponse(pair *token.Pair, u *user.User) *LoginResponse {
	return &LoginResponse{
		AccessToken:  pair.AccessToken,
		RefreshToken: pair.RefreshToken,
		ExpiresIn:    pair.ExpiresIn,
		User: UserSummary{
			ID:        u.ID,
			Username:  u.Username,
			AvatarURL: u.AvatarURL,
		},
	}
}
