package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"

	pc "github.com/JTizzle38/myPokedex/internal/pokecache"
	"github.com/JTizzle38/myPokedex/shared"
)

var baseURL = "https://pokeapi.co/api/v2/"

//PokeAPI - Maps
//GET Request
//Endpoint: https://pokeapi.co/api/v2/location-area
//Returns 20 locations from the pokemon world

func CommandMap(cfg *shared.Config, opts ...any) error {

	fmt.Println("Display the Pokemon map regions! \n ")
	var mapURL string

	//Define which page we are starting on to get requests from
	if cfg.Next == nil {
		mapURL = baseURL + "location-area"
	} else {
		mapURL = *cfg.Next
	}

	//Make the API call
	data, err := GetJSON(mapURL, cfg.Cache)
	if err != nil {
		fmt.Printf("ERROR | api.go | CommandMap():  %s", err)
	}

	//Parse the data from the response
	var page shared.LocationAreaData
	err = json.Unmarshal(data, &page)
	if err != nil {
		return fmt.Errorf("ERROR | api.go | CommandMap(): %s", err)
	}

	//Update the config variable with pagination
	if page.Next != "" {
		cfg.Next = &page.Next
	} else {
		cfg.Next = nil // No more pages
	}

	if page.Previous != "" {
		cfg.Previous = &page.Previous
	} else {
		cfg.Previous = nil // No previous page
	}

	//Print the results (Max 20 locations at a time)
	for _, result := range page.Results {
		fmt.Printf("%s\n", result.Name)
	}

	return nil
}

func CommandMapBack(cfg *shared.Config, opts ...any) error {
	// Check if there's a previous page
	if cfg.Previous == nil {
		fmt.Println("You're on the first page already!")
		return nil
	}

	// Fetch the previous page
	data, err := GetJSON(*cfg.Previous, cfg.Cache)
	if err != nil {
		return fmt.Errorf("failed to fetch location areas: %w", err)
	}

	// Parse the response
	var page shared.LocationAreaData
	err = json.Unmarshal(data, &page)
	if err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	// Update config with pagination URLs
	if page.Next != "" {
		cfg.Next = &page.Next
	} else {
		cfg.Next = nil
	}

	if page.Previous != "" {
		cfg.Previous = &page.Previous
	} else {
		cfg.Previous = nil
	}

	// Print the results
	for _, result := range page.Results {
		fmt.Printf("%s\n", result.Name)
	}

	return nil
}

func CommandExplore(cfg *shared.Config, opts ...any) error {
	//Outputting the current city to be explored
	//This should be the first word followed by the 'explore' command
	//but let's do some quick validation just in case

	if len(opts) == 0 {
		return fmt.Errorf("ERROR | api.go | CommandExplore(): An area name is required from the list output from either the map/mapb commands")
	}

	area, ok := opts[0].(string)
	if !ok {
		return fmt.Errorf("ERROR | api.go | CommandExplore(): The area name must be of type string")
	}

	fmt.Printf("Exploring %s...", area)

	//Creating the full URL to be used based on the received location area
	exploreURL := baseURL + "location-area/" + area

	//Make the API call
	data, err := GetJSON(exploreURL, cfg.Cache)
	if err != nil {
		fmt.Errorf("ERROR | api.go | CommandExplore(): %s", err)
	}

	//Parse the data from the response
	var areaData shared.LocationAreaDetail
	err = json.Unmarshal(data, &areaData)
	if err != nil {
		return fmt.Errorf("ERROR | api.go | CommandExplore(): %s", err)
	}

	//Print the results of all the pokemon in the given location area
	fmt.Println("Found Pokemon: ")
	for _, p := range areaData.PokemonEncounters {
		fmt.Printf("- %s\n", p.Pokemon.Name)
	}

	return nil
}

func CommandCatch(cfg *shared.Config, opts ...any) error {
	if len(opts) == 0 {
		return fmt.Errorf("ERROR | api.go | CommandCatch(): A pokemon name is required in order to use this catch command")
	}

	pokemon, ok := opts[0].(string)
	if !ok {
		return fmt.Errorf("ERROR | api.go | CommandCatch(): The pokemon name must be of type string")
	}

	fmt.Printf("Throwing a Pokeball at %s...", pokemon)

	//Creating the full URL to be used based on the user input pokemon name
	pokemonURL := baseURL + "pokemon/" + pokemon

	//Make the API call
	data, err := GetJSON(pokemonURL, cfg.Cache)
	if err != nil {
		fmt.Errorf("ERROR | api.go | CommandCatch(): %s", err)
	}

	//Parse the data from the response
	var pokemonData shared.PokemonDetail
	err = json.Unmarshal(data, &pokemonData)
	if err != nil {
		return fmt.Errorf("ERROR | api.go | CommandCatch(): %s", err)
	}

	//Calculate the RNG of catching the pokemon
	//User get's 3 chances to catch the mentioned pokemon
	var catchChance int
	for i := 0; i < 3; i++ {
		catchChance = rand.Intn(150)
		fmt.Printf("Throwing a Pokeball at %s...", pokemonData.Name)
		if catchChance >= pokemonData.BaseExperience {
			fmt.Printf("%s was caught!", pokemonData.Name)

			//Saving the captured pokemon to the trainer's pokedex
			cfg.Trainer.Pokedex[pokemonData.Name] = pokemonData
			break
		} else {
			fmt.Printf("Attempting to catch %s again \n", pokemonData.Name)
		}
	}

	//Checks the pokedex to see if the pokemon was captured successfully
	if _, exists := cfg.Trainer.Pokedex[pokemonData.Name]; !exists {
		fmt.Printf("%s escaped!", pokemonData.Name)
	}

	return nil

}

// This is a generic function that can handle any REST API Calls [GET/POST/DELETE/PUT]
func fetchJSONResponse(method, url string, body any, headers map[string]string, cache *pc.Cache) ([]byte, error) {

	var bodyReader io.Reader

	if method == "GET" && cache != nil {
		if dataCache, found := cache.GetEntry(url); found {
			fmt.Println("Using cached data...")
			return dataCache, nil
		}
	}

	// Encode body if present
	if body != nil {
		jsonBytes, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewBuffer(jsonBytes)
	}

	// Build request
	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add headers
	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	// Execute request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %s", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body, %s", err)
	}

	if method == "GET" && cache != nil {
		cache.AddEntry(url, data)
	}

	return data, nil
}

// A wrapper function for making GET requests
func GetJSON(url string, cache *pc.Cache) ([]byte, error) {
	return fetchJSONResponse("GET", url, nil, nil, cache)
}

// A wrapper function for making POST requests
func PostJSON(url string, body any, headers map[string]string, cache *pc.Cache) ([]byte, error) {
	return fetchJSONResponse("POST", url, body, headers, cache)
}

//For reference:
//func MarshalIndent(v interface{}, prefix, indent string) ([]byte, error)
