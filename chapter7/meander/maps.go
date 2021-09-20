package meander

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
)

// Probably, the code should be belongs to other package
var MapAPIKey string

type Geocode struct {
	Lat float64 `json:"latitude"`
	Lng float64 `json:"longitude"`
}

type geocodeResponse struct {
	Results []*Geocode `json:"data"`
}

func (g *Geocode) Public() interface{} {
	return map[string]interface{}{
		"lat": g.Lat,
		"lng": g.Lng,
	}
}

type MapQuery struct {
	Address string
}

func (q *MapQuery) find() (*geocodeResponse, error) {
	u := "http://api.positionstack.com/v1/forward"
	vals := make(url.Values)
	vals.Set("query", q.Address)
	vals.Set("access_key", MapAPIKey)
	fmt.Println(u + "?" + vals.Encode())
	res, err := http.Get(u + "?" + vals.Encode())
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	var response geocodeResponse
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, err
	}
	return &response, nil
}

func (q *MapQuery) Run() []interface{} {
	response, err := q.find()
	location := make([]interface{}, 1)
	if err != nil {
		log.Println("Failed to query geocode:", err)
		return nil
	}
	for i, result := range response.Results {
		location[i] = result
		break
	}
	return location
}
