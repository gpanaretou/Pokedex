package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
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
	}
}

func (cc cliCommand) Exit() {
	os.Exit(0)
}

func (cc cliCommand) Map(c *pokeapi.Config, cache *pokecache.Cache) error {
	mapLocations := pokeapi.MapLocations{}
	entry, ok := cache.Get(c.Next)
	if ok {
		fmt.Println("-In cache !!")
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
		fmt.Println("-In cache !!")
		err := json.Unmarshal(entry, &mapLocations)
		if err != nil {
			return err
		}

	} else {
		data, err := pokeapi.GetMapAreas(c, "mapb")
		if err != nil {
			return err
		}
		cache.Add(c.Previous, data)

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

func main() {
	config := pokeapi.Config{Next: "https://pokeapi.co/api/v2/location-area", Previous: ""}
	cc := cliCommand{}

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

		input := scanner.Text()

		switch input {
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
		case "exit":
			cc.Exit()
		default:
			fmt.Printf("'%s' command does not exist, use 'help' to see a list of commands\n", input)
		}
	}
}
