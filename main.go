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
			Description: "Usage: 'explore [valid city name]` | Allows you to explore all the pokemon available in a specific location area",
			Callback:    internal.CommandExplore,
		},
		"catch": {
			Name:        "catch",
			Description: "Usage: 'catch [name of pokemon]' | Gives you a chance to catch a pokemon based on it's base_experience value and RNG",
			Callback:    internal.CommandCatch,
		},
		"pokedex": {
			Name:        "pokedex",
			Description: "Outputs the current Pokedex list for the trainer",
			Callback:    CommandPokedex,
		},
		"inspect": {
			Name:        "inspect",
			Description: "Inspects more details about the captured pokemon in a trainer's pokedex",
			Callback:    CommandInspect,
		},
	}
}

func main() {

	scanner := bufio.NewScanner(os.Stdin)

	var cmd string
	var area string
	var pokemonName string

	config := &shared.Config{
		Next:     nil,
		Previous: nil,
		Cache:    pc.NewCache(5 * time.Minute),
		Trainer: shared.UserData{
			Name:    "JT",
			Pokedex: make(map[string]shared.PokemonDetail),
		},
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
			cmds["explore"].Callback(config, area)
			fmt.Println()
		case "catch":
			pokemonName = cleanInput(input)[1]
			cmds["catch"].Callback(config, pokemonName)
			fmt.Println()
		case "pokedex":
			cmds["pokedex"].Callback(config, nil)
			fmt.Println()
		case "inspect":
			pokemonName = cleanInput(input)[1]
			cmds["inspect"].Callback(config, pokemonName)
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

func CommandPokedex(cfg *shared.Config, opts ...any) error {
	//Print out current Pokedex count and list of pokemon
	fmt.Printf("Number of Pokemon captured (according to Pokedex): %v \n", len(cfg.Trainer.Pokedex))
	fmt.Printf("Current names of Pokemon: \n")
	counter := 1
	for name := range cfg.Trainer.Pokedex {
		fmt.Printf("%v. %s \n", counter, name)
		counter++
	}

	return nil
}

func CommandInspect(cfg *shared.Config, opts ...any) error {

	if len(opts) == 0 {
		return fmt.Errorf("ERROR | main.go | CommandInspect(): A pokemon name is required in order to use this inspect command")
	}

	pokemon, ok := opts[0].(string)
	if !ok {
		return fmt.Errorf("ERROR | main.go | CommandInspect(): The pokemon name must be of type string")
	}

	//First check to see if the trainer has captured the pokemon that is requesting inspection
	//If not, then print an error like message
	//If so, print out all the details of the pokemon

	if _, exists := cfg.Trainer.Pokedex[pokemon]; !exists {
		fmt.Printf("You have not caught this pokemon yet: %s", pokemon)
		return nil
	} else {
		//Print out the stats of the pokemon inspecting
		fmt.Printf("Name: %s\n", cfg.Trainer.Pokedex[pokemon].Name)
		fmt.Printf("Height: %d \n", cfg.Trainer.Pokedex[pokemon].Height)
		fmt.Printf("Weight: %d\n", cfg.Trainer.Pokedex[pokemon].Weight)
		fmt.Printf("Stats: \n")

		for _, stat := range cfg.Trainer.Pokedex[pokemon].Stats {
			fmt.Printf("- %s: %d\n", stat.Stat.Name, stat.BaseStat)
		}

		fmt.Printf("Types: \n")
		for _, stat := range cfg.Trainer.Pokedex[pokemon].Types {
			fmt.Printf("- %s \n", stat.Type.Name)
		}
	}

	return nil
}
