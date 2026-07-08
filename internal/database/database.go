package database

import (
	"be-technical-test/config"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database struct {
	DB  *gorm.DB
	cfg *config.DatabaseConfig
}

func NewDatabase(cfg *config.DatabaseConfig) (*Database, error) {
	var err error

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.Host,
		cfg.User,
		cfg.Password,
		cfg.Name,
		fmt.Sprintf("%d", cfg.Port),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return &Database{
		DB:  db,
		cfg: cfg}, nil
}
