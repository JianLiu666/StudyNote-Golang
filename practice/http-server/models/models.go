package models

import (
	"database/sql"
	"fmt"
	"httpserver/pkg/setting"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var gormDB *gorm.DB
var sqlDB *sql.DB

type Model struct {
	ID         int   `gorm:"primary_key" json:"id"`
	CreatedOn  int64 `json:"created_on"`
	ModifiedOn int64 `json:"modified_on"`
	DeletedOn  int64 `json:"deleted_on"`
}

func SetUp() {
	var err error
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		setting.DatabaseSetting.User,
		setting.DatabaseSetting.Password,
		setting.DatabaseSetting.Host,
		setting.DatabaseSetting.DBName,
	)

	gormDB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Panicf("failed to open database by gorm: %v", err)
	}

	sqlDB, err = gormDB.DB()
	if err != nil {
		log.Panicf("failed to get sql.DB : %v", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	if err = sqlDB.Ping(); err != nil {
		log.Panicf("failed to ping sql.DB: %v", err)
	}
}

func Close() {
	if err := sqlDB.Close(); err != nil {
		log.Printf("failed to close database: %v", err)
	}
}
