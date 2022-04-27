package postgres

import (
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var Database *Store

type Store struct {
	*UserStore
}

func init() {
	databaseSource := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DATABASE_HOST"), os.Getenv("DATABASE_USER"), os.Getenv("DATABASE_PASSWORD"), os.Getenv("DATABASE_NAME"))
	db, err := sqlx.Open(os.Getenv("DATABASE_DRIVER"), databaseSource)
	if err != nil {
		panic(fmt.Errorf("error while opening a database, error: %w", err))
	}
	if err = db.Ping(); err != nil {
		panic(fmt.Errorf("error while connecting to a database, error: %w", err))
	}
	Database = &Store{
		UserStore: NewUserStore(db),
	}
}
