package shared

type Config struct {
	Next     *string
	Previous *string
}

type CLICommand struct {
	Name        string
	Description string
	Callback    func(*Config) error
}

type LocationAreaData struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}
