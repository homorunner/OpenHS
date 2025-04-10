package test

import (
	"testing"

	"github.com/openhs/internal/engine"
	"github.com/openhs/internal/game"
	"github.com/openhs/internal/game/test"
)

func TestPoisonous(t *testing.T) {
	t.Run("Attack with poisonous effect destroys minions", func(t *testing.T) {
		// Setup
		g := test.CreateTestGame()
		engine := engine.NewEngine(g)
		engine.StartGame()
		player1 := g.Players[0]
		player2 := g.Players[1]

		// Create attacker with poisonous tag
		attackerEntity := test.CreateTestMinionEntity(g, player1,
			test.WithName("Poisonous Minion"),
			test.WithAttack(1),
			test.WithHealth(3),
			test.WithTag(game.TAG_POISONOUS, true))

		// Create defender with high health for opponent
		defenderEntity := test.CreateTestMinionEntity(g, player2,
			test.WithName("Tough Minion"),
			test.WithAttack(2),
			test.WithHealth(10))

		// Add minions to respective players' fields
		player1.Field = append(player1.Field, attackerEntity)
		player2.Field = append(player2.Field, defenderEntity)

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

		// Check that the defender moves to the graveyard when destroyed
		// This happens inside the Attack method, no need to call processDestroyAndUpdateAura separately
		if len(player1.Field) != 1 {
			t.Errorf("Expected attacker's field to have 1 minion, got %d minions", len(player1.Field))
		}

		if len(player2.Field) != 0 {
			t.Errorf("Expected defender's field to be empty, got %d minions", len(player2.Field))
		}

		if len(player2.Graveyard) != 1 || player2.Graveyard[0].Card.Name != "Tough Minion" {
			t.Errorf("Expected Tough Minion to be in defender's graveyard")
		}
	})

	t.Run("Mutual poisonous attack destroys both minions", func(t *testing.T) {
		// Setup
		g := test.CreateTestGame()
		engine := engine.NewEngine(g)
		engine.StartGame()
		player1 := g.Players[0]
		player2 := g.Players[1]

		// Create minions with poisonous for different players
		attackerEntity := test.CreateTestMinionEntity(g, player1,
			test.WithName("Poisonous Attacker"),
			test.WithAttack(1),
			test.WithHealth(2),
			test.WithTag(game.TAG_POISONOUS, true))

		defenderEntity := test.CreateTestMinionEntity(g, player2,
			test.WithName("Poisonous Defender"),
			test.WithAttack(1),
			test.WithHealth(3),
			test.WithTag(game.TAG_POISONOUS, true))

		// Add them to their respective player's field
		player1.Field = append(player1.Field, attackerEntity)
		player2.Field = append(player2.Field, defenderEntity)

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
		// This happens inside the Attack method
		if len(player1.Field) != 0 {
			t.Errorf("Expected attacker's field to be empty, got %d minions", len(player1.Field))
		}

		if len(player2.Field) != 0 {
			t.Errorf("Expected defender's field to be empty, got %d minions", len(player2.Field))
		}

		if len(player1.Graveyard) != 1 {
			t.Errorf("Expected 1 minion in attacker's graveyard, got %d", len(player1.Graveyard))
		}

		if len(player2.Graveyard) != 1 {
			t.Errorf("Expected 1 minion in defender's graveyard, got %d", len(player2.Graveyard))
		}
	})

	t.Run("Poisonous doesn't affect heroes", func(t *testing.T) {
		// Setup
		g := test.CreateTestGame()
		engine := engine.NewEngine(g)
		engine.StartGame()
		player1 := g.Players[0]
		player2 := g.Players[1]

		// Create poisonous minion
		minionEntity := test.CreateTestMinionEntity(g, player1,
			test.WithName("Poisonous Minion"),
			test.WithAttack(2),
			test.WithHealth(2),
			test.WithTag(game.TAG_POISONOUS, true))

		// Add minion to player's field
		player1.Field = append(player1.Field, minionEntity)

		// Get opponent's hero entity
		heroEntity := player2.Hero
		heroEntity.Health = 30

		// Perform attack against opponent's hero
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
