package main

const ballEmoji = "🟢"

// Ball represents the game ball with position and velocity.
type Ball struct {
	Emoji string
	X     int
	Y     int
	XVelo int
	YVelo int
}

// NewBall creates a new ball at position (1, 1) with velocity (1, 1).
func NewBall() *Ball {
	return &Ball{
		Emoji: ballEmoji,
		X:     1,
		Y:     1,
		XVelo: 1,
		YVelo: 1,
	}
}

// UpdatePosition updates the ball's position based on its current velocity.
func (b *Ball) UpdatePosition() {
	b.X += b.XVelo
	b.Y += b.YVelo
}

// BounceWall handles wall collision detection and reverses velocity when hitting walls.
func (b *Ball) BounceWall(maxWidth int, maxHeight int) {
	if b.Y <= 0 || b.Y >= maxHeight {
		b.YVelo *= -1
	}
}

func (b *Ball) Reset(x int, y int, xVelo int, yVelo int) {
	b.X = x
	b.Y = y
	b.XVelo = xVelo
	b.YVelo = yVelo
}

// HasTouched checks if the ball has collided with the given paddle.
func (b *Ball) HasTouched(p Paddle) bool {
	return b.X >= p.X && b.X <= p.X+p.width && b.Y >= p.Y && b.Y <= p.Y+p.height
}

// ReverseX reverses the ball's horizontal velocity.
func (b *Ball) ReverseX() {
	b.XVelo *= -1
}

// ReverseY reverses the ball's vertical velocity.
func (b *Ball) ReverseY() {
	b.YVelo *= -1
}
