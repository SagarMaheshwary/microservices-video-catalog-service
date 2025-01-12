package database

import (
	"fmt"

	"github.com/sagarmaheshwary/microservices-video-catalog-service/internal/config"
	"github.com/sagarmaheshwary/microservices-video-catalog-service/internal/lib/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Conn *gorm.DB

func Connect() {
	var err error
	c := config.Conf.Database

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		c.Host,
		c.Username,
		c.Password,
		c.Database,
		c.Port,
		c.SSLMode,
		c.Timezone,
	)

	Conn, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		logger.Fatal("Database failed to connect on %q: %v", dsn, err)
	}

	logger.Info("Database server connected on %q", dsn)
}

func HealthCheck() bool {
	sqlDB, err := Conn.DB()

	if err != nil {
		logger.Error("DB health check failed! %v", err)

		return false
	}

	if err := sqlDB.Ping(); err != nil {
		logger.Error("DB health check failed! %v", err)

		return false // Database is not reachable
	}

	var result int

	if err := Conn.Raw("SELECT 1").Scan(&result).Error; err != nil {
		logger.Error("DB health check failed! %v", err)

		return false
	}

	return true
}
