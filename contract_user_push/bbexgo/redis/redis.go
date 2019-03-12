package redis

import (
	"bbexgo/config"
	"bbexgo/log"

	// "fmt"
	"strconv"
	"sync"

	"github.com/go-redis/redis"
)

type Pub struct {
	redis.PubSub
}

var (
	Instance *redis.Client
	once     sync.Once
)

func GetInstance() *redis.Client {
	once.Do(func() {
		Instance = crateClient()
		log.Info("redis just do once")
	})
	return Instance
}

func PubMessage(ms interface{}) string {
	var res string
	switch ms.(type) {
	case *redis.Message:
		res = ms.(*redis.Message).Payload
		break
	case *redis.Subscription:
		res = ""
		break
	case *redis.Pong:
		res = ""
		break
	default:
		res = ""
		break
	}
	return res
}

func crateClient() *redis.Client {
	dbIndex := 0
	if config.Get("redis.DBIndex") != "" {
		index, err := strconv.Atoi(config.Get("redis.DBIndex"))
		if err != nil {
			log.Fatal(err)
		}
		dbIndex = int(index)
	}
	return redis.NewClient(&redis.Options{
		Addr:     config.Get("redis.addr"),
		Password: config.Get("redis.password"),
		DB:       dbIndex,
	})

}

func closeClient() {
	Instance.Close()
}
