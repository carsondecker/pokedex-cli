package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/carsondecker/pokedex-cli/pokeapi"
)

type cliCommand struct {
	name        string
	description string
	callback    func([]string) error
}

var commands map[string]cliCommand

// Avoids initialization cycle with commands that reference the commands map
func initCommands() {
	commands = map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"map": {
			name:        "map",
			description: "Shows the next page of location areas",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Shows the previous page of location areas",
			callback:    commandMapB,
		},
		"explore": {
			name:        "explore {area}",
			description: "Finds all the pokemon in given area",
			callback:    commandExplore,
		},
	}
}

func commandExit(_ []string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(_ []string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	for _, v := range commands {
		fmt.Printf("\t%s: %s\n", v.name, v.description)
	}
	return nil
}

func commandMap(_ []string) error {
	data, err := pokeapi.GetNextMapData()
	if err != nil {
		return err
	}
	for _, locationArea := range data.Results {
		fmt.Println(locationArea.Name)
	}
	return nil
}

func commandMapB(_ []string) error {
	data, err := pokeapi.GetPrevMapData()
	if err != nil {
		return err
	}
	for _, locationArea := range data.Results {
		fmt.Println(locationArea.Name)
	}
	return nil
}

func commandExplore(args []string) error {
	if len(args) == 0 {
		return errors.New("explore command requires location as argument")
	}
	data, err := pokeapi.GetAreaData(args[0])
	if err != nil {
		return err
	}
	for _, pokemonEncounter := range data.PokemonEncounters {
		fmt.Printf(" - %s\n", pokemonEncounter.Pokemon.Name)
	}
	return nil
}
