package admin

import (
	"context"
	"errors"
	"time"

	"doctor-go/internal/config"
	"doctor-go/internal/pkg/password"
	"doctor-go/internal/pkg/token"
)

var ErrBadCredentials = errors.New("bad admin credentials")
var ErrAdminDisabled = errors.New("admin disabled")

type Service struct {
	cfg          config.JWTConfig
	repo         *Repository
	tokenService *token.Service
}

func NewService(cfg config.JWTConfig, repo *Repository, tokenService *token.Service) *Service {
	return &Service{cfg: cfg, repo: repo, tokenService: tokenService}
}

func EnsureDefaultAdmin(cfg config.AdminConfig, repo *Repository) error {
	if cfg.DefaultUsername == "" || cfg.DefaultPassword == "" {
		return nil
	}

	exists, err := repo.FindByUsername(cfg.DefaultUsername)
	if err != nil {
		return err
	}
	if exists != nil {
		return repo.EnsureRBAC(exists.ID)
	}

	hash, err := password.Hash(cfg.DefaultPassword)
	if err != nil {
		return err
	}
	admin := &AdminUser{
		Username:     cfg.DefaultUsername,
		PasswordHash: hash,
		Nickname:     cfg.DefaultNickname,
		Status:       1,
	}
	if err := repo.Create(admin); err != nil {
		return err
	}
	return repo.EnsureRBAC(admin.ID)
}

func (s *Service) Login(ctx context.Context, req LoginRequest) (*LoginResponse, error) {
	admin, err := s.repo.FindByUsername(req.Username)
	if err != nil {
		return nil, err
	}
	if admin == nil || !password.Check(admin.PasswordHash, req.Password) {
		return nil, ErrBadCredentials
	}
	if admin.Status != 1 {
		return nil, ErrAdminDisabled
	}

	now := time.Now()
	_ = s.repo.UpdateLastLogin(admin.ID, now)
	admin.LastLoginAt = &now

	pair, err := s.tokenService.GeneratePair(ctx, admin.ID, "admin")
	if err != nil {
		return nil, err
	}
	return &LoginResponse{
		AccessToken:  pair.AccessToken,
		RefreshToken: pair.RefreshToken,
		ExpiresIn:    pair.ExpiresIn,
		Admin: AdminSummary{
			ID:        admin.ID,
			Username:  admin.Username,
			Nickname:  admin.Nickname,
			AvatarURL: admin.AvatarURL,
		},
	}, nil
}

func (s *Service) Refresh(ctx context.Context, refreshToken string) (*token.Pair, error) {
	return s.tokenService.Refresh(ctx, refreshToken, "admin")
}

func (s *Service) Logout(ctx context.Context, accessToken string, refreshToken string) {
	s.tokenService.Logout(ctx, accessToken, refreshToken)
}
