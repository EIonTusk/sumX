package db

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

func Connect() *bun.DB {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	port := os.Getenv("POSTGRES_PORT")
	pdb := os.Getenv("POSTGRES_DB")

	if user == "" || password == "" || port == "" || pdb == "" {
		panic("At least one postgres env var is not set")
	}

    dsn := fmt.Sprintf("postgres://%s:%s@postgres:%s/%s?sslmode=disable", user, password, port, pdb)
    sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))

	db := bun.NewDB(sqldb, pgdialect.New())

    ctx := context.Background()
    if err := db.PingContext(ctx); err != nil {
		panic("DB ping failed")
    }

	models := []any {
		(*XUser)(nil),
		(*Summary)(nil),
	}

	for _, model := range models {
		_, err := db.NewCreateTable().
			Model(model).
			IfNotExists().  // ✅ avoids dropping
			Exec(ctx)
		if err != nil {
			panic(fmt.Errorf("create table failed: %w", err))
		}
	}

	return db
}
