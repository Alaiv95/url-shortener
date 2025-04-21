package storage

const (
	UrlExistsError      = "given url already exists"
	UrlNotProvidedError = "no url provided"
	UrlNotFoundError    = "url not found"
)

type UrlRepo interface {
	SaveUrl(origUrl string, shortUrl string) (string, error)
	Url(shortUrl string) (string, error)
}
