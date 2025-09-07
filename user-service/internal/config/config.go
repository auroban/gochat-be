package config

import (
	"encoding/json"
	"log"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type serverConfig struct {
	Port int `json:"port"`
}

type Config struct {
	Server serverConfig `json:"server"`
}

func Load() *Config {
	_ = godotenv.Load()

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("❌ Failed to read config file: %v", err)
	}

	cfg := &Config{}
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatalf("❌ Unable to parse config: %v", err)
	}

	log.Printf("Loaded config: [%+v]\n", cfg)
	return cfg
}

func (c Config) String() string {
	b, err := json.MarshalIndent(c, "", " ")
	if err != nil {
		log.Printf("Failed to marshal config: %v", err)
		return ""
	}
	return string(b)
}
