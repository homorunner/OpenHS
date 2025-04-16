package tests

import (
	"testing"

	core2025 "github.com/openhs/cards/core2025"
	"github.com/openhs/internal/engine"
	"github.com/openhs/internal/game"
)

var fireballCard *game.Card

func init() {
	(&core2025.Fireball{}).Register(game.GetCardManager())
	fireballCard, _ = game.GetCardManager().CreateCardInstance("Fireball")
}

// TestFireballProperties tests that Fireball has the correct properties
func TestFireballProperties(t *testing.T) {
	// Setup
	g := game.CreateTestGame()
	engine := engine.NewEngine(g)
	engine.StartGame()
	player := g.Players[0]

	// Create a Fireball entity
	entity := game.NewEntity(fireballCard, g, player)

	// Verify the properties
	if entity.Card.Cost != 4 {
		t.Errorf("Expected Fireball cost to be 4, got %d", entity.Card.Cost)
	}
	if entity.Card.Type != game.Spell {
		t.Errorf("Expected Fireball type to be Spell, got %s", entity.Card.Type)
	}
}

// TestFireballEffect tests the main effect of Fireball
func TestFireballEffect(t *testing.T) {
	// Setup
	g := game.CreateTestGame()
	engine := engine.NewEngine(g)
	engine.StartGame()
	
	player1 := g.Players[0]
	player2 := g.Players[1]
	
	// Set player1 mana to enough to cast Fireball
	player1.Mana = 10
	
	// Create a Fireball entity in player1's hand
	fireballEntity := game.NewEntity(fireballCard, g, player1)
	g.AddEntityToHand(player1, fireballEntity, -1)
	
	// Set initial health of player2's hero
	initialHealth := player2.Hero.Health
	
	// Cast Fireball targeting enemy hero
	err := g.PlayCard(player1, len(player1.Hand)-1, player2.Hero, -1, 0)
	if err != nil {
		t.Fatalf("Failed to play Fireball: %v", err)
	}
	
	// Check if 6 damage was dealt to enemy hero
	expectedHealth := initialHealth - 6
	if player2.Hero.Health != expectedHealth {
		t.Errorf("Expected enemy hero health to be %d after Fireball, got %d", 
			expectedHealth, player2.Hero.Health)
	}
	
	// Check if mana was spent
	if player1.Mana != 6 {
		t.Errorf("Expected player mana to be 6 after casting Fireball, got %d", player1.Mana)
	}
	
	// Check if Fireball was moved to graveyard
	if len(player1.Graveyard) != 1 {
		t.Errorf("Expected 1 card in graveyard, got %d", len(player1.Graveyard))
	}
	
	if player1.Graveyard[0] != fireballEntity {
		t.Errorf("Expected Fireball card to be in graveyard")
	}
	
	if fireballEntity.CurrentZone != game.ZONE_GRAVEYARD {
		t.Errorf("Expected Fireball zone to be GRAVEYARD, got %s", fireballEntity.CurrentZone)
	}
}

// TestFireballTargetingMinion tests using Fireball to damage a minion
func TestFireballTargetingMinion(t *testing.T) {
	// Setup
	g := game.CreateTestGame()
	engine := engine.NewEngine(g)
	engine.StartGame()
	
	player1 := g.Players[0]
	player2 := g.Players[1]
	
	// Set player1 mana to enough to cast Fireball
	player1.Mana = 10
	
	// Create a Fireball entity in player1's hand
	fireballEntity := game.NewEntity(fireballCard, g, player1)
	g.AddEntityToHand(player1, fireballEntity, -1)
	
	// Create a test minion on player2's field
	targetMinion := game.CreateTestMinionEntity(g, player2, game.WithHealth(7))
	g.AddEntityToField(player2, targetMinion, -1)
	
	// Cast Fireball targeting the minion
	err := g.PlayCard(player1, len(player1.Hand)-1, targetMinion, -1, 0)
	if err != nil {
		t.Fatalf("Failed to play Fireball: %v", err)
	}
	
	// Check if 6 damage was dealt to minion
	expectedHealth := 7 - 6
	if targetMinion.Health != expectedHealth {
		t.Errorf("Expected minion health to be %d after Fireball, got %d", 
			expectedHealth, targetMinion.Health)
	}
	
	// Check if minion survived (since it had 7 health)
	if targetMinion.CurrentZone != game.ZONE_PLAY {
		t.Errorf("Expected minion to remain on the field, but it's in zone %s", 
			targetMinion.CurrentZone)
	}
	
	// Create another test with a minion that will die
	player1.Mana = 10
	
	// Create a second Fireball entity in player1's hand
	fireballEntity2 := game.NewEntity(fireballCard, g, player1)
	g.AddEntityToHand(player1, fireballEntity2, -1)
	
	// Create a test minion with only 5 health on player2's field
	weakMinion := game.CreateTestMinionEntity(g, player2, game.WithHealth(5))
	g.AddEntityToField(player2, weakMinion, -1)
	
	// Cast Fireball targeting the weak minion
	err = g.PlayCard(player1, len(player1.Hand)-1, weakMinion, -1, 0)
	if err != nil {
		t.Fatalf("Failed to play second Fireball: %v", err)
	}
	
	// Check if the minion's health is reduced to 0 or less
	if weakMinion.Health > 0 {
		t.Errorf("Expected minion health to be <= 0 after Fireball, got %d", 
			weakMinion.Health)
	}

	// Check if the minion is moved to GRAVEYARD
	if weakMinion.CurrentZone != game.ZONE_GRAVEYARD {
		t.Errorf("Expected minion to be in GRAVEYARD, got %s", weakMinion.CurrentZone)
	}
}
