package game

import (
	"testing"
)

func TestBasicAttack(t *testing.T) {
	t.Run("Basic attack between minions", func(t *testing.T) {
		// Setup
		g := CreateTestGame()
		player1 := g.Players[0]
		player2 := g.Players[1]

		// Create attacker and defender entities for different players
		attackerEntity := CreateTestMinionEntity(g, player1, WithName("Test Attacker"), WithAttack(3), WithHealth(4), WithTag(TAG_RUSH, true))
		defenderEntity := CreateTestMinionEntity(g, player2, WithName("Test Defender"), WithAttack(2), WithHealth(5))

		// Add minions to respective player's field
		player1.Field = append(player1.Field, attackerEntity)
		player2.Field = append(player2.Field, defenderEntity)

		// Perform attack
		err := g.Attack(attackerEntity, defenderEntity, false)

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
		g := CreateTestGame()
		player1 := g.Players[0]
		player2 := g.Players[1]

		// Create attacker with zero attack and defender for different players
		attackerEntity := CreateTestMinionEntity(g, player1, WithName("Zero Attack Minion"), WithAttack(0), WithHealth(4))
		defenderEntity := CreateTestMinionEntity(g, player2, WithName("Test Defender"), WithAttack(2), WithHealth(5))

		// Add minions to respective player's field
		player1.Field = append(player1.Field, attackerEntity)
		player2.Field = append(player2.Field, defenderEntity)

		// Perform attack
		err := g.Attack(attackerEntity, defenderEntity, false)

		// Assert
		if err == nil {
			t.Error("Expected an error for zero attack minion, got none")
		}
	})

	t.Run("Attack with null entities", func(t *testing.T) {
		// Setup
		g := CreateTestGame()
		player1 := g.Players[0]
		player2 := g.Players[1]

		// Create valid entities for testing
		validAttacker := CreateTestMinionEntity(g, player1, WithName("Valid Attacker"), WithAttack(1), WithHealth(1))
		validDefender := CreateTestMinionEntity(g, player2, WithName("Valid Defender"), WithAttack(1), WithHealth(1))

		// Add minions to respective player's field
		player1.Field = append(player1.Field, validAttacker)
		player2.Field = append(player2.Field, validDefender)

		// Perform attack with nil attacker
		err := g.Attack(nil, validDefender, false)

		// Assert
		if err == nil {
			t.Error("Expected an error for nil attacker, got none")
		}

		// Perform attack with nil defender
		err = g.Attack(validAttacker, nil, false)

		// Assert
		if err == nil {
			t.Error("Expected an error for nil defender, got none")
		}
	})

	t.Run("Attack with skipValidation", func(t *testing.T) {
		// Setup
		g := CreateTestGame()
		player1 := g.Players[0]
		player2 := g.Players[1]

		// Create zero attack attacker and defender for different players
		attackerEntity := CreateTestMinionEntity(g, player1, WithName("Zero Attack Minion"), WithAttack(0), WithHealth(4))
		defenderEntity := CreateTestMinionEntity(g, player2, WithName("Test Defender"), WithAttack(2), WithHealth(5))

		// Add minions to respective player's field
		player1.Field = append(player1.Field, attackerEntity)
		player2.Field = append(player2.Field, defenderEntity)

		// Perform attack with skipValidation=true
		err := g.Attack(attackerEntity, defenderEntity, true)

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
		g := CreateTestGame()
		player1 := g.Players[0]
		player2 := g.Players[1]

		// Create entities with just enough health to be killed for different players
		attackerEntity := CreateTestMinionEntity(g, player1, WithName("Lethal Attacker"), WithAttack(5), WithHealth(2), WithTag(TAG_RUSH, true))
		defenderEntity := CreateTestMinionEntity(g, player2, WithName("Fragile Defender"), WithAttack(2), WithHealth(2))

		// Add minions to respective player's field
		player1.Field = append(player1.Field, attackerEntity)
		player2.Field = append(player2.Field, defenderEntity)

		// Perform attack
		err := g.Attack(attackerEntity, defenderEntity, false)

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
		g := CreateTestGame()
		g.Phase = MainAction
		player1 := g.Players[0]
		player2 := g.Players[1]

		// Create attacker and defender entities for different players
		attackerEntity := CreateTestMinionEntity(g, player1, WithName("Test Attacker"), WithAttack(3), WithHealth(4))
		defenderEntity := CreateTestMinionEntity(g, player2, WithName("Test Defender"), WithAttack(2), WithHealth(5))

		// Add minions to respective player's field
		player1.Field = append(player1.Field, attackerEntity)
		player2.Field = append(player2.Field, defenderEntity)

		// Perform attack
		_ = g.Attack(attackerEntity, defenderEntity, false)

		// Assert phase returns to MainAction
		if g.Phase != MainAction {
			t.Errorf("Expected game phase to be MainAction after attack, got %v", g.Phase)
		}
	})

	t.Run("Weapon durability decreases on hero attack", func(t *testing.T) {
		g := CreateTestGame()
		player1 := g.Players[0]
		player2 := g.Players[1]

		// Give player a weapon
		player1.Weapon = CreateTestWeaponEntity(g, player1, WithName("Test Weapon"), WithAttack(3), WithHealth(2))

		// Create a defender entity for the opponent
		defenderEntity := CreateTestMinionEntity(g, player2, WithName("Test Defender"), WithAttack(2), WithHealth(5))

		// Add defender to opponent's field
		player2.Field = append(player2.Field, defenderEntity)

		// Set hero's attack value to match weapon's attack
		player1.Hero.Attack = player1.Weapon.Attack

		// Perform attack, skip validation
		err := g.Attack(player1.Hero, defenderEntity, true)

		// Assert
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		// Check weapon durability decreased
		if player1.Weapon.Health != 1 {
			t.Errorf("Expected weapon durability to be 1, got %d", player1.Weapon.Health)
		}

		// Perform another attack
		err = g.Attack(player1.Hero, defenderEntity, true)

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
		g := CreateTestGame()
		player := g.Players[0]

		// Create attacker and defender entities for the same player
		attackerEntity := CreateTestMinionEntity(g, player, WithName("Test Attacker"), WithAttack(3), WithHealth(4), WithTag(TAG_RUSH, true))
		defenderEntity := CreateTestMinionEntity(g, player, WithName("Test Defender"), WithAttack(2), WithHealth(5))

		// Add minions to player's field
		player.Field = append(player.Field, attackerEntity)
		player.Field = append(player.Field, defenderEntity)

		// Attempt normal attack (should fail)
		err := g.Attack(attackerEntity, defenderEntity, false)

		// Assert normal attack fails
		if err == nil {
			t.Error("Expected error when attacking own minion, got none")
		}

		// Attempt attack with skipValidation (should succeed)
		err = g.Attack(attackerEntity, defenderEntity, true)

		// Assert skipValidation attack succeeds
		if err != nil {
			t.Errorf("Expected no error with skipValidation when attacking own minion, got %v", err)
		}

		// Check damage was applied correctly
		if attackerEntity.Health != 2 {
			t.Errorf("Expected attacker health to be 2, got %d", attackerEntity.Health)
		}

		if defenderEntity.Health != 2 {
			t.Errorf("Expected defender health to be 2, got %d", defenderEntity.Health)
		}
	})
}

func TestProcessDestroyAndUpdateAura(t *testing.T) {
	t.Run("Minions with zero health are moved to graveyard", func(t *testing.T) {
		// Setup
		g := CreateTestGame()
		player1 := g.Players[0]
		player2 := g.Players[1]

		// Create entities with zero health
		deadMinion1 := CreateTestMinionEntity(g, player1, WithName("Dead Minion 1"), WithHealth(0))
		deadMinion2 := CreateTestMinionEntity(g, player2, WithName("Dead Minion 2"), WithHealth(0))

		// Add minions to respective player's field
		player1.Field = append(player1.Field, deadMinion1)
		player2.Field = append(player2.Field, deadMinion2)

		// Process deaths
		g.processDestroyAndUpdateAura()

		// Assert minions are moved to graveyard
		if len(player1.Field) != 0 {
			t.Errorf("Expected player1 field to be empty, got %d minions", len(player1.Field))
		}

		if len(player2.Field) != 0 {
			t.Errorf("Expected player2 field to be empty, got %d minions", len(player2.Field))
		}

		if len(player1.Graveyard) != 1 {
			t.Errorf("Expected player1 graveyard to have 1 minion, got %d", len(player1.Graveyard))
		}

		if len(player2.Graveyard) != 1 {
			t.Errorf("Expected player2 graveyard to have 1 minion, got %d", len(player2.Graveyard))
		}
	})

	t.Run("Minions marked for destruction are moved to graveyard", func(t *testing.T) {
		// Setup
		g := CreateTestGame()
		player := g.Players[0]

		// Create minions
		minion1 := CreateTestMinionEntity(g, player, WithName("Healthy Minion"), WithHealth(5))
		minion2 := CreateTestMinionEntity(g, player, WithName("Marked Minion"), WithHealth(5))
		minion2.IsDestroyed = true

		// Add minions to player's field
		player.Field = append(player.Field, minion1, minion2)

		// Process deaths
		g.processDestroyAndUpdateAura()

		// Assert only marked minion is moved to graveyard
		if len(player.Field) != 1 {
			t.Errorf("Expected player field to have 1 minion, got %d", len(player.Field))
		}

		if len(player.Graveyard) != 1 {
			t.Errorf("Expected player graveyard to have 1 minion, got %d", len(player.Graveyard))
		}

		if player.Field[0].Card.Name != "Healthy Minion" {
			t.Errorf("Expected healthy minion to remain on field, got %s", player.Field[0].Card.Name)
		}

		if player.Graveyard[0].Card.Name != "Marked Minion" {
			t.Errorf("Expected marked minion to be in graveyard, got %s", player.Graveyard[0].Card.Name)
		}
	})

	t.Run("Destroyed weapons are moved to graveyard", func(t *testing.T) {
		// Setup
		g := CreateTestGame()
		player := g.Players[0]

		// Give player a destroyed weapon
		player.Weapon = CreateTestWeaponEntity(g, player, WithName("Broken Weapon"), WithHealth(0))

		// Process deaths
		g.processDestroyAndUpdateAura()

		// Assert weapon is moved to graveyard
		if player.Weapon != nil {
			t.Errorf("Expected player weapon to be nil, got %v", player.Weapon)
		}

		if len(player.Graveyard) != 1 {
			t.Errorf("Expected player graveyard to have 1 weapon, got %d", len(player.Graveyard))
		}

		if player.Graveyard[0].Card.Name != "Broken Weapon" {
			t.Errorf("Expected broken weapon to be in graveyard, got %s", player.Graveyard[0].Card.Name)
		}

		if player.Graveyard[0].CurrentZone != ZONE_GRAVEYARD {
			t.Errorf("Expected weapon zone to be ZONE_GRAVEYARD, got %v", player.Graveyard[0].CurrentZone)
		}
	})
}

func TestAttackRestrictions(t *testing.T) {
	t.Run("Rush minions cannot attack heroes on first turn", func(t *testing.T) {
		// Setup
		g := CreateTestGame()
		player1 := g.Players[0]
		player2 := g.Players[1]

		// Create rush minion for player1
		rushMinion := CreateTestMinionEntity(g, player1, WithName("Rush Minion"), WithAttack(3), WithHealth(3), WithTag(TAG_RUSH, true))
		rushMinion.NumTurnInPlay = 0 // Just played this turn

		// Add minion to respective player's field
		player1.Field = append(player1.Field, rushMinion)

		// Try to attack enemy hero
		err := g.Attack(rushMinion, player2.Hero, false)

		// Assert attack is prevented
		if err == nil {
			t.Error("Expected rush minion to be prevented from attacking hero on first turn")
		}
	})

	t.Run("Rush minions can attack heroes after first turn", func(t *testing.T) {
		// Setup
		g := CreateTestGame()
		player1 := g.Players[0]
		player2 := g.Players[1]

		// Create rush minion for player1 that's been in play for a turn
		rushMinion := CreateTestMinionEntity(g, player1, WithName("Rush Minion"), WithAttack(3), WithHealth(3), WithTag(TAG_RUSH, true))
		rushMinion.NumTurnInPlay = 1 // Has been in play for a turn

		// Add minion to respective player's field
		player1.Field = append(player1.Field, rushMinion)

		// Try to attack enemy hero
		err := g.Attack(rushMinion, player2.Hero, false)

		// Assert attack is allowed
		if err != nil {
			t.Errorf("Expected rush minion to be able to attack hero after first turn, got error: %v", err)
		}
	})

	t.Run("Minions become exhausted after attacking", func(t *testing.T) {
		// Setup
		g := CreateTestGame()
		player1 := g.Players[0]
		player2 := g.Players[1]

		// Create minion for player1
		minion := CreateTestMinionEntity(g, player1, WithName("Test Minion"), WithAttack(2), WithHealth(2))
		minion.Exhausted = false

		// Add minion to player's field
		player1.Field = append(player1.Field, minion)

		// Create target for attack
		target := CreateTestMinionEntity(g, player2, WithName("Target Minion"), WithAttack(1), WithHealth(3))
		player2.Field = append(player2.Field, target)

		// Perform attack
		_ = g.Attack(minion, target, false)

		// Assert minion is exhausted
		if !minion.Exhausted {
			t.Error("Expected minion to be exhausted after attacking")
		}

		// Try to attack again
		err := g.Attack(minion, target, false)

		// Assert second attack is prevented
		if err == nil {
			t.Error("Expected second attack to be prevented for exhausted minion")
		}
	})

	t.Run("Windfury minions can attack twice", func(t *testing.T) {
		// Setup
		g := CreateTestGame()
		player1 := g.Players[0]
		player2 := g.Players[1]

		// Create windfury minion for player1
		windfuryMinion := CreateTestMinionEntity(g, player1, WithName("Windfury Minion"), WithAttack(2), WithHealth(4), WithTag(TAG_WINDFURY, true))

		// Add minion to player's field
		player1.Field = append(player1.Field, windfuryMinion)

		// Create target for attack
		target := CreateTestMinionEntity(g, player2, WithName("Target Minion"), WithAttack(1), WithHealth(5))
		player2.Field = append(player2.Field, target)

		// Perform first attack
		err := g.Attack(windfuryMinion, target, false)

		// Assert first attack succeeds
		if err != nil {
			t.Errorf("Expected first attack to succeed, got error: %v", err)
		}

		// Assert minion is not exhausted after first attack
		if windfuryMinion.Exhausted {
			t.Error("Expected windfury minion not to be exhausted after first attack")
		}

		// Perform second attack
		err = g.Attack(windfuryMinion, target, false)

		// Assert second attack succeeds
		if err != nil {
			t.Errorf("Expected second attack to succeed, got error: %v", err)
		}

		// Assert minion is exhausted after second attack
		if !windfuryMinion.Exhausted {
			t.Error("Expected windfury minion to be exhausted after second attack")
		}

		// Try to attack a third time
		err = g.Attack(windfuryMinion, target, false)

		// Assert third attack is prevented
		if err == nil {
			t.Error("Expected third attack to be prevented for exhausted windfury minion")
		}
	})
}
