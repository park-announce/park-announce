package service

import (
	"github.com/park-announce/pa-api/client"
	"github.com/park-announce/pa-api/repository"
)

type SocketService struct {
	redisClient client.IRedisClient
}

type UserService struct {
	userRepository repository.UserRepository
	redisClient    client.IRedisClient
	httpClient     client.IHttpClient
}

type CorporationService struct {
	corporationRepository repository.CorporationRepository
}

func NewSocketService(redisClient client.IRedisClient) SocketService {
	return SocketService{redisClient: redisClient}
}

func NewUserServiceWithHttpClient(redisClient client.IRedisClient, httpClient client.IHttpClient, userRepository repository.UserRepository) UserService {
	return UserService{redisClient: redisClient, httpClient: httpClient, userRepository: userRepository}
}

func NewCorporationService(corporationRepository repository.CorporationRepository) CorporationService {
	return CorporationService{corporationRepository: corporationRepository}
}
