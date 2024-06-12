package config

import (
	"os"
	"path"
	"strconv"
	"time"

	"github.com/gofor-little/env"
	"github.com/sagarmaheshwary/microservices-video-catalog-service/internal/helper"
	"github.com/sagarmaheshwary/microservices-video-catalog-service/internal/lib/log"
)

var conf *Config

type Config struct {
	GRPCServer *grpcServer
	AMQP       *amqp
	Database   *database
	S3         *s3
	GRPCClient *grpcClient
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

type s3 struct {
	Bucket             string
	Region             string
	AccessKey          string
	SecretKey          string
	PresignedUrlExpiry int
}

type grpcClient struct {
	UserServiceurl string
	Timeout        time.Duration
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

	s3UrlExpiry, err := strconv.Atoi(Getenv("AWS_S3_PRESIGNED_URL_EXPIRY", "0"))

	if err != nil {
		log.Error("Invalid AWS_S3_PRESIGNED_URL_EXPIRY value %v", err)
	}

	timeout, err := strconv.Atoi(Getenv("GRPC_CLIENT_TIMEOUT_SECONDS", "5"))

	if err != nil {
		log.Error("Invalid GRPC_CLIENT_TIMEOUT_SECONDS value %v", err)
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
		S3: &s3{
			Bucket:             Getenv("AWS_S3_BUCKET", ""),
			Region:             Getenv("AWS_S3_REGION", ""),
			AccessKey:          Getenv("AWS_S3_ACCESS_KEY", ""),
			SecretKey:          Getenv("AWS_S3_SECRET_KEY", ""),
			PresignedUrlExpiry: s3UrlExpiry,
		},
		GRPCClient: &grpcClient{
			UserServiceurl: Getenv("GRPC_USER_SERVICE_URL", "localhost:5000"),
			Timeout:        time.Duration(timeout) * time.Second,
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

func GetS3() *s3 {
	return conf.S3
}

func GetgrpcClient() *grpcClient {
	return conf.GRPCClient
}

func Getenv(key string, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}

	return defaultVal
}
