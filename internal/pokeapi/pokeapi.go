package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Config struct {
	Next     string
	Previous string
}

func HelloPoke() string {
	return "nice"
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

func GetMapAreas(c *Config, command string) error {
	var requestUrl string
	if command == "map" {
		requestUrl = c.Next
	} else {
		requestUrl = c.Previous
	}

	if command == "mapb" && c.Previous == "" {
		return fmt.Errorf("cannot get previous map areas")
	}

	res, err := http.Get(requestUrl)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	mapLocations := MapLocations{}
	err = json.Unmarshal(data, &mapLocations)
	if err != nil {
		return err
	}

	for _, location := range mapLocations.Results {
		fmt.Println(location.Name)
	}

	c.Next = mapLocations.Next
	if mapLocations.Previous != nil {
		c.Previous = *mapLocations.Previous
	}

	return nil
}
