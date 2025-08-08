/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"collectify/cmd"
	"collectify/internal/config"
	"collectify/internal/db"
	"log"
)

func main() {
	cfg, err := config.InitConfig()
	if err != nil {
		log.Fatalf("Failed to initialize config: %v", err)
	}
	if err = db.InitDB(cfg); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	cmd.Execute()
}
