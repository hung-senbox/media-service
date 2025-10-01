package db

import (
	"fmt"
	"log"
	"media-service/internal/media/model"
	"media-service/pkg/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var MySqlDB *gorm.DB

func ConnectMySQL() {
	d := config.AppConfig.Database.MySQL
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		d.User, d.Password, d.Host, d.Port, d.Name)

	var err error
	MySqlDB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to MySQL: %v", err)
	}

	err = MySqlDB.AutoMigrate(&model.Topic{})
	if err != nil {
		log.Fatalf("AutoMigrate failed: %v", err)
	}

	log.Println("Connected to MySQL and migrated schema")
}
