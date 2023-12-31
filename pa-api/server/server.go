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
		socketHandler := handler.NewSocketHandler(service.NewSocketService(client.NewRedisClientFactory("redis:6379", "").GetRedisClient()))
		socketHandler.HandleSocketConnection(ctx, hub)
	})

	server.POST("/google/oauth2/token", func(ctx *gin.Context) {
		dbClient := dbClientFactory.NewDBClient()
		userHandler := handler.NewUserHandler(service.NewUserServiceWithHttpClient(client.NewRedisClientFactory("redis:6379", "").GetRedisClient(), client.NewHttpClientFactory().GetHttpClient(), repository.NewUserRepository(dbClient)))
		userHandler.HandleOAuth2GoogleToken(ctx)
	})

	server.POST("/preregisterations/google", func(ctx *gin.Context) {
		dbClient := dbClientFactory.NewDBClient()
		userHandler := handler.NewUserHandler(service.NewUserServiceWithHttpClient(client.NewRedisClientFactory("redis:6379", "").GetRedisClient(), client.NewHttpClientFactory().GetHttpClient(), repository.NewUserRepository(dbClient)))
		userHandler.HandleGetGuidForGoogleRegistration(ctx)
	})

	server.POST("/preregisterations", func(ctx *gin.Context) {
		dbClient := dbClientFactory.NewDBClient()
		userHandler := handler.NewUserHandler(service.NewUserServiceWithHttpClient(client.NewRedisClientFactory("redis:6379", "").GetRedisClient(), client.NewHttpClientFactory().GetHttpClient(), repository.NewUserRepository(dbClient)))
		userHandler.HandlePreRegister(ctx)
	})

	server.POST("/registerationverifications", func(ctx *gin.Context) {
		dbClient := dbClientFactory.NewDBClient()
		userHandler := handler.NewUserHandler(service.NewUserServiceWithHttpClient(client.NewRedisClientFactory("redis:6379", "").GetRedisClient(), client.NewHttpClientFactory().GetHttpClient(), repository.NewUserRepository(dbClient)))
		userHandler.HandlePreRegisterVerification(ctx)
	})

	server.POST("/registerations", func(ctx *gin.Context) {
		dbClient := dbClientFactory.NewDBClient()
		userHandler := handler.NewUserHandler(service.NewUserServiceWithHttpClient(client.NewRedisClientFactory("redis:6379", "").GetRedisClient(), client.NewHttpClientFactory().GetHttpClient(), repository.NewUserRepository(dbClient)))
		userHandler.HandleRegister(ctx)
	})

	server.PUT("/corporation/locations/:id", func(ctx *gin.Context) {
		dbClient := dbClientFactory.NewDBClient()
		corporationHandler := handler.NewCorporationHandler(service.NewCorporationService(repository.NewCorporationRepository(dbClient)))
		corporationHandler.HandleCorporationLocationUpdate(ctx)
	})

	server.POST("/corporation/token", func(ctx *gin.Context) {
		dbClient := dbClientFactory.NewDBClient()
		corporationHandler := handler.NewCorporationHandler(service.NewCorporationService(repository.NewCorporationRepository(dbClient)))
		corporationHandler.HandleCorporationToken(ctx)
	})

	server.POST("/corporation/users", func(ctx *gin.Context) {
		dbClient := dbClientFactory.NewDBClient()
		corporationHandler := handler.NewCorporationHandler(service.NewCorporationService(repository.NewCorporationRepository(dbClient)))
		corporationHandler.HandleCorporationUserInsert(ctx)
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
