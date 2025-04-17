package memdb

import (
	"fmt"
	"urlShortener/internal/config"
	"urlShortener/internal/storage"
)

type DB struct {
	data map[string]string
	cfg  *config.Config
}

func New(cfg *config.Config) *DB {
	return &DB{
		data: make(map[string]string),
		cfg:  cfg,
	}
}

func (db *DB) SaveUrl(origUrl string, shortUrl string) (string, error) {
	if origUrl == "" {
		return "", fmt.Errorf(storage.UrlNotProvidedError)
	}

	if _, ok := db.data[origUrl]; ok {
		return "", fmt.Errorf(storage.UrlExistsError)
	}

	db.data[origUrl] = shortUrl

	return shortUrl, nil
}

func (db *DB) OrigUrl(shortUrl string) (string, error) {
	for u, su := range db.data {
		if su == shortUrl {
			return u, nil
		}
	}

	return "", fmt.Errorf(storage.UrlNotFoundError)
}
