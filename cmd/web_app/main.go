package main

import (
	"flag"
	"log"

	"github.com/maxuanquang/social-network/configs"
	_ "github.com/maxuanquang/social-network/docs"
	"github.com/maxuanquang/social-network/internal/app/web_app"
)

//@title			Gin Social Network Service
//	@version		1.0
//	@description	A simple social network management service API in Go using Gin framework
//	@termsOfService	https://maxuanquang.dev/gin-social-network-tos

//	@contact.name	Quang Ma
//	@contact.url	https://www.linkedin.com/in/maxuanquang
//	@contact.email	maxuanquang@gmail.com

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@host
//	@BasePath	/api/v1

//	@securitydefinitions.oauth2.password	OAuth2Password
//	@tokenUrl								https://example.com/oauth/token
//	@scope.read								Grants read access
//	@scope.write							Grants write access
//	@scope.admin							Grants read and write access to administrative information
func main() {
	// Flags
	cfgPath := flag.String("cfg", "config.yml", "Path to config file for this service")

	// Load configs
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
