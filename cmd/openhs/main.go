package main

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/openhs/internal/bootstrap"
	"github.com/openhs/internal/engine"
	"github.com/openhs/internal/game"
	"github.com/openhs/internal/types"
)

func main() {
	// Initialize the application with config
	configPath := filepath.Join("config", "openhs.json")
	if err := bootstrap.Initialize(configPath); err != nil {
		fmt.Printf("Failed to initialize global components: %v\n", err)
		return
	}

	fmt.Println("OpenHS - Hearthstone Simulator Core")
	fmt.Println("Loading sample game...")

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

	// Display game state information
	fmt.Printf("\nGame State:\n")
	fmt.Printf("Current Turn: %d\n", g.CurrentTurn)
	fmt.Printf("Current Phase: %s\n", g.Phase.String())
	fmt.Printf("Current Player Index: %d\n", g.CurrentPlayerIndex)

	// Print the current player's hand
	fmt.Printf("\nCurrent Player's Hand (%d cards):\n", len(g.CurrentPlayer.Hand))
	if len(g.CurrentPlayer.Hand) == 0 {
		fmt.Println("(Empty hand)")
	} else {
		for i, card := range g.CurrentPlayer.Hand {
			cardInfo := []string{
				card.Type.String(),
				card.Name,
			}
			
			// Add cost
			if card.Cost > 0 {
				cardInfo = append(cardInfo, fmt.Sprintf("Cost: %d", card.Cost))
			}
			
			// Add attack and health for minions and weapons
			if card.Type == types.Minion {
				cardInfo = append(cardInfo, fmt.Sprintf("Attack: %d, Health: %d", card.Attack, card.Health))
			} else if card.Type == types.Weapon {
				cardInfo = append(cardInfo, fmt.Sprintf("Attack: %d, Durability: %d", card.Attack, card.Health))
			}
			
			fmt.Printf("  %d. %s\n", i+1, strings.Join(cardInfo, ", "))
		}
	}
}