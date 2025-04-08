package engine

import (
	"testing"

	"github.com/openhs/internal/game"
	"github.com/openhs/internal/game/test"
)

func TestBasicAttack(t *testing.T) {
	t.Run("Basic attack between minions", func(t *testing.T) {
		// Setup
		g := test.CreateTestGame()
		engine := NewEngine(g)
		engine.StartGame()
		player1 := g.Players[0]
		player2 := g.Players[1]

		// Create attacker and defender entities for different players
		attackerEntity := test.CreateTestMinionEntity(player1, test.WithName("Test Attacker"), test.WithAttack(3), test.WithHealth(4))
		defenderEntity := test.CreateTestMinionEntity(player2, test.WithName("Test Defender"), test.WithAttack(2), test.WithHealth(5))

		// Add minions to respective player's field
		player1.Field = append(player1.Field, attackerEntity)
		player2.Field = append(player2.Field, defenderEntity)

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
		g := test.CreateTestGame()
		engine := NewEngine(g)
		engine.StartGame()
		player1 := g.Players[0]
		player2 := g.Players[1]

		// Create attacker with zero attack and defender for different players
		attackerEntity := test.CreateTestMinionEntity(player1, test.WithName("Zero Attack Minion"), test.WithAttack(0), test.WithHealth(4))
		defenderEntity := test.CreateTestMinionEntity(player2, test.WithName("Test Defender"), test.WithAttack(2), test.WithHealth(5))

		// Add minions to respective player's field
		player1.Field = append(player1.Field, attackerEntity)
		player2.Field = append(player2.Field, defenderEntity)

		// Perform attack
		err := engine.Attack(attackerEntity, defenderEntity, false)

		// Assert
		if err == nil {
			t.Error("Expected an error for zero attack minion, got none")
		}
	})

	t.Run("Attack with null entities", func(t *testing.T) {
		// Setup
		g := test.CreateTestGame()
		engine := NewEngine(g)
		engine.StartGame()
		player1 := g.Players[0]
		player2 := g.Players[1]

		// Create valid entities for testing
		validAttacker := test.CreateTestMinionEntity(player1, test.WithName("Valid Attacker"), test.WithAttack(1), test.WithHealth(1))
		validDefender := test.CreateTestMinionEntity(player2, test.WithName("Valid Defender"), test.WithAttack(1), test.WithHealth(1))

		// Add minions to respective player's field
		player1.Field = append(player1.Field, validAttacker)
		player2.Field = append(player2.Field, validDefender)

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
		g := test.CreateTestGame()
		engine := NewEngine(g)
		engine.StartGame()
		player1 := g.Players[0]
		player2 := g.Players[1]

		// Create zero attack attacker and defender for different players
		attackerEntity := test.CreateTestMinionEntity(player1, test.WithName("Zero Attack Minion"), test.WithAttack(0), test.WithHealth(4))
		defenderEntity := test.CreateTestMinionEntity(player2, test.WithName("Test Defender"), test.WithAttack(2), test.WithHealth(5))

		// Add minions to respective player's field
		player1.Field = append(player1.Field, attackerEntity)
		player2.Field = append(player2.Field, defenderEntity)

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
		g := test.CreateTestGame()
		engine := NewEngine(g)
		engine.StartGame()
		player1 := g.Players[0]
		player2 := g.Players[1]

		// Create entities with just enough health to be killed for different players
		attackerEntity := test.CreateTestMinionEntity(player1, test.WithName("Lethal Attacker"), test.WithAttack(5), test.WithHealth(2))
		defenderEntity := test.CreateTestMinionEntity(player2, test.WithName("Fragile Defender"), test.WithAttack(2), test.WithHealth(2))

		// Add minions to respective player's field
		player1.Field = append(player1.Field, attackerEntity)
		player2.Field = append(player2.Field, defenderEntity)

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
		g := test.CreateTestGame()
		g.Phase = game.MainAction
		engine := NewEngine(g)
		engine.StartGame()
		player1 := g.Players[0]
		player2 := g.Players[1]

		// Create attacker and defender entities for different players
		attackerEntity := test.CreateTestMinionEntity(player1, test.WithName("Test Attacker"), test.WithAttack(3), test.WithHealth(4))
		defenderEntity := test.CreateTestMinionEntity(player2, test.WithName("Test Defender"), test.WithAttack(2), test.WithHealth(5))

		// Add minions to respective player's field
		player1.Field = append(player1.Field, attackerEntity)
		player2.Field = append(player2.Field, defenderEntity)

		// Perform attack
		_ = engine.Attack(attackerEntity, defenderEntity, false)

		// Assert phase returns to MainAction
		if g.Phase != game.MainAction {
			t.Errorf("Expected game phase to be MainAction after attack, got %v", g.Phase)
		}
	})

	t.Run("Weapon durability decreases on hero attack", func(t *testing.T) {
		g := test.CreateTestGame()
		engine := NewEngine(g)
		engine.StartGame()
		player1 := g.Players[0]
		player2 := g.Players[1]

		// Give player a weapon
		player1.Weapon = test.CreateTestWeaponEntity(player1, test.WithName("Test Weapon"), test.WithAttack(3), test.WithHealth(2))

		// Create a defender entity for the opponent
		defenderEntity := test.CreateTestMinionEntity(player2, test.WithName("Test Defender"), test.WithAttack(2), test.WithHealth(5))

		// Add defender to opponent's field
		player2.Field = append(player2.Field, defenderEntity)

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
		g := test.CreateTestGame()
		engine := NewEngine(g)
		engine.StartGame()
		player := g.Players[0]

		// Create attacker and defender entities for the same player
		attackerEntity := test.CreateTestMinionEntity(player, test.WithName("Attacker"), test.WithAttack(3), test.WithHealth(10))
		defenderEntity := test.CreateTestMinionEntity(player, test.WithName("Same Player Defender"), test.WithAttack(1), test.WithHealth(30))

		// Add minions to player's field
		player.Field = append(player.Field, attackerEntity, defenderEntity)

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
		g := test.CreateTestGame()
		engine := NewEngine(g)
		engine.StartGame()
		player := g.Players[0]

		// Create minions with zero and negative health
		minion1 := test.CreateTestMinionEntity(player, test.WithName("Dead Minion 1"), test.WithAttack(1), test.WithHealth(0))
		minion2 := test.CreateTestMinionEntity(player, test.WithName("Dead Minion 2"), test.WithAttack(1), test.WithHealth(-1))
		minion3 := test.CreateTestMinionEntity(player, test.WithName("Alive Minion"), test.WithAttack(1), test.WithHealth(2))

		// Add minions to the field
		player.Field = append(player.Field, minion1, minion2, minion3)
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
		g := test.CreateTestGame()
		engine := NewEngine(g)
		engine.StartGame()
		player := g.Players[0]

		// Create a minion marked for destruction
		minion := test.CreateTestMinionEntity(player, test.WithName("Marked Minion"), test.WithAttack(1), test.WithHealth(5))
		minion.IsDestroyed = true

		// Add minion to the field
		player.Field = append(player.Field, minion)

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
		g := test.CreateTestGame()
		engine := NewEngine(g)
		engine.StartGame()
		player := g.Players[0]

		// Create and equip a weapon with zero durability
		weapon := test.CreateTestWeaponEntity(player, test.WithName("Broken Weapon"), test.WithAttack(3), test.WithHealth(0))
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
		g := test.CreateTestGame()
		engine := NewEngine(g)
		engine.StartGame()
		player := g.Players[0]

		// Create and equip a weapon marked for destruction
		weapon := test.CreateTestWeaponEntity(player, test.WithName("Marked Weapon"), test.WithAttack(3), test.WithHealth(2))
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
		g := test.CreateTestGame()
		engine := NewEngine(g)
		engine.StartGame()
		player := g.Players[0]

		// Create minions with various health values
		minion1 := test.CreateTestMinionEntity(player, test.WithName("Dead Minion"), test.WithAttack(1), test.WithHealth(0))
		minion2 := test.CreateTestMinionEntity(player, test.WithName("Alive Minion"), test.WithAttack(1), test.WithHealth(2))

		// Add minions to the field
		player.Field = append(player.Field, minion1, minion2)

		// Also add a weapon
		weapon := test.CreateTestWeaponEntity(player, test.WithName("Broken Weapon"), test.WithAttack(3), test.WithHealth(0))
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
		g := test.CreateTestGame()
		engine := NewEngine(g)
		engine.StartGame()
		player1 := g.Players[0]
		player2 := g.Players[1]

		// Create attacker and defender entities for different players
		attackerEntity := test.CreateTestMinionEntity(player1, test.WithName("Test Attacker"), test.WithAttack(3), test.WithHealth(4))
		defenderEntity := test.CreateTestMinionEntity(player2, test.WithName("Test Defender"), test.WithAttack(2), test.WithHealth(5))

		// Add minions to respective player's field
		player1.Field = append(player1.Field, attackerEntity)
		player2.Field = append(player2.Field, defenderEntity)

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
		g := test.CreateTestGame()
		engine := NewEngine(g)
		engine.StartGame()
		player1 := g.Players[0]
		player2 := g.Players[1]

		// Create attacker and defender entities for different players
		attackerEntity := test.CreateTestMinionEntity(player1, test.WithName("Test Attacker"), test.WithAttack(3), test.WithHealth(4))
		defenderEntity := test.CreateTestMinionEntity(player2, test.WithName("Test Defender"), test.WithAttack(2), test.WithHealth(5))

		// Add minions to respective player's field
		player1.Field = append(player1.Field, attackerEntity)
		player2.Field = append(player2.Field, defenderEntity)

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
		g := test.CreateTestGame()
		engine := NewEngine(g)
		engine.StartGame()
		player := g.Players[0]

		// Create attacker and set it as having already attacked
		attackerEntity := test.CreateTestMinionEntity(player, test.WithName("Test Attacker"), test.WithAttack(3), test.WithHealth(4))
		attackerEntity.NumAttackThisTurn = 1
		attackerEntity.Exhausted = true

		// Add minion to player's field
		player.Field = append(player.Field, attackerEntity)

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
		g := test.CreateTestGame()
		engine := NewEngine(g)
		engine.StartGame()
		player := g.Players[0]
		opponent := g.Players[1]

		// Create a minion entity for the hand
		minionEntity := test.CreateTestMinionEntity(player, test.WithName("Test Minion"), test.WithAttack(2), test.WithHealth(2))
		player.Hand = append(player.Hand, minionEntity)
		player.Mana = 10 // Ensure enough mana

		// Create a defender for opponent
		defenderEntity := test.CreateTestMinionEntity(opponent, test.WithName("Test Defender"), test.WithAttack(1), test.WithHealth(1))
		opponent.Field = append(opponent.Field, defenderEntity)

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
