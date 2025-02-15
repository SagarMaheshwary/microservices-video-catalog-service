package config

import (
	"os"
	"path"
	"strconv"
	"time"

	"github.com/gofor-little/env"
	"github.com/sagarmaheshwary/microservices-video-catalog-service/internal/helper"
	"github.com/sagarmaheshwary/microservices-video-catalog-service/internal/lib/logger"
)

var Conf *Config

type Config struct {
	GRPCServer *GRPCServer
	AMQP       *AMQP
	Database   *Database
	AWS        *AWS
	GRPCClient *GRPCClient
	Prometheus *Prometheus
}

type Database struct {
	Host     string
	Port     string
	Username string
	Password string
	Database string
	SSLMode  string
	Timezone string
}

type GRPCServer struct {
	Host string
	Port int
}

type AMQP struct {
	Host     string
	Port     int
	Username string
	Password string
}

type AWS struct {
	Region               string
	S3Bucket             string
	AccessKey            string
	SecretKey            string
	S3PresignedURLExpiry int
	CloudFrontURL        string
}

type GRPCClient struct {
	UserServiceURL string
	Timeout        time.Duration
}

type Prometheus struct {
	MetricsHost string
	MetricsPort int
}

func Init() {
	envPath := path.Join(helper.GetRootDir(), "..", ".env")

	if _, err := os.Stat(envPath); err == nil {
		if err := env.Load(envPath); err != nil {
			logger.Fatal("Failed to load .env %q: %v", envPath, err)
		}

		logger.Info("Loaded environment variables from %q", envPath)
	} else {
		logger.Info(".env file not found, using system environment variables")
	}

	Conf = &Config{
		GRPCServer: &GRPCServer{
			Host: getEnv("GRPC_HOST", "localhost"),
			Port: getEnvInt("GRPC_PORT", 5003),
		},
		AMQP: &AMQP{
			Host:     getEnv("AMQP_HOST", "localhost"),
			Port:     getEnvInt("AMQP_PORT", 5672),
			Username: getEnv("AMQP_USERNAME", "guest"),
			Password: getEnv("AMQP_PASSWORD", "guest"),
		},
		Database: &Database{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			Username: getEnv("DB_USERNAME", "postgres"),
			Password: getEnv("DB_PASSWORD", "password"),
			Database: getEnv("DB_DATABASE", "microservices_video_catalog_service"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
			Timezone: getEnv("DB_TIMEZONE", "UTC"),
		},
		AWS: &AWS{
			Region:               getEnv("AWS_REGION", ""),
			AccessKey:            getEnv("AWS_ACCESS_KEY", ""),
			SecretKey:            getEnv("AWS_SECRET_KEY", ""),
			S3Bucket:             getEnv("AWS_S3_BUCKET", ""),
			S3PresignedURLExpiry: getEnvInt("AWS_S3_PRESIGNED_URL_EXPIRY", 15),
			CloudFrontURL:        getEnv("AWS_CLOUDFRONT_URL", ""),
		},
		GRPCClient: &GRPCClient{
			UserServiceURL: getEnv("GRPC_USER_SERVICE_URL", "localhost:5000"),
			Timeout:        getEnvDuration("GRPC_CLIENT_TIMEOUT_SECONDS", 5),
		},
		Prometheus: &Prometheus{
			MetricsHost: getEnv("PROMETHEUS_METRICS_HOST", "localhost"),
			MetricsPort: getEnvInt("PROMETHEUS_METRICS_PORT", 5013),
		},
	}
}

func getEnv(key string, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}

	return defaultVal
}

func getEnvInt(key string, defaultVal int) int {
	if val, err := strconv.Atoi(os.Getenv(key)); err == nil {
		return val
	}

	return defaultVal
}

func getEnvDuration(key string, defaultVal time.Duration) time.Duration {
	if val, err := strconv.Atoi(os.Getenv(key)); err == nil {
		return time.Duration(val) * time.Second
	}

	return defaultVal * time.Second
}
