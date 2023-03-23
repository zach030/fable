package cmd

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type FableConfig struct {
	OpenAI    OpenAI    `yaml:"openai"`
	AliyunOSS AliyunOSS `yaml:"aliyunOSS"`
	Milvus    Milvus    `yaml:"milvus"`
	Qdrant    Qdrant    `yaml:"qdrant"`
}

type OpenAI struct {
	APIKey   string `yaml:"apiKey"`
	ProxyURL string `yaml:"proxyURL"`
}

type AliyunOSS struct {
	AliyunOssEndpoint  string `yaml:"aliyunOssEndpoint"`
	AliyunOssAccessKey string `yaml:"aliyunOssAccessKey"`
	AliyunOssSecretKey string `yaml:"aliyunOssSecretKey"`
	Bucket             string `yaml:"bucket"`
}

type Milvus struct {
	Addr       string `yaml:"addr"`
	User       string `yaml:"user"`
	Password   string `yaml:"password"`
	Collection string `yaml:"collection"`
}

type Qdrant struct {
	Addr   string `yaml:"host"`
	ApiKey string `yaml:"apiKey"`
}

func InitCfg() *FableConfig {
	file, err := os.ReadFile("../config/config.yaml")
	if err != nil {
		fmt.Println("read config/config.yaml err:", err)
		return nil
	}
	cfg := &FableConfig{}
	err = yaml.Unmarshal(file, cfg)
	if err != nil {
		fmt.Println("conf Unmarshal err:", err)
		return nil
	}
	return cfg
}
