package main

import (
	"errors"
	"fmt"
	"math/rand/v2"
	"os"

	"github.com/carsondecker/pokedex-cli/pokeapi"
)

type cliCommand struct {
	name        string
	description string
	callback    func([]string) error
}

var commands map[string]cliCommand
var caughtPokemon map[string]pokeapi.PokemonData

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
		"catch": {
			name:        "catch {pokemon name}",
			description: "Attempts to catch the given pokemon",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect {pokemon name}",
			description: "Gets data about the given pokemon if you've caught it",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "Shows all the pokemon in your pokedex",
			callback:    commandPokedex,
		},
	}
	caughtPokemon = make(map[string]pokeapi.PokemonData)
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
		fmt.Printf("\t- %s\n", pokemonEncounter.Pokemon.Name)
	}
	return nil
}

func commandCatch(args []string) error {
	if len(args) == 0 {
		return errors.New("catch command requires pokemon name as argument")
	}

	if _, ok := caughtPokemon[args[0]]; ok {
		return errors.New("you've already caught this pokemon")
	}

	data, err := pokeapi.GetNewPokemonData(args[0])
	if err != nil {
		return err
	}
	fmt.Printf("Throwing a Pokeball at %s...\n", data.Name)
	catchRate := 1.0 / (1.0 + float64(data.BaseExperience)/100.0)
	if catchRate < 0.2 {
		catchRate = 0.2
	}

	isCaught := rand.Float64() < catchRate
	if !isCaught {
		fmt.Printf("%s escaped!\n", data.Name)
		return nil
	}
	fmt.Printf("%s was caught!\n", data.Name)
	caughtPokemon[data.Name] = data
	return nil
}

func commandInspect(args []string) error {
	if len(args) == 0 {
		return errors.New("catch command requires pokemon name as argument")
	}

	data, ok := caughtPokemon[args[0]]
	if !ok {
		return errors.New("that pokemon isn't in your pokedex")
	}

	fmt.Printf("Name: %s\n", data.Name)
	fmt.Printf("Height: %d\n", data.Height)
	fmt.Printf("Weight: %d\n", data.Weight)
	fmt.Println("Stats:")
	for _, pkmnStat := range data.Stats {
		fmt.Printf("\t- %s: %d\n", pkmnStat.Stat.StatName, pkmnStat.BaseStat)
	}
	fmt.Println("Types:")
	for _, pkmnType := range data.Types {
		fmt.Printf("\t- %s\n", pkmnType.Type.TypeName)
	}

	return nil
}

func commandPokedex(args []string) error {
	if len(caughtPokemon) == 0 {
		return errors.New("you haven't caught any pokemon, try using the 'catch' command")
	}

	fmt.Println("Your Pokedex:")
	for name := range caughtPokemon {
		fmt.Printf("\t- %s\n", name)
	}
	return nil
}
