package postgres

import (
	"context"
	"errors"

	"github.com/armbdevelop/urlshorten/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresStorage struct {
    pool *pgxpool.Pool  // или *pgx.Conn
}

func NewPostgresStorage(connString string) (*PostgresStorage, error) {
    pool, err := pgxpool.New(context.Background(), connString)

	if err != nil {
		return nil, err
	}

	return  &PostgresStorage{pool: pool}, nil
}


func (p *PostgresStorage) Save(ctx context.Context, url model.URL) error {
	_, err := p.pool.Exec(ctx, "INSERT INTO urls (original_url, short_url) VALUES ($1, $2)", 
	url.OriginalURL, url.ShortURL)

	return err
}


func (p *PostgresStorage) GetByShort(ctx context.Context, short string) (model.URL, error) {
	var url model.URL

	err := p.pool.QueryRow(ctx, 
		"SELECT id, original_url, short_url, created_at FROM urls WHERE short_url = $1",
        short,
      ).Scan(&url.ID, &url.OriginalURL, &url.ShortURL, &url.CreatedAt)
	
	if errors.Is(err, pgx.ErrNoRows) {
        return model.URL{}, model.ErrNotFound
    }

    if err != nil {
        return model.URL{}, err
    }

    return url, nil

}

func (p *PostgresStorage) GetByOriginal(ctx context.Context, original string) (model.URL, error) {
	var url model.URL

	err := p.pool.QueryRow(ctx, 
		"SELECT id, original_url, short_url, created_at FROM urls WHERE original_url = $1",
        original,
    ).Scan(&url.ID, &url.OriginalURL, &url.ShortURL, &url.CreatedAt)
	
	if errors.Is(err, pgx.ErrNoRows) {
        return model.URL{}, model.ErrNotFound
    }

    if err != nil {
        return model.URL{}, err
    }

    return url, nil
}