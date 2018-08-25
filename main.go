package main

import (
		"github.com/gin-gonic/gin"
	"net/http"
	"io/ioutil"
	"log"
	"fmt"
	"github.com/gin-gonic/gin/binding"
	"strings"
	"encoding/json"
	"github.com/satori/go.uuid"
	"os"
)
func respondWithError(code int, message string,c *gin.Context) {
	resp := map[string]string{"error": message}
	c.JSON(code, resp)
	c.Abort()
}

func DummyMiddleware(c *gin.Context) {
	fmt.Println("Im a test!")
	c.Next()
}

func TokenAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// read from header
		token := c.Request.Header.Get("api_token")

		// read from post body
		//token := c.Request.FormValue("api_token")

		if token == "" {
			respondWithError(401, "API token required", c)
			return
		}
		if token != os.Getenv("API_TOKEN") {
			respondWithError(401, "Invalid API token", c)
			return
		}
		c.Next()
	}
}

func RequestIdMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		u, _ := uuid.NewV4()
		c.Writer.Header().Set("X-Request-Id", u.String())
		c.Next()
	}
}
func main() {
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
		authorized.POST("/search", search)
		// nested group
		//testing := authorized.Group("testing")
		//testing.GET("/analytics", analyticsEndpoint)
	}
	r.Run(":3001")
}

type ReverseRequest struct {
	Lat    string `json:"lat"`
	Lon string `json:"lon"`
}

type Geom struct {
	Type string `json:"type"`
	Coordinates []float64 `json:"-"`
}

type MapIrReverseResponse struct {
	Address string `json:"address"`
	PostalAddress string `json:"postal_address"`
	PostalCompact string `json:"postal_compact"`
	Country string `json:"country"`
	Province string `json:"province"`
	County string `json:"county"`
	City string `json:"city"`
	District string `json:"district"`
	Region string `json:"region"`
	Primary string `json:"primary"`
	Last string `json:"last"`
	Poi string `json:"poi"`
	Plaque string `json:"plaque"`
	PostalCode string `json:"postal_code"`
	Geom Geom `json:"geom"`
}

type Component struct {
	LongName string `json:"long_name"`
	ShortName string `json:"short_name"`
	Type string `json:"type"`
}

type TrafficZone struct {
	Name string `json:"name"`
	InCentral string `json:"in_central"`
	InEvenodd string `json:"in_evenodd"`
}

type Result struct {
	Components []Component `json:"components"`
	Address string `json:"address"`
	Locality string `json:"locality"`
	District string `json:"district"`
	Place string `json:"place"`
	City string `json:"city"`
	Province string `json:"province"`
	TrafficZone TrafficZone `json:"-"`
}

type CedarMapReverseResponse struct {
	Status string `json:"status"`
	Result Result `json:"-"`
}

type Location struct {
	Type string `json:"type"`
	Coordinates []float64 `json:"coordinates"`
}

type SearchRequest struct {
	Text string `json:"text"`
	Location Location `json:"location"`
}

type Message struct {
	status int
	body []byte
}

const (
	MapIrUrl = "https://map.ir/"
	CedarMapUrl = "http://carpino2.api.cedarmaps.com/"
	CedarMapAccessToken = "1b2d5cdc59a35285ac3934f254d2309d6e882000"
	MapIrApiKey = "eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiIsImp0aSI6Ijk1MjhjNjgwMGM5M2I1NmY2NmQ2YjI5ZWVlZjRmZmY3NjZhYjUxODIwNDJhMDE1YTUxOGIyMzZjNzFjNDQ4ZWMzYzRjZTlmNDM3MjFiYjAzIn0.eyJhdWQiOiJteWF3ZXNvbWVhcHAiLCJqdGkiOiI5NTI4YzY4MDBjOTNiNTZmNjZkNmIyOWVlZWY0ZmZmNzY2YWI1MTgyMDQyYTAxNWE1MThiMjM2YzcxYzQ0OGVjM2M0Y2U5ZjQzNzIxYmIwMyIsImlhdCI6MTUzMjM0NDY3OSwibmJmIjoxNTMyMzQ0Njc5LCJleHAiOjE1MzIzNDgyNzksInN1YiI6IiIsInNjb3BlcyI6W119.sUsXA3IQzgU-L-MQPk0XTCSQtbrUtVHWxBQ_ZNTn8VJ6kFcy-X5KogziNk_XNAbLc5E3L80XnQfHQ-54mcgCSOsZ4e7zPpBPbWWMpcQbOgJLJoG8jDGn46L-85aLo1DJNXphGboXILCy9p6AnLpwTkM2u1gBCb6f2FjB7JF1N9wmkU2NHm3ypG7Vg37J3PyCweLBI2l4vCxwVuSZTkHKhyOFXyTW0_Tn5mugRAHaV4ExJ1yMeMbcsfy6M73DOox3YWgSVuzh5hw4bgi37l5AB4eQR0nc71Aqx5NhZFoPs8FRxjM2pP1y52hIlZNIT1m2fBLNzVqv-kjF9gf6aGxo5A"
)

func reverse(c *gin.Context) {
	//var reverseRequest ReverseRequest
	latitude := c.Query("lat")
	longitude := c.Query("lon")

	url := MapIrUrl + fmt.Sprintf("reverse?lat=%v&lon=%v", latitude, longitude)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("accept", "application/json")

	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()

	if res.StatusCode >= 500 {
		// run down time job
	}

	reverseResponse := new(ReverseResponse)
	json.NewDecoder(res.Body).Decode(&reverseResponse)
	c.Header("Content-Type", "application/json; charset=utf-8")

	// TODO (change policy)
	result := string(reverseResponse.Address)

	r := Message{}
	r.body = []byte(result)
	r.status = res.StatusCode

	c.JSON(200, gin.H{
		"status":  r.status,
		"message": string(r.body),
	})
}

func search(c *gin.Context) {
	var sR SearchRequest
	c.ShouldBindWith(&sR, binding.JSON)
	log.Println()
	url := MapIrUrl + "search"
	out, _ := json.Marshal(sR)
	payload := strings.NewReader(string(out))
	req, _ := http.NewRequest("POST", url, payload)
	req.Header.Add("accept", "application/json")
	req.Header.Add("x-api-key", MapIrApiKey)
	req.Header.Add("Content-Type", "application/json")

	res, _ := http.DefaultClient.Do(req)

	if res.StatusCode >= 500 {
		// run down time job
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	c.Header("Content-Type", "application/json; charset=utf-8")
	c.String(200, string(body))
}