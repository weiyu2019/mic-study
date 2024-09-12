package internal

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
)

type RedisConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

var RedisClient *redis.Client

func InitRedis() {
	h := AppConf.RedisConfig.Host
	p := AppConf.RedisConfig.Port
	addr := fmt.Sprintf("%s:%d", h, p)
	fmt.Println("addr", addr)
	RedisClient = redis.NewClient(&redis.Options{
		Addr: addr,
	})
	ping := RedisClient.Ping(context.Background())
	fmt.Println(ping)
	fmt.Println("redis init")
}
