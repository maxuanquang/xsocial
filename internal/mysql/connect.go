package mysql

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDatabase() *gorm.DB {
	// TODO: receive configuration

	const (
		DbServerAddress  = "192.168.0.103:3306"
		DbServerUser     = "quangmx"
		DbServerPassword = "2511"
		DbName           = "social-network"
	)

	// Initialize database
	mysqlConfig := mysql.Config{
		DSN: fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", DbServerUser, DbServerPassword, DbServerAddress, DbName),
	}
	dialector := mysql.New(mysqlConfig)
	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		log.Printf("Failed connecting to db: %v", err)
	}

	return db
}
