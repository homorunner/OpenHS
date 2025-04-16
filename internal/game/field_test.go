package game

import (
	"testing"
)

// TestAddEntityToField tests the AddEntityToField function
func TestAddEntityToField(t *testing.T) {
	g := CreateTestGame()
	player := g.Players[0]

	// Test 1: Add entity to end of field
	testMinion := CreateTestMinionEntity(g, player, WithName("Test Field Minion"))
	success := g.AddEntityToField(player, testMinion, -1)

	if !success {
		t.Error("AddEntityToField returned false when adding to a non-full field")
	}

	if testMinion.CurrentZone != ZONE_PLAY {
		t.Errorf("Entity should be in ZONE_PLAY, got %s", testMinion.CurrentZone)
	}

	if len(player.Field) != 1 {
		t.Errorf("Field should have 1 entity, got %d", len(player.Field))
	}

	if player.Field[0] != testMinion {
		t.Error("Entity not found in expected field position")
	}

	// Test 2: Add entity to specific position in field
	testMinion2 := CreateTestMinionEntity(g, player, WithName("Test Field Minion 2"))
	success = g.AddEntityToField(player, testMinion2, 0)

	if !success {
		t.Error("AddEntityToField returned false when adding to a non-full field")
	}

	if len(player.Field) != 2 {
		t.Errorf("Field should have 2 entities, got %d", len(player.Field))
	}

	if player.Field[0] != testMinion2 {
		t.Error("Entity not found in specified field position")
	}

	// Test 3: Test exhausted status based on charge/rush
	normalMinion := CreateTestMinionEntity(g, player, WithName("Normal Minion"))
	g.AddEntityToField(player, normalMinion, -1)
	if !normalMinion.Exhausted {
		t.Error("Normal minion should be exhausted when played")
	}

	chargeMinion := CreateTestMinionEntity(g, player,
		WithName("Charge Minion"),
		WithTag(TAG_CHARGE, true))
	g.AddEntityToField(player, chargeMinion, -1)
	if chargeMinion.Exhausted {
		t.Error("Charge minion should not be exhausted when played")
	}

	rushMinion := CreateTestMinionEntity(g, player,
		WithName("Rush Minion"),
		WithTag(TAG_RUSH, true))
	g.AddEntityToField(player, rushMinion, -1)
	if rushMinion.Exhausted {
		t.Error("Rush minion should not be exhausted when played")
	}

	// Test 4: Test full field
	player.FieldSize = 7 // Set field size limit
	for i := player.FieldSize - len(player.Field); i > 0; i-- {
		g.AddEntityToField(player, CreateTestMinionEntity(g, player), -1)
	}

	// Try to add to a full field
	fullFieldMinion := CreateTestMinionEntity(g, player, WithName("Full Field Minion"))
	success = g.AddEntityToField(player, fullFieldMinion, -1)
	if success {
		t.Error("AddEntityToField should return false when field is full")
	}

	if fullFieldMinion.CurrentZone != ZONE_NONE {
		t.Errorf("Entity should be in ZONE_NONE when field is full, got %s", fullFieldMinion.CurrentZone)
	}

	if len(player.Field) > player.FieldSize {
		t.Errorf("Field should not exceed maximum size: got %d, max %d",
			len(player.Field), player.FieldSize)
	}
}

// TestAddEntityToHand tests the AddEntityToHand function
func TestAddEntityToHand(t *testing.T) {
	g := CreateTestGame()
	player := g.Players[0]

	// Test 1: Add entity to end of hand
	testCard := CreateTestMinionEntity(g, player, WithName("Test Hand Card"))
	entity, success := g.AddEntityToHand(player, testCard, -1)

	if entity == nil || !success {
		t.Error("AddEntityToHand returned nil or false when adding to a non-full hand")
	}

	if testCard.CurrentZone != ZONE_HAND {
		t.Errorf("Entity should be in ZONE_HAND, got %s", testCard.CurrentZone)
	}

	if len(player.Hand) != 1 {
		t.Errorf("Hand should have 1 entity, got %d", len(player.Hand))
	}

	if player.Hand[0] != testCard {
		t.Error("Entity not found in expected hand position")
	}

	// Test 2: Add entity to specific position in hand
	testCard2 := CreateTestMinionEntity(g, player, WithName("Test Hand Card 2"))
	entity, success = g.AddEntityToHand(player, testCard2, 0)

	if entity == nil || !success {
		t.Error("AddEntityToHand returned nil or false when adding to a non-full hand")
	}

	if len(player.Hand) != 2 {
		t.Errorf("Hand should have 2 entities, got %d", len(player.Hand))
	}

	if player.Hand[0] != testCard2 {
		t.Error("Entity not found in specified hand position")
	}

	// Test 3: Test full hand
	player.HandSize = 5 // Set hand size limit
	for i := player.HandSize - len(player.Hand); i > 0; i-- {
		g.AddEntityToHand(player, CreateTestMinionEntity(g, player), -1)
	}

	// Try to add to a full hand
	fullHandCard := CreateTestMinionEntity(g, player, WithName("Full Hand Card"))
	entity, success = g.AddEntityToHand(player, fullHandCard, -1)
	if entity != nil || success {
		t.Error("AddEntityToHand should return nil and false when hand is full")
	}

	if fullHandCard.CurrentZone != ZONE_REMOVEDFROMGAME {
		t.Errorf("Entity should be in ZONE_REMOVEDFROMGAME when hand is full, got %s", fullHandCard.CurrentZone)
	}

	if len(player.Hand) > player.HandSize {
		t.Errorf("Hand should not exceed maximum size: got %d, max %d",
			len(player.Hand), player.HandSize)
	}
}

// TestMoveFromDeckToHand tests the MoveFromDeckToHand function
func TestMoveFromDeckToHand(t *testing.T) {
	g := CreateTestGame()
	player := g.Players[0]

	// Make sure the deck has cards
	if len(player.Deck) == 0 {
		t.Fatal("Test player deck should not be empty")
	}

	// Save reference to deck card
	deckCard := player.Deck[0]

	// Move card from deck to hand
	entity, success := g.MoveFromDeckToHand(player, 0, -1)

	// Check result
	if entity == nil || !success {
		t.Error("MoveFromDeckToHand returned nil or false when moving to non-full hand")
	}

	if deckCard.CurrentZone != ZONE_HAND {
		t.Errorf("Entity should be in ZONE_HAND, got %s", deckCard.CurrentZone)
	}

	if len(player.Deck) != 9 { // 10 - 1 cards
		t.Errorf("Deck should have 9 entities, got %d", len(player.Deck))
	}

	if len(player.Hand) != 1 {
		t.Errorf("Hand should have 1 entity, got %d", len(player.Hand))
	}

	if player.Hand[0] != deckCard {
		t.Error("Entity not found in hand")
	}

	// Test with full hand
	player.HandSize = 3 // Set small hand size
	for i := player.HandSize - len(player.Hand); i > 0; i-- {
		g.AddEntityToHand(player, CreateTestMinionEntity(g, player), -1)
	}

	// Try to move a card to a full hand
	deckCard = player.Deck[0]
	entity, success = g.MoveFromDeckToHand(player, 0, -1)

	if entity != nil || success {
		t.Error("MoveFromDeckToHand should return nil and false when hand is full")
	}

	if deckCard.CurrentZone != ZONE_REMOVEDFROMGAME {
		t.Errorf("Entity should be in ZONE_REMOVEDFROMGAME when hand is full, got %s", deckCard.CurrentZone)
	}
}

// TestMoveFromHandToField tests the MoveFromHandToField function
func TestMoveFromHandToField(t *testing.T) {
	g := CreateTestGame()
	player := g.Players[0]

	// Add a card to hand
	testCard := CreateTestMinionEntity(g, player, WithName("Test Move Card"))
	g.AddEntityToHand(player, testCard, -1)

	// Move card from hand to field
	success := g.MoveFromHandToField(player, 0, -1)

	// Check result
	if !success {
		t.Error("MoveFromHandToField returned false when moving to non-full field")
	}

	if testCard.CurrentZone != ZONE_PLAY {
		t.Errorf("Entity should be in ZONE_PLAY, got %s", testCard.CurrentZone)
	}

	if len(player.Hand) != 0 {
		t.Errorf("Hand should be empty, got %d cards", len(player.Hand))
	}

	if len(player.Field) != 1 {
		t.Errorf("Field should have 1 entity, got %d", len(player.Field))
	}

	if player.Field[0] != testCard {
		t.Error("Entity not found in field")
	}

	// Test with full field
	player.FieldSize = 3 // Set small field size
	for i := player.FieldSize - len(player.Field); i > 0; i-- {
		g.AddEntityToField(player, CreateTestMinionEntity(g, player), -1)
	}

	// Add another card to hand
	testCard2 := CreateTestMinionEntity(g, player, WithName("Test Move Card 2"))
	g.AddEntityToHand(player, testCard2, -1)

	// Try to move a card to a full field
	success = g.MoveFromHandToField(player, 0, -1)

	if success {
		t.Error("MoveFromHandToField should return false when field is full")
	}

	if testCard2.CurrentZone != ZONE_NONE {
		t.Errorf("Entity should be in ZONE_NONE when field is full, got %s", testCard2.CurrentZone)
	}
}

// TestMoveFromDeckToField tests the MoveFromDeckToField function
func TestMoveFromDeckToField(t *testing.T) {
	g := CreateTestGame()
	player := g.Players[0]

	// Make sure the deck has cards
	if len(player.Deck) == 0 {
		t.Fatal("Test player deck should not be empty")
	}

	// Save reference to deck card
	deckCard := player.Deck[0]

	// Move card from deck to field
	success := g.MoveFromDeckToField(player, 0, -1)

	// Check result
	if !success {
		t.Error("MoveFromDeckToField returned false when moving to non-full field")
	}

	if deckCard.CurrentZone != ZONE_PLAY {
		t.Errorf("Entity should be in ZONE_PLAY, got %s", deckCard.CurrentZone)
	}

	if len(player.Deck) != 9 { // 10 - 1 cards
		t.Errorf("Deck should have 9 entities, got %d", len(player.Deck))
	}

	if len(player.Field) != 1 {
		t.Errorf("Field should have 1 entity, got %d", len(player.Field))
	}

	if player.Field[0] != deckCard {
		t.Error("Entity not found in field")
	}

	// Test with full field
	player.FieldSize = 3 // Set small field size
	for i := player.FieldSize - len(player.Field); i > 0; i-- {
		g.AddEntityToField(player, CreateTestMinionEntity(g, player), -1)
	}

	// Try to move a card to a full field
	deckCard = player.Deck[0]
	success = g.MoveFromDeckToField(player, 0, -1)

	if success {
		t.Error("MoveFromDeckToField should return false when field is full")
	}

	if deckCard.CurrentZone != ZONE_NONE {
		t.Errorf("Entity should be in ZONE_NONE when field is full, got %s", deckCard.CurrentZone)
	}
}

// TestRemoveEntityFromBoard tests the removeEntityFromBoard function
func TestRemoveEntityFromBoard(t *testing.T) {
	g := CreateTestGame()
	player := g.Players[0]

	// Add entities to field
	minion1 := CreateTestMinionEntity(g, player, WithName("Minion 1"))
	minion2 := CreateTestMinionEntity(g, player, WithName("Minion 2"))
	minion3 := CreateTestMinionEntity(g, player, WithName("Minion 3"))

	g.AddEntityToField(player, minion1, -1)
	g.AddEntityToField(player, minion2, -1)
	g.AddEntityToField(player, minion3, -1)

	// Initial field state: [minion1, minion2, minion3]
	if len(player.Field) != 3 {
		t.Fatalf("Field should have 3 entities, got %d", len(player.Field))
	}

	// Remove middle entity
	g.removeEntityFromBoard(player, minion2)

	// Check field size
	if len(player.Field) != 2 {
		t.Errorf("Field should have 2 entities after removal, got %d", len(player.Field))
	}

	// Check remaining entities - the last entity should now be in position 1
	if player.Field[0] != minion1 {
		t.Error("First minion should still be in position 0")
	}

	if player.Field[1] != minion3 {
		t.Error("Last minion should now be in position 1")
	}

	// Remove first entity
	g.removeEntityFromBoard(player, minion1)

	// Check field size
	if len(player.Field) != 1 {
		t.Errorf("Field should have 1 entity after second removal, got %d", len(player.Field))
	}

	// Check remaining entity
	if player.Field[0] != minion3 {
		t.Error("Last minion should be the only one left")
	}

	// Remove last entity
	g.removeEntityFromBoard(player, minion3)

	// Check field is empty
	if len(player.Field) != 0 {
		t.Errorf("Field should be empty after all removals, got %d", len(player.Field))
	}
}
