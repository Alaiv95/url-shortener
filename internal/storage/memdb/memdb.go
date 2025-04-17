package memdb

import (
	"fmt"
	"urlShortener/internal/config"
	"urlShortener/internal/storage"
)

// Storage структура хранилища в памяти
type Storage struct {
	data map[string]string
	cfg  *config.Config
}

// New конструктор инициализации хранилища в памяти
func New(cfg *config.Config) *Storage {
	return &Storage{
		data: make(map[string]string),
		cfg:  cfg,
	}
}

// SaveUrl сохранение новой короткой ссылки в памяти
func (db *Storage) SaveUrl(origUrl string, shortUrl string) (string, error) {
	if origUrl == "" {
		return "", fmt.Errorf(storage.UrlNotProvidedError)
	}

	if _, ok := db.data[shortUrl]; ok {
		return "", fmt.Errorf(storage.UrlExistsError)
	}

	db.data[shortUrl] = origUrl

	return shortUrl, nil
}

// GetUrl получение оригинальной ссылки по короткой
func (db *Storage) GetUrl(shortUrl string) (string, error) {
	if u, ok := db.data[shortUrl]; ok {
		return u, nil
	}

	return "", fmt.Errorf(storage.UrlNotFoundError)
}
