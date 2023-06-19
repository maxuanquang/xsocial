package web_app

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/maxuanquang/social-network/configs"
	v1 "github.com/maxuanquang/social-network/internal/app/web_app/v1"
	client_aap "github.com/maxuanquang/social-network/pkg/client/authen_and_post"
	// client_nf "github.com/maxuanquang/social-network/pkg/client/newsfeed"
	pb_aap "github.com/maxuanquang/social-network/pkg/types/proto/pb/authen_and_post"
	pb_nf "github.com/maxuanquang/social-network/pkg/types/proto/pb/newsfeed"
)

type WebService struct {
	aapClient pb_aap.AuthenticateAndPostClient
	nfClient  pb_nf.NewsfeedClient
}

type WebController struct {
	webService WebService
	router     *gin.Engine
	port       int
}

func (wc *WebController) Run() {
	wc.router.Run(fmt.Sprintf(":%d", wc.port))
}

// NewWebController creates new WebController
func NewWebController(cfg *configs.WebConfig) (*WebController, error) {
	// Intialize webService using grpc clients
	aapClient, err := client_aap.NewClient(cfg.AuthenticateAndPost.Hosts)
	if err != nil {
		return nil, err
	}

	// nfClient, err := client_nf.NewClient(cfg.Newsfeed.Hosts)
	// if err != nil {
	// 	return nil, err
	// }

	webService := WebService{
		aapClient: aapClient,
		// nfClient:  nfClient,
	}

	// Initialize router
	router := gin.Default()
	for _, version := range cfg.APIVersions {
		verXRouter := router.Group(version)
		if version == "v1" {  // TODO: Automate this when a new vision is added
			v1.AddAllRouter(verXRouter, webService.aapClient, webService.nfClient)
		}
	}

	// Create webController
	webController := WebController{
		webService: webService,
		router:     router,
		port:       cfg.Port,
	}

	return &webController, nil
}
