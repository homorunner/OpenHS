package game

import (
	"testing"
)

// TestPlayCard tests the PlayCard function and card playing mechanics
func TestPlayCard(t *testing.T) {
	g := CreateTestGame()
	player := g.Players[0]

	// Add different types of cards to the player's hand for testing
	player.Hand = nil
	g.AddEntityToHand(player, CreateTestMinionEntity(g, player, WithName("Foo")), -1)
	g.AddEntityToHand(player, CreateTestSpellEntity(g, player), -1)
	g.AddEntityToHand(player, CreateTestWeaponEntity(g, player, WithName("Bar")), -1)
	g.AddEntityToHand(player, CreateTestMinionEntity(g, player, WithCost(10)), -1)

	// Setup player resources
	player.Mana = 5
	player.MaxMana = 5

	// Test 1: Play a minion card
	initialHandSize := len(player.Hand)
	initialFieldSize := len(player.Field)
	handIndex := 0 // Test Minion

	err := g.PlayCard(player, handIndex, nil, -1, 0)
	if err != nil {
		t.Fatalf("Failed to play minion card: %v", err)
	}

	// Verify hand decreased by 1
	if len(player.Hand) != initialHandSize-1 {
		t.Fatalf("Expected hand size to decrease by 1, got %d (was %d)",
			len(player.Hand), initialHandSize)
	}

	// Verify field increased by 1
	if len(player.Field) != initialFieldSize+1 {
		t.Fatalf("Expected field size to increase by 1, got %d (was %d)",
			len(player.Field), initialFieldSize)
	}

	// Verify mana was spent
	if player.Mana != 3 {
		t.Fatalf("Expected mana to decrease to 3, got %d", player.Mana)
	}

	// Verify field position - should be at the end when auto-positioned
	if player.Field[len(player.Field)-1].Card.Name != "Foo" {
		t.Fatalf("Expected last field card to be Foo, got %s",
			player.Field[len(player.Field)-1].Card.Name)
	}

	// Test 2: Play a spell card
	initialHandSize = len(player.Hand)
	initialGraveyardSize := len(player.Graveyard)
	handIndex = 0 // Test Spell (after playing the minion, this moved to index 0)

	err = g.PlayCard(player, handIndex, nil, -1, 0)
	if err != nil {
		t.Fatalf("Failed to play spell card: %v", err)
	}

	// Verify hand decreased by 1
	if len(player.Hand) != initialHandSize-1 {
		t.Fatalf("Expected hand size to decrease by 1, got %d (was %d)",
			len(player.Hand), initialHandSize)
	}

	// Verify graveyard increased by 1
	if len(player.Graveyard) != initialGraveyardSize+1 {
		t.Fatalf("Expected graveyard size to increase by 1, got %d (was %d)",
			len(player.Graveyard), initialGraveyardSize)
	}

	// Verify mana was spent
	if player.Mana != 2 {
		t.Fatalf("Expected mana to decrease to 2, got %d", player.Mana)
	}

	// Test 3: Play a weapon card
	initialHandSize = len(player.Hand)
	handIndex = 0 // Test Weapon (after playing the spell, this moved to index 0)

	err = g.PlayCard(player, handIndex, nil, -1, 0)
	if err != nil {
		t.Fatalf("Failed to play weapon card: %v", err)
	}

	// Verify hand decreased by 1
	if len(player.Hand) != initialHandSize-1 {
		t.Fatalf("Expected hand size to decrease by 1, got %d (was %d)",
			len(player.Hand), initialHandSize)
	}

	// Verify player has a weapon equipped
	if player.Weapon == nil {
		t.Fatalf("Expected player to have a weapon equipped")
	}

	// Verify the correct weapon is equipped
	if player.Weapon.Card.Name != "Bar" {
		t.Fatalf("Expected weapon name to be Bar, got %s", player.Weapon.Card.Name)
	}

	// Verify mana was spent
	if player.Mana != 1 {
		t.Fatalf("Expected mana to decrease to 1, got %d", player.Mana)
	}

	// Test 4: Try to play a card with insufficient mana
	// Reset mana to 5
	player.Mana = 5
	handIndex = 0 // Expensive Card (after playing the weapon, this is the only card left)

	err = g.PlayCard(player, handIndex, nil, -1, 0)
	if err == nil {
		t.Fatalf("Expected error when playing card with insufficient mana")
	}

	// Verify mana was not spent
	if player.Mana != 5 {
		t.Fatalf("Expected mana to remain at 5, got %d", player.Mana)
	}

	// Test 5: Try to play a card with invalid hand index
	err = g.PlayCard(player, 10, nil, -1, 0)
	if err == nil {
		t.Fatalf("Expected error when playing card with invalid hand index")
	}
}

// TestPlayCardWithFullField tests playing minions when the field is full
func TestPlayCardWithFullField(t *testing.T) {
	g := CreateTestGame()
	player := g.Players[0]

	// Setup a full field
	player.FieldSize = 7 // Max field size
	for i := 0; i < player.FieldSize; i++ {
		minion := CreateTestMinionEntity(g, player, WithName("Field Minion"))
		g.AddEntityToField(player, minion, -1)
	}

	// Add a minion card to hand
	g.AddEntityToHand(player, CreateTestMinionEntity(g, player), -1)

	player.Mana = 10

	// Try to play minion on full field
	err := g.PlayCard(player, 0, nil, -1, 0)
	if err == nil {
		t.Fatalf("Expected error when playing minion on full field")
	}

	// Verify field size didn't change
	if len(player.Field) != player.FieldSize {
		t.Fatalf("Expected field size to remain %d, got %d", player.FieldSize, len(player.Field))
	}
}

// TestPlayCardWithSpecificPosition tests playing a minion at a specific field position
func TestPlayCardWithSpecificPosition(t *testing.T) {
	g := CreateTestGame()
	player := g.Players[0]

	// Add some minions to the field
	minion1 := CreateTestMinionEntity(g, player, WithName("Field Minion 1"))
	minion2 := CreateTestMinionEntity(g, player, WithName("Field Minion 2"))
	g.AddEntityToField(player, minion1, -1)
	g.AddEntityToField(player, minion2, -1)

	// Add a minion card to hand
	player.Hand = nil
	g.AddEntityToHand(player, CreateTestMinionEntity(g, player, WithName("Test Minion")), -1)

	player.Mana = 10

	// Play minion at position 1 (between the two existing minions)
	err := g.PlayCard(player, 0, nil, 1, 0)
	if err != nil {
		t.Fatalf("Failed to play minion at specific position: %v", err)
	}

	// Verify field size increased
	if len(player.Field) != 3 {
		t.Fatalf("Expected field size to be 3, got %d", len(player.Field))
	}

	// Verify minion is in the correct position
	if player.Field[0].Card.Name != "Field Minion 1" {
		t.Fatalf("Expected minion at position 0 to be 'Field Minion 1', got %s", player.Field[0].Card.Name)
	}
	if player.Field[1].Card.Name != "Test Minion" {
		t.Fatalf("Expected minion at position 1 to be 'Test Minion', got %s", player.Field[1].Card.Name)
	}
	if player.Field[2].Card.Name != "Field Minion 2" {
		t.Fatalf("Expected minion at position 2 to be 'Field Minion 2', got %s", player.Field[2].Card.Name)
	}
}

// TestReplaceWeapon tests replacing an equipped weapon
func TestReplaceWeapon(t *testing.T) {
	g := CreateTestGame()
	player := g.Players[0]

	// Equip initial weapon
	initialWeapon := CreateTestWeaponEntity(g, player,
		WithName("Initial Weapon"),
		WithAttack(2),
		WithHealth(3)) // Weapon health is durability

	// Manually set up the weapon
	initialWeapon.CurrentZone = ZONE_PLAY
	player.Weapon = initialWeapon

	// Add a new weapon to hand
	g.AddEntityToHand(player, CreateTestWeaponEntity(g, player, WithName("New Weapon")), -1)

	player.Mana = 10

	// Play the new weapon
	err := g.PlayCard(player, 0, nil, -1, 0)
	if err != nil {
		t.Fatalf("Failed to play weapon: %v", err)
	}

	// Verify initial weapon was destroyed and sent to graveyard
	if initialWeapon.CurrentZone != ZONE_GRAVEYARD {
		t.Fatalf("Expected initial weapon to be in graveyard, got zone %s", initialWeapon.CurrentZone)
	}

	found := false
	for _, entity := range player.Graveyard {
		if entity == initialWeapon {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("Initial weapon should be in player's graveyard")
	}

	// Verify new weapon is equipped
	if player.Weapon == nil {
		t.Fatalf("Expected player to have a weapon equipped")
	}
	if player.Weapon.Card.Name != "New Weapon" {
		t.Fatalf("Expected weapon name to be 'New Weapon', got %s", player.Weapon.Card.Name)
	}
}
