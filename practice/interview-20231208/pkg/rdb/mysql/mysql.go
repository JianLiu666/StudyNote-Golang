package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"interview20231208/model"
	"interview20231208/pkg/rdb"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type mysqlClient struct {
	gormDB *gorm.DB
	sqlDB  *sql.DB
}

func NewMySqlClient(ctx context.Context, dsn string, connMaxLifetime time.Duration, maxOpenConns, maxIdleConns int) rdb.RDB {
	gormDB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		logrus.Panicf("failed to open database by gorm: %v", err)
	}

	sqlDB, err := gormDB.DB()
	if err != nil {
		logrus.Panicf("failed to get sql.DB: %v", err)
	}

	sqlDB.SetConnMaxLifetime(connMaxLifetime)
	sqlDB.SetMaxOpenConns(maxOpenConns)
	sqlDB.SetMaxIdleConns(maxIdleConns)

	return &mysqlClient{
		gormDB: gormDB,
		sqlDB:  sqlDB,
	}
}

func (c *mysqlClient) Shutdown(ctx context.Context) {
	if err := c.sqlDB.Close(); err != nil {
		logrus.Panicf("failed to close sql.DB : %v", err)
	}
}

func (c *mysqlClient) CreateOrder(ctx context.Context, order *model.Order) {
	fmt.Println(order.ID)
	c.gormDB.WithContext(ctx).Table(rdb.TbOrders).Clauses(clause.OnConflict{UpdateAll: true}).Create(&order)
	fmt.Println(order.ID)
}

func (c *mysqlClient) UpdateOrdersAndCreateTransactionLogs(ctx context.Context, orders map[int]*model.Order, logs []*model.TransactionLog) {
	// TODO
}
