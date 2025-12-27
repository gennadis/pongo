package main

import (
	"log"
	"os"
	"sync"
	"time"

	"github.com/gdamore/tcell/v2"
)

type Ball struct {
	Rune rune
	X    int
	Y    int
}

func NewBall() *Ball {
	return &Ball{
		Rune: '🟢',
		X:    1,
		Y:    1,
	}
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

	x := 0
	for {
		g.Screen.Clear()
		g.Screen.SetContent(x, 10, g.Ball.Rune, nil, style)
		g.Screen.Show()
		x++ // move ball right
		time.Sleep(40 * time.Millisecond)
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
