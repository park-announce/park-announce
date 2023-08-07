package main

import (
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/go-redis/redis"
	_ "github.com/lib/pq"

	"github.com/park-announce/pa-api/background"
	"github.com/park-announce/pa-api/service"
)

func main() {

	httpServerQuit := make(chan os.Signal, 1)
	signal.Notify(httpServerQuit, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	redisLockQuit := make(chan os.Signal, 1)
	signal.Notify(redisLockQuit, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	redisHeartBeatQuit := make(chan os.Signal, 1)
	signal.Notify(redisHeartBeatQuit, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	kafkaConsumerQuit := make(chan os.Signal, 1)
	signal.Notify(kafkaConsumerQuit, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	redisLockObtained := make(chan bool)
	redisLockName := make(chan string)

	rdb := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	hub := service.NewSocketHub()

	go hub.Register()

	wgMain := &sync.WaitGroup{}
	wgMain.Add(5)

	background := background.NewBackgroundOperation()

	go background.GetGlobalInstanceId(wgMain, rdb, redisLockObtained, redisLockQuit, redisLockName)

	//main routine continious running after redis lock obtained
	lockName := <-redisLockName
	<-redisLockObtained

	//send heartbeat message to redis in spesific interval setting expire value for distributed lock key in redis
	go background.RedisHeartBeat(wgMain, rdb, redisHeartBeatQuit, lockName)

	//start server
	go background.StartServer(wgMain, hub, httpServerQuit)

	//consume message from client_response_{{instanceId}} topic and try to send message to related socket client
	//if there is no socket connection, sends this message to dead_letter_message topic
	go background.ConsumeClientResponse(wgMain, rdb, hub, kafkaConsumerQuit)

	//consume message from dead_letter_messages topic and try to send message to related socket client
	go background.ConsumeDeadLetterMessages(wgMain, rdb, hub, kafkaConsumerQuit)

	wgMain.Wait()
}
