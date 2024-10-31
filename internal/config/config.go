package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

// Config - структура для хранения конфигурации
type Config struct {
	Server struct {
		Port string `yaml:"port"`
	} `yaml:"server"`
	Database struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		DBName   string `yaml:"dbname"`
		Sslmode  string `yaml:"sslmode"`
	} `yaml:"database"`
	Logger struct {
		Output   string `yaml:"output"`
		FilePath string `yaml:"filepath"`
		Level    string `yaml:"level"`
	} `yaml:"logger"`
	Jwt struct {
		SecretKey string `yaml:"secret_key"`
	} `yaml:"jwt"`
	SMTP struct {
		Auth_address   string `yaml:"smtp_auth_address"`
		Server_address string `yaml:"smtp_server_address"`
		Name           string `yaml:"smtp_name"`
		Email_address  string `yaml:"smtp_email_address"`
		Email_password string `yaml:"smtp_email_password"`
	} `yaml:"smtp_mail"`
	Token struct {
		AccessTokenTimeDuration  int `yaml:"access_token_duration"`
		RefreshTokenTimeDuration int `yaml:"refresh_token_duration"`
	}
	EnvType string `yaml:"env_type"`
}

// LoadConfig - функция для загрузки конфигурации из YAML файла
func LoadConfig(filename string) *Config {
	data, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalf("cannot read config file: %v", err)
	}
	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		log.Fatalf("cannot unmarshal config data: %v", err)
	}
	return &cfg
}
