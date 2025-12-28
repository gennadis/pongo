# Pongo
A simple Pong game written in Go using the `tcell` for terminal rendering.

## Features
- Classic two-player Pong gameplay in your terminal
- **Scoring system** - First to 3 points wins
- **Game over screen** with winner announcement
- Thread-safe concurrent game loop and input handling
- Graceful shutdown with context cancellation

## Controls
- **Player 1**:
  - `W` - Move up
  - `S` - Move down

- **Player 2**:
  - `Arrow Up` - Move up
  - `Arrow Down` - Move down

- **Quit**: `Ctrl+C` - Exit the game gracefully
