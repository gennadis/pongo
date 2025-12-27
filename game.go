package main

import (
	"os"
	"sync"
	"time"

	"github.com/gdamore/tcell/v2"
)

const gameSpeed = 60

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

		if g.Ball.HasTouched(*g.Paddle1) || g.Ball.HasTouched(*g.Paddle2) {
			g.Ball.ReverseX()
			g.Ball.ReverseY()
		}

		maxWidth, maxHeight := g.Screen.Size()
		g.Ball.BounceWall(maxWidth, maxHeight)
		time.Sleep(gameSpeed * time.Millisecond)
		g.Screen.Show()
	}
}

func (g *Game) HandleKeyPress(wg *sync.WaitGroup, maxHeight int) {
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
