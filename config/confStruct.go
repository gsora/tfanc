package config

// Configuration is a structure holding content of the configuration file
type Configuration struct {
	Targets []Targets `json:"targets"`
}

// Targets is a structure holding a range and a level
type Targets struct {
	Range []int `json:"range"`
	Level int   `json:"level"`
}
