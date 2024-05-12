package config

import (
	"os"
	"path"
	"strconv"

	"github.com/gofor-little/env"
	"github.com/sagarmaheshwary/microservices-upload-service/internal/helper"
	"github.com/sagarmaheshwary/microservices-upload-service/internal/lib/log"
)

var conf *Config

type Config struct {
	GRPCServer *grpcServer
}

type grpcServer struct {
	Host string
	Port int
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

	conf = &Config{
		GRPCServer: &grpcServer{
			Host: Getenv("GRPC_HOST", "localhost"),
			Port: port,
		},
	}
}

func GetgrpcServer() *grpcServer {
	return conf.GRPCServer
}

func Getenv(key string, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}

	return defaultVal
}
