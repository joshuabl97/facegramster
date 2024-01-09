package models

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/rs/zerolog"

	"gopkg.in/yaml.v3"
)

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
	SSLMode  string
}

type Service struct {
	Image       string            `yaml:"image"`
	Restart     string            `yaml:"restart"`
	Environment map[string]string `yaml:"environment"`
	Ports       []string          `yaml:"ports"`
}

type Services struct {
	Version  string             `yaml:"version"`
	Services map[string]Service `yaml:"services"`
}

func DefaultPostgresConfig() PostgresConfig {
	yamlFile, err := os.ReadFile("docker-compose.yaml")
	if err != nil {
		log.Fatalf("Error reading YAML file: %v", err)
	}

	var services Services
	err = yaml.Unmarshal(yamlFile, &services)
	if err != nil {
		log.Fatalf("Error unmarshalling YAML: %v", err)
	}

	cfg := PostgresConfig{
		Host:    "localhost",
		SSLMode: "disable",
	}

	for serviceName, service := range services.Services {
		if serviceName == "db" {

			for key, value := range service.Environment {
				switch key {
				case "POSTGRES_PASSWORD":
					cfg.Password = value
				case "POSTGRES_DB":
					cfg.Database = value
				case "POSTGRES_USER":
					cfg.User = value
				}
			}

			for _, port := range service.Ports {
				ports := strings.Split(port, ":")
				cfg.Port = ports[0]
			}
		}
	}

	return cfg
}

func (cfg PostgresConfig) ConnString() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Database, cfg.SSLMode)
}

// Open will open a SQL connection
// Callers of Open need to make sure that the connection eventually closes
// (run `defer db.Close()`)
func Open(lg *zerolog.Logger, cfg PostgresConfig) (*sql.DB, error) {
	db, err := sql.Open("pgx", cfg.ConnString())
	if err != nil {
		lg.Error().Err(err).Msg("failed to open db")
		return nil, fmt.Errorf("failed to open db: %w", err)
	}

	return db, nil
}
