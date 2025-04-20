# Postgres

### Using Config

```golang
// config.go
pkg config

type Config struct {
    cfgPkg.BaseConfig `mapstructure:",squash"`
    Postgres postgresPkg.Config `mapstructure:"postgres"`
}

// main.go
var conf config.Config
config.LoadConfig("WELLKNOWN", &conf, nil)

pgDb, err := postgresPkg.ConnectRetry(ctx, conf.Postgres, time.Minute)
```

### Using ConnectRetry

`ConnectRetry` is a function that tries to establish a db connection for the given
`maxDuration time.Duration`. It will periodically retry, using a fibonacci backoff
pattern.

For usage example, see above.

### Using MigrateUp

Given the following folder structure:

```
src/
├─ postgres/
│  ├─ postgres.go
│  ├─ migrations/
│  │  ├─ 001_init.down.sql
│  │  ├─ 001_init.up.sql
├─ main.go     
```

First, create a variable of type `embed.FS` and annotate it with `//go:embed relative-path`,
pointing to the folder containing the migration files (no `.` or `..` allowed):
```golang
// postgres/postgres.go
pkg postgres 

//go:embed migrations
var Migrations embed.FS
```

Then, when you established the DB connection (here in `main.go`), pass
the variable and path to `MigrateUP()` 

```golang
// main.go
pkg main

func main() {
	// ...
	pgDb, _ := postgresPkg.ConnectRetry(...)
	postgres.MigrateUP(pgDb, postgres.Migrations, "migrations")
} 
```

The last parameter is the path relative to the embedded file system.