package config

import (
	"os"
	"time"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

var (
	RedisDB                             = 0
	ProjectName           string        = "revolution-fiesta"
	AccessTokenExpiration time.Duration = time.Hour / 4
	YamlConfigPath        string        = "../config.yaml"
	AppConfig             Config
)

type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	Wechat   WechatConfig   `yaml:"wechat"`
	Esign    EsignConfig    `yaml:"esign"`
}

type ServerConfig struct {
	Port string `yaml:"port"`
	Env  string `yaml:"env"`
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
	AppID  string `yaml:"app_id"`
	MchID  string `yaml:"mch_id"`
	APIKey string `yaml:"api_key"`
}

type EsignConfig struct {
	APIURL string `yaml:"api_url"`
	AppID  string `yaml:"app_id"`
	Secret string `yaml:"secret"`
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
