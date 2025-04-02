package game

import (
	"github.com/openhs/internal/config"
	"github.com/openhs/internal/types"
)

// Player represents a player in the game
type Player struct {
	Deck      []types.Card
	Hand      []types.Card
	Field     []types.Card
	Graveyard []types.Card
	Hero      types.Card
	HeroPower types.Card
	Weapon    types.Card

	Mana          int
	MaxMana       int
	TotalMana     int
	FatigueDamage int
	HandSize      int
	HasWeapon     bool
}

// NewPlayer creates a new player from a configuration
func NewPlayer() *Player {
	return &Player{
		Deck:      make([]types.Card, 0),
		Hand:      make([]types.Card, 0),
		Field:     make([]types.Card, 0),
		Graveyard: make([]types.Card, 0),
		HandSize:  10,
		MaxMana:   config.DefaultMaxMana,
		Mana:      config.DefaultStartingMana,
		TotalMana: config.DefaultStartingMana,
	}
}
