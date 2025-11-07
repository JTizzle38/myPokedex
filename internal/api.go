package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/JTizzle38/myPokedex/shared"
)

var baseURL = "https://pokeapi.co/api/v2/"

//PokeAPI - Maps
//GET Request
//Endpoint: https://pokeapi.co/api/v2/location-area
//Returns 20 locations from the pokemon world

func CommandMap(cfg *shared.Config) error {

	fmt.Println("Display the Pokemon map regions! \n ")
	var mapURL string

	//Define which page we are starting on to get requests from
	if cfg.Next == nil {
		mapURL = baseURL + "location-area"
	} else {
		mapURL = *cfg.Next
	}

	//Make the API call
	data, err := GetJSON(mapURL)
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

func CommandMapBack(cfg *shared.Config) error {
	// Check if there's a previous page
	if cfg.Previous == nil {
		fmt.Println("You're on the first page already!")
		return nil
	}

	// Fetch the previous page
	data, err := GetJSON(*cfg.Previous)
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

// This is a generic function that can handle any REST API Calls [GET/POST/DELETE/PUT]
func fetchJSONResponse(method, url string, body any, headers map[string]string) ([]byte, error) {

	var bodyReader io.Reader

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

	defer resp.Body.Close()
	return data, nil
}

// A wrapper function for making GET requests
func GetJSON(url string) ([]byte, error) {
	return fetchJSONResponse("GET", url, nil, nil)
}

// A wrapper function for making POST requests
func PostJSON(url string, body any, headers map[string]string) ([]byte, error) {
	return fetchJSONResponse("POST", url, body, headers)
}

//For reference:
//func MarshalIndent(v interface{}, prefix, indent string) ([]byte, error)
