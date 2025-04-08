package test

import (
	"testing"

	"github.com/openhs/internal/engine"
	"github.com/openhs/internal/game"
	"github.com/openhs/internal/game/test"
)

func TestCharge(t *testing.T) {
	t.Run("Minion with charge can attack the same turn it is played", func(t *testing.T) {
		// Setup
		g := test.CreateTestGame()
		engine := engine.NewEngine(g)
		engine.StartGame()
		player := g.Players[0]
		opponent := g.Players[1]

		// Create a charge minion entity for the hand
		chargeMinionEntity := test.CreateTestMinionEntity(player,
			test.WithName("Charge Minion"),
			test.WithAttack(3),
			test.WithHealth(2),
			test.WithCost(3),
			test.WithTag(game.TAG_CHARGE, true))

		// Add charge minion to player's hand
		player.Hand = []*game.Entity{chargeMinionEntity}
		player.Mana = 10 // Ensure enough mana

		// Create a target minion for the opponent
		targetMinionEntity := test.CreateTestMinionEntity(opponent,
			test.WithName("Target Minion"),
			test.WithAttack(2),
			test.WithHealth(4))

		// Add target minion to opponent's field
		opponent.Field = append(opponent.Field, targetMinionEntity)

		// Play the charge minion
		err := engine.PlayCard(player, 0, nil, -1, 0)
		if err != nil {
			t.Fatalf("Failed to play charge minion: %v", err)
		}

		// Check that the minion is on the field
		if len(player.Field) != 1 {
			t.Fatalf("Expected 1 minion on player's field, got %d", len(player.Field))
		}

		// Get the played charge minion
		playedCharge := player.Field[0]

		// Check that the minion is not exhausted due to charge
		if playedCharge.Exhausted {
			t.Error("Expected charge minion to not be exhausted when played")
		}

		// Attempt to attack with the newly played minion
		err = engine.Attack(playedCharge, targetMinionEntity, false)

		// Assert attack is successful
		if err != nil {
			t.Errorf("Expected attack with charge minion to succeed, got error: %v", err)
		}

		// Verify the attack had an effect
		if targetMinionEntity.Health != 1 { // 4 - 3 = 1
			t.Errorf("Expected target health to be 1 after attack, got %d", targetMinionEntity.Health)
		}

		// Check that the charge minion is exhausted after attacking
		if !playedCharge.Exhausted {
			t.Error("Expected charge minion to be exhausted after attacking")
		}
	})

	t.Run("Regular minion cannot attack the same turn it is played", func(t *testing.T) {
		// Setup
		g := test.CreateTestGame()
		engine := engine.NewEngine(g)
		engine.StartGame()
		player := g.Players[0]
		opponent := g.Players[1]

		// Create a regular minion entity for the hand
		regularMinionEntity := test.CreateTestMinionEntity(player,
			test.WithName("Regular Minion"),
			test.WithAttack(3),
			test.WithHealth(2),
			test.WithCost(3))

		// Add regular minion to player's hand
		player.Hand = append(player.Hand, regularMinionEntity)
		player.Mana = 10 // Ensure enough mana

		// Create a target minion for the opponent
		targetMinionEntity := test.CreateTestMinionEntity(opponent,
			test.WithName("Target Minion"),
			test.WithAttack(2),
			test.WithHealth(4))

		// Add target minion to opponent's field
		opponent.Field = append(opponent.Field, targetMinionEntity)

		// Play the regular minion
		err := engine.PlayCard(player, 0, nil, -1, 0)
		if err != nil {
			t.Fatalf("Failed to play regular minion: %v", err)
		}

		// Check that the minion is on the field
		if len(player.Field) != 1 {
			t.Fatalf("Expected 1 minion on player's field, got %d", len(player.Field))
		}

		// Get the played regular minion
		playedRegular := player.Field[0]

		// Check that the minion is exhausted
		if !playedRegular.Exhausted {
			t.Error("Expected regular minion to be exhausted when played")
		}

		// Attempt to attack with the newly played minion
		err = engine.Attack(playedRegular, targetMinionEntity, false)

		// Assert attack fails
		if err == nil {
			t.Error("Expected attack with regular minion to fail, but it succeeded")
		}
	})
}
