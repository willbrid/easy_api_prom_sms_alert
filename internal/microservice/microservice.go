package microservice

type IMicroservice interface {
	Consume(url string, body string) error
	GetUrlAndBody(dest string, message string) (string, string, error)
}
