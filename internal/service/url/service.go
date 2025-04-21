package url

import (
	"context"
	"log/slog"
	"time"
	serv "urlShortener/internal/service"
	"urlShortener/internal/storage"
)

type Producer interface {
	Produce(msgVal []byte, ctx context.Context)
}

type TimedSetter interface {
	Set(key string, v any, time time.Duration) error
}

type Getter interface {
	Get(key string) ([]byte, error)
}

type service struct {
	repo        storage.UrlRepo
	log         *slog.Logger
	kf          Producer
	cacheSetter TimedSetter
	cacheGetter Getter
}

func New(repo storage.UrlRepo, log *slog.Logger, kf Producer, cacheSetter TimedSetter, getter Getter) serv.UrlService {
	return &service{
		repo:        repo,
		log:         log,
		kf:          kf,
		cacheSetter: cacheSetter,
		cacheGetter: getter,
	}
}
