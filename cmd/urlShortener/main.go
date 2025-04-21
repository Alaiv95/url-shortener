package main

import (
	"urlShortener/internal/server"
)

func main() {
	serv := server.New()

	serv.Start()
}
