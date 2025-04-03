package engine

import (
	"testing"

	"github.com/openhs/internal/game"
)

func TestBasicAttack(t *testing.T) {
	t.Run("Basic attack between minions", func(t *testing.T) {
		// Setup
		g := game.NewGame()
		engine := NewEngine(g)

		// Create a dummy player for entity ownership
		player := game.NewPlayer()

		// Create attacker and defender entities
		attackerCard := &game.Card{
			Name:      "Test Attacker",
			Type:      game.Minion,
			Attack:    3,
			Health:    4,
			MaxHealth: 4,
		}
		attackerEntity := game.NewEntity(attackerCard, player)
		
		defenderCard := &game.Card{
			Name:      "Test Defender",
			Type:      game.Minion,
			Attack:    2,
			Health:    5,
			MaxHealth: 5,
		}
		defenderEntity := game.NewEntity(defenderCard, player)

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
		g := game.NewGame()
		engine := NewEngine(g)

		// Create a dummy player for entity ownership
		player := game.NewPlayer()

		// Create attacker with zero attack
		attackerCard := &game.Card{
			Name:      "Zero Attack Minion",
			Type:      game.Minion,
			Attack:    0,
			Health:    4,
			MaxHealth: 4,
		}
		attackerEntity := game.NewEntity(attackerCard, player)
		
		defenderCard := &game.Card{
			Name:      "Test Defender",
			Type:      game.Minion,
			Attack:    2,
			Health:    5,
			MaxHealth: 5,
		}
		defenderEntity := game.NewEntity(defenderCard, player)

		// Perform attack
		err := engine.Attack(attackerEntity, defenderEntity, false)

		// Assert
		if err == nil {
			t.Error("Expected an error for zero attack minion, got none")
		}
	})

	t.Run("Attack with null entities", func(t *testing.T) {
		// Setup
		g := game.NewGame()
		engine := NewEngine(g)

		// Create a dummy player for entity ownership
		player := game.NewPlayer()

		// Create a valid entity for testing
		validCard := &game.Card{
			Name:      "Valid Minion",
			Type:      game.Minion,
			Attack:    1,
			Health:    1,
			MaxHealth: 1,
		}
		validEntity := game.NewEntity(validCard, player)

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
		g := game.NewGame()
		engine := NewEngine(g)

		// Create a dummy player for entity ownership
		player := game.NewPlayer()

		// Create zero attack attacker that would normally fail validation
		attackerCard := &game.Card{
			Name:      "Zero Attack Minion",
			Type:      game.Minion,
			Attack:    0,
			Health:    4,
			MaxHealth: 4,
		}
		attackerEntity := game.NewEntity(attackerCard, player)
		
		defenderCard := &game.Card{
			Name:      "Test Defender",
			Type:      game.Minion,
			Attack:    2,
			Health:    5,
			MaxHealth: 5,
		}
		defenderEntity := game.NewEntity(defenderCard, player)

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
		g := game.NewGame()
		engine := NewEngine(g)

		// Create a dummy player for entity ownership
		player := game.NewPlayer()

		// Create entities with just enough health to be killed
		attackerCard := &game.Card{
			Name:      "Lethal Attacker",
			Type:      game.Minion,
			Attack:    5,
			Health:    2,
			MaxHealth: 2,
		}
		attackerEntity := game.NewEntity(attackerCard, player)
		
		defenderCard := &game.Card{
			Name:      "Fragile Defender",
			Type:      game.Minion,
			Attack:    2,
			Health:    2,
			MaxHealth: 2,
		}
		defenderEntity := game.NewEntity(defenderCard, player)

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
		g := game.NewGame()
		g.Phase = game.MainAction
		engine := NewEngine(g)

		// Create a dummy player for entity ownership
		player := game.NewPlayer()

		// Create attacker and defender entities
		attackerCard := &game.Card{
			Name:      "Test Attacker",
			Type:      game.Minion,
			Attack:    3,
			Health:    4,
			MaxHealth: 4,
		}
		attackerEntity := game.NewEntity(attackerCard, player)
		
		defenderCard := &game.Card{
			Name:      "Test Defender",
			Type:      game.Minion,
			Attack:    2,
			Health:    5,
			MaxHealth: 5,
		}
		defenderEntity := game.NewEntity(defenderCard, player)

		// Perform attack
		_ = engine.Attack(attackerEntity, defenderEntity, false)

		// Assert phase returns to MainAction
		if g.Phase != game.MainAction {
			t.Errorf("Expected game phase to be MainAction after attack, got %v", g.Phase)
		}
	})

	t.Run("Weapon durability decreases on hero attack", func(t *testing.T) {
		// Setup
		g := game.NewGame()
		engine := NewEngine(g)

		// Create a player with a weapon
		player := game.NewPlayer()
		
		// Create a hero entity
		heroCard := &game.Card{
			Name:      "Test Hero",
			Type:      game.Hero,
			Attack:    0, // Base attack without weapon
			Health:    30,
			MaxHealth: 30,
		}
		heroEntity := game.NewEntity(heroCard, player)
		player.Hero = heroEntity
		heroEntity.Owner = player
		
		// Create and equip a weapon
		weaponCard := &game.Card{
			Name:      "Test Weapon",
			Type:      game.Weapon,
			Attack:    3,
			Health:    2, // Durability for weapons
			MaxHealth: 2,
		}
		weaponEntity := game.NewEntity(weaponCard, player)
		player.Weapon = weaponEntity
		weaponEntity.Owner = player
		
		// Create a defender entity
		defenderCard := &game.Card{
			Name:      "Test Defender",
			Type:      game.Minion,
			Attack:    2,
			Health:    5,
			MaxHealth: 5,
		}
		defenderEntity := game.NewEntity(defenderCard, player)

		// Set hero's attack value to match weapon's attack
		heroEntity.Attack = weaponEntity.Attack

		// Perform attack
		err := engine.Attack(heroEntity, defenderEntity, false)

		// Assert
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		// Check weapon durability decreased
		if player.Weapon.Health != 1 {
			t.Errorf("Expected weapon durability to be 1, got %d", player.Weapon.Health)
		}

		// Perform another attack
		err = engine.Attack(heroEntity, defenderEntity, false)

		// Assert
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		// Check weapon is destroyed after reaching 0 durability
		if player.Weapon != nil {
			t.Errorf("Expected weapon to be destroyed, but it still exists with durability %d", player.Weapon.Health)
		}
	})
} 