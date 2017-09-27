package migrate

import (
	"fmt"

	"github.com/go-pg/migrations"
)

const accountTable = `
CREATE TABLE accounts (
id serial NOT NULL,
created_at timestamp with time zone NOT NULL DEFAULT current_timestamp,
updated_at timestamp with time zone DEFAULT current_timestamp,
last_login timestamp with time zone NOT NULL DEFAULT current_timestamp,
email text NOT NULL UNIQUE,
name text NOT NULL,
active boolean NOT NULL DEFAULT TRUE,
roles text[] NOT NULL DEFAULT '{"user"}',
PRIMARY KEY (id)
)`

const tokenTable = `
CREATE TABLE tokens (
id serial NOT NULL,
created_at timestamp with time zone NOT NULL DEFAULT current_timestamp,
updated_at timestamp with time zone NOT NULL DEFAULT current_timestamp,
account_id int NOT NULL REFERENCES accounts(id),
token text NOT NULL UNIQUE,
expiry timestamp with time zone NOT NULL,
mobile boolean NOT NULL DEFAULT FALSE,
identifier text,
PRIMARY KEY (id)
)`

func init() {
	up := []string{
		accountTable,
		tokenTable,
	}

	down := []string{
		`DROP TABLE tokens`,
		`DROP TABLE accounts`,
	}

	migrations.Register(func(db migrations.DB) error {
		fmt.Println("creating initial tables")
		for _, q := range up {
			_, err := db.Exec(q)
			if err != nil {
				return err
			}
		}
		return nil
	}, func(db migrations.DB) error {
		fmt.Println("dropping initial tables")
		for _, q := range down {
			_, err := db.Exec(q)
			if err != nil {
				return err
			}
		}
		return nil
	})
}
