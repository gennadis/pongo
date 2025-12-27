package main

import (
	"log"

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

	style := tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorBlack)
	screen.SetStyle(style)

	for {
		screen.SetContent(0, 0, 'h', nil, style)
		screen.SetContent(1, 0, 'e', nil, style)
		screen.SetContent(2, 0, 'l', nil, style)
		screen.SetContent(3, 0, 'l', nil, style)
		screen.SetContent(4, 0, 'o', nil, style)
		screen.Show()
	}
}
