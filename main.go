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
	MapIrUrl            = "https://map.ir/"
	CedarMapUrl         = "http://carpino2.api.cedarmaps.com/"
	CedarMapAccessToken = "1b2d5cdc59a35285ac3934f254d2309d6e882000"
	MapIrApiKey         = "eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiIsImp0aSI6Ijk1MjhjNjgwMGM5M2I1NmY2NmQ2YjI5ZWVlZjRmZmY3NjZhYjUxODIwNDJhMDE1YTUxOGIyMzZjNzFjNDQ4ZWMzYzRjZTlmNDM3MjFiYjAzIn0.eyJhdWQiOiJteWF3ZXNvbWVhcHAiLCJqdGkiOiI5NTI4YzY4MDBjOTNiNTZmNjZkNmIyOWVlZWY0ZmZmNzY2YWI1MTgyMDQyYTAxNWE1MThiMjM2YzcxYzQ0OGVjM2M0Y2U5ZjQzNzIxYmIwMyIsImlhdCI6MTUzMjM0NDY3OSwibmJmIjoxNTMyMzQ0Njc5LCJleHAiOjE1MzIzNDgyNzksInN1YiI6IiIsInNjb3BlcyI6W119.sUsXA3IQzgU-L-MQPk0XTCSQtbrUtVHWxBQ_ZNTn8VJ6kFcy-X5KogziNk_XNAbLc5E3L80XnQfHQ-54mcgCSOsZ4e7zPpBPbWWMpcQbOgJLJoG8jDGn46L-85aLo1DJNXphGboXILCy9p6AnLpwTkM2u1gBCb6f2FjB7JF1N9wmkU2NHm3ypG7Vg37J3PyCweLBI2l4vCxwVuSZTkHKhyOFXyTW0_Tn5mugRAHaV4ExJ1yMeMbcsfy6M73DOox3YWgSVuzh5hw4bgi37l5AB4eQR0nc71Aqx5NhZFoPs8FRxjM2pP1y52hIlZNIT1m2fBLNzVqv-kjF9gf6aGxo5A"
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

	releaseMode := os.Getenv("GIN_MODE")
	if releaseMode == "DEBUG" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}
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
