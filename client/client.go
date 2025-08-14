package client

import (
	"fmt"
	"io"
	"log"
	"time"

	"github.com/dailypay/daily-go/pkg/ctxlogger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"
)

type writer struct{}

func (w writer) Printf(string, ...any) {}

func NewClient(_ *ctxlogger.Logger) (*gorm.DB, error) {
	glog := glogger.New(log.New(io.Discard, "\r\n", log.LstdFlags), glogger.Config{})

	db, err := gorm.Open(postgres.Open(
		"postgres://root:password@localhost:5400/postgres",
	), &gorm.Config{Logger: glog})
	if err != nil {
		return nil, fmt.Errorf("connect: %w", err)
	}

	db.AutoMigrate(&FeatureFlag{})

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("getting *sql.DB object: %w", err)
	}

	maxOpenConnections := 20
	maxIdleConnections := 5
	maxConnectionLifetime := time.Hour

	sqlDB.SetMaxOpenConns(maxOpenConnections)
	sqlDB.SetMaxIdleConns(maxIdleConnections)
	sqlDB.SetConnMaxLifetime(maxConnectionLifetime)

	return db, nil
}
