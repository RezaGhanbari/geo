package main

type Geom struct {
	Type        string   `json:"type"`
	Coordinates []string `json:"coordinates"`
}

type YReverseResponse struct {
	Address        string `json:"address"`
	PostalAddress  string `json:"postal_address"`
	AddressCompact string `json:"address_compact"`
	Country        string `json:"country"`
	Province       string `json:"province"`
	County         string `json:"county"`
	City           string `json:"city"`
	District       string `json:"district"`
	Region         string `json:"region"`
	Primary        string `json:"primary"`
	Last           string `json:"last"`
	Poi            string `json:"poi"`
	Plaque         string `json:"plaque"`
	PostalCode     string `json:"postal_code"`
	Geom           Geom   `json:"-"`
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

type XMapReverseResponse struct {
	Status string `json:"status"`
	Result Result `json:"result"`
}

type Location struct {
	Type        string    `json:"type"`
	Coordinates []float64 `json:"coordinates"`
}

type Message struct {
	Status int `json:"status"`
	Body   []byte `json:"body"`
}

type BB struct {
	NE string `json:"ne"`
	SW string `json:"sw"`
}

type XSearchLocation struct {
	BB     BB     `json:"bb"`
	Center string `json:"center"`
}

type XSearchComponent struct {
	Country    string   `json:"country"`
	Province   string   `json:"province"`
	City       string   `json:"city"`
	Districts  []string `json:"districts"`
	Localities []string `json:"localities"`
}

type XSearchResult struct {
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

type YCoordinate struct {
	Lat string `json:"lat"`
	Lon string `json:"lon"`
}

type YValue struct {
	Text       string          `json:"text"`
	Title      string          `json:"title"`
	Address    string          `json:"address"`
	Province   string          `json:"province"`
	City       string          `json:"city"`
	Type       string          `json:"type"`
	FClass     string          `json:"FClass"`
	Coordinate YCoordinate `json:"-"`
}

type GeorgeSearchResponse struct {
	Result []SearchResponse `json:"result"`
}

type ErrorLogger struct {
	Timestamp string `json:"timestamp"`
	Url       string `json:"url"`
	Status    string `json:"status"`
	ClientId  string `json:"client_id"`
	GcmToken  string `json:"gcm_token"`
}
