package main

import (
	"aegix/internal/httpserver"
	"log"
)

func main() {
	err := httpserver.Start(); if err != nil {
		log.Fatal(err)
	}
}
