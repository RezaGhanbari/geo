package main

import (
	"fmt"
	"time"
	"os"
	"github.com/garyburd/redigo/redis"
	log "github.com/Sirupsen/logrus"

)

func PanicOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%v: $v", msg, err)
		panic(fmt.Sprintf("%v: %v", msg, err))
	}
}

func newRedisPool(server, password string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			if password != "" {
				if _, err := c.Do("AUTH", password); err != nil {
					c.Close()
					return nil, err
				}
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}


func RedisPing() {
	serverUrl := os.Getenv("REDIS_URL")
	password := os.Getenv("REDIS_PASSWORD")
	RedisPool = newRedisPool(serverUrl, password)
	c := RedisPool.Get()
	defer c.Close()
	pong, err := redis.String(c.Do("PING"))
	PanicOnError(err, "Cannot ping Redis")
	log.Infof("Redis Ping: %s", pong)
}

func RedisGet(str string) string{
	serverUrl := os.Getenv("REDIS_URL")
	password := os.Getenv("REDIS_PASSWORD")
	RedisPool = newRedisPool(serverUrl, password)
	c := RedisPool.Get()
	defer c.Close()

	ret, err := redis.String(c.Do("GET", str))
	if err != nil {
		panic(err)
	}
	return ret
}

func RedisSet(key string, value string) string {
	serverUrl := os.Getenv("REDIS_URL")
	password := os.Getenv("REDIS_PASSWORD")
	RedisPool = newRedisPool(serverUrl, password)
	c := RedisPool.Get()
	defer c.Close()

	ret, err := redis.String(c.Do("SET", key, value))
	if err != nil {
		panic(err)
	}
	return ret
}

func RedisSetObject(obj string, fileds ... ErrorObject) string {
	serverUrl := os.Getenv("REDIS_URL")
	password := os.Getenv("REDIS_PASSWORD")
	RedisPool = newRedisPool(serverUrl, password)
	c := RedisPool.Get()
	defer c.Close()

	ret, err := redis.String(c.Do("HMSET", obj, ))
	if err != nil {
		panic(err)
	}
	return ret
}