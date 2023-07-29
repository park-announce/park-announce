package service

import (
	"github.com/park-announce/pa-api/client"
	"github.com/park-announce/pa-api/repository"
)

type SocketService struct {
}

type UserService struct {
	userRepository repository.UserRepository
	redisClient    client.IRedisClient
	httpClient     client.IHttpClient
}

func NewSocketService() SocketService {
	return SocketService{}
}

func NewUserServiceWithHttpClient(redisClient client.IRedisClient, httpClient client.IHttpClient, userRepository repository.UserRepository) UserService {
	return UserService{redisClient: redisClient, httpClient: httpClient, userRepository: userRepository}
}
