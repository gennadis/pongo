package main

const ballEmoji = "🟢"

type Ball struct {
	Emoji string
	X     int
	Y     int
	XVelo int
	YVelo int
}

func NewBall() *Ball {
	return &Ball{
		Emoji: ballEmoji,
		X:     1,
		Y:     1,
		XVelo: 1,
		YVelo: 1,
	}
}

func (b *Ball) UpdatePosition() {
	b.X += b.XVelo
	b.Y += b.YVelo
}

func (b *Ball) BounceWall(maxWidth int, maxHeight int) {
	if b.X <= 0 || b.X >= maxWidth {
		b.XVelo *= -1
	}

	if b.Y <= 0 || b.Y >= maxHeight {
		b.YVelo *= -1
	}
}

func (b *Ball) BouncePaddle(maxWidth int, maxHeight int) {
	if b.X <= 0 || b.X >= maxWidth {
		b.ReverseX()
	}

	if b.Y <= 0 || b.Y >= maxHeight {
		b.ReverseY()
	}
}

func (b *Ball) HasTouched(p Paddle) bool {
	return b.X >= p.X && b.X <= p.X+p.width && b.Y >= p.Y && b.Y <= p.Y+p.height
}

func (b *Ball) ReverseX() {
	b.XVelo *= -1
}

func (b *Ball) ReverseY() {
	b.YVelo *= -1
}
