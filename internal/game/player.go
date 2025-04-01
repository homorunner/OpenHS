package game

import (
	"github.com/openhs/internal/config"
	"github.com/openhs/internal/types"
)

// Player represents a player in the game
type Player struct {
	Mana      int
	MaxMana   int
	Deck      []types.Card
	Hand      []types.Card
	Board     []types.Card
	Hero      types.Card
	HeroPower types.Card
}

// NewPlayer creates a new player from a configuration
func NewPlayer(config config.PlayerConfig) *Player {
	return &Player{
		Deck:  make([]types.Card, 0),
		Hand:  make([]types.Card, 0),
		Board: make([]types.Card, 0),
	}
}
