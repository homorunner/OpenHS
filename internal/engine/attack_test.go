package engine

import (
	"testing"

	"github.com/openhs/internal/game"
)

func TestBasicAttack(t *testing.T) {
	t.Run("Basic attack between minions", func(t *testing.T) {
		// Setup
		g := game.CreateTestGame()
		engine := NewEngine(g)
		engine.StartGame()
		player1 := g.Players[0]
		player2 := g.Players[1]

		// Create attacker and defender entities for different players
		attackerEntity := game.CreateTestMinionEntity(g, player1, game.WithName("Test Attacker"), game.WithAttack(3), game.WithHealth(4), game.WithTag(game.TAG_RUSH, true))
		defenderEntity := game.CreateTestMinionEntity(g, player2, game.WithName("Test Defender"), game.WithAttack(2), game.WithHealth(5))

		// Add minions to respective player's field
		engine.AddEntityToField(player1, attackerEntity, -1)
		engine.AddEntityToField(player2, defenderEntity, -1)

		// Perform attack
		err := engine.Attack(attackerEntity, defenderEntity, false)

		// Assert
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		// Check health values after attack
		if attackerEntity.Health != 2 {
			t.Errorf("Expected attacker health to be 2, got %d", attackerEntity.Health)
		}

		if defenderEntity.Health != 2 {
			t.Errorf("Expected defender health to be 2, got %d", defenderEntity.Health)
		}
	})

	t.Run("Attack with zero attack minion", func(t *testing.T) {
		// Setup
		g := game.CreateTestGame()
		engine := NewEngine(g)
		engine.StartGame()
		player1 := g.Players[0]
		player2 := g.Players[1]

		// Create attacker with zero attack and defender for different players
		attackerEntity := game.CreateTestMinionEntity(g, player1, game.WithName("Zero Attack Minion"), game.WithAttack(0), game.WithHealth(4))
		defenderEntity := game.CreateTestMinionEntity(g, player2, game.WithName("Test Defender"), game.WithAttack(2), game.WithHealth(5))

		// Add minions to respective player's field
		engine.AddEntityToField(player1, attackerEntity, -1)
		engine.AddEntityToField(player2, defenderEntity, -1)

		// Perform attack
		err := engine.Attack(attackerEntity, defenderEntity, false)

		// Assert
		if err == nil {
			t.Error("Expected an error for zero attack minion, got none")
		}
	})

	t.Run("Attack with null entities", func(t *testing.T) {
		// Setup
		g := game.CreateTestGame()
		engine := NewEngine(g)
		engine.StartGame()
		player1 := g.Players[0]
		player2 := g.Players[1]

		// Create valid entities for testing
		validAttacker := game.CreateTestMinionEntity(g, player1, game.WithName("Valid Attacker"), game.WithAttack(1), game.WithHealth(1))
		validDefender := game.CreateTestMinionEntity(g, player2, game.WithName("Valid Defender"), game.WithAttack(1), game.WithHealth(1))

		// Add minions to respective player's field
		engine.AddEntityToField(player1, validAttacker, -1)
		engine.AddEntityToField(player2, validDefender, -1)

		// Perform attack with nil attacker
		err := engine.Attack(nil, validDefender, false)

		// Assert
		if err == nil {
			t.Error("Expected an error for nil attacker, got none")
		}

		// Perform attack with nil defender
		err = engine.Attack(validAttacker, nil, false)

		// Assert
		if err == nil {
			t.Error("Expected an error for nil defender, got none")
		}
	})

	t.Run("Attack with skipValidation", func(t *testing.T) {
		// Setup
		g := game.CreateTestGame()
		engine := NewEngine(g)
		engine.StartGame()
		player1 := g.Players[0]
		player2 := g.Players[1]

		// Create zero attack attacker and defender for different players
		attackerEntity := game.CreateTestMinionEntity(g, player1, game.WithName("Zero Attack Minion"), game.WithAttack(0), game.WithHealth(4))
		defenderEntity := game.CreateTestMinionEntity(g, player2, game.WithName("Test Defender"), game.WithAttack(2), game.WithHealth(5))

		// Add minions to respective player's field
		engine.AddEntityToField(player1, attackerEntity, -1)
		engine.AddEntityToField(player2, defenderEntity, -1)

		// Perform attack with skipValidation=true
		err := engine.Attack(attackerEntity, defenderEntity, true)

		// Assert
		if err != nil {
			t.Errorf("Expected no error with skipValidation, got %v", err)
		}

		// Since attacker has 0 attack, defender shouldn't take damage
		if defenderEntity.Health != 5 {
			t.Errorf("Expected defender health to be 5, got %d", defenderEntity.Health)
		}

		// Attacker should still take damage
		if attackerEntity.Health != 2 {
			t.Errorf("Expected attacker health to be 2, got %d", attackerEntity.Health)
		}
	})

	t.Run("Attack that kills both entities", func(t *testing.T) {
		// Setup
		g := game.CreateTestGame()
		engine := NewEngine(g)
		engine.StartGame()
		player1 := g.Players[0]
		player2 := g.Players[1]

		// Create entities with just enough health to be killed for different players
		attackerEntity := game.CreateTestMinionEntity(g, player1, game.WithName("Lethal Attacker"), game.WithAttack(5), game.WithHealth(2), game.WithTag(game.TAG_RUSH, true))
		defenderEntity := game.CreateTestMinionEntity(g, player2, game.WithName("Fragile Defender"), game.WithAttack(2), game.WithHealth(2))

		// Add minions to respective player's field
		engine.AddEntityToField(player1, attackerEntity, -1)
		engine.AddEntityToField(player2, defenderEntity, -1)

		// Perform attack
		err := engine.Attack(attackerEntity, defenderEntity, false)

		// Assert
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		// Check both entities should be at 0 health or less
		if attackerEntity.Health > 0 {
			t.Errorf("Expected attacker to be dead, got health %d", attackerEntity.Health)
		}

		if defenderEntity.Health > 0 {
			t.Errorf("Expected defender to be dead, got health %d", defenderEntity.Health)
		}
	})

	t.Run("Game phase changes during attack", func(t *testing.T) {
		// Setup
		g := game.CreateTestGame()
		g.Phase = game.MainAction
		engine := NewEngine(g)
		engine.StartGame()
		player1 := g.Players[0]
		player2 := g.Players[1]

		// Create attacker and defender entities for different players
		attackerEntity := game.CreateTestMinionEntity(g, player1, game.WithName("Test Attacker"), game.WithAttack(3), game.WithHealth(4))
		defenderEntity := game.CreateTestMinionEntity(g, player2, game.WithName("Test Defender"), game.WithAttack(2), game.WithHealth(5))

		// Add minions to respective player's field
		engine.AddEntityToField(player1, attackerEntity, -1)
		engine.AddEntityToField(player2, defenderEntity, -1)

		// Perform attack
		_ = engine.Attack(attackerEntity, defenderEntity, false)

		// Assert phase returns to MainAction
		if g.Phase != game.MainAction {
			t.Errorf("Expected game phase to be MainAction after attack, got %v", g.Phase)
		}
	})

	t.Run("Weapon durability decreases on hero attack", func(t *testing.T) {
		g := game.CreateTestGame()
		engine := NewEngine(g)
		engine.StartGame()
		player1 := g.Players[0]
		player2 := g.Players[1]

		// Give player a weapon
		player1.Weapon = game.CreateTestWeaponEntity(g, player1, game.WithName("Test Weapon"), game.WithAttack(3), game.WithHealth(2))

		// Create a defender entity for the opponent
		defenderEntity := game.CreateTestMinionEntity(g, player2, game.WithName("Test Defender"), game.WithAttack(2), game.WithHealth(5))

		// Add defender to opponent's field
		engine.AddEntityToField(player2, defenderEntity, -1)

		// Set hero's attack value to match weapon's attack
		player1.Hero.Attack = player1.Weapon.Attack

		// Perform attack, skip validation
		err := engine.Attack(player1.Hero, defenderEntity, true)

		// Assert
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		// Check weapon durability decreased
		if player1.Weapon.Health != 1 {
			t.Errorf("Expected weapon durability to be 1, got %d", player1.Weapon.Health)
		}

		// Perform another attack
		err = engine.Attack(player1.Hero, defenderEntity, true)

		// Assert
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		// Check weapon is destroyed after reaching 0 durability
		if player1.Weapon != nil {
			t.Errorf("Expected weapon to be destroyed, but it still exists with durability %d", player1.Weapon.Health)
		}
	})

	t.Run("Player can attack own minions only when skipValidation is true", func(t *testing.T) {
		// Setup
		g := game.CreateTestGame()
		engine := NewEngine(g)
		engine.StartGame()
		player := g.Players[0]

		// Create attacker and defender entities for the same player
		attackerEntity := game.CreateTestMinionEntity(g, player, game.WithName("Attacker"), game.WithAttack(3), game.WithHealth(10))
		defenderEntity := game.CreateTestMinionEntity(g, player, game.WithName("Same Player Defender"), game.WithAttack(1), game.WithHealth(30))

		// Add minions to player's field
		engine.AddEntityToField(player, attackerEntity, -1)
		engine.AddEntityToField(player, defenderEntity, -1)

		// Attempt to attack own minion with validation (should fail)
		err := engine.Attack(attackerEntity, defenderEntity, false)
		if err == nil {
			t.Error("Expected an error when attacking own minion with validation enabled, got none")
		}

		// Initial health check
		initialDefenderHealth := defenderEntity.Health
		initialAttackerHealth := attackerEntity.Health

		// Attempt to attack own minion without validation (should succeed)
		err = engine.Attack(attackerEntity, defenderEntity, true)
		if err != nil {
			t.Errorf("Expected no error when attacking own minion with skipValidation=true, got: %v", err)
		}

		// Verify exact damage calculation
		expectedHealth := initialDefenderHealth - attackerEntity.Attack
		if defenderEntity.Health != expectedHealth {
			t.Errorf("Expected defender health to be %d, got %d",
				expectedHealth, defenderEntity.Health)
		}

		// Verify attacker also took damage
		expectedAttackerHealth := initialAttackerHealth - defenderEntity.Attack
		if attackerEntity.Health != expectedAttackerHealth {
			t.Errorf("Expected attacker health to be %d, got %d",
				expectedAttackerHealth, attackerEntity.Health)
		}
	})
}

func TestProcessDestroyAndUpdateAura(t *testing.T) {
	t.Run("Minions with zero or negative health are moved to graveyard", func(t *testing.T) {
		// Setup
		g := game.CreateTestGame()
		engine := NewEngine(g)
		engine.StartGame()
		player := g.Players[0]

		// Create minions with zero and negative health
		minion1 := game.CreateTestMinionEntity(g, player, game.WithName("Dead Minion 1"), game.WithAttack(1), game.WithHealth(0))
		minion2 := game.CreateTestMinionEntity(g, player, game.WithName("Dead Minion 2"), game.WithAttack(1), game.WithHealth(-1))
		minion3 := game.CreateTestMinionEntity(g, player, game.WithName("Alive Minion"), game.WithAttack(1), game.WithHealth(2))

		// Add minions to the field
		engine.AddEntityToField(player, minion1, -1)
		engine.AddEntityToField(player, minion2, -1)
		engine.AddEntityToField(player, minion3, -1)
		initialFieldSize := len(player.Field)

		// Process deaths
		engine.processDestroyAndUpdateAura()

		// Check that the dead minions were removed from the field
		if len(player.Field) != initialFieldSize-2 {
			t.Errorf("Expected field size to be %d, got %d", initialFieldSize-2, len(player.Field))
		}

		// Check that the dead minions were added to the graveyard
		if len(player.Graveyard) != 2 {
			t.Errorf("Expected graveyard size to be 2, got %d", len(player.Graveyard))
		}

		// Check that the living minion is still on the field
		if len(player.Field) != 1 || player.Field[0].Card.Name != "Alive Minion" {
			t.Errorf("Expected only 'Alive Minion' to remain on the field")
		}
	})

	t.Run("Minions marked as destroyed are moved to graveyard", func(t *testing.T) {
		// Setup
		g := game.CreateTestGame()
		engine := NewEngine(g)
		engine.StartGame()
		player := g.Players[0]

		// Create a minion marked for destruction
		minion := game.CreateTestMinionEntity(g, player, game.WithName("Marked Minion"), game.WithAttack(1), game.WithHealth(5))
		minion.IsDestroyed = true

		// Add minion to the field
		engine.AddEntityToField(player, minion, -1)

		// Process deaths
		engine.processDestroyAndUpdateAura()

		// Check that the minion was removed from the field
		if len(player.Field) != 0 {
			t.Errorf("Expected field to be empty, got %d minions", len(player.Field))
		}

		// Check that the minion was added to the graveyard
		if len(player.Graveyard) != 1 {
			t.Errorf("Expected graveyard size to be 1, got %d", len(player.Graveyard))
		}
	})

	t.Run("Weapons with zero or negative durability are destroyed", func(t *testing.T) {
		// Setup
		g := game.CreateTestGame()
		engine := NewEngine(g)
		engine.StartGame()
		player := g.Players[0]

		// Create and equip a weapon with zero durability
		weapon := game.CreateTestWeaponEntity(g, player, game.WithName("Broken Weapon"), game.WithAttack(3), game.WithHealth(0))
		player.Weapon = weapon

		// Process deaths
		engine.processDestroyAndUpdateAura()

		// Check that the weapon was destroyed
		if player.Weapon != nil {
			t.Errorf("Expected weapon to be destroyed, but it still exists")
		}

		// Check that the weapon was added to the graveyard
		if len(player.Graveyard) != 1 {
			t.Errorf("Expected graveyard size to be 1, got %d", len(player.Graveyard))
		}
	})

	t.Run("Weapons marked as destroyed are removed", func(t *testing.T) {
		// Setup
		g := game.CreateTestGame()
		engine := NewEngine(g)
		engine.StartGame()
		player := g.Players[0]

		// Create and equip a weapon marked for destruction
		weapon := game.CreateTestWeaponEntity(g, player, game.WithName("Marked Weapon"), game.WithAttack(3), game.WithHealth(2))
		weapon.IsDestroyed = true
		player.Weapon = weapon

		// Process deaths
		engine.processDestroyAndUpdateAura()

		// Check that the weapon was destroyed
		if player.Weapon != nil {
			t.Errorf("Expected weapon to be destroyed, but it still exists")
		}

		// Check that the weapon was added to the graveyard
		if len(player.Graveyard) != 1 {
			t.Errorf("Expected graveyard size to be 1, got %d", len(player.Graveyard))
		}
	})

	t.Run("Process continues until no more entities die", func(t *testing.T) {
		// Setup - we need to ensure the code can handle a cascading effect
		// For now we're just testing it runs without errors, as processReborn is empty
		g := game.CreateTestGame()
		engine := NewEngine(g)
		engine.StartGame()
		player := g.Players[0]

		// Create minions with various health values
		minion1 := game.CreateTestMinionEntity(g, player, game.WithName("Dead Minion"), game.WithAttack(1), game.WithHealth(0))
		minion2 := game.CreateTestMinionEntity(g, player, game.WithName("Alive Minion"), game.WithAttack(1), game.WithHealth(2))

		// Add minions to the field
		engine.AddEntityToField(player, minion1, -1)
		engine.AddEntityToField(player, minion2, -1)

		// Also add a weapon
		weapon := game.CreateTestWeaponEntity(g, player, game.WithName("Broken Weapon"), game.WithAttack(3), game.WithHealth(0))
		player.Weapon = weapon

		// Process deaths
		engine.processDestroyAndUpdateAura()

		// Verify final state
		if len(player.Field) != 1 {
			t.Errorf("Expected 1 minion on field, got %d", len(player.Field))
		}

		if player.Weapon != nil {
			t.Errorf("Expected weapon to be destroyed")
		}

		if len(player.Graveyard) != 2 {
			t.Errorf("Expected 2 entities in graveyard, got %d", len(player.Graveyard))
		}
	})
}

func TestAttackRestrictions(t *testing.T) {
	t.Run("Entity cannot attack when exhausted", func(t *testing.T) {
		// Setup
		g := game.CreateTestGame()
		engine := NewEngine(g)
		engine.StartGame()
		player1 := g.Players[0]
		player2 := g.Players[1]

		// Create attacker and defender entities for different players
		attackerEntity := game.CreateTestMinionEntity(g, player1, game.WithName("Test Attacker"), game.WithAttack(3), game.WithHealth(4))
		defenderEntity := game.CreateTestMinionEntity(g, player2, game.WithName("Test Defender"), game.WithAttack(2), game.WithHealth(5))

		// Add minions to respective player's field
		engine.AddEntityToField(player1, attackerEntity, -1)
		engine.AddEntityToField(player2, defenderEntity, -1)

		// Set attacker as exhausted
		attackerEntity.Exhausted = true

		// Attempt to attack
		err := engine.Attack(attackerEntity, defenderEntity, false)

		// Assert that attack is rejected due to exhaustion
		if err == nil {
			t.Error("Expected an error when attacking with exhausted entity, got none")
		}
	})

	t.Run("Entity cannot attack more than once per turn", func(t *testing.T) {
		// Setup
		g := game.CreateTestGame()
		engine := NewEngine(g)
		engine.StartGame()
		player1 := g.Players[0]
		player2 := g.Players[1]

		// Create attacker and defender entities for different players
		attackerEntity := game.CreateTestMinionEntity(g, player1, game.WithName("Test Attacker"), game.WithAttack(3), game.WithHealth(4))
		defenderEntity := game.CreateTestMinionEntity(g, player2, game.WithName("Test Defender"), game.WithAttack(2), game.WithHealth(5))

		// Add minions to respective player's field
		engine.AddEntityToField(player1, attackerEntity, -1)
		engine.AddEntityToField(player2, defenderEntity, -1)

		// Set attacker as ready to attack
		attackerEntity.Exhausted = false

		// First attack should succeed
		err := engine.Attack(attackerEntity, defenderEntity, false)
		if err != nil {
			t.Errorf("Expected first attack to succeed, got error: %v", err)
		}

		// Check that NumAttackThisTurn was incremented
		if attackerEntity.NumAttackThisTurn != 1 {
			t.Errorf("Expected NumAttackThisTurn to be 1, got %d", attackerEntity.NumAttackThisTurn)
		}

		// Second attack should fail
		err = engine.Attack(attackerEntity, defenderEntity, false)
		if err == nil {
			t.Error("Expected an error when attacking twice in one turn, got none")
		}
	})

	t.Run("Entity attack counters are reset at turn start", func(t *testing.T) {
		// Setup
		g := game.CreateTestGame()
		engine := NewEngine(g)
		engine.StartGame()
		player := g.Players[0]

		// Create attacker and set it as having already attacked
		attackerEntity := game.CreateTestMinionEntity(g, player, game.WithName("Test Attacker"), game.WithAttack(3), game.WithHealth(4))
		attackerEntity.NumAttackThisTurn = 1
		attackerEntity.Exhausted = true

		// Add minion to player's field
		engine.AddEntityToField(player, attackerEntity, -1)

		// End the turn twice to reset the attack counters
		engine.EndPlayerTurn()
		engine.EndPlayerTurn()

		// Check that NumAttackThisTurn was reset
		if attackerEntity.NumAttackThisTurn != 0 {
			t.Errorf("Expected NumAttackThisTurn to be reset to 0, got %d", attackerEntity.NumAttackThisTurn)
		}

		// Check that Exhausted was set to false
		if attackerEntity.Exhausted {
			t.Error("Expected Exhausted to be reset to false")
		}
	})

	t.Run("Newly played minions are exhausted", func(t *testing.T) {
		// Setup
		g := game.CreateTestGame()
		engine := NewEngine(g)
		engine.StartGame()
		player := g.Players[0]
		opponent := g.Players[1]

		// Create a minion entity for the hand
		minionEntity := game.CreateTestMinionEntity(g, player, game.WithName("Test Minion"), game.WithAttack(2), game.WithHealth(2))
		player.Hand = append(player.Hand, minionEntity)
		player.Mana = 10 // Ensure enough mana

		// Create a defender for opponent
		defenderEntity := game.CreateTestMinionEntity(g, opponent, game.WithName("Test Defender"), game.WithAttack(1), game.WithHealth(1))
		engine.AddEntityToField(opponent, defenderEntity, -1)

		// Play the minion
		err := engine.PlayCard(player, 0, nil, -1, 0)
		if err != nil {
			t.Fatalf("Failed to play minion: %v", err)
		}

		// Check that the minion is on the field
		if len(player.Field) != 1 {
			t.Fatalf("Expected 1 minion on field, got %d", len(player.Field))
		}

		// Check that the minion is exhausted
		playedMinion := player.Field[0]
		if !playedMinion.Exhausted {
			t.Error("Expected newly played minion to be exhausted")
		}

		// Attempt to attack with the newly played minion
		err = engine.Attack(playedMinion, defenderEntity, false)
		if err == nil {
			t.Error("Expected an error when attacking with newly played minion, got none")
		}
	})
}

// poisonous test cases moved to poisonous_test.go
