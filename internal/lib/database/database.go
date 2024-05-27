package database

import (
	"fmt"

	"github.com/sagarmaheshwary/microservices-video-catalog-service/internal/config"
	"github.com/sagarmaheshwary/microservices-video-catalog-service/internal/lib/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	c := config.GetDatabase()

	dns := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		c.Host,
		c.Username,
		c.Password,
		c.Database,
		c.Port,
		c.SSLMode,
		c.Timezone,
	)

	var err error

	DB, err = gorm.Open(postgres.Open(dns), &gorm.Config{})

	if err != nil {
		log.Error("Database failed to connect on %q: %v", dns, err)
	}

	DB.Raw("SELECT 1+1")

	log.Info("Database server connected on %q", dns)
}
