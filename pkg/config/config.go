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
	Active string        `yaml:"active"` // "mysql" or "mongodb"
	MySQL  MySQLConfig   `yaml:"mysql"`
	Mongo  MongoDBConfig `yaml:"mongodb"`
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

type AppConfigStruct struct {
	Server   ServerConfig     `yaml:"server"`
	Database DatabaseConfig   `yaml:"database"`
	Consul   ConsulConfig     `yaml:"consul"`
	Zap      ZapConfig        `mapstructure:"zap"`
	Registry Registry         `mapstructure:"registry" validate:"required"`
	App      AppConfiguration `mapstructure:"app"`
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
