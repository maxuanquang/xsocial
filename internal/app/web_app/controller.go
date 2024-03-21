package web_app

import (
	// "context"
	"fmt"

	// "net/http"
	"net/http/pprof"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/maxuanquang/social-network/configs"
	"github.com/maxuanquang/social-network/internal/app/web_app/service"
	v1 "github.com/maxuanquang/social-network/internal/app/web_app/v1"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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
	// Init web services
	webService, err := service.NewWebService(cfg)
	if err != nil {
		return nil, err
	}

	// Init router
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Cookie"},
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		ExposeHeaders: []string{"Set-Cookie"},
	}))

	// Serve frontend static files
	router.Use(static.Serve("/", static.LocalFile("./web/build", true)))

	for _, version := range cfg.APIVersions {
		verXRouter := router.Group(fmt.Sprint("/api/" + version))
		// TODO: Automate applying new API version
		if version == "v1" {
			v1.AddAllRouter(verXRouter, webService)
		}
	}

	// Init other support tools
	initSwagger(router)
	initPprof(router)
	initPrometheus(router)

	return &WebController{
		webService: *webService,
		router:     router,
		port:       cfg.Port,
	}, nil
}

func initSwagger(router *gin.Engine) {
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

func initPprof(router *gin.Engine) {
	router.GET("/debug/pprof/", func(context *gin.Context) {
		pprof.Index(context.Writer, context.Request)
	})
	router.GET("/debug/pprof/profile", func(context *gin.Context) {
		pprof.Profile(context.Writer, context.Request)
	})
	router.GET("/debug/pprof/trace", func(context *gin.Context) {
		pprof.Trace(context.Writer, context.Request)
	})
}

func initPrometheus(router *gin.Engine) {
	handler := promhttp.Handler()
	router.GET("/metrics", func(context *gin.Context) {
		handler.ServeHTTP(context.Writer, context.Request)
	})
}
