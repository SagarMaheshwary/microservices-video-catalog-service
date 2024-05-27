package config

import (
	"os"
	"path"
	"strconv"

	"github.com/gofor-little/env"
	"github.com/sagarmaheshwary/microservices-video-catalog-service/internal/helper"
	"github.com/sagarmaheshwary/microservices-video-catalog-service/internal/lib/log"
)

var conf *Config

type Config struct {
	GRPCServer *grpcServer
	AMQP       *amqp
	Database   *database
}

type database struct {
	Host     string
	Port     string
	Username string
	Password string
	Database string
	SSLMode  string
	Timezone string
}

type grpcServer struct {
	Host string
	Port int
}

type amqp struct {
	Host     string
	Port     int
	Username string
	Password string
}

func Init() {
	envPath := path.Join(helper.RootDir(), "..", ".env")

	if err := env.Load(envPath); err != nil {
		log.Fatal("Failed to load .env %q: %v", envPath, err)
	}

	log.Info("Loaded %q", envPath)

	port, err := strconv.Atoi(Getenv("GRPC_PORT", "5000"))

	if err != nil {
		log.Error("Invalid GRPC_PORT value %v", err)
	}

	amqpPort, err := strconv.Atoi(Getenv("AMQP_PORT", "5672"))

	if err != nil {
		log.Error("Invalid AMQP_PORT value %v", err)
	}

	conf = &Config{
		GRPCServer: &grpcServer{
			Host: Getenv("GRPC_HOST", "localhost"),
			Port: port,
		},
		AMQP: &amqp{
			Host:     Getenv("AMQP_HOST", "localhost"),
			Port:     amqpPort,
			Username: Getenv("AMQP_USERNAME", "guest"),
			Password: Getenv("AMQP_PASSWORD", "guest"),
		},
		Database: &database{
			Host:     Getenv("DB_HOST", "localhost"),
			Port:     Getenv("DB_PORT", "5432"),
			Username: Getenv("DB_USERNAME", "postgres"),
			Password: Getenv("DB_PASSWORD", "password"),
			Database: Getenv("DB_DATABASE", "microservices_video_catalog_service"),
			SSLMode:  Getenv("DB_SSLMODE", "disable"),
			Timezone: Getenv("DB_TIMEZONE", "UTC"),
		},
	}
}

func GetgrpcServer() *grpcServer {
	return conf.GRPCServer
}

func Getamqp() *amqp {
	return conf.AMQP
}

func GetDatabase() *database {
	return conf.Database
}

func Getenv(key string, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}

	return defaultVal
}
