package config

import (
	"fmt"
	"os"

	"github.com/goccy/go-yaml"
	"github.com/joho/godotenv"
)

type Service struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type Database struct {
	Host         string `yaml:"host"`
	Port         string `yaml:"port"`
	Name         string `yaml:"name"`
	SSLMode      string `yaml:"sslmode"`
	MaxOpenConns int    `yaml:"max_open_connections"`
	MaxIdleConns int    `yaml:"max_idle_connections"`

	User     string `yaml:"-"`
	Password string `yaml:"-"`
}

type Config struct {
	Service  Service  `yaml:"service"`
	Database Database `yaml:"database"`
}

func (d *Database) GetDSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		d.Host, d.Port, d.User, d.Password, d.Name, d.SSLMode,
	)
}

func LoadConfig(configPath, envPath string) (*Config, error) {
	if err := godotenv.Load(envPath); err != nil {
		return nil, fmt.Errorf("failed to load env file %s: %s", envPath, err.Error())
	}

	file, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file %s: %s", configPath, err.Error())
	}

	var config Config

	err = yaml.Unmarshal(file, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to parse YAML config %s: %s", configPath, err.Error())
	}

	config.Database.User = os.Getenv("DB_USER")
	config.Database.Password = os.Getenv("DB_PASSWORD")

	return &config, nil
}
