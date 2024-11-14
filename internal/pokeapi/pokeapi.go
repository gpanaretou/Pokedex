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

type Area struct {
	ID                int    `json:"id"`
	Name              string `json:"name"`
	GameIndex         int    `json:"game_index"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
		VersionDetails []struct {
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
			MaxChance        int `json:"max_chance"`
			EncounterDetails []struct {
				MinLevel        int   `json:"min_level"`
				MaxLevel        int   `json:"max_level"`
				ConditionValues []any `json:"condition_values"`
				Chance          int   `json:"chance"`
				Method          struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				} `json:"method"`
			} `json:"encounter_details"`
		} `json:"version_details"`
	} `json:"pokemon_encounters"`
}

type Pokemon struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	BaseExperience int    `json:"base_experience"`
	Height         int    `json:"height"`
	IsDefault      bool   `json:"is_default"`
	Weight         int    `json:"weight"`
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

func ExploreArea(area string) ([]byte, error) {
	baseUrl := "https://pokeapi.co/api/v2/location-area/"
	requestUrl := baseUrl + area

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

func GetPokemonExperience(pokemon string) ([]byte, error) {
	baseUrl := "https://pokeapi.co/api/v2/pokemon/"
	requestUrl := baseUrl + pokemon

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
