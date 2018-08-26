package main

import (
	"github.com/gin-gonic/gin"
	"fmt"
	"net/http"
	"encoding/json"
	"geo/Redis"
	"time"
	"encoding/binary"
	"bytes"
	"github.com/satori/go.uuid"
)

func reverse(c *gin.Context) {
	//var reverseRequest ReverseRequest
	latitude := c.Query("lat")
	longitude := c.Query("lon")

	url := CedarMapUrl + fmt.Sprintf("v1/geocode/cedarmaps.streets/%v,%v?access_token=%v",
		latitude, longitude, CedarMapAccessToken)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("accept", "application/json")

	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()

	if res.StatusCode >= 100 {
		var bin_buf bytes.Buffer
		errorObject := ErrorObject{Status:res.StatusCode, Timestamp:time.Now(), Url:c.Request.URL.Path}
		binary.Write(&bin_buf, binary.BigEndian, errorObject)
		u, _ := uuid.NewV4()
		Redis.Set(u.String(), bin_buf.Bytes())

	}

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
	r.body = []byte(result)
	r.status = res.StatusCode

	c.JSON(200, gin.H{
		"result": string(r.body),
	})
}

func search(c *gin.Context) {
	name := c.Query("name")
	latitude := c.Query("lat")
	longitude := c.Query("lon")
	distance := c.Query("distance")

	url := CedarMapUrl + fmt.
		Sprintf("v1/geocode/cedarmaps.streets/%v.json?access_token=%v&location=%v,%v&distance=%v",
			name, CedarMapAccessToken, latitude, longitude, distance)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("accept", "application/json")

	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()

	if res.StatusCode >= 500 {
		//
		// run down time job
	}

	cedarSearchResponse := new(CedarMapSearchResponse)
	json.NewDecoder(res.Body).Decode(&cedarSearchResponse)

	georgeSearchResponse := GeorgeSearchResponse{}
	resultString := make([]string, 0)
	for _, value := range cedarSearchResponse.Results {
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

			//locality part
			//if len(value.Components.Localities) > 0 {
			//	if localityName := value.Components.Localities[0]; localityName != "" {
			//		if len(volunteerValue) > 0 && volunteerValue != "" {
			//			volunteerValue += ", "
			//		}
			//		volunteerValue += localityName
			//	}
			//}

			// address part
			if len(volunteerValue) > 0 && volunteerValue != "" {
				volunteerValue += ", "
			}
			volunteerValue += value.Address

			resultString = append(resultString, volunteerValue)
		} else {
			continue
		}
	}
	georgeSearchResponse.Result = resultString
	c.Header("Content-Type", "application/json; charset=utf-8")
	c.JSON(res.StatusCode, georgeSearchResponse)
}