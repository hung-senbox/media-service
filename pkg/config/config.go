package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type ServerConfig struct {
	Port string `yaml:"port"`
}

type DatabaseConfig struct {
	Active     string           `yaml:"active"` // "mysql" or "mongodb"
	MySQL      MySQLConfig      `yaml:"mysql"`
	Mongo      MongoDBConfig    `yaml:"mongodb"`
	Redis      RedisConfig      `yaml:"redis"`
	RedisCache RedisCacheConfig `yaml:"redis_cache"`
}

type MySQLConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
}

type MongoDBConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
}

type RedisCacheConfig struct {
	TTLSeconds int    `yaml:"ttlSeconds"`
	Host       string `yaml:"host"`
	Port       string `yaml:"port"`
	Password   string `yaml:"password"`
	DB         int    `yaml:"db"`
}

type ConsulConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type ZapConfig struct {
	Development bool   `mapstructure:"development"`
	Caller      bool   `mapstructure:"caller"`
	Stacktrace  string `mapstructure:"stacktrace"`
	Cores       struct {
		Console struct {
			Type     string `mapstructure:"type"`
			Level    string `mapstructure:"level"`
			Encoding string `mapstructure:"encoding"`
		} `mapstructure:"console"`
	} `mapstructure:"cores"`
}

type AppConfiguration struct {
	Name        string    `mapstructure:"name"`
	Version     string    `mapstructure:"version"`
	Environment string    `mapstructure:"environment"`
	API         APIConfig `mapstructure:"api"`
}

type APIConfig struct {
	Rest RestConfig `mapstructure:"rest"`
}

type RestConfig struct {
	Host    string        `mapstructure:"host"`
	Port    string        `mapstructure:"port"`
	Setting SettingConfig `mapstructure:"setting"`
}
type SettingConfig struct {
	Debug               bool     `mapstructure:"debug"`
	DebugErrorsResponse bool     `mapstructure:"debugErrorsResponse"`
	IgnoreLogUrls       []string `mapstructure:"ignoreLogUrls"`
}

type Registry struct {
	Host string `mapstructure:"host" validate:"required"`
}

type RedisConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

// ---------------- S3 configuration ----------------
type SenboxFormSubmitBucket struct {
	Domain               string `yaml:"domain"`
	Region               string `yaml:"region"`
	BucketName           string `yaml:"bucket_name"`
	AccessKey            string `yaml:"access_key"`
	SecretKey            string `yaml:"secret_key"`
	CloudfrontKeyGroupID string `yaml:"cloudfront_key_group_id"`
	CloudfrontKeyPath    string `yaml:"cloudfront_key_path"`
}

type S3 struct {
	SenboxFormSubmitBucket SenboxFormSubmitBucket `yaml:"senbox-form-submit-bucket"`
}

// ---------------- S3 configuration ----------------

type AppConfigStruct struct {
	Server   ServerConfig     `yaml:"server"`
	Database DatabaseConfig   `yaml:"database"`
	Consul   ConsulConfig     `yaml:"consul"`
	Zap      ZapConfig        `mapstructure:"zap"`
	Registry Registry         `mapstructure:"registry" validate:"required"`
	App      AppConfiguration `mapstructure:"app"`
	S3       S3               `yaml:"s3"`
}

var AppConfig *AppConfigStruct

func LoadConfig(filePath string) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Failed to read config file: %v", err)
	}

	AppConfig = &AppConfigStruct{}
	err = yaml.Unmarshal(data, AppConfig)
	if err != nil {
		log.Fatalf("Failed to unmarshal config: %v", err)
	}

	log.Println("Config loaded successfully")
}
