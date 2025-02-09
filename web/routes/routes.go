package routes

import (
	"context"
	"gobp/web/middleware"
	assests "gobp/web/ui/assets"
	"gobp/web/ui/handlers"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Initialize(ctx context.Context, l *zap.Logger) (router *gin.Engine) {
	l.Sugar().Info("Initializing logger")

	router = gin.Default()
	router.Use(gin.Recovery())
	//Assests and Tailwind
	router.StaticFS("/assets", http.FS(assests.AssestFS))

	router.LoadHTMLGlob("web/ui/templates/*")

	//secure group
	rSecure := router.Group("/sec")

	rSecure.Use(middleware.ContextMiddleware(ctx))
	rSecure.GET("/home", handlers.IndexHandler)
	// router.GET("/", func(c *gin.Context) { c.Redirect(http.StatusMovedPermanently, "sec/home") })
	rSecure.POST("/checkmail", handlers.IndexHandler)
	router.GET("/test", handlers.IndexHandler) // without middleware
	router.GET("/", handlers.IndexHandler)
	router.GET("/about", handlers.AboutHandler)
	router.GET("/message", handlers.MessageHandler)

	//auth group sets the context and calls auth middleware
	rAuth := router.Group("/auth")
	rAuth.Use(middleware.ContextMiddleware(ctx), middleware.AuthMiddleware(ctx))
	rAuth.POST("/gobp/deactivate/:id/:isactive", handlers.IndexHandler)

	for _, route := range router.Routes() {
		l.Sugar().Infof("Route: %s %s", route.Method, route.Path)
	}

	return router
}
