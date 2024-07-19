package main

type Config struct {
	Teams             []TeamConfig `json:"teams"`
	Currencies        []string     `json:"currencies"`
	LockLengthSeconds int          `json:"lock_length_seconds"`
}

type TeamConfig struct {
	Name           string `json:"name"`
	Secret         string `json:"secret"`
	Admin          bool   `json:"admin"`
	InitialBalance *[]int `json:"initial_balance"`
}
