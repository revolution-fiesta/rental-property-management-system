package config

import (
	"fmt"
	"time"
	//"gopkg.in/yaml.v2"
	"os"

	"gopkg.in/yaml.v3"
)

var (
	RedisDB                             = 0
	ProjectName           string        = "revolution-fiesta"
	AccessTokenExpiration time.Duration = time.Hour / 4
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

var AppConfig Config

func LoadConfig() {
	file, err := os.Open("config/config.yaml")
	if err != nil {
		panic(fmt.Sprintf("Error loading config: %v", err))
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&AppConfig)
	if err != nil {
		panic(fmt.Sprintf("Error decoding config: %v", err))
	}
}
