package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	initCommands()
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
		commands[inputWords[0]].callback()
	}
}

func cleanInput(text string) []string {
	return strings.Fields(strings.TrimSpace(strings.ToLower(text)))
}
