package tests

import (
	"testing"

	core2025 "github.com/openhs/cards/core2025"
	"github.com/openhs/internal/engine"
	"github.com/openhs/internal/game"
)

var arcaneIntellectCard *game.Card

func init() {
	(&core2025.ArcaneIntellect{}).Register(game.GetCardManager())
	arcaneIntellectCard, _ = game.GetCardManager().CreateCardInstance("Arcane Intellect")
}

// TestArcaneIntellectProperties tests that Arcane Intellect has the correct properties
func TestArcaneIntellectProperties(t *testing.T) {
	// Setup
	g := game.CreateTestGame()
	engine := engine.NewEngine(g)
	engine.StartGame()
	player := g.Players[0]

	// Create an Arcane Intellect entity
	entity := game.NewEntity(arcaneIntellectCard, g, player)

	// Verify the properties
	if entity.Card.Cost != 3 {
		t.Errorf("Expected Arcane Intellect cost to be 3, got %d", entity.Card.Cost)
	}
	if entity.Card.Type != game.Spell {
		t.Errorf("Expected Arcane Intellect type to be Spell, got %s", entity.Card.Type)
	}
}

// TestArcaneIntellectEffect tests the main effect of Arcane Intellect (drawing 2 cards)
func TestArcaneIntellectEffect(t *testing.T) {
	// Setup
	g := game.CreateTestGame()
	engine := engine.NewEngine(g)
	engine.StartGame()
	
	player := g.Players[0]
	
	// Set player mana to enough to cast Arcane Intellect
	player.Mana = 10
	
	// Create an Arcane Intellect entity in player's hand
	arcaneIntellectEntity := game.NewEntity(arcaneIntellectCard, g, player)
	g.AddEntityToHand(player, arcaneIntellectEntity, -1)
	
	// Count cards in deck and hand before casting
	initialDeckCount := len(player.Deck)
	initialHandCount := len(player.Hand)
	
	// Cast Arcane Intellect
	err := g.PlayCard(player, len(player.Hand)-1, nil, -1, 0)
	if err != nil {
		t.Fatalf("Failed to play Arcane Intellect: %v", err)
	}
	
	// Check if 2 cards were drawn
	expectedDeckCount := initialDeckCount - 2
	if len(player.Deck) != expectedDeckCount {
		t.Errorf("Expected deck count to be %d after Arcane Intellect, got %d", 
			expectedDeckCount, len(player.Deck))
	}
	
	// Check if hand has 1 more card (+2 drawn, -1 played)
	expectedHandCount := initialHandCount + 2 - 1
	if len(player.Hand) != expectedHandCount {
		t.Errorf("Expected hand count to be %d after Arcane Intellect, got %d", 
			expectedHandCount, len(player.Hand))
	}
	
	// Check if mana was spent
	if player.Mana != 7 {
		t.Errorf("Expected player mana to be 7 after casting Arcane Intellect, got %d", player.Mana)
	}
	
	// Check if Arcane Intellect was moved to graveyard
	if len(player.Graveyard) != 1 {
		t.Errorf("Expected 1 card in graveyard, got %d", len(player.Graveyard))
	}
	
	if player.Graveyard[0] != arcaneIntellectEntity {
		t.Errorf("Expected Arcane Intellect card to be in graveyard")
	}
	
	if arcaneIntellectEntity.CurrentZone != game.ZONE_GRAVEYARD {
		t.Errorf("Expected Arcane Intellect zone to be GRAVEYARD, got %s", arcaneIntellectEntity.CurrentZone)
	}
}

// TestArcaneIntellectEmptyDeck tests what happens when there aren't enough cards in the deck
func TestArcaneIntellectEmptyDeck(t *testing.T) {
	// Setup
	g := game.CreateTestGame()
	engine := engine.NewEngine(g)
	engine.StartGame()
	
	player := g.Players[0]
	
	// Set player mana to enough to cast Arcane Intellect
	player.Mana = 10
	
	// Record the hero's initial health
	initialHeroHealth := player.Hero.Health
	
	// Empty the deck except for 1 card
	for len(player.Deck) > 1 {
		player.Deck = player.Deck[:len(player.Deck)-1]
	}
	
	// Create an Arcane Intellect entity in player's hand
	arcaneIntellectEntity := game.NewEntity(arcaneIntellectCard, g, player)
	g.AddEntityToHand(player, arcaneIntellectEntity, -1)
	
	// Count cards in hand before casting
	initialHandCount := len(player.Hand)
	
	// Cast Arcane Intellect
	err := g.PlayCard(player, len(player.Hand)-1, nil, -1, 0)
	if err != nil {
		t.Fatalf("Failed to play Arcane Intellect: %v", err)
	}
	
	// Check if deck is now empty
	if len(player.Deck) != 0 {
		t.Errorf("Expected deck to be empty after Arcane Intellect, got %d cards", len(player.Deck))
	}
	
	// Check if hand has correct number of cards (+1 drawn, -1 played)
	expectedHandCount := initialHandCount
	if len(player.Hand) != expectedHandCount {
		t.Errorf("Expected hand count to be %d after Arcane Intellect, got %d", 
			expectedHandCount, len(player.Hand))
	}
	
	// Check if Arcane Intellect was moved to graveyard
	if len(player.Graveyard) != 1 {
		t.Errorf("Expected 1 card in graveyard, got %d", len(player.Graveyard))
	}
	
	// Check fatigue damage (1 damage for the first fatigue)
	expectedHealthAfterFatigue := initialHeroHealth - 1
	if player.Hero.Health != expectedHealthAfterFatigue {
		t.Errorf("Expected hero health to be %d after taking fatigue damage, got %d", 
			expectedHealthAfterFatigue, player.Hero.Health)
	}
	
	// Check that fatigue counter was increased
	if player.FatigueDamage != 1 {
		t.Errorf("Expected fatigue damage counter to be 1, got %d", player.FatigueDamage)
	}
} 