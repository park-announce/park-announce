package client

type IHttpClientFactory interface {
	GetHttpClient() IHttpClient
}

type HttpClientFactory struct {
}

func NewHttpClientFactory() IHttpClientFactory {
	return &HttpClientFactory{}
}

func (f *HttpClientFactory) GetHttpClient() IHttpClient {
	return NewHttpClient()
}

type IRedisClientFactory interface {
	GetRedisClient() IRedisClient
}

type RedisClientFactory struct {
	address  string
	password string
}

func NewRedisClientFactory(address string, password string) IRedisClientFactory {
	return &RedisClientFactory{address: address, password: password}
}

func (f *RedisClientFactory) GetRedisClient() IRedisClient {
	return NewRedisClient(f.address, f.password)
}
