package app

import (
	"doctor-go/internal/config"
	"doctor-go/internal/infrastructure/logger"
	"doctor-go/internal/infrastructure/mysql"
	"doctor-go/internal/infrastructure/oss"
	"doctor-go/internal/infrastructure/redis"
	"doctor-go/internal/modules/admin"
)

type Server struct {
	cfg *config.Config
	dep *Dependencies
}

type Dependencies struct {
	MySQL  *mysql.DB
	Redis  *redis.Client
	OSS    *oss.Client
	Logger *logger.Logger
}

func NewServer(cfg *config.Config) (*Server, error) {
	log, err := logger.New(cfg.App.Env)
	if err != nil {
		return nil, err
	}

	db, err := mysql.New(cfg.MySQL)
	if err != nil {
		return nil, err
	}
	if err := admin.EnsureDefaultAdmin(cfg.Admin, admin.NewRepository(db)); err != nil {
		return nil, err
	}

	cache := redis.New(cfg.Redis)
	storage, err := oss.New(cfg.OSS)
	if err != nil {
		return nil, err
	}

	return &Server{
		cfg: cfg,
		dep: &Dependencies{
			MySQL:  db,
			Redis:  cache,
			OSS:    storage,
			Logger: log,
		},
	}, nil
}

func (s *Server) Run() error {
	router := NewRouter(s.cfg, s.dep)
	return router.Run(s.cfg.App.Addr)
}
