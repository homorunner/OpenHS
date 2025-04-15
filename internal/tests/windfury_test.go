package test

import (
	"testing"

	"github.com/openhs/internal/engine"
	"github.com/openhs/internal/game"
)

func TestWindfury(t *testing.T) {
	t.Run("Minion with windfury can attack twice", func(t *testing.T) {
		// Setup
		g := game.CreateTestGame()
		engine := engine.NewEngine(g)
		engine.StartGame()
		player := g.Players[0]
		opponent := g.Players[1]

		// Create attacker with windfury
		attackerEntity := game.CreateTestMinionEntity(g, player,
			game.WithName("Windfury Minion"),
			game.WithAttack(3),
			game.WithHealth(4),
			game.WithTag(game.TAG_WINDFURY, true))

		// Create two defender entities for the opponent
		defender1 := game.CreateTestMinionEntity(g, opponent,
			game.WithName("Defender 1"),
			game.WithAttack(2),
			game.WithHealth(3))
		defender2 := game.CreateTestMinionEntity(g, opponent,
			game.WithName("Defender 2"),
			game.WithAttack(1),
			game.WithHealth(5))

		// Add minions to respective player fields
		engine.AddEntityToField(player, attackerEntity, -1)
		engine.AddEntityToField(opponent, defender1, -1)
		engine.AddEntityToField(opponent, defender2, -1)

		// Make sure minion isn't exhausted
		attackerEntity.Exhausted = false

		// First attack should succeed
		err := engine.Attack(attackerEntity, defender1, false)
		if err != nil {
			t.Errorf("Expected first attack to succeed, got error: %v", err)
		}

		// Check that NumAttackThisTurn was incremented
		if attackerEntity.NumAttackThisTurn != 1 {
			t.Errorf("Expected NumAttackThisTurn to be 1, got %d", attackerEntity.NumAttackThisTurn)
		}

		// Check that attacker is not exhausted after first attack
		if attackerEntity.Exhausted {
			t.Errorf("Expected attacker to not be exhausted after first attack")
		}

		// Second attack should also succeed
		err = engine.Attack(attackerEntity, defender2, false)
		if err != nil {
			t.Errorf("Expected second attack to succeed, got error: %v", err)
		}

		// Check that NumAttackThisTurn was incremented again
		if attackerEntity.NumAttackThisTurn != 2 {
			t.Errorf("Expected NumAttackThisTurn to be 2, got %d", attackerEntity.NumAttackThisTurn)
		}

		// Check that attacker is exhausted after second attack
		if !attackerEntity.Exhausted {
			t.Errorf("Expected attacker to be exhausted after second attack")
		}

		// Third attack should fail
		err = engine.Attack(attackerEntity, defender1, false)
		if err == nil {
			t.Errorf("Expected third attack to fail, but it succeeded")
		}
	})

	t.Run("Normal minion can only attack once", func(t *testing.T) {
		// Setup
		g := game.CreateTestGame()
		engine := engine.NewEngine(g)
		engine.StartGame()
		player := g.Players[0]
		opponent := g.Players[1]

		// Create normal attacker without windfury
		attackerEntity := game.CreateTestMinionEntity(g, player,
			game.WithName("Normal Minion"),
			game.WithAttack(3),
			game.WithHealth(4))

		// Create two defender entities for the opponent
		defender1 := game.CreateTestMinionEntity(g, opponent,
			game.WithName("Defender 1"),
			game.WithAttack(2),
			game.WithHealth(3))
		defender2 := game.CreateTestMinionEntity(g, opponent,
			game.WithName("Defender 2"),
			game.WithAttack(1),
			game.WithHealth(5))

		// Add minions to respective player fields
		engine.AddEntityToField(player, attackerEntity, -1)
		engine.AddEntityToField(opponent, defender1, -1)
		engine.AddEntityToField(opponent, defender2, -1)

		// Make sure minion isn't exhausted
		attackerEntity.Exhausted = false

		// First attack should succeed
		err := engine.Attack(attackerEntity, defender1, false)
		if err != nil {
			t.Errorf("Expected first attack to succeed, got error: %v", err)
		}

		// Check that NumAttackThisTurn was incremented
		if attackerEntity.NumAttackThisTurn != 1 {
			t.Errorf("Expected NumAttackThisTurn to be 1, got %d", attackerEntity.NumAttackThisTurn)
		}

		// Check that attacker is exhausted after first attack
		if !attackerEntity.Exhausted {
			t.Errorf("Expected attacker to be exhausted after first attack")
		}

		// Second attack should fail
		err = engine.Attack(attackerEntity, defender2, false)
		if err == nil {
			t.Errorf("Expected second attack to fail, but it succeeded")
		}
	})

	t.Run("Hero with windfury weapon can attack twice", func(t *testing.T) {
		// Setup
		g := game.CreateTestGame()
		engine := engine.NewEngine(g)
		engine.StartGame()
		player := g.Players[0]
		opponent := g.Players[1]

		// Create a windfury weapon and equip it
		weapon := game.CreateTestWeaponEntity(g, player,
			game.WithName("Windfury Weapon"),
			game.WithAttack(2),
			game.WithHealth(2),
			game.WithTag(game.TAG_WINDFURY, true))

		player.Weapon = weapon

		// Set hero's attack value to match weapon's attack
		player.Hero.Attack = weapon.Attack

		// Create two defender entities for the opponent
		defender1 := game.CreateTestMinionEntity(g, opponent,
			game.WithName("Defender 1"),
			game.WithAttack(2),
			game.WithHealth(3))
		defender2 := game.CreateTestMinionEntity(g, opponent,
			game.WithName("Defender 2"),
			game.WithAttack(1),
			game.WithHealth(5))

		// Add minions to opponent's field
		engine.AddEntityToField(opponent, defender1, -1)
		engine.AddEntityToField(opponent, defender2, -1)

		// Make sure hero isn't exhausted
		player.Hero.Exhausted = false

		// First attack should succeed
		err := engine.Attack(player.Hero, defender1, false)
		if err != nil {
			t.Errorf("Expected first attack to succeed, got error: %v", err)
		}

		// Check that NumAttackThisTurn was incremented
		if player.Hero.NumAttackThisTurn != 1 {
			t.Errorf("Expected NumAttackThisTurn to be 1, got %d", player.Hero.NumAttackThisTurn)
		}

		// Check that hero is not exhausted after first attack
		if player.Hero.Exhausted {
			t.Errorf("Expected hero to not be exhausted after first attack")
		}

		// Second attack should also succeed
		err = engine.Attack(player.Hero, defender2, false)
		if err != nil {
			t.Errorf("Expected second attack to succeed, got error: %v", err)
		}

		// Check that NumAttackThisTurn was incremented again
		if player.Hero.NumAttackThisTurn != 2 {
			t.Errorf("Expected NumAttackThisTurn to be 2, got %d", player.Hero.NumAttackThisTurn)
		}

		// Check that hero is exhausted after second attack
		if !player.Hero.Exhausted {
			t.Errorf("Expected hero to be exhausted after second attack")
		}

		// Third attack should fail
		err = engine.Attack(player.Hero, defender1, false)
		if err == nil {
			t.Errorf("Expected third attack to fail, but it succeeded")
		}
	})
}
