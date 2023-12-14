package mysql

import (
	"context"
	"database/sql"
	"interview20231208/model"
	"interview20231208/pkg/e"
	"interview20231208/pkg/rdb"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var _ rdb.RDB = (*mysqlClient)(nil)

type mysqlClient struct {
	gormDB *gorm.DB
	sqlDB  *sql.DB
}

func NewMySqlClient(ctx context.Context, dsn string, connMaxLifetime time.Duration, maxOpenConns, maxIdleConns int) *mysqlClient {
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
	c.gormDB.WithContext(ctx).Table(rdb.TbOrders).Create(order)
}

func (c *mysqlClient) UpdateOrdersAndCreateTransactionLogs(ctx context.Context, orders map[int]*model.Order, logs []*model.TransactionLog) {
	c.gormDB.Transaction(func(tx *gorm.DB) error {
		for _, order := range orders {
			if err := tx.Table(rdb.TbOrders).Where("id = ?", order.ID).Update("status", order.Status).Error; err != nil {
				return err
			}
		}

		for _, log := range logs {
			if err := tx.Table(rdb.TbTransactionLogs).Create(log).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (c *mysqlClient) GetOrders(ctx context.Context, opts *model.OrderQueryOpts) []*model.Order {
	result := []*model.Order{}

	tx := c.gormDB.Table(rdb.TbOrders).Select("*")

	if opts.UserID != -1 {
		tx = tx.Where("userId = ?", opts.UserID)
	}
	if opts.Status != e.STATUS_PLACEHOLDER {
		tx = tx.Where("status = ?", opts.Status)
	}
	if !opts.StartTimestamp.IsZero() && !opts.EndTimestamp.IsZero() && opts.StartTimestamp.Unix() <= opts.EndTimestamp.Unix() {
		tx = tx.Where("timestamp >= ? AND timestamp <= ?", opts.StartTimestamp, opts.EndTimestamp)
	}

	if err := tx.Find(&result).Error; err != nil {
		logrus.Error(err)
		return []*model.Order{}
	}

	return result
}

func (c *mysqlClient) GetTransactionLogs(ctx context.Context, opts *model.TransactionLogQueryOpts) []*model.TransactionLog {
	result := []*model.TransactionLog{}

	tx := c.gormDB.Table(rdb.TbTransactionLogs).Select("*")

	if opts.BuyerOrderID != -1 {
		tx = tx.Where("buyerOrderId = ?", opts.BuyerOrderID)
	}
	if opts.SellerOrderID != -1 {
		tx = tx.Where("sellerOrderId = ?", opts.SellerOrderID)
	}
	if !opts.StartTimestamp.IsZero() && !opts.EndTimestamp.IsZero() && opts.StartTimestamp.Unix() <= opts.EndTimestamp.Unix() {
		tx = tx.Where("timestamp >= ? AND timestamp <= ?", opts.StartTimestamp, opts.EndTimestamp)
	}

	if err := tx.Find(&result).Error; err != nil {
		logrus.Error(err)
		return []*model.TransactionLog{}
	}

	return result
}
