package rdb

import (
	sqle "github.com/dolthub/go-mysql-server"
	"github.com/dolthub/go-mysql-server/memory"
	"github.com/dolthub/go-mysql-server/server"
	"github.com/dolthub/go-mysql-server/sql"
	"github.com/dolthub/go-mysql-server/sql/information_schema"
	"github.com/dolthub/go-mysql-server/sql/types"
)

type sqlServer struct {
	engine *sqle.Engine
	server *server.Server
}

func NewSqlServer(address string) *sqlServer {
	engine := sqle.NewDefault(
		memory.NewDBProvider(
			createFakeDatabase(),
			information_schema.NewInformationSchemaDatabase(),
		))

	config := server.Config{
		Protocol: "tcp",
		Address:  address,
	}
	server, err := server.NewDefaultServer(config, engine)
	if err != nil {
		panic(err)
	}

	return &sqlServer{
		engine: engine,
		server: server,
	}
}

func (f *sqlServer) Enable() {
	go func() {
		if err := f.server.Start(); err != nil {
			panic(err)
		}
	}()
}

func createFakeDatabase() *memory.Database {
	db := memory.NewDatabase("trading")
	db.EnablePrimaryKeyIndexes()

	db.AddTable(TbOrders, createOrdersTable(db))
	db.AddTable(TbTransactionLogs, createTransactionLogsTable(db))

	return db
}

func createOrdersTable(db *memory.Database) *memory.Table {
	table := memory.NewTable(TbOrders, sql.NewPrimaryKeySchema(sql.Schema{
		{Name: "id", Type: types.Int32, Nullable: false, Source: TbOrders, AutoIncrement: true, PrimaryKey: true},
		{Name: "userId", Type: types.Int32, Nullable: false, Source: TbOrders},
		{Name: "roleType", Type: types.Int8, Nullable: false, Source: TbOrders},
		{Name: "orderType", Type: types.Int8, Nullable: false, Source: TbOrders},
		{Name: "durationType", Type: types.Int8, Nullable: false, Source: TbOrders},
		{Name: "price", Type: types.Int32, Nullable: false, Source: TbOrders},
		{Name: "quantity", Type: types.Int32, Nullable: false, Source: TbOrders},
		{Name: "status", Type: types.Int8, Nullable: false, Source: TbOrders},
		{Name: "timestamp", Type: types.DatetimeMaxPrecision, Nullable: false, Source: TbOrders},
	}), db.GetForeignKeyCollection())

	return table
}

func createTransactionLogsTable(db *memory.Database) *memory.Table {
	table := memory.NewTable(TbTransactionLogs, sql.NewPrimaryKeySchema(sql.Schema{
		{Name: "id", Type: types.Int32, Nullable: false, Source: TbTransactionLogs, AutoIncrement: true, PrimaryKey: true},
		{Name: "buyerOrderId", Type: types.Int32, Nullable: false, Source: TbTransactionLogs},
		{Name: "sellerOrderId", Type: types.Int32, Nullable: false, Source: TbTransactionLogs},
		{Name: "price", Type: types.Int32, Nullable: false, Source: TbTransactionLogs},
		{Name: "quantity", Type: types.Int32, Nullable: false, Source: TbTransactionLogs},
		{Name: "timestamp", Type: types.DatetimeMaxPrecision, Nullable: false, Source: TbTransactionLogs},
	}), db.GetForeignKeyCollection())

	return table
}
