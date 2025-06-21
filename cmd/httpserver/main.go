package main

import (
	"log"
	"othnx/internal/infrastructure/httpserver"
)

func main() {
	err := httpserver.Start()
	if err != nil {
		log.Fatal(err)
	}
}
