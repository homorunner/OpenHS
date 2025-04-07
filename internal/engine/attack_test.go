package engine

import (
	"testing"

	"github.com/openhs/internal/game"
)

func TestBasicAttack(t *testing.T) {
	t.Run("Basic attack between minions", func(t *testing.T) {
		// Setup
		g := createTestGame()
		engine := NewEngine(g)
		player := g.Players[0]

		// Create attacker and defender entities
		attackerEntity := createTestMinionEntity(player, withName("Test Attacker"), withAttack(3), withHealth(4))
		defenderEntity := createTestMinionEntity(player, withName("Test Defender"), withAttack(2), withHealth(5))

		// Add minions to player's field
		player.Field = append(player.Field, attackerEntity, defenderEntity)

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
		g := createTestGame()
		engine := NewEngine(g)
		player := g.Players[0]

		// Create attacker with zero attack
		attackerEntity := createTestMinionEntity(player, withName("Zero Attack Minion"), withAttack(0), withHealth(4))
		defenderEntity := createTestMinionEntity(player, withName("Test Defender"), withAttack(2), withHealth(5))

		// Add minions to player's field
		player.Field = append(player.Field, attackerEntity, defenderEntity)

		// Perform attack
		err := engine.Attack(attackerEntity, defenderEntity, false)

		// Assert
		if err == nil {
			t.Error("Expected an error for zero attack minion, got none")
		}
	})

	t.Run("Attack with null entities", func(t *testing.T) {
		// Setup
		g := createTestGame()
		engine := NewEngine(g)
		player := g.Players[0]

		// Create a valid entity for testing
		validEntity := createTestMinionEntity(player, withName("Valid Minion"), withAttack(1), withHealth(1))

		// Perform attack with nil attacker
		err := engine.Attack(nil, validEntity, false)

		// Assert
		if err == nil {
			t.Error("Expected an error for nil attacker, got none")
		}

		// Perform attack with nil defender
		err = engine.Attack(validEntity, nil, false)

		// Assert
		if err == nil {
			t.Error("Expected an error for nil defender, got none")
		}
	})

	t.Run("Attack with skipValidation", func(t *testing.T) {
		// Setup
		g := createTestGame()
		engine := NewEngine(g)
		player := g.Players[0]

		// Create zero attack attacker that would normally fail validation
		attackerEntity := createTestMinionEntity(player, withName("Zero Attack Minion"), withAttack(0), withHealth(4))
		defenderEntity := createTestMinionEntity(player, withName("Test Defender"), withAttack(2), withHealth(5))

		// Add minions to player's field
		player.Field = append(player.Field, attackerEntity, defenderEntity)

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
		g := createTestGame()
		engine := NewEngine(g)
		player := g.Players[0]

		// Create entities with just enough health to be killed
		attackerEntity := createTestMinionEntity(player, withName("Lethal Attacker"), withAttack(5), withHealth(2))
		defenderEntity := createTestMinionEntity(player, withName("Fragile Defender"), withAttack(2), withHealth(2))

		// Add minions to player's field
		player.Field = append(player.Field, attackerEntity, defenderEntity)

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
		g := createTestGame()
		g.Phase = game.MainAction
		engine := NewEngine(g)
		player := g.Players[0]

		// Create attacker and defender entities
		attackerEntity := createTestMinionEntity(player, withName("Test Attacker"), withAttack(3), withHealth(4))
		defenderEntity := createTestMinionEntity(player, withName("Test Defender"), withAttack(2), withHealth(5))

		// Add minions to player's field
		player.Field = append(player.Field, attackerEntity, defenderEntity)

		// Perform attack
		_ = engine.Attack(attackerEntity, defenderEntity, false)

		// Assert phase returns to MainAction
		if g.Phase != game.MainAction {
			t.Errorf("Expected game phase to be MainAction after attack, got %v", g.Phase)
		}
	})

	t.Run("Weapon durability decreases on hero attack", func(t *testing.T) {
		g := createTestGame()
		engine := NewEngine(g)
		player := g.Players[0]
		
		// Give player a weapon
		player.Weapon = createTestWeaponEntity(player, withName("Test Weapon"), withAttack(3), withHealth(2))
		
		// Create a defender entity
		defenderEntity := createTestMinionEntity(player, withName("Test Defender"), withAttack(2), withHealth(5))

		// Add defender to player's field
		player.Field = append(player.Field, defenderEntity)

		// Set hero's attack value to match weapon's attack
		player.Hero.Attack = player.Weapon.Attack

		// Perform attack
		err := engine.Attack(player.Hero, defenderEntity, false)

		// Assert
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		// Check weapon durability decreased
		if player.Weapon.Health != 1 {
			t.Errorf("Expected weapon durability to be 1, got %d", player.Weapon.Health)
		}

		// Perform another attack
		err = engine.Attack(player.Hero, defenderEntity, false)

		// Assert
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		// Check weapon is destroyed after reaching 0 durability
		if player.Weapon != nil {
			t.Errorf("Expected weapon to be destroyed, but it still exists with durability %d", player.Weapon.Health)
		}
	})

	t.Run("Attack with poisonous effect destroys minions", func(t *testing.T) {
		// Setup
		g := createTestGame()
		engine := NewEngine(g)
		player := g.Players[0]
		
		// Create attacker with poisonous tag
		attackerEntity := createTestMinionEntity(player, 
			withName("Poisonous Minion"), 
			withAttack(1), 
			withHealth(3),
			withTag(game.TAG_POISONOUS, true))
		
		// Create defender with high health
		defenderEntity := createTestMinionEntity(player, 
			withName("Tough Minion"), 
			withAttack(2), 
			withHealth(10))
			
		// Add minions to player's field
		player.Field = append(player.Field, attackerEntity, defenderEntity)
			
		// Perform attack
		err := engine.Attack(attackerEntity, defenderEntity, false)
		
		// Assert
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		
		// Attacker should take defender's damage
		if attackerEntity.Health != 1 {
			t.Errorf("Expected attacker health to be 1, got %d", attackerEntity.Health)
		}
		
		// Defender should be marked as destroyed even though it has health remaining
		if defenderEntity.Health != 9 {
			t.Errorf("Expected defender health to be 9, got %d", defenderEntity.Health)
		}
		
		if !defenderEntity.IsDestroyed {
			t.Errorf("Expected defender to be marked as destroyed due to poisonous")
		}
		
		// Check that the defender was moved to the graveyard after processDestroyAndUpdateAura
		if len(player.Field) != 0 {
			t.Errorf("Expected field to be empty, got %d minions", len(player.Field))
		}
		
		if len(player.Graveyard) != 1 || player.Graveyard[0].Card.Name != "Tough Minion" {
			t.Errorf("Expected Tough Minion to be in graveyard")
		}
	})
	
	t.Run("Mutual poisonous attack destroys both minions", func(t *testing.T) {
		// Setup
		g := createTestGame()
		engine := NewEngine(g)
		player := g.Players[0]
		
		// Create both minions with poisonous
		attackerEntity := createTestMinionEntity(player, 
			withName("Poisonous Attacker"), 
			withAttack(1), 
			withHealth(2),
			withTag(game.TAG_POISONOUS, true))
			
		defenderEntity := createTestMinionEntity(player, 
			withName("Poisonous Defender"), 
			withAttack(1), 
			withHealth(3),
			withTag(game.TAG_POISONOUS, true))
			
		// Add them to the field to check graveyard logic
		player.Field = append(player.Field, attackerEntity, defenderEntity)
		
		// Perform attack
		err := engine.Attack(attackerEntity, defenderEntity, false)
		
		// Assert
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		
		// Both should take damage
		if attackerEntity.Health != 1 { // 2 - 1 = 1
			t.Errorf("Expected attacker health to be 1, got %d", attackerEntity.Health)
		}
		
		if defenderEntity.Health != 2 { // 3 - 1 = 2
			t.Errorf("Expected defender health to be 2, got %d", defenderEntity.Health)
		}
		
		// Both should be marked as destroyed due to poisonous
		if !attackerEntity.IsDestroyed {
			t.Errorf("Expected attacker to be marked as destroyed due to poisonous")
		}
		
		if !defenderEntity.IsDestroyed {
			t.Errorf("Expected defender to be marked as destroyed due to poisonous")
		}
		
		// Check that both were moved to the graveyard
		if len(player.Field) != 0 {
			t.Errorf("Expected field to be empty, got %d minions", len(player.Field))
		}
		
		if len(player.Graveyard) != 2 {
			t.Errorf("Expected 2 minions in graveyard, got %d", len(player.Graveyard))
		}
	})
	
	t.Run("Poisonous doesn't affect heroes", func(t *testing.T) {
		// Setup
		g := createTestGame()
		engine := NewEngine(g)
		player := g.Players[0]
		
		// Create poisonous minion
		minionEntity := createTestMinionEntity(player, 
			withName("Poisonous Minion"), 
			withAttack(2), 
			withHealth(2),
			withTag(game.TAG_POISONOUS, true))
			
		// Add minion to player's field
		player.Field = append(player.Field, minionEntity)
			
		// Get hero entity
		heroEntity := player.Hero
		heroEntity.Health = 30
		
		// Perform attack against hero
		err := engine.Attack(minionEntity, heroEntity, false)
		
		// Assert
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		
		// Hero should take damage but not be destroyed by poisonous
		if heroEntity.Health != 28 { // 30 - 2 = 28
			t.Errorf("Expected hero health to be 28, got %d", heroEntity.Health)
		}
		
		if heroEntity.IsDestroyed {
			t.Errorf("Hero should not be affected by poisonous")
		}
	})
}

func TestProcessDestroyAndUpdateAura(t *testing.T) {
	t.Run("Minions with zero or negative health are moved to graveyard", func(t *testing.T) {
		// Setup
		g := createTestGame()
		engine := NewEngine(g)
		player := g.Players[0]

		// Create minions with zero and negative health
		minion1 := createTestMinionEntity(player, withName("Dead Minion 1"), withAttack(1), withHealth(0))
		minion2 := createTestMinionEntity(player, withName("Dead Minion 2"), withAttack(1), withHealth(-1))
		minion3 := createTestMinionEntity(player, withName("Alive Minion"), withAttack(1), withHealth(2))
		
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
		g := createTestGame()
		engine := NewEngine(g)
		player := g.Players[0]
		
		// Create a minion marked for destruction
		minion := createTestMinionEntity(player, withName("Marked Minion"), withAttack(1), withHealth(5))
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
		g := createTestGame()
		engine := NewEngine(g)
		player := g.Players[0]
		
		// Create and equip a weapon with zero durability
		weapon := createTestWeaponEntity(player, withName("Broken Weapon"), withAttack(3), withHealth(0))
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
		g := createTestGame()
		engine := NewEngine(g)
		player := g.Players[0]
		
		// Create and equip a weapon marked for destruction
		weapon := createTestWeaponEntity(player, withName("Marked Weapon"), withAttack(3), withHealth(2))
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
		g := createTestGame()
		engine := NewEngine(g)
		player := g.Players[0]
		
		// Create minions with various health values
		minion1 := createTestMinionEntity(player, withName("Dead Minion"), withAttack(1), withHealth(0))
		minion2 := createTestMinionEntity(player, withName("Alive Minion"), withAttack(1), withHealth(2))
		
		// Add minions to the field
		player.Field = append(player.Field, minion1, minion2)
		
		// Also add a weapon
		weapon := createTestWeaponEntity(player, withName("Broken Weapon"), withAttack(3), withHealth(0))
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