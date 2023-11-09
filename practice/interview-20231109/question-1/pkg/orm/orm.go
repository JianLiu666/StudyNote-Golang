package orm

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func NewConnect() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("./sqlite.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}
	if err := sqlDB.Ping(); err != nil {
		panic(err)
	}

	return db
}
