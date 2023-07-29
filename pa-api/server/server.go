package server

import (
	"github.com/gin-gonic/gin"
	"github.com/park-announce/pa-api/client"
	"github.com/park-announce/pa-api/factory"
	"github.com/park-announce/pa-api/handler"
	"github.com/park-announce/pa-api/middleware"
	"github.com/park-announce/pa-api/repository"
	"github.com/park-announce/pa-api/service"
)

func NewServer(hub *service.SocketHub) *gin.Engine {

	gin.SetMode(gin.ReleaseMode)
	server := gin.New()
	factory.InitFactoryList()

	dbClientFactory := repository.NewDbClientFactory("postgres", "postgres://park_announce:PosgresDb1591*@db/padb?sslmode=disable")
	AddDefaultMiddlewaresToEngine(server)

	server.GET("/socket/connect", func(ctx *gin.Context) {
		socketHandler := handler.NewSocketHandler(service.NewSocketService())
		socketHandler.HandleSocketConnection(ctx, hub)
	})

	server.POST("/google/oauth2/code", func(ctx *gin.Context) {
		dbClient := dbClientFactory.NewDBClient()
		userHandler := handler.NewUserHandler(service.NewUserServiceWithHttpClient(client.NewRedisClientFactory("redis:6379", "").GetRedisClient(), client.NewHttpClientFactory().GetHttpClient(), repository.NewUserRepository(dbClient)))
		userHandler.HandleOAuth2GoogleCode(ctx)
	})

	server.POST("/google/oauth2/token", func(ctx *gin.Context) {
		dbClient := dbClientFactory.NewDBClient()
		userHandler := handler.NewUserHandler(service.NewUserServiceWithHttpClient(client.NewRedisClientFactory("redis:6379", "").GetRedisClient(), client.NewHttpClientFactory().GetHttpClient(), repository.NewUserRepository(dbClient)))
		userHandler.HandleOAuth2GoogleToken(ctx)
	})

	server.POST("/google/oauth2/register", func(ctx *gin.Context) {
		dbClient := dbClientFactory.NewDBClient()
		userHandler := handler.NewUserHandler(service.NewUserServiceWithHttpClient(client.NewRedisClientFactory("redis:6379", "").GetRedisClient(), client.NewHttpClientFactory().GetHttpClient(), repository.NewUserRepository(dbClient)))
		userHandler.HandleOAuth2GoogleRegister(ctx)
	})

	return server
}

func AddDefaultMiddlewaresToEngine(server *gin.Engine) {
	//engine.Use(secure.Secure(secure.Options))
	server.Use(gin.Logger())
	server.Use(gin.Recovery())
	server.Use(middleware.CORSMiddleware())
	server.Use(middleware.UseUserMiddleware())
}
