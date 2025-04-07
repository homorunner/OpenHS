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
		player := g.Players[0]

		// Create attacker with poisonous tag
		attackerEntity := test.CreateTestMinionEntity(player,
			test.WithName("Poisonous Minion"),
			test.WithAttack(1),
			test.WithHealth(3),
			test.WithTag(game.TAG_POISONOUS, true))

		// Create defender with high health
		defenderEntity := test.CreateTestMinionEntity(player,
			test.WithName("Tough Minion"),
			test.WithAttack(2),
			test.WithHealth(10))

		// Add minions to player's field
		player.Field = append(player.Field, attackerEntity, defenderEntity)

		// Perform attack, skip validation
		err := engine.Attack(attackerEntity, defenderEntity, true)

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
		if len(player.Field) != 1 {
			t.Errorf("Expected field to have 1 minion, got %d minions", len(player.Field))
		}

		if len(player.Graveyard) != 1 || player.Graveyard[0].Card.Name != "Tough Minion" {
			t.Errorf("Expected Tough Minion to be in graveyard")
		}
	})

	t.Run("Mutual poisonous attack destroys both minions", func(t *testing.T) {
		// Setup
		g := test.CreateTestGame()
		engine := engine.NewEngine(g)
		engine.StartGame()
		player := g.Players[0]

		// Create both minions with poisonous
		attackerEntity := test.CreateTestMinionEntity(player,
			test.WithName("Poisonous Attacker"),
			test.WithAttack(1),
			test.WithHealth(2),
			test.WithTag(game.TAG_POISONOUS, true))

		defenderEntity := test.CreateTestMinionEntity(player,
			test.WithName("Poisonous Defender"),
			test.WithAttack(1),
			test.WithHealth(3),
			test.WithTag(game.TAG_POISONOUS, true))

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
		g := test.CreateTestGame()
		engine := engine.NewEngine(g)
		engine.StartGame()
		player := g.Players[0]

		// Create poisonous minion
		minionEntity := test.CreateTestMinionEntity(player,
			test.WithName("Poisonous Minion"),
			test.WithAttack(2),
			test.WithHealth(2),
			test.WithTag(game.TAG_POISONOUS, true))

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
