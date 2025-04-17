package tests

import (
	"testing"

	cards "github.com/openhs/cards/classic"
	"github.com/openhs/internal/engine"
	"github.com/openhs/internal/game"
)

var waterElementalCard *game.Card

func init() {
	(&cards.WaterElemental{}).Register(game.GetCardManager())
	waterElementalCard, _ = game.GetCardManager().CreateCardInstance("Water Elemental")
}

// TestWaterElementalProperties tests that Water Elemental has the correct properties
func TestWaterElementalProperties(t *testing.T) {
	// Setup
	g := game.CreateTestGame()
	engine := engine.NewEngine(g)
	engine.StartGame()
	player := g.Players[0]

	// Create a Water Elemental entity
	entity := game.NewEntity(waterElementalCard, g, player)

	// Verify the properties
	if entity.Card.Cost != 4 {
		t.Errorf("Expected Water Elemental cost to be 4, got %d", entity.Card.Cost)
	}
	if entity.Card.Attack != 3 {
		t.Errorf("Expected Water Elemental attack to be 3, got %d", entity.Card.Attack)
	}
	if entity.Card.Health != 6 {
		t.Errorf("Expected Water Elemental health to be 6, got %d", entity.Card.Health)
	}
	if entity.Card.Type != game.Minion {
		t.Errorf("Expected Water Elemental type to be Minion, got %s", entity.Card.Type)
	}
}

// TestWaterElementalFreezeEffect tests that Water Elemental freezes targets it damages
func TestWaterElementalFreezeEffect(t *testing.T) {
	// Setup
	g := game.CreateTestGame()
	engine := engine.NewEngine(g)
	engine.StartGame()

	player1 := g.Players[0]
	player2 := g.Players[1]

	// Set player1 mana to enough to cast Water Elemental
	player1.Mana = 10

	// Create a Water Elemental entity in player1's hand
	waterElementalEntity := game.NewEntity(waterElementalCard, g, player1)
	g.AddEntityToHand(player1, waterElementalEntity, -1)

	// Play Water Elemental
	err := g.PlayCard(player1, len(player1.Hand)-1, nil, -1, 0)
	if err != nil {
		t.Fatalf("Failed to play Water Elemental: %v", err)
	}

	// Check if Water Elemental is on the field
	if len(player1.Field) != 1 {
		t.Fatalf("Expected 1 minion on the field, got %d", len(player1.Field))
	}
	if player1.Field[0] != waterElementalEntity {
		t.Errorf("Expected Water Elemental to be on the field")
	}

	// Create a target minion for player2
	targetMinion := game.CreateTestMinionEntity(g, player2,
		game.WithName("Target Minion"),
		game.WithAttack(2),
		game.WithHealth(4))
	g.AddEntityToField(player2, targetMinion, 0)

	// End turn for player1
	engine.EndPlayerTurn()

	// End turn for player2
	engine.EndPlayerTurn()

	// Now player1 can attack with Water Elemental
	if waterElementalEntity.Exhausted {
		t.Errorf("Expected Water Elemental to be ready to attack")
	}

	// Attack the target minion with Water Elemental
	err = g.Attack(waterElementalEntity, targetMinion, false)
	if err != nil {
		t.Fatalf("Failed to attack with Water Elemental: %v", err)
	}

	// Check if target minion is frozen
	if !game.HasTag(targetMinion.Tags, game.TAG_FROZEN) {
		t.Errorf("Expected target minion to be frozen after being damaged by Water Elemental")
	}

	// End turn for both players to ready Water Elemental again
	engine.EndPlayerTurn()
	engine.EndPlayerTurn()

	// Attack the enemy hero with Water Elemental
	err = g.Attack(waterElementalEntity, player2.Hero, false)
	if err != nil {
		t.Fatalf("Failed to attack enemy hero with Water Elemental: %v", err)
	}

	// Check if enemy hero is frozen
	if !game.HasTag(player2.Hero.Tags, game.TAG_FROZEN) {
		t.Errorf("Expected enemy hero to be frozen after being damaged by Water Elemental")
	}
}

// TestWaterElementalFreezePersistence tests that a frozen target remains frozen until the end of its next turn
func TestWaterElementalFreezePersistence(t *testing.T) {
	// Setup
	g := game.CreateTestGame()
	engine := engine.NewEngine(g)
	engine.StartGame()

	player1 := g.Players[0]
	player2 := g.Players[1]

	// Set player1 mana to enough to cast Water Elemental
	player1.Mana = 10

	// Create a Water Elemental entity in player1's hand
	waterElementalEntity := game.NewEntity(waterElementalCard, g, player1)
	g.AddEntityToHand(player1, waterElementalEntity, -1)

	// Play Water Elemental
	err := g.PlayCard(player1, len(player1.Hand)-1, nil, -1, 0)
	if err != nil {
		t.Fatalf("Failed to play Water Elemental: %v", err)
	}

	// Create a target minion for player2
	targetMinion := game.CreateTestMinionEntity(g, player2,
		game.WithName("Target Minion"),
		game.WithAttack(2),
		game.WithHealth(4))
	g.AddEntityToField(player2, targetMinion, 0)

	// End turn for player1
	engine.EndPlayerTurn()

	// End turn for player2
	engine.EndPlayerTurn()

	// Attack the target minion with Water Elemental
	err = g.Attack(waterElementalEntity, targetMinion, false)
	if err != nil {
		t.Fatalf("Failed to attack with Water Elemental: %v", err)
	}

	// Check if target minion is frozen
	if !game.HasTag(targetMinion.Tags, game.TAG_FROZEN) {
		t.Errorf("Expected target minion to be frozen after being damaged by Water Elemental")
	}

	// End turn for player1
	engine.EndPlayerTurn()

	// Check if target minion is still frozen on player2's turn
	if !game.HasTag(targetMinion.Tags, game.TAG_FROZEN) {
		t.Errorf("Expected target minion to still be frozen on its controller's turn")
	}

	// Try to attack with the frozen minion - should fail
	err = g.Attack(targetMinion, player1.Hero, false)
	if err == nil {
		t.Errorf("Expected attack to fail due to freeze, but it succeeded")
	}

	// End turn for player2
	engine.EndPlayerTurn()

	// Check if target minion is unfrozen at the end of player2's next turn
	// If the minion was not exhausted at the end of the turn, it should be unfrozen
	if targetMinion.Exhausted {
		// If exhausted, should still be frozen
		if !game.HasTag(targetMinion.Tags, game.TAG_FROZEN) {
			t.Errorf("Expected exhausted minion to still be frozen")
		}
	} else {
		// If not exhausted, should be unfrozen
		if game.HasTag(targetMinion.Tags, game.TAG_FROZEN) {
			t.Errorf("Expected non-exhausted minion to be unfrozen after its controller's turn")
		}
	}
}

// TestNoFreezeWhenNotDamagedByWaterElemental tests that entities are not frozen when damaged by other sources
func TestNoFreezeWhenNotDamagedByWaterElemental(t *testing.T) {
	// Setup
	g := game.CreateTestGame()
	engine := engine.NewEngine(g)
	engine.StartGame()

	player1 := g.Players[0]
	player2 := g.Players[1]

	// Set player1 mana to enough to cast Water Elemental
	player1.Mana = 10

	// Create a Water Elemental entity in player1's hand
	waterElementalEntity := game.NewEntity(waterElementalCard, g, player1)
	g.AddEntityToHand(player1, waterElementalEntity, -1)

	// Play Water Elemental
	err := g.PlayCard(player1, len(player1.Hand)-1, nil, -1, 0)
	if err != nil {
		t.Fatalf("Failed to play Water Elemental: %v", err)
	}

	// Create another minion for player1
	otherMinion := game.CreateTestMinionEntity(g, player1,
		game.WithName("Other Minion"),
		game.WithAttack(2),
		game.WithHealth(2))
	g.AddEntityToField(player1, otherMinion, 1)

	// Create a target minion for player2
	targetMinion := game.CreateTestMinionEntity(g, player2,
		game.WithName("Target Minion"),
		game.WithAttack(2),
		game.WithHealth(4))
	g.AddEntityToField(player2, targetMinion, 0)

	// End turn for player1
	engine.EndPlayerTurn()

	// End turn for player2
	engine.EndPlayerTurn()

	// Attack the target minion with otherMinion (not Water Elemental)
	err = g.Attack(otherMinion, targetMinion, false)
	if err != nil {
		t.Fatalf("Failed to attack with other minion: %v", err)
	}

	// Check that target minion is NOT frozen
	if game.HasTag(targetMinion.Tags, game.TAG_FROZEN) {
		t.Errorf("Target minion should not be frozen when damaged by a non-Water Elemental minion")
	}
}
