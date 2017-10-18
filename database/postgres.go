package database

import (
	"log"
	"time"

	"github.com/spf13/viper"

	"github.com/go-pg/pg"
)

// DBConn returns a postgres connection pool.
func DBConn() (*pg.DB, error) {

	opts, err := pg.ParseURL(viper.GetString("database_url"))
	if err != nil {
		return nil, err
	}

	db := pg.Connect(opts)
	if err := checkConn(db); err != nil {
		return nil, err
	}

	if viper.GetBool("db_debug") {
		db.OnQueryProcessed(func(event *pg.QueryProcessedEvent) {
			query, err := event.FormattedQuery()
			if err != nil {
				panic(err)
			}
			log.Printf("%s %s\n", time.Since(event.StartTime), query)
		})
	}

	return db, nil
}

func checkConn(db *pg.DB) error {
	var n int
	if _, err := db.QueryOne(pg.Scan(&n), "SELECT 1"); err != nil {
		return err
	}
	return nil
}
