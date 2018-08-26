package main

import (
	"time"
	"github.com/garyburd/redigo/redis"
)

type ReverseRequest struct {
	Lat string `json:"lat"`
	Lon string `json:"lon"`
}

type Geom struct {
	Type        string    `json:"type"`
	Coordinates []float64 `json:"-"`
}

type MapIrReverseResponse struct {
	Address       string `json:"address"`
	PostalAddress string `json:"postal_address"`
	PostalCompact string `json:"postal_compact"`
	Country       string `json:"country"`
	Province      string `json:"province"`
	County        string `json:"county"`
	City          string `json:"city"`
	District      string `json:"district"`
	Region        string `json:"region"`
	Primary       string `json:"primary"`
	Last          string `json:"last"`
	Poi           string `json:"poi"`
	Plaque        string `json:"plaque"`
	PostalCode    string `json:"postal_code"`
	Geom          Geom   `json:"geom"`
}

type Component struct {
	LongName  string `json:"long_name"`
	ShortName string `json:"short_name"`
	Type      string `json:"type"`
}

type TrafficZone struct {
	Name      string `json:"name"`
	InCentral string `json:"in_central"`
	InEvenodd string `json:"in_evenodd"`
}

type Result struct {
	Components  []Component `json:"components"`
	Address     string      `json:"address"`
	Locality    string      `json:"locality"`
	District    string      `json:"district"`
	Place       string      `json:"place"`
	City        string      `json:"city"`
	Province    string      `json:"province"`
	TrafficZone TrafficZone `json:"-"`
}

type CedarMapReverseResponse struct {
	Status string `json:"status"`
	Result Result `json:"result"`
}

type Location struct {
	Type        string    `json:"type"`
	Coordinates []float64 `json:"coordinates"`
}

type SearchRequest struct {
	Text     string   `json:"text"`
	Location Location `json:"location"`
}

type Message struct {
	status int
	body   []byte
}

type BB struct {
	NE string `json:"ne"`
	SW string `json:"sw"`
}

type CedarSearchLocation struct {
	BB     BB     `json:"bb"`
	Center string `json:"center"`
}

type CedarSearchComponent struct {
	Country    string   `json:"country"`
	Province   string   `json:"province"`
	City       string   `json:"city"`
	Districts  []string `json:"districts"`
	Localities []string `json:"localities"`
}

type CedarSearchResult struct {
	Id         int                  `json:"id"`
	Name       string               `json:"name"`
	NameEn     string               `json:"name_en"`
	Type       string               `json:"type"`
	Location   CedarSearchLocation  `json:"location"`
	Address    string               `json:"address"`
	Components CedarSearchComponent `json:"components"`
}

type CedarMapSearchResponse struct {
	Status  string              `json:"status"`
	Results []CedarSearchResult `json:"results"`
}

type GeorgeSearchResponse struct {
	Result []string `json:"result"`
}

type ErrorObject struct {
	Timestamp time.Time `json:"timestamp,omitempty"`
	Url string
	Status int
}

type RedisStore struct {
	pool              *redis.Pool
	defaultExpiration time.Duration
}
