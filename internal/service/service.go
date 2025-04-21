package service

import "context"

type UrlService interface {
	SaveUrl(origUrl string, domain string, ctx context.Context) (string, error)
	Url(key string, ctx context.Context) (string, error)
}
