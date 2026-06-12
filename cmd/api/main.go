package main

import (
	"fmt"
	"log"

	landcorg "github.com/LandcLi/LandcOrg"
)

func main() {
	server, err := landcorg.InitDefault("config.yml")
	if err != nil {
		log.Fatalf("Failed to initialize: %v", err)
	}

	fmt.Println("LandcOrg starting on :8083")
	if err := server.Run(); err != nil {
		log.Fatalf("Server stopped with error: %v", err)
	}
}
