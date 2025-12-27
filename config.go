package main

// Config holds game configuration parameters.
type Config struct {
	GameSpeed   int
	PaddleSize  int
	PaddleSpeed int
}

// NewConfig returns a new Config with default values.
func NewConfig() *Config {
	return &Config{
		GameSpeed:   40,
		PaddleSize:  10,
		PaddleSpeed: 3,
	}
}
