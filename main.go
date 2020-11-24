package main

import (
	"encoding/json"
	"github.com/cnjack/throttle"
	"github.com/gin-gonic/gin"
	"net"
	"net/http"
	"os"
	"strconv"
	"github.com/gin-contrib/cache"
	"github.com/gin-contrib/cache/persistence"
	"time"
)

//var (
//	RedisPool *redis.Pool
//)

const (
	MapIrUrl            = ""
	CedarMapUrl         = ""
	CedarMapAccessToken = ""
	MapIrApiKey         = ""
	requestId = "X-Request-Id"
)

type MessageError struct {
	Status int    `json:"status"`
	Body   string `json:"body"`
}

func main() {

	redisPool := Init()
	store := persistence.NewRedisCacheWithPool(redisPool, time.Minute)

	serverUrl := os.Getenv("SERVER")
	port := os.Getenv("PORT")
	if port == "" {
		port = "3001"
	}
	Limit := os.Getenv("LIMIT")

	ThrottlingLimit, err := strconv.ParseUint(Limit, 10, 64)
	if err != nil {
		ThrottlingLimit = 100
	}

	r := gin.Default()

	//releaseMode := os.Getenv("GIN_MODE")
	//if releaseMode == "DEBUG" {
	//gin.SetMode(gin.ReleaseMode)
	//} else {
	gin.SetMode(gin.DebugMode)
	//}
	r.Use(gin.Logger())
	r.Use(TokenAuthMiddleware())

	r.Use(gin.Recovery())

	authorized := r.Group("/")
	authorized.Use(RequestIdMiddleware())

	errorResp := MessageError{}
	errorResp.Body = "Too many requests, try later."
	errorResp.Status = 429

	outError, _ := json.Marshal(errorResp)
	authorized.Use(throttle.Policy(&throttle.Quota{
		Limit: ThrottlingLimit,
		Within: time.Minute,
	}, &throttle.Options{
		StatusCode: 429,
		Message: string(outError),
		IdentificationFunction: func(req *http.Request) string {
			if rI := req.Header.Get(requestId); rI != "" {
				return rI
			}

			ip, _, err := net.SplitHostPort(req.RemoteAddr)
			if err != nil {
				panic(err.Error())
			}
			return ip
		},
		KeyPrefix: "T-GEO",
		Disabled: false,
	}))
	{
		authorized.GET("/reverse", cache.CachePage(store, 10 * time.Second, reverse))
		authorized.GET("/search", cache.CachePage(store, 10 * time.Second, search))
	}
	r.Run(serverUrl + ":" + port)
}
