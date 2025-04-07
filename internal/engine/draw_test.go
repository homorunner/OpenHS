package engine

import (
	"testing"

	"github.com/openhs/internal/game"
)

// TestDrawCard tests the DrawCard helper function
func TestDrawCard(t *testing.T) {
	g := createTestGame()
	e := NewEngine(g)
	e.StartGame()

	// Test 1: Normal card draw
	player := g.Players[0]
	initialHandSize := len(player.Hand)
	initialDeckSize := len(player.Deck)

	drawnEntity := e.DrawCard(player)

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
	emptyPlayer := game.NewPlayer()
	emptyPlayer.Hero = createTestHeroEntity(emptyPlayer, withName("Test Hero"), withHealth(30))

	// Try to draw from empty deck
	drawn := e.DrawCard(emptyPlayer)

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
	e.DrawCard(emptyPlayer)

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
	g := createTestGame()
	e := NewEngine(g)
	e.StartGame()
	player := g.Players[0]

	// Add some specific cards to test
	player.Deck = append(player.Deck, createTestMinionEntity(player, withName("Special Card 1")))
	player.Deck = append(player.Deck, createTestMinionEntity(player, withName("Special Card 2")))

	initialHandSize := len(player.Hand)
	initialDeckSize := len(player.Deck)

	// Test 1: Draw a specific card
	cardToDraw := "Special Card 1"
	drawnEntity := e.DrawSpecificCard(player, cardToDraw)

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
	nonExistentCard := e.DrawSpecificCard(player, "Non-existent Card")

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

// TestAddCardToHand tests the full hand discard logic
func TestAddCardToHand(t *testing.T) {
	g := createTestGame()
	e := NewEngine(g)
	e.StartGame()

	// Create a player with a nearly full hand
	player := game.NewPlayer()
	player.HandSize = 3 // Small hand size for testing

	// Add cards until one away from full
	addToHand(player, createTestMinionEntity(player, withName("Hand Card 1")))
	addToHand(player, createTestMinionEntity(player, withName("Hand Card 2")))

	// Test 1: Add card to hand with space available
	newEntity := createTestMinionEntity(player, withName("New Card"))
	success := e.AddCardToHand(player, newEntity)

	// Verify card was added successfully
	if !success {
		t.Fatalf("Expected AddCardToHand to return true when space is available")
	}

	if len(player.Hand) != 3 {
		t.Fatalf("Expected hand size to be 3, got %d", len(player.Hand))
	}

	if player.Hand[2].Card.Name != "New Card" {
		t.Fatalf("Expected last card in hand to be 'New Card', got %s", player.Hand[2].Card.Name)
	}

	// Test 2: Add card to full hand
	overflowEntity := createTestMinionEntity(player, withName("Overflow Card"))
	success = e.AddCardToHand(player, overflowEntity)

	// Verify card was not added
	if success {
		t.Fatalf("Expected AddCardToHand to return false when hand is full")
	}

	if len(player.Hand) != 3 {
		t.Fatalf("Expected hand size to remain 3, got %d", len(player.Hand))
	}
}
