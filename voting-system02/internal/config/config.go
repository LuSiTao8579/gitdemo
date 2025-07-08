package config

import (
	"log"
	"os"
	_ "strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerAddress string
	DataFilePath  string
	TokenSecret   string
	TokenDuration time.Duration
}

func LoadConfig() *Config {
	// 加载.env文件
	_ = godotenv.Load()

	// 设置默认值
	cfg := &Config{
		ServerAddress: ":8080",
		DataFilePath:  "data.json",
		TokenSecret:   "default-secret",
	}

	// 从环境变量读取配置
	if addr := os.Getenv("SERVER_ADDRESS"); addr != "" {
		cfg.ServerAddress = addr
	}

	if path := os.Getenv("DATA_FILE_PATH"); path != "" {
		cfg.DataFilePath = path
	}

	if secret := os.Getenv("TOKEN_SECRET"); secret != "" {
		cfg.TokenSecret = secret
	}

	if durationStr := os.Getenv("TOKEN_DURATION"); durationStr != "" {
		duration, err := time.ParseDuration(durationStr)
		if err != nil {
			log.Printf("Invalid TOKEN_DURATION: %v, using default", err)
			cfg.TokenDuration = 24 * time.Hour
		} else {
			cfg.TokenDuration = duration
		}
	}

	return cfg
}
