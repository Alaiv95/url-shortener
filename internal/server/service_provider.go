package server

import (
	"log/slog"
	"os"
	"urlShortener/internal/api"
	"urlShortener/internal/config"
	"urlShortener/internal/kafka"
	"urlShortener/internal/redis"
	"urlShortener/internal/service"
	"urlShortener/internal/service/url"
	"urlShortener/internal/storage"
	"urlShortener/internal/storage/pg"
)

const (
	local = "local"
	prod  = "prod"
)

type serviceProvider struct {
	redis      *redis.Client
	kafka      *kafka.Client
	urlService service.UrlService
	urlRepo    storage.UrlRepo
	config     *config.Config
	urlApi     *api.API
	log        *slog.Logger
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (p *serviceProvider) Config() *config.Config {
	if p.config == nil {
		cfg := config.MustLoad()
		p.config = cfg
		p.Logger().Debug("Config initialized")
	}

	return p.config
}

func (p *serviceProvider) Redis() *redis.Client {
	if p.redis == nil {
		red := redis.New(&p.Config().Redis)
		p.redis = red
		p.Logger().Debug("Redis initialized")
	}

	return p.redis
}

func (p *serviceProvider) Logger() *slog.Logger {
	if p.log == nil {
		logger := setupLogger(p.Config().Env)
		p.log = logger
		p.Logger().Debug("Logger initialized")
	}

	return p.log
}

func (p *serviceProvider) Kafka() *kafka.Client {
	if p.kafka == nil {
		kf, err := kafka.New(p.Config().Kafka.Address, p.Config().Kafka.UrlTopic, p.Logger())
		if err != nil {
			p.Logger().Error("error initializing Kafka client")
			os.Exit(1)
		}

		p.kafka = kf
		p.Logger().Debug("Kafka initialized")
	}

	return p.kafka
}

func (p *serviceProvider) UrlRepo() storage.UrlRepo {
	if p.urlRepo == nil {
		repo, err := pg.New(p.Config().StoragePath)
		if err != nil {
			p.Logger().Error("error initializing UrlRepo")
			os.Exit(1)
		}

		p.urlRepo = repo
		p.Logger().Debug("Storage initialized")
	}

	return p.urlRepo
}

func (p *serviceProvider) UrlService() service.UrlService {
	if p.urlService == nil {
		p.urlService = url.New(p.UrlRepo(), p.Logger(), p.Kafka(), p.Redis(), p.Redis())
		p.Logger().Debug("UrlService initialized")
	}

	return p.urlService
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case local:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))
	case prod:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}))
	}

	return log
}
