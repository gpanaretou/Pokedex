package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"os"
	"strings"
	"time"

	"github.com/gpanaretou/Pokedex/internal/pokeapi"
	"github.com/gpanaretou/Pokedex/internal/pokecache"
)

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

func (cc cliCommand) Help() map[string]cliCommand {
	return map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    nil,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    nil,
		},
		"map": {
			name:        "map",
			description: "Show 20 locations, subsequent uses show the next 20 locations",
			callback:    nil,
		},
		"mapb": {
			name:        "mapb",
			description: "Show the previous 20 locations",
			callback:    nil,
		},
		"explore": {
			name:        "explore",
			description: "Explore a location and view the available pokemon there",
			callback:    nil,
		},
	}
}

func (cc cliCommand) Exit() {
	os.Exit(0)
}

func (cc cliCommand) Map(c *pokeapi.Config, cache *pokecache.Cache) error {
	mapLocations := pokeapi.MapLocations{}
	entry, ok := cache.Get(c.Next)
	if ok {
		err := json.Unmarshal(entry, &mapLocations)
		if err != nil {
			return err
		}

	} else {
		data, err := pokeapi.GetMapAreas(c, "map")
		if err != nil {
			return err
		}
		cache.Add(c.Next, data)

		err = json.Unmarshal(data, &mapLocations)
		if err != nil {
			return err
		}
	}

	c.Next = mapLocations.Next
	if mapLocations.Previous != nil {
		c.Previous = *mapLocations.Previous
	}

	for _, location := range mapLocations.Results {
		fmt.Println(location.Name)
	}

	return nil
}

func (cc cliCommand) Mapb(c *pokeapi.Config, cache *pokecache.Cache) error {
	mapLocations := pokeapi.MapLocations{}
	entry, ok := cache.Get(c.Previous)
	if ok {
		err := json.Unmarshal(entry, &mapLocations)
		if err != nil {
			return err
		}

	} else {
		data, err := pokeapi.GetMapAreas(c, "mapb")
		if err != nil {
			return err
		}

		err = json.Unmarshal(data, &mapLocations)
		if err != nil {
			return err
		}
		cache.Add(c.Previous, data)
	}

	c.Next = mapLocations.Next

	if mapLocations.Previous != nil {
		c.Previous = *mapLocations.Previous
	}

	for _, location := range mapLocations.Results {
		fmt.Println(location.Name)
	}

	return nil
}

func (cc cliCommand) Explore(area string, cache *pokecache.Cache) error {
	baseUrl := "https://pokeapi.co/api/v2/location-area/"

	areaInfo := pokeapi.Area{}
	requestUrl := baseUrl + area
	_, ok := cache.Get(requestUrl)
	if ok {
		fmt.Println("area explore in cahce!")
		return nil
	} else {
		data, err := pokeapi.ExploreArea(area)
		if err != nil {
			return nil
		}

		err = json.Unmarshal(data, &areaInfo)
		if err != nil {
			return err
		}
		cache.Add(requestUrl, data)
	}

	for _, value := range areaInfo.PokemonEncounters {
		fmt.Printf("  - %s\n", value.Pokemon.Name)
	}
	return nil
}

func (cc cliCommand) Catch(pokemon string, cache *pokecache.Cache) (Pokemon pokeapi.Pokemon, caught bool, err error) {
	baseCaptureChance := rand.Float32() * 1000
	pokemonInfo := pokeapi.Pokemon{}

	entry, ok := cache.Get(pokemon)
	if ok {
		err := json.Unmarshal(entry, &pokemonInfo)
		if err != nil {
			return pokemonInfo, false, err
		}
	} else {
		data, err := pokeapi.GetPokemonExperience(pokemon)
		if err != nil {
			return pokemonInfo, false, err
		}

		err = json.Unmarshal(data, &pokemonInfo)
		if err != nil {
			return pokemonInfo, false, err
		}
		cache.Add(pokemon, data)
	}

	isCaught := false

	fmt.Printf("Throwning a pokeball at %s...\n", pokemonInfo.Name)
	if baseCaptureChance > float32(pokemonInfo.BaseExperience) {
		isCaught = true
		fmt.Printf("%s was caught!\n", pokemonInfo.Name)
	} else {
		fmt.Printf("%s escaped!\n", pokemonInfo.Name)
	}

	return pokemonInfo, isCaught, nil
}

func main() {
	config := pokeapi.Config{Next: "https://pokeapi.co/api/v2/location-area", Previous: ""}
	cc := cliCommand{}
	pokedex := map[string]pokeapi.Pokemon{}

	t := time.Duration(time.Minute)
	cache := pokecache.NewCache(t * 5)

	for {

		fmt.Print("Pokedex > ")
		scanner := bufio.NewScanner(os.Stdin)
		err := scanner.Scan()

		if !err {
			fmt.Println("Something went wrong when trying to read from value.")
			os.Exit(1)
		}

		args := strings.Split(scanner.Text(), " ")
		command := args[0]

		switch command {
		case "help":
			fmt.Println("Usage:")
			commands := cc.Help()
			for _, command := range commands {
				fmt.Printf("%s: %s\n", command.name, command.description)
			}
		case "map":
			err := cc.Map(&config, &cache)
			if err != nil {
				fmt.Println(err)
			}
		case "mapb":
			err := cc.Mapb(&config, &cache)
			if err != nil {
				fmt.Println(err)
			}
		case "explore":
			if len(args) < 2 {
				fmt.Println("incorrect usage of 'explore', try 'explore area'")
				continue
			}
			err := cc.Explore(args[1], &cache)
			if err != nil {
				fmt.Println(err)
			}
		case "catch":
			if len(args) < 2 {
				fmt.Println("incorrect usage of 'catch', try 'explore pokemon'")
				continue
			}
			name := strings.ToLower(args[1])
			Pokemon, isCaught, err := cc.Catch(name, &cache)
			if err != nil {
				fmt.Printf("error: %v\n", err)
			}

			_, ok := pokedex[name]
			if !ok && isCaught {
				fmt.Printf("%v added to pokedex!\n", name)
				pokedex[name] = Pokemon
			}

		case "pokedex":
			fmt.Println("Your Pokedex:")
			for pokemon, _ := range pokedex {
				fmt.Printf("\t- %v\n", pokemon)
			}

		case "exit":
			cc.Exit()
		default:
			fmt.Printf("'%s' command does not exist, use 'help' to see a list of commands\n", command)
		}
	}
}
