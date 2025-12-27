package main

import "strings"

const paddleSize = 12

type Paddle struct {
	width  int
	height int
	X      int
	Y      int
	YVelo  int
	String string
}

func NewPaddle(x int, y int) *Paddle {
	return &Paddle{
		width:  1,
		height: paddleSize,
		X:      x,
		Y:      y,
		YVelo:  3,
		String: strings.Repeat("-", paddleSize),
	}
}

func (p *Paddle) MoveUp() {
	if p.Y > 0 {
		p.Y -= p.YVelo
	}
}

func (p *Paddle) MoveDown(maxHeight int) {
	if p.Y < maxHeight-p.height {
		p.Y += p.YVelo
	}
}
