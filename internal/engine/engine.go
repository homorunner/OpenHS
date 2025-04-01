package engine

import (
	"github.com/openhs/internal/game"
)

// Engine handles the game rules and mechanics
type Engine struct {
	game *game.Game
}

// NewEngine creates a new game engine
func NewEngine(g *game.Game) *Engine {
	return &Engine{
		game: g,
	}
}
