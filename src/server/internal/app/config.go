package app

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	JWT      JWTConfig
	CORS     CORSConfig
	Docs     DocsConfig
}

type DocsConfig struct {
	SpecPath string
}

type RedisConfig struct {
	URL string
}

type ServerConfig struct {
	Port int
	Mode string
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
	SSLMode  string
}

func (d DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		d.Host, d.Port, d.User, d.Password, d.Name, d.SSLMode,
	)
}

type JWTConfig struct {
	Secret          string
	ExpirationHours int
	Issuer          string
}

type CORSConfig struct {
	AllowedOrigins []string
}

// LoadDotenv looks for a .env file in common locations (cwd, parent dirs up
// to 4 levels, and next to the running binary) and loads it into the process
// environment. Existing env vars are not overwritten. Returns the path that
// was loaded, or "" if no file was found.
func LoadDotenv() string {
	candidates := []string{".env"}
	if exe, err := os.Executable(); err == nil {
		candidates = append(candidates, filepath.Join(filepath.Dir(exe), ".env"))
	}
	if cwd, err := os.Getwd(); err == nil {
		dir := cwd
		for i := 0; i < 4; i++ {
			candidates = append(candidates, filepath.Join(dir, ".env"))
			parent := filepath.Dir(dir)
			if parent == dir {
				break
			}
			dir = parent
		}
	}
	for _, p := range candidates {
		if _, err := os.Stat(p); err == nil {
			if err := godotenv.Load(p); err == nil {
				return p
			}
		}
	}
	return ""
}

func LoadConfig() (*Config, error) {
	LoadDotenv()

	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath("./configs")
	v.AddConfigPath("../../configs")

	v.SetEnvPrefix("")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	_ = v.ReadInConfig()

	bindEnv(v)

	cfg := &Config{
		Server: ServerConfig{
			Port: v.GetInt("server.port"),
			Mode: v.GetString("server.mode"),
		},
		Database: DatabaseConfig{
			Host:     v.GetString("database.host"),
			Port:     v.GetInt("database.port"),
			User:     v.GetString("database.user"),
			Password: v.GetString("database.password"),
			Name:     v.GetString("database.name"),
			SSLMode:  v.GetString("database.sslmode"),
		},
		JWT: JWTConfig{
			Secret:          v.GetString("jwt.secret"),
			ExpirationHours: v.GetInt("jwt.expirationHours"),
			Issuer:          v.GetString("jwt.issuer"),
		},
		CORS: CORSConfig{
			AllowedOrigins: v.GetStringSlice("cors.allowedOrigins"),
		},
		Docs: DocsConfig{
			SpecPath: v.GetString("docs.specPath"),
		},
		Redis: RedisConfig{
			URL: v.GetString("redis.url"),
		},
	}

	if cfg.Docs.SpecPath == "" {
		cfg.Docs.SpecPath = "./api/openapi.yaml"
	}

	if cfg.JWT.Secret == "" || cfg.JWT.Secret == "change-me-in-env" {
		return nil, fmt.Errorf("JWT_SECRET must be set in environment")
	}

	return cfg, nil
}

func bindEnv(v *viper.Viper) {
	_ = v.BindEnv("server.port", "PORT")
	_ = v.BindEnv("server.mode", "APP_MODE")

	_ = v.BindEnv("database.host", "DB_HOST")
	_ = v.BindEnv("database.port", "DB_PORT")
	_ = v.BindEnv("database.user", "DB_USER")
	_ = v.BindEnv("database.password", "DB_PASSWORD")
	_ = v.BindEnv("database.name", "DB_NAME")
	_ = v.BindEnv("database.sslmode", "DB_SSLMODE")

	_ = v.BindEnv("jwt.secret", "JWT_SECRET")
	_ = v.BindEnv("jwt.expirationHours", "JWT_EXPIRATION_HOURS")
	_ = v.BindEnv("jwt.issuer", "JWT_ISSUER")

	_ = v.BindEnv("docs.specPath", "DOCS_SPEC_PATH")

	_ = v.BindEnv("redis.url", "REDIS_URL")
}
