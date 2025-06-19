package main

import (
	"log"
	"othnx/internal/httpserver"
)

func main() {
	err := httpserver.Start()
	if err != nil {
		log.Fatal(err)
	}
}
