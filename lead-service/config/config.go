package config

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
)

type ServerConfig struct {
	Port int `yaml:"port"`
}

type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
	SSLMode  string `yaml:"sslmode"`
	Action   string `yaml:"action"`
}

type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	Grpc     GrpcConfig     `yaml:"grpc"`
}

type GrpcConfig struct {
	EducationService struct {
		Address string `yaml:"address"`
	} `yaml:"education_service"`
	UserService struct {
		Address string `yaml:"address"`
	} `yaml:"user_service"`
}

func LoadConfig() (*Config, error) {
	file, err := os.Open("config/config.yaml")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := yaml.Unmarshal(bytes, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
