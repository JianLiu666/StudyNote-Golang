package models

import (
	"database/sql"
	"fmt"
	"httpserver/pkg/setting"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var gormDB *gorm.DB
var sqlDB *sql.DB

type Model struct {
	ID         int `gorm:"primary_key" json:"id"`
	CreatedOn  int `json:"created_on"`
	ModifiedOn int `json:"modified_on"`
	DeletedOn  int `json:"deleted_on"`
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

	gormDB.Callback().Create().Before("gorm:before_create").Register("update_timestamp", updateTimeStampForCreateCallback)
	gormDB.Callback().Update().Before("gorm:before_update").Register("Update_timestamp", updateTimeStampForUpdateCallback)
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

func updateTimeStampForCreateCallback(db *gorm.DB) {
	if db.Error != nil || db.DryRun || db.Statement.Schema == nil {
		return
	}

	nowTime := time.Now().Unix()

	createTimeField := db.Statement.Schema.LookUpField("CreatedOn")
	if createTimeField == nil {
		return
	}

	if createTimeField.DefaultValue == "" {
		db.Statement.SetColumn("CreatedOn", nowTime)
	}

	modifyTimeField := db.Statement.Schema.LookUpField("ModifiedOn")
	if modifyTimeField == nil {
		return
	}

	if modifyTimeField.DefaultValue == "" {
		db.Statement.SetColumn("ModifiedOn", nowTime)
	}
}

func updateTimeStampForUpdateCallback(db *gorm.DB) {
	if db.Error != nil || db.DryRun || db.Statement.Schema == nil {
		return
	}

	modifyTimeField := db.Statement.Schema.LookUpField("ModifiedOn")
	if modifyTimeField == nil {
		return
	}

	if modifyTimeField.DefaultValue == "" && !db.Statement.Changed("ModifiedOn") {
		db.Statement.SetColumn("ModifiedOn", time.Now().Unix())
	}
}
