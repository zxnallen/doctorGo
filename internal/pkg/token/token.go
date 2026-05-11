package token

import (
	"context"
	"strconv"
	"time"

	"doctor-go/internal/config"
	"doctor-go/internal/infrastructure/redis"
	appJWT "doctor-go/internal/pkg/jwt"
)

type Pair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}

type Service struct {
	cfg   config.JWTConfig
	cache *redis.Client
}

func NewService(cfg config.JWTConfig, cache *redis.Client) *Service {
	return &Service{cfg: cfg, cache: cache}
}

func (s *Service) GeneratePair(ctx context.Context, userID uint64, role string) (*Pair, error) {
	access, err := appJWT.GenerateWithRoleAndType(s.cfg.Secret, userID, role, "access", s.cfg.ExpireSeconds)
	if err != nil {
		return nil, err
	}
	refresh, err := appJWT.GenerateWithRoleAndType(s.cfg.Secret, userID, role, "refresh", s.cfg.RefreshExpireSeconds)
	if err != nil {
		return nil, err
	}
	claims, err := appJWT.Parse(s.cfg.Secret, refresh)
	if err != nil {
		return nil, err
	}
	if s.cache != nil {
		key := refreshKey(role, userID, claims.ID)
		_ = s.cache.Set(ctx, key, "1", time.Duration(s.cfg.RefreshExpireSeconds)*time.Second).Err()
	}
	return &Pair{AccessToken: access, RefreshToken: refresh, ExpiresIn: s.cfg.ExpireSeconds}, nil
}

func (s *Service) Refresh(ctx context.Context, refreshToken string, expectedRole string) (*Pair, error) {
	claims, err := appJWT.Parse(s.cfg.Secret, refreshToken)
	if err != nil {
		return nil, err
	}
	if claims.Role != expectedRole || claims.TokenType != "refresh" {
		return nil, appJWT.ErrInvalidToken
	}
	if s.cache != nil {
		exists, err := s.cache.Exists(ctx, refreshKey(claims.Role, claims.UserID, claims.ID)).Result()
		if err != nil || exists == 0 {
			return nil, appJWT.ErrInvalidToken
		}
		_ = s.cache.Del(ctx, refreshKey(claims.Role, claims.UserID, claims.ID)).Err()
	}
	return s.GeneratePair(ctx, claims.UserID, claims.Role)
}

func (s *Service) Logout(ctx context.Context, accessToken string, refreshToken string) {
	if s.cache == nil {
		return
	}
	if accessToken != "" {
		if claims, err := appJWT.Parse(s.cfg.Secret, accessToken); err == nil && claims.ID != "" {
			ttl := time.Until(claims.ExpiresAt.Time)
			if ttl > 0 {
				_ = s.cache.Set(ctx, blacklistKey(claims.ID), "1", ttl).Err()
			}
		}
	}
	if refreshToken != "" {
		if claims, err := appJWT.Parse(s.cfg.Secret, refreshToken); err == nil && claims.ID != "" {
			_ = s.cache.Del(ctx, refreshKey(claims.Role, claims.UserID, claims.ID)).Err()
		}
	}
}

func (s *Service) IsBlacklisted(ctx context.Context, jti string) bool {
	if s.cache == nil || jti == "" {
		return false
	}
	exists, err := s.cache.Exists(ctx, blacklistKey(jti)).Result()
	return err == nil && exists > 0
}

func refreshKey(role string, userID uint64, jti string) string {
	return "auth:refresh:" + role + ":" + strconv.FormatUint(userID, 10) + ":" + jti
}

func blacklistKey(jti string) string {
	return "auth:blacklist:" + jti
}
