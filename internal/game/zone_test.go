package game

import (
	"testing"
)

// TestEntityZoneTracking tests that entity zones are correctly tracked across various operations
func TestEntityZoneTracking(t *testing.T) {
	g := CreateTestGame()

	player := g.Players[0]

	// Test 1: Verify new entities have zone set to NONE
	newEntity := CreateTestMinionEntity(g, player)
	if newEntity.CurrentZone != ZONE_NONE {
		t.Errorf("New entity should have zone NONE, got %s", newEntity.CurrentZone)
	}

	// Test 2: Verify entities created in deck have zone set to DECK
	deckEntity := player.Deck[0]
	if deckEntity.CurrentZone != ZONE_DECK {
		t.Errorf("Deck entity should have zone DECK, got %s", deckEntity.CurrentZone)
	}

	// Test 3: Verify zone changes when drawing cards
	drawnEntity := g.DrawCard(player)
	if drawnEntity.CurrentZone != ZONE_HAND {
		t.Errorf("Drawn entity should have zone HAND, got %s", drawnEntity.CurrentZone)
	}

	// Verify entity is no longer in deck
	for _, entity := range player.Deck {
		if entity == drawnEntity {
			t.Errorf("Drawn entity should no longer be in deck")
		}
	}

	// Test 4: Verify zone changes when playing minions
	// We need to ensure we have enough mana to play the card
	player.Mana = 10

	// Add a test card to hand with known cost
	testCard := CreateTestMinionEntity(g, player, WithName("Test Zone Minion"), WithCost(2))
	g.AddEntityToHand(player, testCard, -1)

	// Verify it has HAND zone
	if testCard.CurrentZone != ZONE_HAND {
		t.Errorf("Test card should have zone HAND, got %s", testCard.CurrentZone)
	}

	// Play the card (last card in hand)
	handIndex := len(player.Hand) - 1
	err := g.PlayCard(player, handIndex, nil, -1, 0)
	if err != nil {
		t.Fatalf("Failed to play test card: %v", err)
	}

	// Verify it has PLAY zone
	if testCard.CurrentZone != ZONE_PLAY {
		t.Errorf("Played minion should have zone PLAY, got %s", testCard.CurrentZone)
	}

	// Test 5: Verify zone changes when playing spells
	// Add a test spell to hand
	testSpell := CreateTestSpellEntity(g, player, WithName("Test Zone Spell"), WithCost(1))
	g.AddEntityToHand(player, testSpell, -1)

	// Play the spell
	handIndex = len(player.Hand) - 1
	err = g.PlayCard(player, handIndex, nil, -1, 0)
	if err != nil {
		t.Fatalf("Failed to play test spell: %v", err)
	}

	// Verify it went to GRAVEYARD
	if testSpell.CurrentZone != ZONE_GRAVEYARD {
		t.Errorf("Played spell should have zone GRAVEYARD, got %s", testSpell.CurrentZone)
	}

	// Test 6: Verify zone changes when equipping weapons
	// Add a test weapon to hand
	testWeapon := CreateTestWeaponEntity(g, player, WithName("Test Zone Weapon"), WithCost(1))
	g.AddEntityToHand(player, testWeapon, -1)

	// Play the weapon
	handIndex = len(player.Hand) - 1
	err = g.PlayCard(player, handIndex, nil, -1, 0)
	if err != nil {
		t.Fatalf("Failed to play test weapon: %v", err)
	}

	// Verify it has PLAY zone
	if testWeapon.CurrentZone != ZONE_PLAY {
		t.Errorf("Equipped weapon should have zone PLAY, got %s", testWeapon.CurrentZone)
	}

	// Add another weapon to hand
	replacementWeapon := CreateTestWeaponEntity(g, player, WithName("Replacement Weapon"), WithCost(1))
	g.AddEntityToHand(player, replacementWeapon, -1)

	// Play the second weapon
	handIndex = len(player.Hand) - 1
	err = g.PlayCard(player, handIndex, nil, -1, 0)
	if err != nil {
		t.Fatalf("Failed to play replacement weapon: %v", err)
	}

	// Verify first weapon went to GRAVEYARD
	if testWeapon.CurrentZone != ZONE_GRAVEYARD {
		t.Errorf("Replaced weapon should have zone GRAVEYARD, got %s", testWeapon.CurrentZone)
	}

	// Verify new weapon is in PLAY
	if replacementWeapon.CurrentZone != ZONE_PLAY {
		t.Errorf("New weapon should have zone PLAY, got %s", replacementWeapon.CurrentZone)
	}

	// Test 7: Verify zone changes when entities die
	// Get a reference to a minion on the field
	minion := testCard

	// Kill the minion by setting health to 0
	minion.Health = 0

	// Process deaths
	g.ProcessGraveyard()

	// Verify minion went to GRAVEYARD
	if minion.CurrentZone != ZONE_GRAVEYARD {
		t.Errorf("Dead minion should have zone GRAVEYARD, got %s", minion.CurrentZone)
	}

	// Test 8: Verify hero in PLAY zone
	if player.Hero.CurrentZone != ZONE_PLAY {
		t.Errorf("Hero should have zone PLAY, got %s", player.Hero.CurrentZone)
	}

	// Test 9: Test handling of full hand (card gets discarded)
	// Fill player's hand to capacity
	player.HandSize = 3 // Set small hand size for testing
	player.Hand = nil   // Clear hand

	for i := 0; i < player.HandSize; i++ {
		filler := CreateTestMinionEntity(g, player, WithName("Filler Card"))
		g.AddEntityToHand(player, filler, -1)
	}

	// Try to draw with full hand
	if len(player.Deck) > 0 {
		// Save the entity that should be drawn
		toBurned := player.Deck[len(player.Deck)-1]

		// Attempt to draw with full hand
		drawnEntity = g.DrawCard(player)

		// Verify null is returned
		if drawnEntity != nil {
			t.Errorf("Drawing with full hand should return nil")
		}

		// Verify the card that would be drawn has zone REMOVEDFROMGAME
		if toBurned.CurrentZone != ZONE_REMOVEDFROMGAME {
			t.Errorf("Burned card should have zone REMOVEDFROMGAME, got %s", toBurned.CurrentZone)
		}
	}
}

// TestZoneTrackingDuringHeroReplacement tests that entity zones are properly updated when replacing heroes
func TestZoneTrackingDuringHeroReplacement(t *testing.T) {
	g := CreateTestGame()

	player := g.Players[0]

	// Store reference to original hero
	originalHero := player.Hero

	// Verify original hero is in PLAY zone
	if originalHero.CurrentZone != ZONE_PLAY {
		t.Errorf("Original hero should have zone PLAY, got %s", originalHero.CurrentZone)
	}

	// Create a hero card to replace the current one
	newHero := CreateTestHeroEntity(g, player,
		WithName("Replacement Hero"),
		WithHealth(15))

	// Add to hand
	g.AddEntityToHand(player, newHero, -1)

	// Verify it's in HAND zone
	if newHero.CurrentZone != ZONE_HAND {
		t.Errorf("New hero card should have zone HAND, got %s", newHero.CurrentZone)
	}

	// Play the hero card
	handIndex := len(player.Hand) - 1
	err := g.PlayCard(player, handIndex, nil, -1, 0)
	if err != nil {
		t.Fatalf("Failed to play hero card: %v", err)
	}

	// Verify new hero is in PLAY zone and is now the player's hero
	if newHero.CurrentZone != ZONE_PLAY {
		t.Errorf("New hero should have zone PLAY, got %s", newHero.CurrentZone)
	}
	if player.Hero != newHero {
		t.Errorf("Player's hero should be the new hero")
	}

	// Verify original hero is in GRAVEYARD
	if originalHero.CurrentZone != ZONE_GRAVEYARD {
		t.Errorf("Original hero should have zone GRAVEYARD, got %s", originalHero.CurrentZone)
	}
}

// TestZoneTrackingDuringCombat tests that entity zones are properly updated during combat
func TestZoneTrackingDuringCombat(t *testing.T) {
	g := CreateTestGame()

	player1 := g.Players[0]
	player2 := g.Players[1]

	// Create attacking minion for player 1
	attacker := CreateTestMinionEntity(g, player1,
		WithName("Attacker"),
		WithAttack(5),
		WithHealth(5))
	g.AddEntityToField(player1, attacker, -1)
	attacker.Exhausted = false // Allow it to attack

	// Create defending minion for player 2 that will die from the attack
	defender := CreateTestMinionEntity(g, player2,
		WithName("Defender"),
		WithAttack(2),
		WithHealth(3))
	g.AddEntityToField(player2, defender, -1)

	// Verify both are in PLAY zone
	if attacker.CurrentZone != ZONE_PLAY {
		t.Errorf("Attacker should have zone PLAY, got %s", attacker.CurrentZone)
	}
	if defender.CurrentZone != ZONE_PLAY {
		t.Errorf("Defender should have zone PLAY, got %s", defender.CurrentZone)
	}

	// Execute attack
	err := g.ProcessAttack(attacker, defender)
	if err != nil {
		t.Fatalf("Failed to process attack: %v", err)
	}

	// Verify defender died and went to GRAVEYARD
	if defender.CurrentZone != ZONE_GRAVEYARD {
		t.Errorf("Dead defender should have zone GRAVEYARD, got %s", defender.CurrentZone)
	}

	// Verify attacker survived and remained in PLAY
	if attacker.CurrentZone != ZONE_PLAY {
		t.Errorf("Surviving attacker should have zone PLAY, got %s", attacker.CurrentZone)
	}

	// Verify defender is in player2's graveyard
	found := false
	for _, entity := range player2.Graveyard {
		if entity == defender {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Defender should be in player2's graveyard")
	}

	// Test mutual destruction scenario
	// Create two minions that will kill each other
	minion1 := CreateTestMinionEntity(g, player1,
		WithName("Minion1"),
		WithAttack(4),
		WithHealth(3))
	g.AddEntityToField(player1, minion1, -1)
	minion1.Exhausted = false // Allow it to attack

	minion2 := CreateTestMinionEntity(g, player2,
		WithName("Minion2"),
		WithAttack(3),
		WithHealth(4))
	g.AddEntityToField(player2, minion2, -1)

	// Execute attack
	err = g.ProcessAttack(minion1, minion2)
	if err != nil {
		t.Fatalf("Failed to process attack in mutual destruction scenario: %v", err)
	}

	// Verify both died and went to GRAVEYARD
	if minion1.CurrentZone != ZONE_GRAVEYARD {
		t.Errorf("Minion1 should have zone GRAVEYARD, got %s", minion1.CurrentZone)
	}
	if minion2.CurrentZone != ZONE_GRAVEYARD {
		t.Errorf("Minion2 should have zone GRAVEYARD, got %s", minion2.CurrentZone)
	}

	// Test weapon destruction
	// Create a weapon for player1
	weapon := CreateTestWeaponEntity(g, player1,
		WithName("Testing Weapon"),
		WithAttack(2),
		WithHealth(1)) // 1 durability

	player1.Weapon = weapon
	player1.Hero.Attack = 1
	weapon.CurrentZone = ZONE_PLAY

	// Have the hero attack to use the weapon
	err = g.ProcessAttack(player1.Hero, player2.Hero)
	if err != nil {
		t.Fatalf("Failed to process hero attack: %v", err)
	}

	// Verify weapon was destroyed and went to GRAVEYARD
	if weapon.CurrentZone != ZONE_GRAVEYARD {
		t.Errorf("Destroyed weapon should have zone GRAVEYARD, got %s", weapon.CurrentZone)
	}

	// Verify weapon is in player1's graveyard
	found = false
	for _, entity := range player1.Graveyard {
		if entity == weapon {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Weapon should be in player1's graveyard")
	}

	// Test poisonous effect
	// Create a poisonous minion
	poisonous := CreateTestMinionEntity(g, player1,
		WithName("Poisonous Minion"),
		WithAttack(1),
		WithHealth(1))
	poisonous.Tags = append(poisonous.Tags, NewTag(TAG_POISONOUS, true))
	g.AddEntityToField(player1, poisonous, -1)
	poisonous.Exhausted = false

	// Create a big target
	bigTarget := CreateTestMinionEntity(g, player2,
		WithName("Big Target"),
		WithAttack(1),
		WithHealth(10))
	g.AddEntityToField(player2, bigTarget, -1)

	// Execute attack
	err = g.ProcessAttack(poisonous, bigTarget)
	if err != nil {
		t.Fatalf("Failed to process poisonous attack: %v", err)
	}

	// Verify poisonous killed the big minion
	if bigTarget.CurrentZone != ZONE_GRAVEYARD {
		t.Errorf("Poisoned target should have zone GRAVEYARD, got %s", bigTarget.CurrentZone)
	}
}
