package migrate

import (
	"fmt"
	"os"

	"github.com/dhax/go-base/database"
	"github.com/go-pg/migrations"
	"github.com/go-pg/pg"
)

// Migrate runs go-pg migrations
func Migrate(args []string) {
	db, err := database.DBConn()
	if err != nil {
		panic(err)
	}

	err = db.RunInTransaction(func(tx *pg.Tx) error {
		oldVersion, newVersion, err := migrations.Run(tx, args...)
		if err != nil {
			return err
		}
		if newVersion != oldVersion {
			fmt.Printf("migrated from version %d to %d\n", oldVersion, newVersion)
		} else {
			fmt.Printf("version is %d\n", oldVersion)
		}
		return nil
	})
	if err != nil {
		exitf(err.Error())
	}

}

// Reset runs reverts all migrations to version 0 and then applies all migrations to latest
func Reset() {
	db, err := database.DBConn()
	if err != nil {
		exitf(err.Error())
	}

	version, err := migrations.Version(db)
	if err != nil {
		exitf(err.Error())
	}

	err = db.RunInTransaction(func(tx *pg.Tx) error {
		for version != 0 {
			oldVersion, newVersion, err := migrations.Run(tx, "down")
			if err != nil {
				return err
			}
			fmt.Printf("migrated from version %d to %d\n", oldVersion, newVersion)
			version = newVersion
		}
		return nil
	})
	if err != nil {
		exitf(err.Error())
	}
}

func errorf(s string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, s+"\n", args...)
}

func exitf(s string, args ...interface{}) {
	errorf(s, args...)
	os.Exit(1)
}
