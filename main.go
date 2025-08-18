package main

import (
	"collectify/cmd"
	"collectify/internal/cli"
	"collectify/internal/config"
	"collectify/internal/db"
	"collectify/internal/service"
	"embed"
	"log"
)

//go:embed web/build
var frontendEmbedFS embed.FS

func main() {
	cfg, err := config.InitConfig()
	if err != nil {
		log.Fatalf("Failed to initialize config: %v", err)
	}
	if err = db.InitDB(cfg); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// 如果回收站未启用，则清空回收站，避免唯一约束冲突
	if !cfg.RecycleBin.Enable {
		err = service.ClearRecycleBin()
		if err != nil {
			log.Fatalf("Failed to clear recycle bin: %v", err)
		}
	}

	cli.SetFrontendFS(frontendEmbedFS)

	cmd.Execute()
}
