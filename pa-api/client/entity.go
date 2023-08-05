package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/go-redis/redis"
)

type IHttpClient interface {
	Get(url string) IHttpClient
	PostJson(url string, data interface{}) IHttpClient
	PutJson(url string, data interface{}) IHttpClient
	PostUrlEncoded(url string, data url.Values) IHttpClient
	EndStruct(response interface{}) error
}

type HttpClient struct {
	client   *http.Client
	response *http.Response
	err      error
	inTrx    bool
}

func (client *HttpClient) Get(url string) IHttpClient {

	if client.inTrx {
		panic(errors.New("client is in trx"))
	}

	client.inTrx = true

	res, err := client.client.Get(url)
	if err != nil {
		client.err = errors.New(fmt.Sprintf("error in Get-> %s", err.Error()))
		return client
	}

	client.response = res
	return client
}

func (client *HttpClient) PostJson(url string, data interface{}) IHttpClient {

	if client.inTrx {
		panic(errors.New("client is in trx"))
	}

	client.inTrx = true

	dataBytes, err := json.Marshal(data)

	if err != nil {
		client.err = errors.New(fmt.Sprintf("error in PostJson marshal-> %s", err.Error()))
		return client
	}

	reader := strings.NewReader(string(dataBytes))

	res, err := client.client.Post(url, "application/json", reader)

	if err != nil {
		client.err = errors.New(fmt.Sprintf("error in Post-> %s", err.Error()))
		return client
	}

	client.response = res
	return client
}

func (client *HttpClient) PutJson(url string, data interface{}) IHttpClient {

	if client.inTrx {
		panic(errors.New("client is in trx"))
	}

	client.inTrx = true

	dataBytes, err := json.Marshal(data)

	if err != nil {
		client.err = errors.New(fmt.Sprintf("error in PutJson marshal-> %s", err.Error()))
		return client
	}

	reader := strings.NewReader(string(dataBytes))

	request, err := http.NewRequest(http.MethodPut, url, reader)

	if err != nil {
		client.err = errors.New(fmt.Sprintf("error in PutJson http.NewRequest-> %s", err.Error()))
		return client
	}
	request.Header.Set("Content-Type", "application/json")

	res, err := client.client.Do(request)

	if err != nil {
		client.err = errors.New(fmt.Sprintf("error in Put-> %s", err.Error()))
		return client
	}

	client.response = res
	return client
}

func (client *HttpClient) PostUrlEncoded(url string, data url.Values) IHttpClient {

	if client.inTrx {
		panic(errors.New("client is in trx"))
	}

	client.inTrx = true

	reader := strings.NewReader(string(data.Encode()))

	res, err := client.client.Post(url, "application/x-www-form-urlencoded", reader)

	if err != nil {
		client.err = errors.New(fmt.Sprintf("error in response decoding-> %s", err.Error()))
		return client
	}

	client.response = res
	return client
}

func (client *HttpClient) EndStruct(response interface{}) error {

	if !client.inTrx {
		panic(errors.New("client is not in trx"))
	}

	client.inTrx = false

	if client.err != nil {
		return client.err
	}

	client.err = nil

	if client.response.StatusCode != http.StatusOK {
		defer client.response.Body.Close()
		resp, err := ioutil.ReadAll(client.response.Body)

		if err != nil {
			fmt.Println("ioutil.ReadAll error ->", err.Error())
		} else {
			fmt.Println("client.response.Body ->", string(resp))
		}

		return errors.New(fmt.Sprintf("response.StatusCode is not OK -> %d", client.response.StatusCode))
	}

	defer client.response.Body.Close()

	err := json.NewDecoder(client.response.Body).Decode(response)

	if err != nil {
		return errors.New(fmt.Sprintf("error in response decoding-> %s", err.Error()))
	}

	return err
}

func NewHttpClient() IHttpClient {
	return &HttpClient{client: &http.Client{}}
}

type IRedisClient interface {
	HMGet(key string, fields ...string) (string, error)
	HMSet(key string, fields map[string]interface{}) (string, error)
	HDel(key string, fields ...string) (int64, error)
	HIncrBy(key, field string, incr int64) (int64, error)
	HGet(key, field string) (string, error)
	HSet(key, field string, value interface{}) (bool, error)
	Set(key string, value interface{}, expiration time.Duration) (string, error)
	Get(key string) (string, error)
}

type RedisClient struct {
	client *redis.Client
}

func NewRedisClient(address string, password string) IRedisClient {
	return &RedisClient{client: redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password, // no password set
		DB:       0,        // use default DB
	})}
}

func (c *RedisClient) HMGet(key string, fields ...string) (string, error) {
	var result string = ""
	response, err := c.client.HMGet(key, fields...).Result()
	if err == nil && response != nil && len(response) > 0 && response[0] != nil {
		result = response[0].(string)
	}

	return result, err
}

func (c *RedisClient) HMSet(key string, fields map[string]interface{}) (string, error) {
	var result string = ""

	response, err := c.client.HMSet(key, fields).Result()

	if err == nil {
		result = response
	}

	return result, err
}

func (c *RedisClient) HDel(key string, fields ...string) (int64, error) {
	var result int64 = 0

	response, err := c.client.HDel(key, fields...).Result()

	if err == nil {
		result = response
	}

	return result, err
}

func (c *RedisClient) HIncrBy(key, field string, incr int64) (int64, error) {
	var result int64 = 0

	response, err := c.client.HIncrBy(key, field, incr).Result()

	if err == nil {
		result = response
	}

	return result, err
}

func (c *RedisClient) HGet(key string, field string) (string, error) {
	var result string = ""
	result, err := c.client.HGet(key, field).Result()

	if err == redis.Nil {
		return result, nil
	}

	return result, err
}

func (c *RedisClient) HSet(key string, field string, value interface{}) (bool, error) {
	var result bool = false

	result, err := c.client.HSet(key, field, value).Result()

	return result, err
}

func (c *RedisClient) Get(key string) (string, error) {
	var result string = ""
	result, err := c.client.Get(key).Result()

	if err == redis.Nil {
		return result, nil
	}

	return result, err
}

func (c *RedisClient) Set(key string, value interface{}, expiration time.Duration) (string, error) {
	var result string = ""

	result, err := c.client.Set(key, value, expiration).Result()

	return result, err
}
