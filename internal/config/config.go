package config

import (
	"flag"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	DBUrl        string `yaml:"db_url"`
	GrinexAPIUrl string `yaml:"grinex_api_url"`
	Port         string `yaml:"port"`
}

func Load() *Config {
	cfg := &Config{}
	file, err := os.Open("config.yml")
	if err == nil {
		defer func() {
			if err := file.Close(); err != nil {
				log.Printf("Error closing file: %v", err)
			}
		}()
		decoder := yaml.NewDecoder(file)
		if err := decoder.Decode(cfg); err != nil {
			log.Fatalf("failed to decode config.yml: %v", err)
		}
	}

	if env := os.Getenv("DB_URL"); env != "" {
		cfg.DBUrl = env
	}
	if env := os.Getenv("API_URL"); env != "" {
		cfg.GrinexAPIUrl = env
	}
	if env := os.Getenv("PORT"); env != "" {
		cfg.Port = env
	}

	dbUrl := flag.String("db-url", cfg.DBUrl, "PostgreSQL connection URL")
	apiUrl := flag.String("api-url", cfg.GrinexAPIUrl, "Grinex API base URL")
	port := flag.String("port", cfg.Port, "Port for gRPC server")

	flag.Parse()

	cfg.DBUrl = *dbUrl
	cfg.GrinexAPIUrl = *apiUrl
	cfg.Port = *port

	return cfg
}
