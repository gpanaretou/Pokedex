package pokeapi

import (
	"fmt"
	"io"
	"net/http"
)

type Config struct {
	Next     string
	Previous string
}

type MapLocations struct {
	Count    int     `json:"count"`
	Next     string  `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

func GetMapAreas(c *Config, command string) ([]byte, error) {
	var requestUrl string
	if command == "map" {
		requestUrl = c.Next
	} else {
		requestUrl = c.Previous
	}

	if command == "mapb" && c.Previous == "" {
		return []byte{}, fmt.Errorf("cannot get previous map areas")
	}

	res, err := http.Get(requestUrl)
	if err != nil {
		return []byte{}, err
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return []byte{}, err
	}

	return data, nil
}
