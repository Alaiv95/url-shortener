package storage

const (
	UrlExistsError      = "given url already exists"
	UrlNotProvidedError = "no url provided"
	UrlNotFoundError    = "url not found"
)

type Storage interface {
	SaveUrl(origUrl string, shortUrl string) (string, error)
	OrigUrl(shortUrl string) (string, error)
}
