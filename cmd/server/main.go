package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	"github.com/openhs/internal/bootstrap"
	"github.com/openhs/internal/engine"
	"github.com/openhs/internal/game"
)

// GameState represents the simplified game state sent to the frontend
type GameState struct {
	CurrentTurn        int                  `json:"currentTurn"`
	Phase              string               `json:"phase"`
	CurrentPlayerIndex int                  `json:"currentPlayerIndex"`
	Players            []*SimplifiedPlayer  `json:"players"`
	Actions            []string             `json:"availableActions"`
}

// SimplifiedPlayer represents the player state for the frontend
type SimplifiedPlayer struct {
	Hero      *SimplifiedEntity   `json:"hero"`
	Hand      []*SimplifiedEntity `json:"hand"`
	Field     []*SimplifiedEntity `json:"field"`
	Mana      int                 `json:"mana"`
	TotalMana int                 `json:"totalMana"`
	Weapon    *SimplifiedEntity   `json:"weapon,omitempty"`
}

// SimplifiedEntity represents a card entity for the frontend
type SimplifiedEntity struct {
	Name        string   `json:"name"`
	Attack      int      `json:"attack"`
	Health      int      `json:"health"`
	Cost        int      `json:"cost"`
	Type        string   `json:"type"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
	CanAttack   bool     `json:"canAttack"`
}

var (
	gameEngine *engine.Engine
	gameObj    *game.Game
)

func main() {
	// Initialize the application with config
	configPath := filepath.Join("config", "openhs.json")
	if err := bootstrap.Initialize(configPath); err != nil {
		fmt.Printf("Failed to initialize global components: %v\n", err)
		return
	}

	// Load the sample game
	gameManager := game.GetGameManager()
	g, err := gameManager.LoadGameByID("sample_game")
	if err != nil {
		fmt.Printf("Failed to load sample game: %v\n", err)
		return
	}

	// Create a new engine
	e := engine.NewEngine(g)
	
	// Start the game
	err = e.StartGame()
	if err != nil {
		fmt.Printf("Failed to start game: %v\n", err)
		return
	}

	gameEngine = e
	gameObj = g

	// Set up HTTP handlers
	http.HandleFunc("/api/game", gameStateHandler)
	http.HandleFunc("/api/action", actionHandler)
	
	// Serve static files
	fs := http.FileServer(http.Dir("frontend/static"))
	http.Handle("/", fs)

	// Start the server
	fmt.Println("Starting web server on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func gameStateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	gameState := convertGameState(gameObj)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(gameState)
}

func actionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var action struct {
		Type      string `json:"type"`
		CardIndex int    `json:"cardIndex"`
		Position  int    `json:"position"`
		Target    int    `json:"target"`
	}

	if err := json.NewDecoder(r.Body).Decode(&action); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var err error
	switch action.Type {
	case "playCard":
		// Use proper method available in engine with the correct parameters
		err = gameEngine.PlayCard(gameObj.CurrentPlayer, action.CardIndex, nil, action.Position, 0)
	case "attack":
		// Fix: Use proper Attack method signature
		var attacker, target *game.Entity
		
		// Get attacker from current player's field
		if action.CardIndex >= 0 && action.CardIndex < len(gameObj.CurrentPlayer.Field) {
			attacker = gameObj.CurrentPlayer.Field[action.CardIndex]
		}
		
		// Get target from opponent's field
		opponent := gameObj.Players[1-gameObj.CurrentPlayerIndex]
		if action.Target >= 0 && action.Target < len(opponent.Field) {
			target = opponent.Field[action.Target]
		}
		
		if attacker != nil && target != nil {
			err = gameEngine.Attack(attacker, target, false)
		} else {
			err = fmt.Errorf("invalid attacker or target")
		}
	case "endTurn":
		err = gameEngine.EndPlayerTurn()
	default:
		http.Error(w, "Invalid action type", http.StatusBadRequest)
		return
	}

	if err != nil {
		http.Error(w, fmt.Sprintf("Action failed: %v", err), http.StatusBadRequest)
		return
	}

	// Sleep briefly to simulate game processing
	time.Sleep(300 * time.Millisecond)

	// Return updated game state
	gameState := convertGameState(gameObj)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(gameState)
}

func convertGameState(g *game.Game) *GameState {
	gameState := &GameState{
		CurrentTurn:        g.CurrentTurn,
		Phase:              g.Phase.String(),
		CurrentPlayerIndex: g.CurrentPlayerIndex,
		Players:            make([]*SimplifiedPlayer, len(g.Players)),
		Actions:            []string{"playCard", "attack", "endTurn"},
	}

	// Convert both players
	for i, player := range g.Players {
		simplifiedPlayer := &SimplifiedPlayer{
			Hero: &SimplifiedEntity{
				Name:   player.Hero.Card.Name,
				Health: player.Hero.Health,
				Type:   "Hero",
			},
			Hand:      make([]*SimplifiedEntity, len(player.Hand)),
			Field:     make([]*SimplifiedEntity, len(player.Field)),
			Mana:      player.Mana,
			TotalMana: player.TotalMana,
		}

		// Convert hand
		for j, card := range player.Hand {
			simplifiedPlayer.Hand[j] = &SimplifiedEntity{
				Name:        card.Card.Name,
				Attack:      card.Attack,
				Health:      card.Health,
				Cost:        card.Card.Cost,
				Type:        card.Card.Type.String(),
				Description: card.Card.Description,
				Tags:        convertTagsToString(card.Tags),
			}
		}

		// Convert field
		for j, card := range player.Field {
			canAttack := !card.Exhausted && card.NumAttackThisTurn < getMaxAttacksPerTurn(card.Tags)
			simplifiedPlayer.Field[j] = &SimplifiedEntity{
				Name:        card.Card.Name,
				Attack:      card.Attack,
				Health:      card.Health,
				Cost:        card.Card.Cost,
				Type:        card.Card.Type.String(),
				Description: card.Card.Description,
				Tags:        convertTagsToString(card.Tags),
				CanAttack:   canAttack && i == g.CurrentPlayerIndex,
			}
		}

		// Add weapon if exists
		if player.Weapon != nil {
			simplifiedPlayer.Weapon = &SimplifiedEntity{
				Name:        player.Weapon.Card.Name,
				Attack:      player.Weapon.Attack,
				Health:      player.Weapon.Health,
				Cost:        player.Weapon.Card.Cost,
				Type:        "Weapon",
				Description: player.Weapon.Card.Description,
			}
		}

		gameState.Players[i] = simplifiedPlayer
	}

	return gameState
}

// Fix: Change to handle game.Tag type properly
func convertTagsToString(tags []game.Tag) []string {
	result := make([]string, len(tags))
	for i, tag := range tags {
		// Convert tag to string based on the tag value
		switch tag.Type {
		case game.TAG_TAUNT:
			result[i] = "Taunt"
		case game.TAG_DIVINE_SHIELD:
			result[i] = "Divine Shield"
		case game.TAG_CHARGE:
			result[i] = "Charge"
		case game.TAG_RUSH:
			result[i] = "Rush"
		case game.TAG_WINDFURY:
			result[i] = "Windfury"
		case game.TAG_POISONOUS:
			result[i] = "Poisonous"
		default:
			result[i] = fmt.Sprintf("Tag(%d)", int(tag.Type))
		}
	}
	return result
}

// Fix: Use proper TagType check
func getMaxAttacksPerTurn(tags []game.Tag) int {
	for _, tag := range tags {
		if tag.Type == game.TAG_WINDFURY {
			return 2
		}
	}
	return 1
} 