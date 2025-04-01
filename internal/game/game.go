package game

import (
	"github.com/openhs/internal/card"
	"github.com/openhs/internal/config"
	"github.com/openhs/internal/logger"
)

type Game struct {
	Players     []*Player
	CurrentTurn int
	Phase       GamePhase
}

type GamePhase int

const (
	StartGame GamePhase = iota
	StartTurn
	Draw
	Play
	Combat
	EndTurn
	EndGame
)

func NewGame() *Game {
	return &Game{
		Players:     make([]*Player, 0),
		CurrentTurn: 0,
		Phase:       StartGame,
	}
}

// LoadGame creates a new game from a configuration
func LoadGame(config *config.GameConfig) (*Game, error) {
	g := NewGame()
	
	// Create players based on configuration
	for _, playerConfig := range config.Players {
		player := NewPlayer(playerConfig)
		
		// Load hero card
		cardManager := card.GetCardManager()
		heroCard, err := cardManager.CreateCard(playerConfig.Hero)
		if err != nil {
			logger.Error("Failed to load hero card: " + err.Error())
			return nil, err
		}
		player.Hero = *heroCard
		
		// Load deck cards
		for _, cardName := range playerConfig.Deck {
			cardInstance, err := cardManager.CreateCard(cardName)
			if err != nil {
				logger.Error("Failed to load card: " + err.Error())
				return nil, err
			}
			player.Deck = append(player.Deck, *cardInstance)
		}
		
		g.Players = append(g.Players, player)
	}
	
	return g, nil
}
