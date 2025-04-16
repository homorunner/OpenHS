package game

import (
	"testing"
)

// TestDrawCard tests the Game.DrawCard function
func TestDrawCard(t *testing.T) {
	g := CreateTestGame()
	// No need for StartGame, the test game is already set up

	// Test 1: Normal card draw
	player := g.Players[0]
	initialHandSize := len(player.Hand)
	initialDeckSize := len(player.Deck)

	drawnEntity := g.DrawCard(player)

	// Verify hand increased by 1
	if len(player.Hand) != initialHandSize+1 {
		t.Fatalf("Expected hand size to increase by 1, got %d (was %d)",
			len(player.Hand), initialHandSize)
	}

	// Verify deck decreased by 1
	if len(player.Deck) != initialDeckSize-1 {
		t.Fatalf("Expected deck size to decrease by 1, got %d (was %d)",
			len(player.Deck), initialDeckSize)
	}

	// Verify the card drawn is the last one from the deck
	lastCardName := "Test Card" // The name we used when creating test cards
	if player.Hand[len(player.Hand)-1].Card.Name != lastCardName {
		t.Fatalf("Expected drawn card name to be %s, got %s",
			lastCardName, player.Hand[len(player.Hand)-1].Card.Name)
	}

	// Verify drawn card reference is returned
	if drawnEntity == nil {
		t.Fatal("Expected drawn entity to be returned, got nil")
	}
	if drawnEntity.Card.Name != lastCardName {
		t.Fatalf("Expected drawn card name to be %s, got %s",
			lastCardName, drawnEntity.Card.Name)
	}

	// Test 2: Drawing from an empty deck should trigger fatigue damage
	// Create a test hero card and entity
	emptyPlayer := NewPlayer()
	emptyPlayer.Hero = CreateTestHeroEntity(g, emptyPlayer, WithName("Test Hero"), WithHealth(30))

	// Try to draw from empty deck
	drawn := g.DrawCard(emptyPlayer)

	// Verify nil is returned for an empty deck
	if drawn != nil {
		t.Fatalf("Expected nil to be returned when drawing from empty deck, got %v", drawn)
	}

	// Verify fatigue counter increases
	if emptyPlayer.FatigueDamage != 1 {
		t.Fatalf("Expected fatigue damage to be 1, got %d", emptyPlayer.FatigueDamage)
	}

	// Verify hero took fatigue damage
	if emptyPlayer.Hero.Health != 29 {
		t.Fatalf("Expected hero to take 1 fatigue damage, health is %d instead of 29", emptyPlayer.Hero.Health)
	}

	// Draw again to verify fatigue damage increases
	g.DrawCard(emptyPlayer)

	// Verify fatigue damage increased
	if emptyPlayer.FatigueDamage != 2 {
		t.Fatalf("Expected fatigue damage to be 2, got %d", emptyPlayer.FatigueDamage)
	}

	// Verify hero took 2 more fatigue damage
	if emptyPlayer.Hero.Health != 27 {
		t.Fatalf("Expected hero to take 2 fatigue damage, health is %d instead of 27", emptyPlayer.Hero.Health)
	}
}

// TestDrawSpecificCard tests drawing a specific card from the deck
func TestDrawSpecificCard(t *testing.T) {
	g := CreateTestGame()
	// No need for StartGame, the test game is already set up
	player := g.Players[0]

	// Add some specific cards to test
	player.Deck = append(player.Deck, CreateTestMinionEntity(g, player, WithName("Special Card 1")))
	player.Deck = append(player.Deck, CreateTestMinionEntity(g, player, WithName("Special Card 2")))

	initialHandSize := len(player.Hand)
	initialDeckSize := len(player.Deck)

	// Test 1: Draw a specific card
	cardToDraw := "Special Card 1"
	drawnEntity := g.DrawSpecificCard(player, cardToDraw)

	// Verify hand increased by 1
	if len(player.Hand) != initialHandSize+1 {
		t.Fatalf("Expected hand size to increase by 1, got %d (was %d)",
			len(player.Hand), initialHandSize)
	}

	// Verify deck decreased by 1
	if len(player.Deck) != initialDeckSize-1 {
		t.Fatalf("Expected deck size to decrease by 1, got %d (was %d)",
			len(player.Deck), initialDeckSize)
	}

	// Verify the specific card was drawn
	if player.Hand[len(player.Hand)-1].Card.Name != cardToDraw {
		t.Fatalf("Expected drawn card name to be %s, got %s",
			cardToDraw, player.Hand[len(player.Hand)-1].Card.Name)
	}

	// Verify correct card reference is returned
	if drawnEntity == nil {
		t.Fatalf("Expected drawn entity to be returned, got nil")
	}
	if drawnEntity.Card.Name != cardToDraw {
		t.Fatalf("Expected returned card name to be %s, got %s",
			cardToDraw, drawnEntity.Card.Name)
	}

	// Test 2: Drawing a card that doesn't exist
	nonExistentCard := g.DrawSpecificCard(player, "Non-existent Card")

	// Verify hand size didn't change
	if len(player.Hand) != initialHandSize+1 {
		t.Fatalf("Expected hand size to remain %d, got %d",
			initialHandSize+1, len(player.Hand))
	}

	// Verify nil is returned for non-existent card
	if nonExistentCard != nil {
		t.Fatalf("Expected nil to be returned when drawing non-existent card, got %v", nonExistentCard)
	}
}

// TestAddEntityToHand tests the Game.AddEntityToHand function
func TestAddEntityToHand(t *testing.T) {
	g := CreateTestGame()
	player := g.Players[0]

	// Test 1: Adding to hand normally
	entity := CreateTestMinionEntity(g, player)
	entity.CurrentZone = ZONE_DECK

	result, ok := g.AddEntityToHand(player, entity, -1)

	if !ok {
		t.Errorf("Expected AddEntityToHand to return true")
	}

	if result != entity {
		t.Errorf("Expected returned entity to be the same as input entity")
	}

	if entity.CurrentZone != ZONE_HAND {
		t.Errorf("Entity should have zone HAND, got %s", entity.CurrentZone)
	}

	// Test 2: Adding to hand with full hand
	// Fill the hand
	for i := 0; i < player.HandSize; i++ {
		player.Hand = append(player.Hand, CreateTestMinionEntity(g, player))
	}

	// Create an entity that should not be able to fit in hand
	overflowEntity := CreateTestMinionEntity(g, player)
	overflowEntity.CurrentZone = ZONE_DECK

	result, ok = g.AddEntityToHand(player, overflowEntity, -1)

	if ok {
		t.Errorf("Expected AddEntityToHand to return false with full hand")
	}

	if result != nil {
		t.Errorf("Expected nil result with full hand")
	}

	if overflowEntity.CurrentZone != ZONE_REMOVEDFROMGAME {
		t.Errorf("Entity should have zone REMOVEDFROMGAME, got %s", overflowEntity.CurrentZone)
	}
}
