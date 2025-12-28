package main

import (
	"context"
	"sync"
	"time"

	"github.com/gdamore/tcell/v2"
)

// Game holds the game state including screen, entities, and synchronization primitives.
type Game struct {
	mu      sync.RWMutex
	ctx     context.Context
	cancel  context.CancelFunc
	config  *Config
	Screen  tcell.Screen
	Ball    *Ball
	Paddle1 *Paddle
	Paddle2 *Paddle
}

// NewGame creates a new Game with the given context, screen, entities, and configuration.
func NewGame(ctx context.Context, cancel context.CancelFunc, cfg *Config, scr tcell.Screen, b *Ball, p1 *Paddle, p2 *Paddle) *Game {
	return &Game{
		ctx:     ctx,
		cancel:  cancel,
		config:  cfg,
		Screen:  scr,
		Ball:    b,
		Paddle1: p1,
		Paddle2: p2,
	}
}

// Run is the main game loop that updates and renders the game state.
// It runs in a separate goroutine and respects context cancellation for graceful shutdown.
func (g *Game) Run() {
	screenStyle := tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorWhite)
	paddleStyle := tcell.StyleDefault.Background(tcell.ColorWhite).Foreground(tcell.ColorWhite)
	g.Screen.SetStyle(screenStyle)

	ticker := time.NewTicker(time.Duration(g.config.GameSpeed) * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-g.ctx.Done():
			return
		case <-ticker.C:
			g.mu.Lock()
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
			}

			maxWidth, maxHeight := g.Screen.Size()
			g.Ball.BounceWall(maxWidth, maxHeight)

			if g.Ball.X <= 0 {
				g.Ball.Reset(maxWidth/2, maxHeight/2, -1, 1)
			}

			if g.Ball.X >= maxWidth {
				g.Ball.Reset(maxWidth/2, maxHeight/2, 1, 1)
			}

			g.Screen.Show()
			g.mu.Unlock()
		}
	}
}

// HandleKeyPress processes keyboard input and updates paddle positions.
// It runs in a separate goroutine and respects context cancellation for graceful shutdown.
func (g *Game) HandleKeyPress() {
	for {
		select {
		case <-g.ctx.Done():
			return
		default:
			event := g.Screen.PollEvent()
			switch ev := event.(type) {
			case *tcell.EventKey:
				g.mu.Lock()
				_, maxHeight := g.Screen.Size()
				switch {
				// player one
				case ev.Rune() == 'w':
					g.Paddle1.MoveUp()
				case ev.Rune() == 's':
					g.Paddle1.MoveDown(maxHeight)
				// player two
				case ev.Key() == tcell.KeyUp:
					g.Paddle2.MoveUp()
				case ev.Key() == tcell.KeyDown:
					g.Paddle2.MoveDown(maxHeight)
				// quit
				case ev.Key() == tcell.KeyCtrlC:
					g.mu.Unlock()
					g.cancel()
					return
				}
				g.mu.Unlock()
			case *tcell.EventResize:
				g.Screen.Sync()
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
