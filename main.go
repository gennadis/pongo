package main

import (
	"log"
	"os"
	"sync"
	"time"

	"github.com/gdamore/tcell/v2"
)

type Game struct {
	Screen tcell.Screen
}

func NewGame(scr tcell.Screen) *Game {
	return &Game{
		Screen: scr,
	}
}

func (g *Game) Run() {
	style := tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorWhite)
	g.Screen.SetStyle(style)

	x := 0
	for {
		g.Screen.Clear()
		g.Screen.SetContent(x, 10, 'h', nil, style)
		g.Screen.SetContent(x+1, 10, 'e', nil, style)
		g.Screen.SetContent(x+2, 10, 'l', nil, style)
		g.Screen.SetContent(x+3, 10, 'l', nil, style)
		g.Screen.SetContent(x+4, 10, 'o', nil, style)
		g.Screen.Show()
		x++ // move text right
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

	game := NewGame(scr)
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
