package main

import (
	"log"
	"sync"

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

	width, height := screen.Size()
	ball := NewBall()
	paddle1 := NewPaddle(2, height/2-3)
	paddle2 := NewPaddle(width-3, height/2-3)
	game := NewGame(screen, ball, paddle1, paddle2)

	go game.Run()

	wg := sync.WaitGroup{}
	wg.Add(1)
	go game.HandleKeyPress(&wg, height)
	wg.Wait()
}
