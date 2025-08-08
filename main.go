/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"collectify/cmd"
	"collectify/internal/config"
	"log"
)

func main() {
	if err := config.InitConfig(); err != nil {
		log.Fatalf("Failed to initialize config: %v", err)
	}

	cmd.Execute()
}
