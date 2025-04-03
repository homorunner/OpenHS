package game

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/openhs/internal/logger"
)

type Game struct {
	Players            []*Player
	CurrentPlayer      *Player
	CurrentTurn        int
	CurrentPlayerIndex int
	Phase              GamePhase
}

type GamePhase int

const (
	InvalidPhase GamePhase = iota
	BeginFirst
	BeginShuffle
	BeginDraw
	BeginMulligan
	MainBegin
	MainReady
	MainResource
	MainDraw
	MainStart
	MainAction
	MainCombat
	MainEnd
	MainNext
	FinalWrapup
	FinalGameover
	MainCleanup
	MainStartTriggers
	MainSetActionStepType
	MainPreAction
	MainPostAction
)

// String returns a string representation of the GamePhase
func (p GamePhase) String() string {
	switch p {
	case InvalidPhase:
		return "Invalid Phase"
	case BeginFirst:
		return "Begin First"
	case BeginShuffle:
		return "Begin Shuffle"
	case BeginDraw:
		return "Begin Draw"
	case BeginMulligan:
		return "Begin Mulligan"
	case MainBegin:
		return "Main Begin"
	case MainReady:
		return "Main Ready"
	case MainResource:
		return "Main Resource"
	case MainDraw:
		return "Main Draw"
	case MainStart:
		return "Main Start"
	case MainAction:
		return "Main Action"
	case MainCombat:
		return "Main Combat"
	case MainEnd:
		return "Main End"
	case MainNext:
		return "Main Next"
	case FinalWrapup:
		return "Final Wrapup"
	case FinalGameover:
		return "Final Gameover"
	case MainCleanup:
		return "Main Cleanup"
	case MainStartTriggers:
		return "Main Start Triggers"
	case MainSetActionStepType:
		return "Main Set Action Step Type"
	case MainPreAction:
		return "Main Pre Action"
	case MainPostAction:
		return "Main Post Action"
	default:
		return fmt.Sprintf("Unknown Phase (%d)", int(p))
	}
}

func NewGame() *Game {
	return &Game{
		Players:     make([]*Player, 0),
		CurrentTurn: 0,
		Phase:       InvalidPhase,
	}
}

// LoadGame creates a new game from a configuration
func LoadGame(config *GameConfig) (*Game, error) {
	g := NewGame()

	// Create players based on configuration
	for _, playerConfig := range config.Players {
		player := NewPlayer()

		// Load hero card
		cardManager := GetCardManager()
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

// GameConfig represents the configuration for a game
type GameConfig struct {
	Players []PlayerConfig `json:"players"`
}

// PlayerConfig represents the configuration for a player
type PlayerConfig struct {
	Hero string   `json:"hero"`
	Deck []string `json:"deck"`
}

// LoadGameConfig loads a game configuration from a JSON file
func LoadGameConfig(configPath string) (*GameConfig, error) {
	// Read the JSON file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	// Parse the JSON data
	var config GameConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
