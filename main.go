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
		"explore": {
			Name:        "explore",
			Description: "Allows you to explore all the pokemon available in a specific location area",
			Callback:    internal.CommandExplore,
		},
	}
}

func main() {

	scanner := bufio.NewScanner(os.Stdin)

	var cmd string
	var area string

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

		cmd = cleanInput(input)[0]
		switch cmd {
		case "help":
			cmds["help"].Callback(config, nil)
			fmt.Println()
		case "map":
			cmds["map"].Callback(config, nil)
			fmt.Println()
		case "mapb":
			cmds["mapb"].Callback(config, nil)
			fmt.Println()
		case "explore":
			area = cleanInput(input)[1]
			fmt.Println("JTK - ", area)
			cmds["explore"].Callback(config, area)
			fmt.Println()
		case "exit":
			cmds["exit"].Callback(config, nil)
			fmt.Println()
		default:
			fmt.Println("An invalid command was received. Please type 'help' to see a list of valid commands")
			fmt.Println()
		}

		fmt.Printf("Your command was: %s\n", cmd)
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
func CommandHelp(cfg *shared.Config, opts ...any) error {
	fmt.Println("Welcome to the Pokedex! \nUsage: \n ")

	for _, cmd := range cmds {
		fmt.Printf("%s: %s \n", cmd.Name, cmd.Description)
	}

	return nil
}

// Exits the program
func CommandExit(cfg *shared.Config, opts ...any) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}
