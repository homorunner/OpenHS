package tests

import (
	"testing"

	core2025 "github.com/openhs/cards/core2025"
	"github.com/openhs/internal/engine"
	"github.com/openhs/internal/game"
)

var frostboltCard *game.Card

func init() {
	(&core2025.Frostbolt{}).Register(game.GetCardManager())
	frostboltCard, _ = game.GetCardManager().CreateCardInstance("Frostbolt")
}

// TestFrostboltProperties tests that Frostbolt has the correct properties
func TestFrostboltProperties(t *testing.T) {
	// Setup
	g := game.CreateTestGame()
	engine := engine.NewEngine(g)
	engine.StartGame()
	player := g.Players[0]

	// Create a Frostbolt entity
	entity := game.NewEntity(frostboltCard, g, player)

	// Verify the properties
	if entity.Card.Cost != 2 {
		t.Errorf("Expected Frostbolt cost to be 2, got %d", entity.Card.Cost)
	}
	if entity.Card.Type != game.Spell {
		t.Errorf("Expected Frostbolt type to be Spell, got %s", entity.Card.Type)
	}
}

// TestFrostboltEffect tests the main effect of Frostbolt (dealing 3 damage and freezing)
func TestFrostboltEffect(t *testing.T) {
	// Setup
	g := game.CreateTestGame()
	engine := engine.NewEngine(g)
	engine.StartGame()

	player1 := g.Players[0]
	player2 := g.Players[1]

	// Set player mana to enough to cast Frostbolt
	player1.Mana = 10

	// Create a Frostbolt entity in player's hand
	frostboltEntity := game.NewEntity(frostboltCard, g, player1)
	g.AddEntityToHand(player1, frostboltEntity, -1)

	// Create a target minion for player2
	targetMinion := game.CreateTestMinionEntity(g, player2,
		game.WithName("Target Minion"),
		game.WithHealth(5))
	g.AddEntityToField(player2, targetMinion, 0)

	// Record initial health
	initialHealth := targetMinion.Health

	// Cast Frostbolt on the minion
	err := g.PlayCard(player1, len(player1.Hand)-1, targetMinion, -1, 0)
	if err != nil {
		t.Fatalf("Failed to play Frostbolt: %v", err)
	}

	// Check if damage was applied
	expectedHealth := initialHealth - 3
	if targetMinion.Health != expectedHealth {
		t.Errorf("Expected target health to be %d after Frostbolt, got %d",
			expectedHealth, targetMinion.Health)
	}

	// Check if target was frozen
	if !game.HasTag(targetMinion.Tags, game.TAG_FROZEN) {
		t.Errorf("Expected target to be frozen after Frostbolt")
	}

	// Check if mana was spent
	if player1.Mana != 8 {
		t.Errorf("Expected player mana to be 8 after casting Frostbolt, got %d", player1.Mana)
	}

	// Check if Frostbolt was moved to graveyard
	if len(player1.Graveyard) != 1 {
		t.Errorf("Expected 1 card in graveyard, got %d", len(player1.Graveyard))
	}

	if player1.Graveyard[0] != frostboltEntity {
		t.Errorf("Expected Frostbolt card to be in graveyard")
	}
}

// TestFrostboltHero tests using Frostbolt on a hero
func TestFrostboltHero(t *testing.T) {
	// Setup
	g := game.CreateTestGame()
	engine := engine.NewEngine(g)
	engine.StartGame()

	player1 := g.Players[0]
	player2 := g.Players[1]

	// Set player mana to enough to cast Frostbolt
	player1.Mana = 10

	// Create a Frostbolt entity in player's hand
	frostboltEntity := game.NewEntity(frostboltCard, g, player1)
	g.AddEntityToHand(player1, frostboltEntity, -1)

	// Record initial hero health
	initialHealth := player2.Hero.Health

	// Cast Frostbolt on enemy hero
	err := g.PlayCard(player1, len(player1.Hand)-1, player2.Hero, -1, 0)
	if err != nil {
		t.Fatalf("Failed to play Frostbolt on hero: %v", err)
	}

	// Check if damage was applied
	expectedHealth := initialHealth - 3
	if player2.Hero.Health != expectedHealth {
		t.Errorf("Expected hero health to be %d after Frostbolt, got %d",
			expectedHealth, player2.Hero.Health)
	}

	// Check if hero was frozen
	if !game.HasTag(player2.Hero.Tags, game.TAG_FROZEN) {
		t.Errorf("Expected hero to be frozen after Frostbolt")
	}
}
