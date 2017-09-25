package main

import (
	"fmt"

	"github.com/go-pg/migrations"
)

func init() {
	migrations.Register(func(db migrations.DB) error {
		fmt.Println("seeding my_table...")
		_, err := db.Exec(`INSERT INTO my_table VALUES (1)`)
		return err
	}, func(db migrations.DB) error {
		fmt.Println("truncating my_table...")
		_, err := db.Exec(`TRUNCATE my_table`)
		return err
	})
}
