package engine

import (
	"testing"

	"github.com/openhs/internal/game/test"
)

// TestHealCard tests the HealCard function
func TestHealCard(t *testing.T) {
	g := test.CreateTestGame()
	e := NewEngine(g)
	e.StartGame()

	// Test 1: Healing with valid amount
	player := g.Players[0]
	entity := test.CreateTestMinionEntity(g, player, test.WithHealth(20))
	entity.Health = 10 // Set current health to 10

	e.Heal(nil, entity, 5)

	if entity.Health != 15 {
		t.Errorf("Expected health to be 15 after healing 5, got %d", entity.Health)
	}

	// Test 2: Healing beyond max health should cap at max health
	e.Heal(nil, entity, 10)
	if entity.Health != entity.MaxHealth {
		t.Errorf("Expected health to be capped at max health %d, got %d", entity.MaxHealth, entity.Health)
	}

	// Test 3: Healing a card already at max health
	e.Heal(nil, entity, 5)
	if entity.Health != entity.MaxHealth {
		t.Errorf("Expected health to remain at max health %d, got %d", entity.MaxHealth, entity.Health)
	}
}
