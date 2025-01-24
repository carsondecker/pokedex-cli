package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/carsondecker/pokedex-cli/pokeapi"
)

func main() {
	initCommands()
	pokeapi.CreateCache(5)

	mapConfig := pokeapi.Config{
		Next: "https://pokeapi.co/api/v2/location-area/?offset=0&limit=20",
		Prev: "",
	}

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		if scanner.Text() == "" {
			fmt.Println("Please enter a command.")
			continue
		}
		inputWords := cleanInput(scanner.Text())
		if _, ok := commands[inputWords[0]]; !ok {
			fmt.Println("Unknown command.")
			continue
		}
		err := commands[inputWords[0]].callback(&mapConfig)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	}
}

func cleanInput(text string) []string {
	return strings.Fields(strings.TrimSpace(strings.ToLower(text)))
}
