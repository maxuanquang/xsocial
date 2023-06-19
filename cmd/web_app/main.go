package main

import (
	"flag"
	"log"
	"github.com/maxuanquang/social-network/configs"
	"github.com/maxuanquang/social-network/internal/app/web_app"
)

func main() {
	// Flags
	cfgPath := flag.String("conf", "configs/files/test.yml", "Path to config file for this service")

	// Load configurations
	cfg, err := configs.GetWebConfig(*cfgPath)
	if err != nil {
		log.Fatalf("failed to parse config: %v", err)
		return
	}

	// Create new web controller
	web_controller, err := web_app.NewWebController(cfg)
	if err != nil {
		log.Fatalf("failed to create web controller: %v", err)
		return
	}

	// Serve
	web_controller.Run()
}