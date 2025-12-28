package main

import (
	"context"
	"strconv"
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
	Player1 *Player
	Player2 *Player
}

// NewGame creates a new Game with the given context, screen, entities, and configuration.
func NewGame(ctx context.Context, cancel context.CancelFunc, cfg *Config, scr tcell.Screen, b *Ball, p1 *Player, p2 *Player) *Game {
	return &Game{
		ctx:     ctx,
		cancel:  cancel,
		config:  cfg,
		Screen:  scr,
		Ball:    b,
		Player1: p1,
		Player2: p2,
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

			maxWidth, maxHeight := g.Screen.Size()

			// ball
			drawSprite(
				g.Screen,
				g.Ball.X, g.Ball.Y,
				g.Ball.X, g.Ball.Y,
				screenStyle, g.Ball.Emoji,
			)
			// player one paddle
			drawSprite(
				g.Screen,
				g.Player1.Paddle.X, g.Player1.Paddle.Y,
				g.Player1.Paddle.X+g.Player1.Paddle.width, g.Player1.Paddle.Y+g.Player1.Paddle.height,
				paddleStyle, g.Player1.Paddle.String,
			)
			// player one scrore
			drawSprite(g.Screen,
				(maxWidth/2)-5, 1,
				1, 1,
				screenStyle, strconv.Itoa(g.Player1.Score),
			)
			// player two paddle
			drawSprite(
				g.Screen,
				g.Player2.Paddle.X, g.Player2.Paddle.Y,
				g.Player2.Paddle.X+g.Player2.Paddle.width, g.Player2.Paddle.Y+g.Player2.Paddle.height,
				paddleStyle, g.Player2.Paddle.String,
			)
			// player two score
			drawSprite(g.Screen,
				(maxWidth/2)+5, 1,
				1, 1,
				screenStyle, strconv.Itoa(g.Player2.Score),
			)

			if g.Ball.HasTouched(*g.Player1.Paddle) || g.Ball.HasTouched(*g.Player2.Paddle) {
				g.Ball.ReverseX()
			}

			g.Ball.BounceWall(maxWidth, maxHeight)

			if g.Ball.X <= 0 {
				g.Player2.Score++
				g.Ball.Reset(maxWidth/2, maxHeight/2, -1, 1)
			}

			if g.Ball.X >= maxWidth {
				g.Player1.Score++
				g.Ball.Reset(maxWidth/2, maxHeight/2, 1, 1)
			}

			if g.GameOver() {
				drawSprite(g.Screen,
					(maxWidth/2)-4, 7,
					(maxWidth/2)+5, 7,
					screenStyle, "Game Over",
				)

				drawSprite(g.Screen,
					(maxWidth/2)-8, 11,
					(maxWidth/2)+5, 7,
					screenStyle, g.DeclareWinner()+" Wins!",
				)
				g.Screen.Show()
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
					g.Player1.Paddle.MoveUp()
				case ev.Rune() == 's':
					g.Player1.Paddle.MoveDown(maxHeight)
				// player two
				case ev.Key() == tcell.KeyUp:
					g.Player2.Paddle.MoveUp()
				case ev.Key() == tcell.KeyDown:
					g.Player2.Paddle.MoveDown(maxHeight)
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

func (g *Game) GameOver() bool {
	return g.Player1.Score == 3 || g.Player2.Score == 3
}

func (g *Game) DeclareWinner() string {
	if !g.GameOver() {
		return "No winner"
	}

	if g.Player1.Score > g.Player2.Score {
		return "Player 1"
	} else {
		return "Player 2"
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
