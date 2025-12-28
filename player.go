package main

// Player holds the player state including score and paddle.
type Player struct {
	Score  int
	Paddle *Paddle
}

// NewPlayer creates a new Player with the given score and paddle.
func NewPlayer(p *Paddle) *Player {
	return &Player{
		Score:  0,
		Paddle: p,
	}
}
