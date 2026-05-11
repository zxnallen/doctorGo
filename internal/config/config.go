package config

import (
	"errors"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	App   AppConfig   `mapstructure:"app"`
	Admin AdminConfig `mapstructure:"admin"`
	MySQL MySQLConfig `mapstructure:"mysql"`
	Redis RedisConfig `mapstructure:"redis"`
	JWT   JWTConfig   `mapstructure:"jwt"`
	OSS   OSSConfig   `mapstructure:"oss"`
}

type AppConfig struct {
	Name           string   `mapstructure:"name"`
	Env            string   `mapstructure:"env"`
	Addr           string   `mapstructure:"addr"`
	TrustedProxies []string `mapstructure:"trusted_proxies"`
}

type AdminConfig struct {
	DefaultUsername string `mapstructure:"default_username"`
	DefaultPassword string `mapstructure:"default_password"`
	DefaultNickname string `mapstructure:"default_nickname"`
}

type MySQLConfig struct {
	DSN string `mapstructure:"dsn"`
}

type RedisConfig struct {
	Addr     string `mapstructure:"addr"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type JWTConfig struct {
	Secret               string `mapstructure:"secret"`
	ExpireSeconds        int    `mapstructure:"expire_seconds"`
	RefreshExpireSeconds int    `mapstructure:"refresh_expire_seconds"`
}

type OSSConfig struct {
	Endpoint        string `mapstructure:"endpoint"`
	AccessKeyID     string `mapstructure:"access_key_id"`
	AccessKeySecret string `mapstructure:"access_key_secret"`
	Bucket          string `mapstructure:"bucket"`
	BaseURL         string `mapstructure:"base_url"`
}

func Load() (*Config, error) {
	env := getEnv("APP_ENV", "local")
	cfg := defaultConfig(env)

	v := viper.New()
	v.SetConfigName("config." + env)
	v.SetConfigType("yaml")
	v.AddConfigPath("./configs")
	v.AddConfigPath("../configs")
	v.AddConfigPath("../../configs")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()
	bindEnvs(v)

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	applyEnvOverrides(&cfg)
	if err := validate(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func getEnv(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func getEnvInt(key string, fallback int) int {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return parsed
}

func defaultConfig(env string) Config {
	return Config{
		App: AppConfig{
			Name:           "doctor-go",
			Env:            env,
			Addr:           ":8080",
			TrustedProxies: []string{"127.0.0.1", "::1", "172.16.0.0/12", "10.0.0.0/8"},
		},
		Admin: AdminConfig{
			DefaultUsername: "admin",
			DefaultPassword: "admin123456",
			DefaultNickname: "系统管理员",
		},
		MySQL: MySQLConfig{
			DSN: "root:root@tcp(127.0.0.1:3306)/doctor_go?charset=utf8mb4&parseTime=True&loc=Local",
		},
		Redis: RedisConfig{
			Addr:     "127.0.0.1:6379",
			Password: "",
			DB:       0,
		},
		JWT: JWTConfig{
			Secret:               "change-me-in-production",
			ExpireSeconds:        86400,
			RefreshExpireSeconds: 604800,
		},
	}
}

func bindEnvs(v *viper.Viper) {
	keys := []string{
		"app.name",
		"app.env",
		"app.addr",
		"app.trusted_proxies",
		"admin.default_username",
		"admin.default_password",
		"admin.default_nickname",
		"mysql.dsn",
		"redis.addr",
		"redis.password",
		"redis.db",
		"jwt.secret",
		"jwt.expire_seconds",
		"jwt.refresh_expire_seconds",
		"oss.endpoint",
		"oss.access_key_id",
		"oss.access_key_secret",
		"oss.bucket",
		"oss.base_url",
	}
	for _, key := range keys {
		_ = v.BindEnv(key)
	}
}

func applyEnvOverrides(cfg *Config) {
	cfg.App.Env = getEnv("APP_ENV", cfg.App.Env)
	cfg.App.Name = getEnv("APP_NAME", cfg.App.Name)
	cfg.App.Addr = getEnv("APP_ADDR", cfg.App.Addr)
	if value := os.Getenv("APP_TRUSTED_PROXIES"); value != "" {
		cfg.App.TrustedProxies = splitList(value)
	}
	cfg.Admin.DefaultUsername = getEnv("ADMIN_DEFAULT_USERNAME", cfg.Admin.DefaultUsername)
	cfg.Admin.DefaultPassword = getEnv("ADMIN_DEFAULT_PASSWORD", cfg.Admin.DefaultPassword)
	cfg.Admin.DefaultNickname = getEnv("ADMIN_DEFAULT_NICKNAME", cfg.Admin.DefaultNickname)
	cfg.MySQL.DSN = getEnv("MYSQL_DSN", cfg.MySQL.DSN)
	cfg.Redis.Addr = getEnv("REDIS_ADDR", cfg.Redis.Addr)
	cfg.Redis.Password = getEnv("REDIS_PASSWORD", cfg.Redis.Password)
	cfg.Redis.DB = getEnvInt("REDIS_DB", cfg.Redis.DB)
	cfg.JWT.Secret = getEnv("JWT_SECRET", cfg.JWT.Secret)
	cfg.JWT.ExpireSeconds = getEnvInt("JWT_EXPIRE_SECONDS", cfg.JWT.ExpireSeconds)
	cfg.JWT.RefreshExpireSeconds = getEnvInt("JWT_REFRESH_EXPIRE_SECONDS", cfg.JWT.RefreshExpireSeconds)
	cfg.OSS.Endpoint = getEnv("OSS_ENDPOINT", cfg.OSS.Endpoint)
	cfg.OSS.AccessKeyID = getEnv("OSS_ACCESS_KEY_ID", cfg.OSS.AccessKeyID)
	cfg.OSS.AccessKeySecret = getEnv("OSS_ACCESS_KEY_SECRET", cfg.OSS.AccessKeySecret)
	cfg.OSS.Bucket = getEnv("OSS_BUCKET", cfg.OSS.Bucket)
	cfg.OSS.BaseURL = getEnv("OSS_BASE_URL", cfg.OSS.BaseURL)
}

func splitList(value string) []string {
	parts := strings.Split(value, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

func validate(cfg *Config) error {
	if cfg.App.Env != "prod" {
		return nil
	}
	if cfg.MySQL.DSN == "" {
		return errors.New("MYSQL_DSN is required in prod")
	}
	if cfg.Redis.Addr == "" {
		return errors.New("REDIS_ADDR is required in prod")
	}
	if cfg.JWT.Secret == "" || strings.HasPrefix(cfg.JWT.Secret, "change-me") {
		return errors.New("JWT_SECRET is required in prod")
	}
	return nil
}
