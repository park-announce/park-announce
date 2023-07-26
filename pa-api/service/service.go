package service

import (
	"github.com/park-announce/pa-api/client"
)

type SocketService struct {
}

type UserService struct {
	redisClient client.IRedisClient
	httpClient  client.IHttpClient
}

func NewSocketService() SocketService {
	return SocketService{}
}

func NewUserServiceWithHttpClient(redisClient client.IRedisClient, httpClient client.IHttpClient) UserService {
	return UserService{redisClient: redisClient, httpClient: httpClient}
}
