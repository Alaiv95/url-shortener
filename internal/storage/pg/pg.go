package pg

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"urlShortener/internal/storage"
)

// Storage структура хранилища в памяти
type Storage struct {
	db *sql.DB
}

// New конструктор инициализации хранилища в памяти
func New(path string) (*Storage, error) {
	const op = "storage.pg.New"

	db, err := sql.Open("postgres", path)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	stmt, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS url(
		id SERIAL PRIMARY KEY,
		slug VARCHAR(255) UNIQUE NOT NULL,
		url VARCHAR(255) NOT NULL);
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	stmt2, err := db.Prepare(`CREATE INDEX IF NOT EXISTS slug_idx ON url(slug)`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt2.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

// SaveUrl сохранение новой короткой ссылки в памяти
func (s *Storage) SaveUrl(origUrl string, shortUrl string) (string, error) {
	if origUrl == "" {
		return "", fmt.Errorf(storage.UrlNotProvidedError)
	}

	stmt, err := s.db.Prepare("INSERT INTO url(url, slug) VALUES ($1, $2)")
	if err != nil {
		return "", err
	}

	res, err := stmt.Exec(origUrl, shortUrl)
	if err != nil {
		return "", fmt.Errorf(storage.UrlExistsError)
	}

	_ = res

	return shortUrl, nil
}

// Url получение оригинальной ссылки по короткой
func (s *Storage) Url(shortUrl string) (string, error) {
	stmt, err := s.db.Prepare("SELECT url FROM url WHERE slug = ?")
	if err != nil {
		return "", err
	}

	var resUrl string

	err = stmt.QueryRow(shortUrl).Scan(&resUrl)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf(storage.UrlNotFoundError)
		}
		return "", err
	}

	return resUrl, nil
}
