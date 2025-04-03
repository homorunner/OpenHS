package engine

import (
	"testing"

	"github.com/openhs/internal/game"
	"github.com/openhs/internal/types"
)

func TestAttack(t *testing.T) {
	t.Run("Basic attack between minions", func(t *testing.T) {
		// Setup
		g := game.NewGame()
		engine := NewEngine(g)

		// Create attacker and defender cards
		attacker := createTestMinion().WithAttack(3).WithHealth(4)
		defender := createTestMinion().WithAttack(2).WithHealth(5)

		// Perform attack
		err := engine.Attack(&attacker, &defender, false)

		// Assert
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		// Check health values after attack
		if attacker.Health != 2 {
			t.Errorf("Expected attacker health to be 2, got %d", attacker.Health)
		}

		if defender.Health != 2 {
			t.Errorf("Expected defender health to be 2, got %d", defender.Health)
		}
	})

	t.Run("Attack with zero attack minion", func(t *testing.T) {
		// Setup
		g := game.NewGame()
		engine := NewEngine(g)

		// Create attacker with zero attack
		attacker := createTestMinion().WithAttack(0).WithHealth(4)
		defender := createTestMinion().WithAttack(2).WithHealth(5)

		// Perform attack
		err := engine.Attack(&attacker, &defender, false)

		// Assert
		if err == nil {
			t.Error("Expected an error for zero attack minion, got none")
		}
	})

	t.Run("Attack with null cards", func(t *testing.T) {
		// Setup
		g := game.NewGame()
		engine := NewEngine(g)

		// Perform attack with nil attacker
		err := engine.Attack(nil, &types.Card{}, false)

		// Assert
		if err == nil {
			t.Error("Expected an error for nil attacker, got none")
		}

		// Perform attack with nil defender
		err = engine.Attack(&types.Card{Attack: 1}, nil, false)

		// Assert
		if err == nil {
			t.Error("Expected an error for nil defender, got none")
		}
	})

	t.Run("Attack with skipValidation", func(t *testing.T) {
		// Setup
		g := game.NewGame()
		engine := NewEngine(g)

		// Create zero attack attacker that would normally fail validation
		attacker := createTestMinion().WithAttack(0).WithHealth(4)
		defender := createTestMinion().WithAttack(2).WithHealth(5)

		// Perform attack with skipValidation=true
		err := engine.Attack(&attacker, &defender, true)

		// Assert
		if err != nil {
			t.Errorf("Expected no error with skipValidation, got %v", err)
		}

		// Since attacker has 0 attack, defender shouldn't take damage
		if defender.Health != 5 {
			t.Errorf("Expected defender health to be 5, got %d", defender.Health)
		}

		// Attacker should still take damage
		if attacker.Health != 2 {
			t.Errorf("Expected attacker health to be 2, got %d", attacker.Health)
		}
	})

	t.Run("Attack that kills both cards", func(t *testing.T) {
		// Setup
		g := game.NewGame()
		engine := NewEngine(g)

		// Create cards with just enough health to be killed
		attacker := createTestMinion().WithAttack(5).WithHealth(2)
		defender := createTestMinion().WithAttack(2).WithHealth(2)

		// Perform attack
		err := engine.Attack(&attacker, &defender, false)

		// Assert
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		// Check both cards should be at 0 health or less
		if attacker.Health > 0 {
			t.Errorf("Expected attacker to be dead, got health %d", attacker.Health)
		}

		if defender.Health > 0 {
			t.Errorf("Expected defender to be dead, got health %d", defender.Health)
		}
	})

	t.Run("Game phase changes during attack", func(t *testing.T) {
		// Setup
		g := game.NewGame()
		g.Phase = game.MainAction
		engine := NewEngine(g)

		attacker := createTestMinion().WithAttack(3).WithHealth(4)
		defender := createTestMinion().WithAttack(2).WithHealth(5)

		// Perform attack
		_ = engine.Attack(&attacker, &defender, false)

		// Assert phase returns to MainAction
		if g.Phase != game.MainAction {
			t.Errorf("Expected game phase to be MainAction after attack, got %v", g.Phase)
		}
	})
} 