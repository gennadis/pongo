package main

import (
	"context"
	"log"

	"github.com/gdamore/tcell/v2"
)

func main() {
	config := NewConfig()

	screen, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("new screen: %+v", err)
	}

	if err := screen.Init(); err != nil {
		log.Fatalf("screen init: %+v", err)
	}
	defer func() {
		screen.Fini()
	}()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	width, height := screen.Size()
	ball := NewBall()
	paddle1 := NewPaddle(0, height/2-3, config)
	paddle2 := NewPaddle(width-1, height/2-3, config)
	game := NewGame(ctx, cancel, config, screen, ball, paddle1, paddle2)

	go game.Run()
	go game.HandleKeyPress()

	<-ctx.Done()
}
