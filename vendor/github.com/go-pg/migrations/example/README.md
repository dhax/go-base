# Example

You need to create database `pg_migrations_example` before running this example.

```bash
> psql -c "CREATE DATABASE pg_migrations_example"
CREATE DATABASE

> go run *.go init
version is 0

> go run *.go version
version is 0

> go run *.go
creating table my_table...
adding id column...
seeding my_table...
migrated from version 0 to 3

> go run *.go version
version is 3

> go run *.go reset
truncating my_table...
dropping id column...
dropping table my_table...
migrated from version 3 to 0

> go run *.go
creating table my_table...
adding id column...
seeding my_table...
migrated from version 0 to 3

> go run *.go down
truncating my_table...
migrated from version 3 to 2

> go run *.go version
version is 2

> go run *.go set_version 1
migrated from version 2 to 1

> go run *.go create add email to users
created migration 4_add_email_to_users.go
```

## Transactions

If you'd want to wrap the whole run in a big transaction, which may be the case if you have multi-statement migrations, the code in `main.go` should be slightly modified:

```go
var oldVersion, newVersion int64

err := db.RunInTransaction(func(tx *pg.Tx) (err error) {
    oldVersion, newVersion, err = migrations.Run(tx, flag.Args()...)
    return
})
if err != nil {
    exitf(err.Error())
}
```
