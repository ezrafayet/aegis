package main

import (
	"aegis/internal/infrastructure/httpserver"
	"log"
)

func main() {
	if err := httpserver.Start(); err != nil {
		log.Fatal(err)
	}
}
