package main

import (
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gdamore/tcell/v2"
)

const (
	ballEmoji  = "🟢"
	paddleSize = 12
	gameSpeed  = 60
)

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

func (b *Ball) hasTouched(p Paddle) bool {
	return b.X >= p.X && b.X <= p.X+p.width && b.Y >= p.Y && b.Y <= p.Y+p.height
}

func (b *Ball) reverseX() {
	b.XVelo *= -1
}

func (b *Ball) reverseY() {
	b.YVelo *= -1
}

func (b *Ball) BouncePaddle(maxWidth int, maxHeight int) {
	if b.X <= 0 || b.X >= maxWidth {
		b.reverseX()
	}

	if b.Y <= 0 || b.Y >= maxHeight {
		b.reverseY()
	}
}

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

type Game struct {
	Screen  tcell.Screen
	Ball    *Ball
	Paddle1 *Paddle
	Paddle2 *Paddle
}

func NewGame(scr tcell.Screen, b *Ball, p1 *Paddle, p2 *Paddle) *Game {
	return &Game{
		Screen:  scr,
		Ball:    b,
		Paddle1: p1,
		Paddle2: p2,
	}
}

func (g *Game) Run() {
	screenStyle := tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorWhite)
	paddleStyle := tcell.StyleDefault.Background(tcell.ColorWhite).Foreground(tcell.ColorWhite)
	g.Screen.SetStyle(screenStyle)

	for {
		g.Screen.Clear()
		g.Ball.UpdatePosition()
		// ball
		drawSprite(
			g.Screen,
			g.Ball.X, g.Ball.Y,
			g.Ball.X, g.Ball.Y,
			screenStyle, g.Ball.Emoji,
		)
		// player one
		drawSprite(
			g.Screen,
			g.Paddle1.X, g.Paddle1.Y,
			g.Paddle1.X+g.Paddle1.width, g.Paddle1.Y+g.Paddle1.height,
			paddleStyle, g.Paddle1.String,
		)
		// player two
		drawSprite(
			g.Screen,
			g.Paddle2.X, g.Paddle2.Y,
			g.Paddle2.X+g.Paddle2.width, g.Paddle2.Y+g.Paddle2.height,
			paddleStyle, g.Paddle2.String,
		)

		if g.Ball.hasTouched(*g.Paddle1) || g.Ball.hasTouched(*g.Paddle2) {
			g.Ball.reverseX()
			g.Ball.reverseY()
		}

		maxWidth, maxHeight := g.Screen.Size()
		g.Ball.BounceWall(maxWidth, maxHeight)
		time.Sleep(gameSpeed * time.Millisecond)
		g.Screen.Show()
	}
}

func (g *Game) handleKeyPress(wg *sync.WaitGroup, maxHeight int) {
	defer wg.Done()
	for {
		switch event := g.Screen.PollEvent().(type) {
		case *tcell.EventKey:
			switch {
			// player one
			case event.Rune() == 'w':
				g.Paddle1.MoveUp()
			case event.Rune() == 's':
				g.Paddle1.MoveDown(maxHeight)
			// player two
			case event.Key() == tcell.KeyUp:
				g.Paddle2.MoveUp()
			case event.Key() == tcell.KeyDown:
				g.Paddle2.MoveDown(maxHeight)
			// quit
			case event.Key() == tcell.KeyCtrlC:
				g.Screen.Fini()
				os.Exit(0)
			}
		case *tcell.EventResize:
			g.Screen.Sync()
		}
	}
}

func main() {
	scr, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("new screen: %+v", err)
	}

	if err := scr.Init(); err != nil {
		log.Fatalf("screen init: %+v", err)
	}

	ball := NewBall()

	width, height := scr.Size()
	p1 := NewPaddle(2, height/2-3)
	p2 := NewPaddle(width-3, height/2-3)

	game := NewGame(scr, ball, p1, p2)
	go game.Run()

	wg := sync.WaitGroup{}
	wg.Add(1)
	go game.handleKeyPress(&wg, height)
	wg.Wait()
}

func drawSprite(s tcell.Screen, x1, y1, x2, y2 int, style tcell.Style, text string) {
	row := y1
	col := x1

	for _, r := range []rune(text) {
		s.SetContent(col, row, r, nil, style)
		col++
		if col >= x2 {
			row++
			col = x1
		}
		if row > y2 {
			break
		}
	}
}
