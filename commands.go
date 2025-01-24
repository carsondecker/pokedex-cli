package main

import (
	"fmt"
	"os"

	"github.com/carsondecker/pokedex-cli/pokeapi"
)

type cliCommand struct {
	name        string
	description string
	callback    func() error
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
	}
}

func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp() error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	for k, v := range commands {
		fmt.Printf("\t%s: %s\n", k, v.description)
	}
	return nil
}

func commandMap() error {
	data, err := pokeapi.GetNextMapData()
	if err != nil {
		return err
	}
	for _, locationArea := range data.Results {
		fmt.Println(locationArea.Name)
	}
	return nil
}

func commandMapB() error {
	data, err := pokeapi.GetPrevMapData()
	if err != nil {
		return err
	}
	for _, locationArea := range data.Results {
		fmt.Println(locationArea.Name)
	}
	return nil
}
