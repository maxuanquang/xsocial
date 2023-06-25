package web_app

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/maxuanquang/social-network/configs"
	v1 "github.com/maxuanquang/social-network/internal/app/web_app/v1"
	"github.com/maxuanquang/social-network/internal/app/web_app/service"
)

type WebController struct {
	webService service.WebService
	router     *gin.Engine
	port       int
}

func (wc *WebController) Run() {
	wc.router.Run(fmt.Sprintf(":%d", wc.port))
}

// NewWebController creates new WebController
func NewWebController(cfg *configs.WebConfig) (*WebController, error) {
	// Intialize webService
	webService, err := service.NewWebService(cfg)
	if err != nil {
		return nil, err
	}

	// Initialize router
	router := gin.Default()
	for _, version := range cfg.APIVersions {
		verXRouter := router.Group(version)
		if version == "v1" { // TODO: Automate this when a new vision is added
			v1.AddAllRouter(verXRouter, webService)
		}
	}

	// Create webController
	webController := WebController{
		webService: *webService,
		router:     router,
		port:       cfg.Port,
	}

	return &webController, nil
}
