package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/JTizzle38/myPokedex/internal"
	pc "github.com/JTizzle38/myPokedex/internal/pokecache"
	"github.com/JTizzle38/myPokedex/shared"
)

var cmds map[string]shared.CLICommand

func init() {
	cmds = map[string]shared.CLICommand{
		"help": {
			Name:        "help",
			Description: "Displays a help message",
			Callback:    CommandHelp,
		},
		"exit": {
			Name:        "exit",
			Description: "Exit the Pokedex",
			Callback:    CommandExit,
		},
		"map": {
			Name:        "map",
			Description: "Displays 20 location areas of the Pokemon world",
			Callback:    internal.CommandMap,
		},
		"mapb": {
			Name:        "mapb",
			Description: "Displays the previous 20 location areas of the Pokemon world",
			Callback:    internal.CommandMapBack,
		},
	}
}

func main() {

	scanner := bufio.NewScanner(os.Stdin)

	var word string

	config := &shared.Config{
		Next:     nil,
		Previous: nil,
		Cache:    pc.NewCache(5 * time.Minute),
	}

	for {
		fmt.Print("Pokedex > ")
		if !scanner.Scan() {
			break
		}
		input := scanner.Text()

		switch input {
		case "help":
			word = cleanInput(input)[0]
			cmds["help"].Callback(config)
			fmt.Println()
		case "map":
			word = cleanInput(input)[0]
			cmds["map"].Callback(config)
			fmt.Println()
		case "mapb":
			word = cleanInput(input)[0]
			cmds["mapb"].Callback(config)
			fmt.Println()
		case "exit":
			word = cleanInput(input)[0]
			cmds["exit"].Callback(config)
			fmt.Println()
		default:
			word = cleanInput(input)[0]
			fmt.Println()
		}

		fmt.Printf("Your command was: %s\n", word)
		fmt.Println()

		if err := scanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "reading standard input:", err)
		}
	}

}

// Splits the user's input into individual words
// Also lower cases each word
func cleanInput(text string) []string {
	output := strings.Fields(text)
	for i, word := range output {
		output[i] = strings.ToLower(word)
	}
	return output
}

// Displays the help message and available commands for the Pokedex
func CommandHelp(cfg *shared.Config) error {
	fmt.Println("Welcome to the Pokedex! \nUsage: \n ")

	for _, cmd := range cmds {
		fmt.Printf("%s: %s \n", cmd.Name, cmd.Description)
	}

	return nil
}

// Exits the program
func CommandExit(cfg *shared.Config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}
