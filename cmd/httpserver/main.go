package main

import (
	"auth302/internal/infrastructure/httpserver"
	"log"
)

func main() {
	if err := httpserver.Start(); err != nil {
		log.Fatal(err)
	}
}
