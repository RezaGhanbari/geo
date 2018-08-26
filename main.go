package main

import (
	"github.com/gin-gonic/gin"
	"os"

	"github.com/garyburd/redigo/redis"
)

var (
	RedisPool *redis.Pool
)

const (
	MapIrUrl            = "https://map.ir/"
	CedarMapUrl         = "http://carpino2.api.cedarmaps.com/"
	CedarMapAccessToken = "1b2d5cdc59a35285ac3934f254d2309d6e882000"
	MapIrApiKey         = "eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiIsImp0aSI6Ijk1MjhjNjgwMGM5M2I1NmY2NmQ2YjI5ZWVlZjRmZmY3NjZhYjUxODIwNDJhMDE1YTUxOGIyMzZjNzFjNDQ4ZWMzYzRjZTlmNDM3MjFiYjAzIn0.eyJhdWQiOiJteWF3ZXNvbWVhcHAiLCJqdGkiOiI5NTI4YzY4MDBjOTNiNTZmNjZkNmIyOWVlZWY0ZmZmNzY2YWI1MTgyMDQyYTAxNWE1MThiMjM2YzcxYzQ0OGVjM2M0Y2U5ZjQzNzIxYmIwMyIsImlhdCI6MTUzMjM0NDY3OSwibmJmIjoxNTMyMzQ0Njc5LCJleHAiOjE1MzIzNDgyNzksInN1YiI6IiIsInNjb3BlcyI6W119.sUsXA3IQzgU-L-MQPk0XTCSQtbrUtVHWxBQ_ZNTn8VJ6kFcy-X5KogziNk_XNAbLc5E3L80XnQfHQ-54mcgCSOsZ4e7zPpBPbWWMpcQbOgJLJoG8jDGn46L-85aLo1DJNXphGboXILCy9p6AnLpwTkM2u1gBCb6f2FjB7JF1N9wmkU2NHm3ypG7Vg37J3PyCweLBI2l4vCxwVuSZTkHKhyOFXyTW0_Tn5mugRAHaV4ExJ1yMeMbcsfy6M73DOox3YWgSVuzh5hw4bgi37l5AB4eQR0nc71Aqx5NhZFoPs8FRxjM2pP1y52hIlZNIT1m2fBLNzVqv-kjF9gf6aGxo5A"
)

func main() {
	port := os.Getenv("PORT")
	//mapName := os.Getenv("MAP_NAME")
	r := gin.Default()
	//TODO (GIN_MODE=release)
	r.Use(gin.Logger())
	r.Use(RequestIdMiddleware())
	r.Use(TokenAuthMiddleware())

	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	r.Use(gin.Recovery())

	authorized := r.Group("/")
	authorized.Use(TokenAuthMiddleware())
	{
		authorized.GET("/reverse", reverse)
		authorized.GET("/search", search)
	}
	r.Run(port)
}
