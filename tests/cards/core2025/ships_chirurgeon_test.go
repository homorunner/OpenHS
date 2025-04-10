package tests

import (
	"testing"

	core2025 "github.com/openhs/internal/cards/core2025"
	"github.com/openhs/internal/engine"
	"github.com/openhs/internal/game"
	"github.com/openhs/internal/game/test"
)

var card *game.Card

func init() {
	(&core2025.ShipsChirurgeon{}).Register(game.GetCardManager())
	card, _ = game.GetCardManager().CreateCardInstance("Ship's Chirurgeon")
}

// TestShipsChirurgeonProperties tests that Ship's Chirurgeon has the correct properties
func TestShipsChirurgeonProperties(t *testing.T) {
	// Setup
	g := test.CreateTestGame()
	engine := engine.NewEngine(g)
	engine.StartGame()
	player := g.Players[0]

	// Create a Ship's Chirurgeon entity
	entity := game.NewEntity(card, g, player)

	// Verify the properties
	if entity.Card.Cost != 1 {
		t.Errorf("Expected Ship's Chirurgeon cost to be 1, got %d", entity.Card.Cost)
	}
	if entity.Attack != 1 {
		t.Errorf("Expected Ship's Chirurgeon attack to be 1, got %d", entity.Attack)
	}
	if entity.Health != 2 {
		t.Errorf("Expected Ship's Chirurgeon health to be 2, got %d", entity.Health)
	}
}

// TestShipsChirurgeonEffect tests the main effect of Ship's Chirurgeon
func TestShipsChirurgeonEffect(t *testing.T) {
	t.Run("Ship's Chirurgeon gives +1 health to friendly minions when summoned", func(t *testing.T) {
		// Setup
		g := test.CreateTestGame()
		engine := engine.NewEngine(g)
		engine.StartGame()
		player := g.Players[0]

		// Create Ship's Chirurgeon entity with its real name
		entity := game.NewEntity(card, g, player)

		// Create a test minion to be played after Ship's Chirurgeon
		testMinion := test.CreateTestMinionEntity(g, player,
			test.WithName("Test Minion"),
			test.WithCost(1),
			test.WithAttack(2),
			test.WithHealth(3))

		// Add cards to player's hand
		player.Hand = []*game.Entity{entity, testMinion}
		player.Mana = 10 // Ensure enough mana

		// Play Ship's Chirurgeon
		err := engine.PlayCard(player, 0, nil, -1, 0)
		if err != nil {
			t.Fatalf("Failed to play Ship's Chirurgeon: %v", err)
		}

		// Check that Ship's Chirurgeon is on the field
		if len(player.Field) != 1 {
			t.Fatalf("Expected 1 minion on player's field, got %d", len(player.Field))
		}

		// Play the test minion
		err = engine.PlayCard(player, 0, nil, -1, 0)
		if err != nil {
			t.Fatalf("Failed to play test minion: %v", err)
		}

		// Check that the test minion is on the field
		if len(player.Field) != 2 {
			t.Fatalf("Expected 2 minions on player's field, got %d", len(player.Field))
		}

		// Get the played test minion
		playedMinion := player.Field[1]

		// Verify that the test minion received +1 health
		if playedMinion.Health != 4 {
			t.Errorf("Expected test minion health to be 4 (3 + 1), got %d", playedMinion.Health)
		}
		if playedMinion.MaxHealth != 4 {
			t.Errorf("Expected test minion max health to be 4 (3 + 1), got %d", playedMinion.MaxHealth)
		}
	})

	t.Run("Ship's Chirurgeon doesn't give +1 health to opponent's minions", func(t *testing.T) {
		// Setup
		g := test.CreateTestGame()
		engine := engine.NewEngine(g)
		engine.StartGame()
		player := g.Players[0]
		opponent := g.Players[1]

		// Create Ship's Chirurgeon entity
		entity := game.NewEntity(card, g, player)

		// Create a test minion for the opponent
		opponentMinion := test.CreateTestMinionEntity(g, opponent,
			test.WithName("Opponent Minion"),
			test.WithCost(1),
			test.WithAttack(2),
			test.WithHealth(3))

		// Add Ship's Chirurgeon to player's hand
		player.Hand = []*game.Entity{entity}
		player.Mana = 10 // Ensure enough mana

		// Add opponent's minion to their hand
		opponent.Hand = []*game.Entity{opponentMinion}
		opponent.Mana = 10 // Ensure enough mana

		// Play Ship's Chirurgeon
		err := engine.PlayCard(player, 0, nil, -1, 0)
		if err != nil {
			t.Fatalf("Failed to play Ship's Chirurgeon: %v", err)
		}

		// Play opponent's minion
		g.CurrentPlayerIndex = 1
		g.CurrentPlayer = g.Players[1]
		err = engine.PlayCard(opponent, 0, nil, -1, 0)
		if err != nil {
			t.Fatalf("Failed to play opponent's minion: %v", err)
		}

		// Get the played opponent's minion
		opponentPlayedMinion := opponent.Field[0]

		// Verify that opponent's minion did NOT receive +1 health (should stay at 3)
		if opponentPlayedMinion.Health != 3 {
			t.Errorf("Expected opponent minion health to be 3 (unchanged), got %d", opponentPlayedMinion.Health)
		}
		if opponentPlayedMinion.MaxHealth != 3 {
			t.Errorf("Expected opponent minion max health to be 3 (unchanged), got %d", opponentPlayedMinion.MaxHealth)
		}
	})

	t.Run("Ship's Chirurgeon effect doesn't apply to itself", func(t *testing.T) {
		// Setup
		g := test.CreateTestGame()
		engine := engine.NewEngine(g)
		engine.StartGame()
		player := g.Players[0]

		// Create Ship's Chirurgeon entity
		entity := game.NewEntity(card, g, player)

		// Add Ship's Chirurgeon to player's hand
		player.Hand = []*game.Entity{entity}
		player.Mana = 10 // Ensure enough mana

		// Play Ship's Chirurgeon
		err := engine.PlayCard(player, 0, nil, -1, 0)
		if err != nil {
			t.Fatalf("Failed to play Ship's Chirurgeon: %v", err)
		}

		// Get the played Ship's Chirurgeon
		playedChirurgeon := player.Field[0]

		// Verify that Ship's Chirurgeon did NOT give itself +1 health (should stay at 2)
		if playedChirurgeon.Health != 2 {
			t.Errorf("Expected Ship's Chirurgeon health to be 2 (unchanged), got %d", playedChirurgeon.Health)
		}
		if playedChirurgeon.MaxHealth != 2 {
			t.Errorf("Expected Ship's Chirurgeon max health to be 2 (unchanged), got %d", playedChirurgeon.MaxHealth)
		}
	})

	t.Run("Multiple Ship's Chirurgeons give multiple +1 health bonuses", func(t *testing.T) {
		// Setup
		g := test.CreateTestGame()
		engine := engine.NewEngine(g)
		engine.StartGame()
		player := g.Players[0]

		// Create two Ship's Chirurgeon entities
		entity1 := game.NewEntity(card, g, player)
		entity2 := game.NewEntity(card, g, player)

		// Create a test minion
		testMinion := test.CreateTestMinionEntity(g, player,
			test.WithName("Test Minion"),
			test.WithCost(1),
			test.WithAttack(2),
			test.WithHealth(3))

		// Add cards to player's hand
		player.Hand = []*game.Entity{entity1, entity2, testMinion}
		player.Mana = 10 // Ensure enough mana

		// Play first Ship's Chirurgeon
		err := engine.PlayCard(player, 0, nil, -1, 0)
		if err != nil {
			t.Fatalf("Failed to play first Ship's Chirurgeon: %v", err)
		}

		// Play second Ship's Chirurgeon
		err = engine.PlayCard(player, 0, nil, -1, 0)
		if err != nil {
			t.Fatalf("Failed to play second Ship's Chirurgeon: %v", err)
		}

		// Play the test minion
		err = engine.PlayCard(player, 0, nil, -1, 0)
		if err != nil {
			t.Fatalf("Failed to play test minion: %v", err)
		}

		// Get the played test minion
		playedMinion := player.Field[2]

		// Verify that the test minion received +2 health (1 from each Ship's Chirurgeon)
		if playedMinion.Health != 5 {
			t.Errorf("Expected test minion health to be 5 (3 + 1 + 1), got %d", playedMinion.Health)
		}
		if playedMinion.MaxHealth != 5 {
			t.Errorf("Expected test minion max health to be 5 (3 + 1 + 1), got %d", playedMinion.MaxHealth)
		}
	})
}
