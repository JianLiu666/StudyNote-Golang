package article

import (
	"sync"

	"gorm.io/gorm"
)

var instance *service
var once sync.Once

type service struct {
	db *gorm.DB
}

func Init(db *gorm.DB) {
	once.Do(func() {
		instance = &service{
			db: db,
		}
	})
}
