package database

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"randomMeetsProject/config"
	"randomMeetsProject/internal/models/sql_models"
)

func DB() (*gorm.DB, error) {
	cfg := config.Config{}
	db, err := gorm.Open(postgres.Open(cfg.DbUrl()), &gorm.Config{})
	if err != nil {
		log.Fatal("Error connecting to database")
		return nil, err
	}
	return db, nil
}

func InitDB() error {
	db, err := DB()
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	tables := []interface{}{&sql_models.User{}, &sql_models.MeetUp{}}

	err = db.AutoMigrate(tables...)
	if err != nil {
		return fmt.Errorf("failed to auto migrate tables: %w\nTables: %+v", err, tables)
	}

	fmt.Println("Database migrated successfully")
	return nil
}
