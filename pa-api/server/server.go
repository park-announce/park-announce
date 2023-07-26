package server

import (
	"github.com/gin-gonic/gin"
	"github.com/park-announce/pa-api/client"
	"github.com/park-announce/pa-api/handler"
	"github.com/park-announce/pa-api/middleware"
	"github.com/park-announce/pa-api/service"
)

func NewServer(hub *service.SocketHub) *gin.Engine {

	gin.SetMode(gin.ReleaseMode)
	server := gin.New()
	AddDefaultMiddlewaresToEngine(server)

	server.GET("/socket/connect", func(ctx *gin.Context) {

		colorHandler := handler.NewSocketHandler(service.NewSocketService())
		colorHandler.HandleSocketConnection(ctx, hub)
	})

	server.POST("/google/oauth2/token", func(ctx *gin.Context) {
		userHandler := handler.NewUserHandler(service.NewUserServiceWithHttpClient(client.NewRedisClientFactory("redis:6379", "").GetRedisClient(), client.NewHttpClientFactory().GetHttpClient()))
		userHandler.HandleOAuth2Google(ctx)
	})

	return server
}

func AddDefaultMiddlewaresToEngine(server *gin.Engine) {
	//engine.Use(secure.Secure(secure.Options))
	server.Use(gin.Logger())
	server.Use(gin.Recovery())
	server.Use(middleware.UseUserMiddleware())
}
