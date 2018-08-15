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
)

func main() {
	// Creates a router without any middleware by default
	r := gin.Default()

	// Global middleware
	// Logger middleware will write the logs to gin.DefaultWriter even if you set with GIN_MODE=release.
	// By default
	// gin.DefaultWriter = os.Stdout
	r.Use(gin.Logger())

	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	r.Use(gin.Recovery())

	r.GET("/reverse", reverse)
	r.POST("/search/", search)
	// Authorization group
	// authorized := r.Group("/", AuthRequired())
	// exactly the same as:
	//authorized := r.Group("/")
	// per group middleware! in this case we use the custom created
	// AuthRequired() middleware just in the "authorized" group.
	//authorized.Use(AuthRequired())
	//{
	//	authorized.POST("/login", loginEndpoint)
	//	authorized.POST("/submit", submitEndpoint)
	//	authorized.POST("/read", readEndpoint)
	//
		// nested group
		//testing := authorized.Group("testing")
		//testing.GET("/analytics", analyticsEndpoint)
	//}

	// Listen and serve on 0.0.0.0:8080
	r.Run(":3001")
}

type ReverseRequest struct {
	Lat    string `json:"lat"`
	Lon string `json:"lon"`
}

type Location struct {
	Type string `json:"type"`
	Coordinates []float64 `json:"coordinates"`
}

type SearchRequest struct {
	Text string `json:"text"`
	Location Location `json:"location"`
}

const (
	MapIrUrl = "https://map.ir/"
	MapIrApiKey = "eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiIsImp0aSI6Ijk1MjhjNjgwMGM5M2I1NmY2NmQ2YjI5ZWVlZjRmZmY3NjZhYjUxODIwNDJhMDE1YTUxOGIyMzZjNzFjNDQ4ZWMzYzRjZTlmNDM3MjFiYjAzIn0.eyJhdWQiOiJteWF3ZXNvbWVhcHAiLCJqdGkiOiI5NTI4YzY4MDBjOTNiNTZmNjZkNmIyOWVlZWY0ZmZmNzY2YWI1MTgyMDQyYTAxNWE1MThiMjM2YzcxYzQ0OGVjM2M0Y2U5ZjQzNzIxYmIwMyIsImlhdCI6MTUzMjM0NDY3OSwibmJmIjoxNTMyMzQ0Njc5LCJleHAiOjE1MzIzNDgyNzksInN1YiI6IiIsInNjb3BlcyI6W119.sUsXA3IQzgU-L-MQPk0XTCSQtbrUtVHWxBQ_ZNTn8VJ6kFcy-X5KogziNk_XNAbLc5E3L80XnQfHQ-54mcgCSOsZ4e7zPpBPbWWMpcQbOgJLJoG8jDGn46L-85aLo1DJNXphGboXILCy9p6AnLpwTkM2u1gBCb6f2FjB7JF1N9wmkU2NHm3ypG7Vg37J3PyCweLBI2l4vCxwVuSZTkHKhyOFXyTW0_Tn5mugRAHaV4ExJ1yMeMbcsfy6M73DOox3YWgSVuzh5hw4bgi37l5AB4eQR0nc71Aqx5NhZFoPs8FRxjM2pP1y52hIlZNIT1m2fBLNzVqv-kjF9gf6aGxo5A"

)

func reverse(c *gin.Context) {
	//var reverseRequest ReverseRequest
	latitude := c.Query("lat")
	longitude := c.Query("lon")
	log.Println(latitude)
	log.Println(longitude)

	url := MapIrUrl + fmt.Sprintf("reverse?lat=%v&lon=%v", latitude, longitude)

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("accept", "application/json")
	req.Header.Add("x-api-key", MapIrApiKey)

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	c.Header("Content-Type", "application/json; charset=utf-8")
	c.String(200, string(body))
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

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	c.Header("Content-Type", "application/json; charset=utf-8")
	c.String(200, string(body))
}