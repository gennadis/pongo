package main

import (
	"log"
	"os"
	"sync"
	"time"

	"github.com/gdamore/tcell/v2"
)

func main() {
	screen, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("new screen: %+v", err)
	}

	if err := screen.Init(); err != nil {
		log.Fatalf("screen init: %+v", err)
	}

	style := tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorWhite)
	screen.SetStyle(style)

	wg := sync.WaitGroup{}
	wg.Add(2)
	go runGameLoop(&wg, screen, style)
	go handleKeyPress(&wg, screen)
	wg.Wait()
}

func runGameLoop(wg *sync.WaitGroup, screen tcell.Screen, style tcell.Style) {
	defer wg.Done()
	x := 0
	for {
		screen.Clear()
		screen.SetContent(x, 10, 'h', nil, style)
		screen.SetContent(x+1, 10, 'e', nil, style)
		screen.SetContent(x+2, 10, 'l', nil, style)
		screen.SetContent(x+3, 10, 'l', nil, style)
		screen.SetContent(x+4, 10, 'o', nil, style)
		screen.Show()
		x++ // move text right
		time.Sleep(40 * time.Millisecond)
	}
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
