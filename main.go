package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/gpanaretou/Pokedex/internal/pokeapi"
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

func (cc cliCommand) Map(c *pokeapi.Config) error {
	err := pokeapi.GetMapAreas(c, "map")
	if err != nil {
		return err
	}
	return nil
}

func (cc cliCommand) Mapb(c *pokeapi.Config) error {
	err := pokeapi.GetMapAreas(c, "mapb")
	if err != nil {
		return err
	}
	return nil
}

func main() {
	config := pokeapi.Config{Next: "https://pokeapi.co/api/v2/location-area", Previous: ""}
	cc := cliCommand{}

	for {

		fmt.Print("Pokedex > ")
		scanner := bufio.NewScanner(os.Stdin)
		err := scanner.Scan()

		if !err {
			fmt.Println("Something went wrong when trying to read from value.")
			os.Exit(1)
		}

		if scanner.Text() == "help" {
			fmt.Println()
			fmt.Println("Usage:")
			commands := cc.Help()
			for _, command := range commands {
				fmt.Printf("%s: %s\n", command.name, command.description)
			}
			fmt.Println()
		}
		if scanner.Text() == "map" {
			err := cc.Map(&config)
			if err != nil {
				fmt.Println(err)
			}
		}
		if scanner.Text() == "mapb" {
			err := cc.Mapb(&config)
			if err != nil {
				fmt.Println(err)
			}
		}

		if scanner.Text() == "exit" {
			cc.Exit()
		}
	}
}
