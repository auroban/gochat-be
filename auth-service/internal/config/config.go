package config

import (
	"encoding/json"
	"log"

	"github.com/spf13/viper"
)

type serverConfig struct {
	Port int `json:"port"`
}

type loggerConfig struct {
	Level string `json:"level"`
}

type dbConfig struct {
	Host     string `json:"Host"`
	Port     int    `json:"Port"`
	User     string `json:"User"`
	Password string `json:"Password"`
	DBName   string `json:"DBName"`
	Schema   string `json:"Schema"`
	SSLMode  string `json:"SSLMode"`
}

type Config struct {
	Server   serverConfig `json:"Server"`
	Database dbConfig     `json:"Database"`
	Log      loggerConfig `json:"Log"`
}

func Load() *Config {
	viper.AutomaticEnv()
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./configs")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("❌ Failed to read config file: %v", err)
	}

	cfg := &Config{}
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatalf("❌ Unable to parse config: %v", err)
	}

	cfg.Database.Host = viper.GetString("DB_HOST")
	cfg.Database.Port = viper.GetInt("DB_PORT")
	cfg.Database.User = viper.GetString("DB_USER")
	cfg.Database.Password = viper.GetString("DB_PASSWORD")
	cfg.Database.DBName = viper.GetString("DB_NAME")
	cfg.Database.Schema = viper.GetString("DB_SCHEMA")
	cfg.Database.SSLMode = viper.GetString("DB_SSLMODE")

	validateEnvVariables(*cfg)
	return cfg
}

func validateEnvVariables(cfg Config) {
	missing := []string{}

	if cfg.Database.Host == "" {
		missing = append(missing, "DB_HOST")
	}
	if cfg.Database.Port == 0 {
		missing = append(missing, "DB_PORT")
	}
	if cfg.Database.User == "" {
		missing = append(missing, "DB_USER")
	}
	if cfg.Database.Password == "" {
		missing = append(missing, "DB_PASSWORD")
	}
	if cfg.Database.DBName == "" {
		missing = append(missing, "DB_NAME")
	}
	if cfg.Database.Schema == "" {
		missing = append(missing, "DB_SCHEMA")
	}
	if cfg.Database.SSLMode == "" {
		missing = append(missing, "DB_SSLMODE")
	}

	if len(missing) > 0 {
		log.Fatalf("❌ Missing required environment variables: %v", missing)
	}

}

func (c Config) String() string {
	b, err := json.MarshalIndent(c, "", " ")
	if err != nil {
		log.Printf("Failed to marshal config: %v", err)
		return ""
	}
	return string(b)
}
