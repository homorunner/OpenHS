package engine

import (
	"testing"

	"github.com/openhs/internal/game"
)

// TestHealCard tests the HealCard function
func TestHealCard(t *testing.T) {
	g := createTestGame()
	e := NewEngine(g)

	// Test 1: Healing with valid amount
	card := &game.Card{
		Health:    10,
		MaxHealth: 20,
	}

	e.HealCard(card, 5)

	if card.Health != 15 {
		t.Errorf("Expected health to be 15 after healing 5, got %d", card.Health)
	}

	// Test 2: Healing beyond max health should cap at max health
	e.HealCard(card, 10)
	if card.Health != card.MaxHealth {
		t.Errorf("Expected health to be capped at max health %d, got %d", card.MaxHealth, card.Health)
	}

	// Test 3: Healing a card already at max health
	e.HealCard(card, 5)
	if card.Health != card.MaxHealth {
		t.Errorf("Expected health to remain at max health %d, got %d", card.MaxHealth, card.Health)
	}
} 