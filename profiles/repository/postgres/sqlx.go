package postgres

import (
	"context"
	"datapreparation/profiles"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/spf13/viper"
)

type sqlxDB struct {
	db *sqlx.DB
}

func NewSQLX(h, p, u, dbname, pass string) (profiles.Repository, error) {
	connectionString := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", h, p, u, dbname, pass)
	db, err := sqlx.Connect("postgres", connectionString)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}

	return &sqlxDB{db}, nil
}

func (repo *sqlxDB) Save(ctx context.Context, key string, encrypted [][]byte) error {
	q := fmt.Sprintf(`
		INSERT INTO %s (ICCID, chunks) VALUES (?, ?) ON CONFLICT DO NOTHING
	`, viper.GetString("TABLE NAME"))
	q = repo.db.Rebind(q)

	_, err := repo.db.Exec(q, key, pq.Array(encrypted))
	if err != nil {
		return err
	}

	return nil
}
