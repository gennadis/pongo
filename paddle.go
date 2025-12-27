package main

import "strings"

// Paddle represents a player's paddle with position and movement speed.
type Paddle struct {
	width  int
	height int
	X      int
	Y      int
	YVelo  int
	String string
}

// NewPaddle creates a new paddle at the given position using the provided configuration.
func NewPaddle(x int, y int, config *Config) *Paddle {
	return &Paddle{
		width:  1,
		height: config.PaddleSize,
		X:      x,
		Y:      y,
		YVelo:  config.PaddleSpeed,
		String: strings.Repeat("-", config.PaddleSize),
	}
}

// MoveUp moves the paddle upward, respecting the top boundary.
func (p *Paddle) MoveUp() {
	if p.Y > 0 {
		p.Y -= p.YVelo
	}
}

// MoveDown moves the paddle downward, respecting the bottom boundary.
func (p *Paddle) MoveDown(maxHeight int) {
	if p.Y < maxHeight-p.height {
		p.Y += p.YVelo
	}
}
