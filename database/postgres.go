// Package database implements postgres connection and queries.
package database

import (
	"context"
	"database/sql"

	"github.com/spf13/viper"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
)

// DBConn returns a postgres connection pool.
func DBConn() (*bun.DB, error) {
	viper.SetDefault("db_network", "tcp")
	viper.SetDefault("db_addr", "localhost:5432")
	viper.SetDefault("db_user", "postgres")
	viper.SetDefault("db_password", "postgres")
	viper.SetDefault("db_database", "postgres")

	dsn := "postgres://" + viper.GetString("db_user") + ":" + viper.GetString("db_password") + "@" + viper.GetString("db_addr") + "/" + viper.GetString("db_database") + "?sslmode=disable"

	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))

	db := bun.NewDB(sqldb, pgdialect.New())

	if err := checkConn(db); err != nil {
		return nil, err
	}

	if viper.GetBool("db_debug") {
		db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))
	}

	return db, nil
}

func checkConn(db *bun.DB) error {
	var n int
	return db.NewSelect().ColumnExpr("1").Scan(context.Background(), &n)
}
