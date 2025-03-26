package config

import (
	"crypto/rsa"
	"fmt"
	"os"
	"time"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

var (
	RedisDB                             = 0
	ProjectName           string        = "revolution-fiesta"
	AccessTokenExpiration time.Duration = time.Hour / 4
	YamlConfigPath        string        = "./config.yaml"
	AppConfig             Config
	PrivateKey            *rsa.PrivateKey
)

type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	Wechat   WechatConfig   `yaml:"wechat"`
}

type ServerConfig struct {
	Port        string `yaml:"port"`
	Env         string `yaml:"env"`
	CrtFilePath string `yaml:"crt"`
	KeyFilePath string `yaml:"key"`
}

type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
	SSLMode  string `yaml:"sslmode"`
}

type WechatConfig struct {
	AppID     string `yaml:"app_id"`
	AppSecret string `yaml:"app_secret"`
}

func LoadConfig() error {
	file, err := os.Open(YamlConfigPath)
	if err != nil {
		return errors.Wrapf(err, "error loading config")
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&AppConfig)
	if err != nil {
		return errors.Wrapf(err, "Error decoding config")
	}

	return nil
}

func GetPostgresDsn() string {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		AppConfig.Database.Host,
		AppConfig.Database.User,
		AppConfig.Database.Password,
		AppConfig.Database.Name,
		AppConfig.Database.Port,
		AppConfig.Database.SSLMode)
	fmt.Println(dsn)
	return dsn
}
