package main

type Config struct {
	Language          string       `json:"language"`
	Teams             []TeamConfig `json:"teams"`
	ResourceNames     []string     `json:"resource_names"`
	GemNames          []string     `json:"gem_names"`
	LockLengthSeconds int          `json:"lock_length_seconds"`
}

type TeamConfig struct {
	Name           string          `json:"name"`
	Secret         string          `json:"secret,omitempty"`
	Admin          bool            `json:"admin,omitempty"`
	InitialBalance *InitialBalance `json:"initial_balance,omitempty"`
}

type InitialBalance struct {
	Resources []int `json:"resources"`
	Gems      []int `json:"gems"`
}

func stripConfigSecret(cfg Config) Config {
	teams := make([]TeamConfig, len(cfg.Teams))
	for i := range cfg.Teams {
		teams[i] = cfg.Teams[i]
		teams[i].Secret = ""
		teams[i].InitialBalance = nil
	}
	cfg.Teams = teams
	return cfg
}
