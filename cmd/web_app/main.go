package main

import (
	"flag"
	"github.com/maxuanquang/social-network/configs"
	"github.com/maxuanquang/social-network/internal/app/web_app"
	"log"
	_ "github.com/maxuanquang/social-network/docs"
)

// @title 			Gin Social Network Service
// @version 		1.0
// @description 	A simple social network management service API in Go using Gin framework
func main() {
	// Flags
	cfgPath := flag.String("conf", "config.yml", "Path to config file for this service")

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
