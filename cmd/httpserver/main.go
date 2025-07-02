package main

import (
	"aegis/internal/infrastructure/httpserver"
	"log"
)

func main() {
	err := httpserver.Start()
	if err != nil {
		log.Fatal(err)
	}
}
