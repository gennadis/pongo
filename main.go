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
	paddleSize = 8
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

func (b *Ball) Bounce(maxWidth int, maxHeight int) {
	if b.X <= 0 || b.X >= maxWidth {
		b.XVelo *= -1
	}

	if b.Y <= 0 || b.Y >= maxHeight {
		b.YVelo *= -1
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

		drawSprite(
			g.Screen,
			g.Ball.X, g.Ball.Y,
			g.Ball.X, g.Ball.Y,
			screenStyle, g.Ball.Emoji,
		)

		drawSprite(
			g.Screen,
			g.Paddle1.X, g.Paddle1.Y,
			g.Paddle1.X+g.Paddle1.width, g.Paddle1.Y+g.Paddle1.height,
			paddleStyle, g.Paddle1.String,
		)

		drawSprite(
			g.Screen,
			g.Paddle2.X, g.Paddle2.Y,
			g.Paddle2.X+g.Paddle2.width, g.Paddle2.Y+g.Paddle2.height,
			paddleStyle, g.Paddle2.String,
		)

		maxWidth, maxHeight := g.Screen.Size()
		g.Ball.Bounce(maxWidth, maxHeight)

		time.Sleep(40 * time.Millisecond)
		g.Screen.Show()
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
	go handleKeyPress(&wg, scr)
	wg.Wait()
}

func handleKeyPress(wg *sync.WaitGroup, screen tcell.Screen) {
	defer wg.Done()
	for {
		switch event := screen.PollEvent().(type) {
		case *tcell.EventResize:
			screen.Sync()
		case *tcell.EventKey:
			if event.Key() == tcell.KeyCtrlC {
				screen.Fini()
				os.Exit(0)
			}
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
