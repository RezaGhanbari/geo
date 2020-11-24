package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/satori/go.uuid"
	"net/http"
	"os"
	"time"
	"strings"
)

var mapName = os.Getenv("MAP_NAME")

func reverse(c *gin.Context) {
	if mapName == "" {
		mapName = "X"
	}

	latitude := c.Query("lat")
	longitude := c.Query("lon")
	ClientId := c.Request.Header.Get("x-client-id")

	if ClientId == "" {
		respondWithError(400, "credentials are not provided.", c)
		return
	}

	var url string
	if mapName == "X" {
		url = CedarMapUrl + fmt.Sprintf("...",
			latitude, longitude, CedarMapAccessToken)
	} else if mapName == "X" {
		url = MapIrUrl + fmt.Sprintf("...", latitude, longitude)
	}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("accept", "application/json")
	if mapName == "X" {
		req.Header.Add("x-api-key", MapIrApiKey)
	}
	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()

	// set log to redis if third party service is down.
	if res.StatusCode >= 500 {
		s := strings.Split(ClientId, ":")
		cId, gcmToken := s[0], s[1]
		errorObject := ErrorLogger{time.Now().String(), c.Request.URL.Path,
			res.Status, cId, gcmToken}

		b, _ := json.Marshal(errorObject)

		u, _ := uuid.NewV4()

		err := Set("ERROR_" + u.String(), b)
		if err != nil {
			panic(err)
		}
	}

	if mapName == "X" {
		reverseResponse := new(CedarMapReverseResponse)
		json.NewDecoder(res.Body).Decode(&reverseResponse)
		c.Header("Content-Type", "application/json; charset=utf-8")

		var mainAddress string
		if district := reverseResponse.Result.District; district != "" {

			// city part
			if city := reverseResponse.Result.City; city != "" {
				mainAddress += city
			}

			// district part
			if district := reverseResponse.Result.District; district != "" {
				if len(mainAddress) > 0 && mainAddress != "" {
					mainAddress += ", "
				}
				mainAddress += district
			}

			// place part
			if place := reverseResponse.Result.Place; place != "" {
				if len(mainAddress) > 0 && mainAddress != "" {
					mainAddress += ", "
				}
				mainAddress += place
			}

			// locality part
			if localityName := reverseResponse.Result.Locality; localityName != "" {
				if len(mainAddress) > 0 && mainAddress != "" {
					mainAddress += ", "
				}
				mainAddress += localityName
			}

			// address part
			if address := reverseResponse.Result.Address; address != "" {
				if len(mainAddress) > 0 && mainAddress != "" {
					mainAddress += ", "
				}
				mainAddress += address
			}
		}
		result := mainAddress
		r := Message{}
		r.Body = []byte(result)
		r.Status = res.StatusCode
		c.JSON(r.Status, gin.H{
			"result": string(r.Body),
		})
	} else if mapName == "Y" {
		reverseResponse := new(MapIrReverseResponse)
		json.NewDecoder(res.Body).Decode(&reverseResponse)
		c.Header("Content-Type", "application/json; charset=utf-8")
		r := Message{}
		r.Body = []byte(reverseResponse.AddressCompact)
		r.Status = res.StatusCode
		c.JSON(r.Status, gin.H{
			"result": string(r.Body),
		})
	}
}

type SearchResponse struct {
	Address string `json:"address"`
}

func search(c *gin.Context) {
	if mapName == "" {
		mapName = "Y"
	}
	name := c.Query("name")
	latitude := c.Query("lat")
	longitude := c.Query("lon")
	distance := c.Query("distance")
	clientId := c.Request.Header.Get("x-client-id")
	if clientId == "" {
		respondWithError(400, "Credentials are not provided.", c)
		return
	}

	if mapName == "X" {
		url := CedarMapUrl + fmt.
			Sprintf("...",
				name, CedarMapAccessToken, latitude, longitude, distance)

		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Add("accept", "application/json")

		res, _ := http.DefaultClient.Do(req)
		defer res.Body.Close()

		if res.StatusCode >= 500 {
			s := strings.Split(clientId, ":")
			cId, gcmToken := s[0], s[1]
			errorObject := ErrorLogger{time.Now().String(), c.Request.URL.Path,
				res.Status, cId, gcmToken}
			b, _ := json.Marshal(errorObject)
			u, _ := uuid.NewV4()
			err := Set("ERROR_" + u.String(), b)
			if err != nil {
				panic(err)
			}
		}
		cedarSearchResponse := new(CedarMapSearchResponse)
		json.NewDecoder(res.Body).Decode(&cedarSearchResponse)
		georgeSearchResponse := GeorgeSearchResponse{}

		resultString := make([]SearchResponse, 0)

		for _, value := range cedarSearchResponse.Results {
			searchResponse := SearchResponse{}
			if value.Address != "" {
				volunteerValue := ""
				// city part
				if city := value.Components.City; city != "" {
					volunteerValue += city
				}

				// district part
				if len(value.Components.Districts) > 0 {
					if district := value.Components.Districts[0]; district != "" {
						if len(volunteerValue) > 0 && volunteerValue != "" {
							volunteerValue += ", "
						}
						volunteerValue += district
					}
				}

				// address part
				if len(volunteerValue) > 0 && volunteerValue != "" {
					volunteerValue += ", "
				}
				volunteerValue += value.Address
				searchResponse.Address = volunteerValue

				resultString = append(resultString, searchResponse)
			} else {
				continue
			}
		}
		georgeSearchResponse.Result = resultString
		c.Header("Content-Type", "application/json; charset=utf-8")
		c.JSON(res.StatusCode, georgeSearchResponse)
	} else if mapName == "X" {
		r := Message{}
		r.Body = []byte("Not implemented")
		r.Status = 501

		c.JSON(r.Status, gin.H{
			"result": string(r.Body),
		})
	}
}
