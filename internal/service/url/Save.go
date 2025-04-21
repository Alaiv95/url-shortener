package url

import (
	"context"
	"fmt"
	"time"
	"urlShortener/internal/lib/base62"
	"urlShortener/internal/lib/numGen"
)

func (s *service) SaveUrl(origUrl string, domain string, ctx context.Context) (string, error) {
	s.log.Info("Start saving short url for:  ", "url", origUrl)

	slug := base62.ConvertNum(numGen.Generate())

	resp, err := s.repo.SaveUrl(origUrl, slug)
	if err != nil {
		s.log.Error("Error saving url", "err", err.Error())
		return "", err
	}

	// form short url with current domain
	shortUrl := fmt.Sprintf("%s/api/v1/url/%s", domain, resp)

	go func() {
		// produce message to kafka that url created
		s.kf.Produce([]byte(shortUrl), ctx)

		// set url to cache for 12 hr
		err = s.cacheSetter.Set(slug, origUrl, time.Hour*12)
		if err != nil {
			s.log.Error("Error saving url to cache", "err", err.Error())
		}
	}()

	return shortUrl, nil
}
