package game

import (
	"testing"
)

// TestDealDamage tests the Game.DealDamage function
func TestDealDamage(t *testing.T) {
	g := CreateTestGame()

	// Test 1: Dealing damage to a minion
	player := g.Players[0]
	entity := CreateTestMinionEntity(g, player, WithHealth(20))
	entity.Health = 20

	g.DealDamage(nil, entity, 5)

	if entity.Health != 15 {
		t.Errorf("Expected health to be 15 after dealing 5 damage, got %d", entity.Health)
	}

	// Test 2: Dealing zero damage should have no effect
	healthBefore := entity.Health
	g.DealDamage(nil, entity, 0)
	if entity.Health != healthBefore {
		t.Errorf("Expected health to remain %d after dealing 0 damage, got %d", healthBefore, entity.Health)
	}

	// Test 3: Dealing negative damage should have no effect
	healthBefore = entity.Health
	g.DealDamage(nil, entity, -5)
	if entity.Health != healthBefore {
		t.Errorf("Expected health to remain %d after dealing negative damage, got %d", healthBefore, entity.Health)
	}
}

// TestHeal tests the Game.Heal function
func TestHeal(t *testing.T) {
	g := CreateTestGame()

	// Test 1: Healing with valid amount
	player := g.Players[0]
	entity := CreateTestMinionEntity(g, player, WithHealth(20))
	entity.Health = 10 // Set current health to 10

	g.Heal(nil, entity, 5)

	if entity.Health != 15 {
		t.Errorf("Expected health to be 15 after healing 5, got %d", entity.Health)
	}

	// Test 2: Healing beyond max health should cap at max health
	g.Heal(nil, entity, 10)
	if entity.Health != entity.MaxHealth {
		t.Errorf("Expected health to be capped at max health %d, got %d", entity.MaxHealth, entity.Health)
	}

	// Test 3: Healing a card already at max health
	g.Heal(nil, entity, 5)
	if entity.Health != entity.MaxHealth {
		t.Errorf("Expected health to remain at max health %d, got %d", entity.MaxHealth, entity.Health)
	}
}

// TestSetHealth tests the Game.SetHealth function
func TestSetHealth(t *testing.T) {
	g := CreateTestGame()

	// Test: Setting health should update both current and max health
	player := g.Players[0]
	entity := CreateTestMinionEntity(g, player, WithHealth(5))
	
	g.SetHealth(entity, 10)
	
	if entity.Health != 10 || entity.MaxHealth != 10 {
		t.Errorf("Expected both Health and MaxHealth to be 10, got Health=%d, MaxHealth=%d", 
			entity.Health, entity.MaxHealth)
	}
} 