package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/openhs/internal/bootstrap"
	"github.com/openhs/internal/engine"
	"github.com/openhs/internal/game"
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

	// Start CLI game loop
	scanner := bufio.NewScanner(os.Stdin)
	running := true

	for running {
		// Display game state information
		displayGameState(g)

		fmt.Println("\nCommands:")
		fmt.Println("  p <card_number> [<position>] - Play a card from your hand")
		fmt.Println("  e - End your turn")
		fmt.Println("  q - Quit the game")
		fmt.Print("\nEnter command: ")

		if !scanner.Scan() {
			break
		}

		command := scanner.Text()
		parts := strings.Fields(command)

		if len(parts) == 0 {
			continue
		}

		switch parts[0] {
		case "p":
			handlePlayCard(e, g, parts)
		case "e":
			e.EndPlayerTurn()
		case "q":
			running = false
			fmt.Println("Thanks for playing!")
		default:
			fmt.Println("Unknown command")
		}
	}
}

func displayGameState(g *game.Game) {
	fmt.Printf("\n=== Game State ===\n")
	fmt.Printf("Current Turn: %d\n", g.CurrentTurn)
	fmt.Printf("Current Phase: %s\n", g.Phase.String())
	fmt.Printf("Current Player: %s (%s)\n", g.CurrentPlayer.Hero.Name, []string{"First", "Second"}[g.CurrentPlayerIndex])
	fmt.Printf("Player Mana: %d/%d\n", g.CurrentPlayer.Mana, g.CurrentPlayer.TotalMana)
	fmt.Printf("Player Health: %d\n", g.CurrentPlayer.Hero.Health)
	if g.CurrentPlayer.FatigueDamage > 0 {
		fmt.Printf("Player Fatigue Damage: %d\n", g.CurrentPlayer.FatigueDamage)
	}

	// Print the current player's field
	fmt.Printf("\nCurrent Player's Field (%d cards):\n", len(g.CurrentPlayer.Field))
	if len(g.CurrentPlayer.Field) == 0 {
		fmt.Println("(Empty field)")
	} else {
		for i, card := range g.CurrentPlayer.Field {
			fmt.Printf("  %d. %s (%d/%d)\n", i+1, card.Name, card.Attack, card.Health)
		}
	}

	// Print weapon if equipped
	if g.CurrentPlayer.HasWeapon {
		fmt.Printf("\nEquipped Weapon: %s (%d/%d)\n",
			g.CurrentPlayer.Weapon.Name,
			g.CurrentPlayer.Weapon.Attack,
			g.CurrentPlayer.Weapon.Health)
	}

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
			if card.Type == game.Minion {
				cardInfo = append(cardInfo, fmt.Sprintf("Attack: %d, Health: %d", card.Attack, card.Health))
			} else if card.Type == game.Weapon {
				cardInfo = append(cardInfo, fmt.Sprintf("Attack: %d, Durability: %d", card.Attack, card.Health))
			}

			fmt.Printf("  %d. %s\n", i+1, strings.Join(cardInfo, ", "))
		}
	}
}

func handlePlayCard(e *engine.Engine, g *game.Game, parts []string) {
	if len(parts) < 2 {
		fmt.Println("Error: Please specify a card number")
		return
	}

	// Parse card index
	cardNum, err := strconv.Atoi(parts[1])
	if err != nil {
		fmt.Println("Error: Invalid card number")
		return
	}

	// Adjust for 0-indexed hand
	handIndex := cardNum - 1

	if handIndex < 0 || handIndex >= len(g.CurrentPlayer.Hand) {
		fmt.Println("Error: Card number out of range")
		return
	}

	// Get the card to play
	card := g.CurrentPlayer.Hand[handIndex]

	// Default position is -1 (auto-position)
	position := -1

	// If position is specified and the card is a minion
	if len(parts) >= 3 && card.Type == game.Minion {
		pos, err := strconv.Atoi(parts[2])
		if err == nil && pos > 0 && pos <= len(g.CurrentPlayer.Field)+1 {
			position = pos - 1
		}
	}

	// For now, we don't handle targeting or choose one effects
	var target *game.Card = nil
	chooseOne := 0

	// Play the card
	err = e.PlayCard(g.CurrentPlayer, handIndex, target, position, chooseOne)
	if err != nil {
		fmt.Printf("Error playing card: %v\n", err)
		return
	}

	fmt.Printf("Played %s successfully!\n", card.Name)
}
