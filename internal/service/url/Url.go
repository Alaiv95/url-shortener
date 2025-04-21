package url

import (
	"context"
	"encoding/json"
)

func (s *service) Url(key string, ctx context.Context) (string, error) {
	var origUrl string
	_ = ctx

	if cached, err := s.cacheGetter.Get(key); err == nil {
		err = json.Unmarshal(cached, &origUrl)
		if err != nil {
			s.log.Error("Error parsing value from cache", "err", err.Error())
			return "", err
		}
	} else {
		origUrl, err = s.repo.Url(key)

		if err != nil {
			s.log.Error("Error getting url", "err", err.Error())
			return "", err
		}
	}

	return origUrl, nil
}
