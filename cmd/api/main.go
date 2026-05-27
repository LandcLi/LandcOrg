package main

import (
	"fmt"
	"log"

	landcorg "github.com/LandcLi/LandcOrg"
)

func main() {
	engine, err := landcorg.InitDefault("config.yml")
	if err != nil {
		log.Fatalf("Failed to initialize: %v", err)
	}

	port := ":8083"
	fmt.Printf("LandcOrg starting on %s\n", port)
	if err := engine.Run(port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
