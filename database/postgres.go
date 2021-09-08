// Package database implements postgres connection and queries.
package database

import (
	"log"

	"github.com/spf13/viper"

	"github.com/go-pg/pg"
)

// DBConn returns a postgres connection pool.
func DBConn() (*pg.DB, error) {
	viper.SetDefault("db_network", "tcp")
	viper.SetDefault("db_addr", "localhost:5432")
	viper.SetDefault("db_user", "postgres")
	viper.SetDefault("db_password", "postgres")
	viper.SetDefault("db_database", "postgres")

	db := pg.Connect(&pg.Options{
		Network:  viper.GetString("db_network"),
		Addr:     viper.GetString("db_addr"),
		User:     viper.GetString("db_user"),
		Password: viper.GetString("db_password"),
		Database: viper.GetString("db_database"),
	})

	if err := checkConn(db); err != nil {
		return nil, err
	}

	if viper.GetBool("db_debug") {
		db.AddQueryHook(&logSQL{})
	}

	return db, nil
}

type logSQL struct{}

func (l *logSQL) BeforeQuery(e *pg.QueryEvent) {}

func (l *logSQL) AfterQuery(e *pg.QueryEvent) {
	query, err := e.FormattedQuery()
	if err != nil {
		panic(err)
	}
	log.Println(query)
}

func checkConn(db *pg.DB) error {
	var n int
	_, err := db.QueryOne(pg.Scan(&n), "SELECT 1")
	return err
}
