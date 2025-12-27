package main

import (
	"log"
	"os"
	"sync"
	"time"

	"github.com/gdamore/tcell/v2"
)

const ballEmoji = '🟢'

type Ball struct {
	Rune  rune
	X     int
	Y     int
	XVelo int
	YVelo int
}

func NewBall() *Ball {
	return &Ball{
		Rune:  ballEmoji,
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

type Game struct {
	Screen tcell.Screen
	Ball   Ball
}

func NewGame(scr tcell.Screen, b Ball) *Game {
	return &Game{
		Screen: scr,
		Ball:   b,
	}
}

func (g *Game) Run() {
	style := tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorWhite)
	g.Screen.SetStyle(style)

	for {
		g.Screen.Clear()
		g.Ball.UpdatePosition()
		g.Screen.SetContent(g.Ball.X, g.Ball.Y, g.Ball.Rune, nil, style)
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
	game := NewGame(scr, *ball)
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
