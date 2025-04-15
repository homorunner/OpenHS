package engine

import (
	"testing"

	"github.com/openhs/internal/game"
)

// TestEntityZoneTracking tests that entity zones are correctly tracked across various operations
func TestEntityZoneTracking(t *testing.T) {
	g := game.CreateTestGame()
	e := NewEngine(g)
	e.StartGame()

	player := g.Players[0]

	// Test 1: Verify new entities have zone set to NONE
	newEntity := game.CreateTestMinionEntity(g, player)
	if newEntity.CurrentZone != game.ZONE_NONE {
		t.Errorf("New entity should have zone NONE, got %s", newEntity.CurrentZone)
	}

	// Test 2: Verify entities created in deck have zone set to DECK
	deckEntity := player.Deck[0]
	if deckEntity.CurrentZone != game.ZONE_DECK {
		t.Errorf("Deck entity should have zone DECK, got %s", deckEntity.CurrentZone)
	}

	// Test 3: Verify zone changes when drawing cards
	drawnEntity := e.DrawCard(player)
	if drawnEntity.CurrentZone != game.ZONE_HAND {
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
	testCard := game.CreateTestMinionEntity(g, player, game.WithName("Test Zone Minion"), game.WithCost(2))
	e.AddEntityToHand(player, testCard, -1)

	// Verify it has HAND zone
	if testCard.CurrentZone != game.ZONE_HAND {
		t.Errorf("Test card should have zone HAND, got %s", testCard.CurrentZone)
	}

	// Play the card (last card in hand)
	handIndex := len(player.Hand) - 1
	err := e.PlayCard(player, handIndex, nil, -1, 0)
	if err != nil {
		t.Fatalf("Failed to play test card: %v", err)
	}

	// Verify it has PLAY zone
	if testCard.CurrentZone != game.ZONE_PLAY {
		t.Errorf("Played minion should have zone PLAY, got %s", testCard.CurrentZone)
	}

	// Test 5: Verify zone changes when playing spells
	// Add a test spell to hand
	testSpell := game.CreateTestSpellEntity(g, player, game.WithName("Test Zone Spell"), game.WithCost(1))
	e.AddEntityToHand(player, testSpell, -1)

	// Play the spell
	handIndex = len(player.Hand) - 1
	err = e.PlayCard(player, handIndex, nil, -1, 0)
	if err != nil {
		t.Fatalf("Failed to play test spell: %v", err)
	}

	// Verify it went to GRAVEYARD
	if testSpell.CurrentZone != game.ZONE_GRAVEYARD {
		t.Errorf("Played spell should have zone GRAVEYARD, got %s", testSpell.CurrentZone)
	}

	// Test 6: Verify zone changes when equipping weapons
	// Add a test weapon to hand
	testWeapon := game.CreateTestWeaponEntity(g, player, game.WithName("Test Zone Weapon"), game.WithCost(1))
	e.AddEntityToHand(player, testWeapon, -1)

	// Play the weapon
	handIndex = len(player.Hand) - 1
	err = e.PlayCard(player, handIndex, nil, -1, 0)
	if err != nil {
		t.Fatalf("Failed to play test weapon: %v", err)
	}

	// Verify it has PLAY zone
	if testWeapon.CurrentZone != game.ZONE_PLAY {
		t.Errorf("Equipped weapon should have zone PLAY, got %s", testWeapon.CurrentZone)
	}

	// Add another weapon to hand
	replacementWeapon := game.CreateTestWeaponEntity(g, player, game.WithName("Replacement Weapon"), game.WithCost(1))
	e.AddEntityToHand(player, replacementWeapon, -1)

	// Play the second weapon
	handIndex = len(player.Hand) - 1
	err = e.PlayCard(player, handIndex, nil, -1, 0)
	if err != nil {
		t.Fatalf("Failed to play replacement weapon: %v", err)
	}

	// Verify first weapon went to GRAVEYARD
	if testWeapon.CurrentZone != game.ZONE_GRAVEYARD {
		t.Errorf("Replaced weapon should have zone GRAVEYARD, got %s", testWeapon.CurrentZone)
	}

	// Verify new weapon is in PLAY
	if replacementWeapon.CurrentZone != game.ZONE_PLAY {
		t.Errorf("New weapon should have zone PLAY, got %s", replacementWeapon.CurrentZone)
	}

	// Test 7: Verify zone changes when entities die
	// Get a reference to a minion on the field
	minion := testCard

	// Kill the minion by setting health to 0
	minion.Health = 0

	// Process deaths
	e.processGraveyard()

	// Verify minion went to GRAVEYARD
	if minion.CurrentZone != game.ZONE_GRAVEYARD {
		t.Errorf("Dead minion should have zone GRAVEYARD, got %s", minion.CurrentZone)
	}

	// Test 8: Verify hero in PLAY zone
	if player.Hero.CurrentZone != game.ZONE_PLAY {
		t.Errorf("Hero should have zone PLAY, got %s", player.Hero.CurrentZone)
	}

	// Test 9: Test handling of full hand (card gets discarded)
	// Fill player's hand to capacity
	player.HandSize = 3 // Set small hand size for testing
	player.Hand = nil   // Clear hand

	for i := 0; i < player.HandSize; i++ {
		filler := game.CreateTestMinionEntity(g, player, game.WithName("Filler Card"))
		e.AddEntityToHand(player, filler, -1)
	}

	// Try to draw with full hand
	if len(player.Deck) > 0 {
		// Save the entity that should be drawn
		toBurned := player.Deck[len(player.Deck)-1]

		// Attempt to draw with full hand
		drawnEntity = e.DrawCard(player)

		// Verify null is returned
		if drawnEntity != nil {
			t.Errorf("Drawing with full hand should return nil")
		}

		// Verify the card that would be drawn has zone REMOVEDFROMGAME
		if toBurned.CurrentZone != game.ZONE_REMOVEDFROMGAME {
			t.Errorf("Burned card should have zone REMOVEDFROMGAME, got %s", toBurned.CurrentZone)
		}
	}
}

// TestZoneTrackingDuringHeroReplacement tests that entity zones are properly updated when replacing heroes
func TestZoneTrackingDuringHeroReplacement(t *testing.T) {
	g := game.CreateTestGame()
	e := NewEngine(g)
	e.StartGame()

	player := g.Players[0]

	// Store reference to original hero
	originalHero := player.Hero

	// Verify original hero is in PLAY zone
	if originalHero.CurrentZone != game.ZONE_PLAY {
		t.Errorf("Original hero should have zone PLAY, got %s", originalHero.CurrentZone)
	}

	// Create a hero card to replace the current one
	newHero := game.CreateTestHeroEntity(g, player,
		game.WithName("Replacement Hero"),
		game.WithHealth(15))

	// Add to hand
	e.AddEntityToHand(player, newHero, -1)

	// Verify it's in HAND zone
	if newHero.CurrentZone != game.ZONE_HAND {
		t.Errorf("New hero card should have zone HAND, got %s", newHero.CurrentZone)
	}

	// Give player enough mana
	player.Mana = 10

	// Play the hero card (last card in hand)
	handIndex := len(player.Hand) - 1
	err := e.PlayCard(player, handIndex, nil, -1, 0)
	if err != nil {
		t.Fatalf("Failed to play hero card: %v", err)
	}

	// Verify new hero is in PLAY zone
	if newHero.CurrentZone != game.ZONE_PLAY {
		t.Errorf("New hero should have zone PLAY, got %s", newHero.CurrentZone)
	}

	// Verify original hero is in GRAVEYARD zone
	if originalHero.CurrentZone != game.ZONE_GRAVEYARD {
		t.Errorf("Original hero should have zone GRAVEYARD, got %s", originalHero.CurrentZone)
	}

	// Verify original hero is in player's graveyard
	found := false
	for _, entity := range player.Graveyard {
		if entity == originalHero {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Original hero should be in player's graveyard")
	}

	// Verify player.Hero is set to the new hero
	if player.Hero != newHero {
		t.Errorf("Player.Hero should be set to the new hero")
	}
}

// TestZoneTrackingDuringCombat tests that entity zones are properly updated during combat
func TestZoneTrackingDuringCombat(t *testing.T) {
	g := game.CreateTestGame()
	e := NewEngine(g)
	e.StartGame()

	player1 := g.Players[0]
	player2 := g.Players[1]

	// Create attacking minion for player 1
	attacker := game.CreateTestMinionEntity(g, player1,
		game.WithName("Attacker"),
		game.WithAttack(5),
		game.WithHealth(5))
	e.AddEntityToField(player1, attacker, -1)
	attacker.Exhausted = false // Allow it to attack

	// Create defending minion for player 2 that will die from the attack
	defender := game.CreateTestMinionEntity(g, player2,
		game.WithName("Defender"),
		game.WithAttack(2),
		game.WithHealth(3))
	e.AddEntityToField(player2, defender, -1)

	// Verify both are in PLAY zone
	if attacker.CurrentZone != game.ZONE_PLAY {
		t.Errorf("Attacker should have zone PLAY, got %s", attacker.CurrentZone)
	}
	if defender.CurrentZone != game.ZONE_PLAY {
		t.Errorf("Defender should have zone PLAY, got %s", defender.CurrentZone)
	}

	// Execute attack
	err := e.ProcessAttack(attacker, defender)
	if err != nil {
		t.Fatalf("Failed to process attack: %v", err)
	}

	// Verify defender died and went to GRAVEYARD
	if defender.CurrentZone != game.ZONE_GRAVEYARD {
		t.Errorf("Dead defender should have zone GRAVEYARD, got %s", defender.CurrentZone)
	}

	// Verify attacker survived and remained in PLAY
	if attacker.CurrentZone != game.ZONE_PLAY {
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
	minion1 := game.CreateTestMinionEntity(g, player1,
		game.WithName("Minion1"),
		game.WithAttack(4),
		game.WithHealth(3))
	e.AddEntityToField(player1, minion1, -1)
	minion1.Exhausted = false // Allow it to attack

	minion2 := game.CreateTestMinionEntity(g, player2,
		game.WithName("Minion2"),
		game.WithAttack(3),
		game.WithHealth(4))
	e.AddEntityToField(player2, minion2, -1)

	// Execute attack
	err = e.ProcessAttack(minion1, minion2)
	if err != nil {
		t.Fatalf("Failed to process attack in mutual destruction scenario: %v", err)
	}

	// Verify both died and went to GRAVEYARD
	if minion1.CurrentZone != game.ZONE_GRAVEYARD {
		t.Errorf("Minion1 should have zone GRAVEYARD, got %s", minion1.CurrentZone)
	}
	if minion2.CurrentZone != game.ZONE_GRAVEYARD {
		t.Errorf("Minion2 should have zone GRAVEYARD, got %s", minion2.CurrentZone)
	}

	// Test weapon destruction
	// Create a weapon for player1
	weapon := game.CreateTestWeaponEntity(g, player1,
		game.WithName("Testing Weapon"),
		game.WithAttack(2),
		game.WithHealth(1)) // 1 durability

	player1.Weapon = weapon
	player1.Hero.Attack = 1
	weapon.CurrentZone = game.ZONE_PLAY

	// Have the hero attack to use the weapon
	err = e.ProcessAttack(player1.Hero, player2.Hero)
	if err != nil {
		t.Fatalf("Failed to process hero attack: %v", err)
	}

	// Verify weapon was destroyed and went to GRAVEYARD
	if weapon.CurrentZone != game.ZONE_GRAVEYARD {
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
	poisonous := game.CreateTestMinionEntity(g, player1,
		game.WithName("Poisonous Minion"),
		game.WithAttack(1),
		game.WithHealth(1))
	poisonous.Tags = append(poisonous.Tags, game.NewTag(game.TAG_POISONOUS, true))
	e.AddEntityToField(player1, poisonous, -1)
	poisonous.Exhausted = false

	// Create a big target
	bigTarget := game.CreateTestMinionEntity(g, player2,
		game.WithName("Big Target"),
		game.WithAttack(1),
		game.WithHealth(10))
	e.AddEntityToField(player2, bigTarget, -1)

	// Execute attack
	err = e.ProcessAttack(poisonous, bigTarget)
	if err != nil {
		t.Fatalf("Failed to process poisonous attack: %v", err)
	}

	// Verify poisonous killed the big minion
	if bigTarget.CurrentZone != game.ZONE_GRAVEYARD {
		t.Errorf("Poisoned target should have zone GRAVEYARD, got %s", bigTarget.CurrentZone)
	}
}
