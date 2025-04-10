package test

import (
	"testing"

	"github.com/openhs/internal/engine"
	"github.com/openhs/internal/game"
	"github.com/openhs/internal/game/test"
)

func TestRush(t *testing.T) {
	t.Run("Minion with rush can attack a minion the same turn it is played", func(t *testing.T) {
		// Setup
		g := test.CreateTestGame()
		engine := engine.NewEngine(g)
		engine.StartGame()
		player := g.Players[0]
		opponent := g.Players[1]

		// Create a rush minion entity for the hand
		rushMinionEntity := test.CreateTestMinionEntity(g, player,
			test.WithName("Rush Minion"),
			test.WithAttack(3),
			test.WithHealth(2),
			test.WithCost(3),
			test.WithTag(game.TAG_RUSH, true))

		// Add rush minion to player's hand
		player.Hand = []*game.Entity{rushMinionEntity}
		player.Mana = 10 // Ensure enough mana

		// Create a target minion for the opponent
		targetMinionEntity := test.CreateTestMinionEntity(g, opponent,
			test.WithName("Target Minion"),
			test.WithAttack(2),
			test.WithHealth(4))

		// Add target minion to opponent's field
		opponent.Field = append(opponent.Field, targetMinionEntity)

		// Play the rush minion
		err := engine.PlayCard(player, 0, nil, -1, 0)
		if err != nil {
			t.Fatalf("Failed to play rush minion: %v", err)
		}

		// Check that the minion is on the field
		if len(player.Field) != 1 {
			t.Fatalf("Expected 1 minion on player's field, got %d", len(player.Field))
		}

		// Get the played rush minion
		playedRush := player.Field[0]

		// Check that the minion is not exhausted due to rush
		if playedRush.Exhausted {
			t.Error("Expected rush minion to not be exhausted when played")
		}

		// Check that NumTurnInPlay is 0 (first turn)
		if playedRush.NumTurnInPlay != 0 {
			t.Errorf("Expected NumTurnInPlay to be 0, got %d", playedRush.NumTurnInPlay)
		}

		// Attempt to attack a minion with the newly played rush minion
		err = engine.Attack(playedRush, targetMinionEntity, false)

		// Assert attack is successful
		if err != nil {
			t.Errorf("Expected attack with rush minion to succeed against a minion, got error: %v", err)
		}

		// Verify the attack had an effect
		if targetMinionEntity.Health != 1 { // 4 - 3 = 1
			t.Errorf("Expected target health to be 1 after attack, got %d", targetMinionEntity.Health)
		}

		// Check that the rush minion is exhausted after attacking
		if !playedRush.Exhausted {
			t.Error("Expected rush minion to be exhausted after attacking")
		}
	})

	t.Run("Minion with rush cannot attack a hero the same turn it is played", func(t *testing.T) {
		// Setup
		g := test.CreateTestGame()
		engine := engine.NewEngine(g)
		engine.StartGame()
		player := g.Players[0]
		opponent := g.Players[1]

		// Create a rush minion entity for the hand
		rushMinionEntity := test.CreateTestMinionEntity(g, player,
			test.WithName("Rush Minion"),
			test.WithAttack(3),
			test.WithHealth(2),
			test.WithCost(3),
			test.WithTag(game.TAG_RUSH, true))

		// Add rush minion to player's hand
		player.Hand = []*game.Entity{rushMinionEntity}
		player.Mana = 10 // Ensure enough mana

		// Play the rush minion
		err := engine.PlayCard(player, 0, nil, -1, 0)
		if err != nil {
			t.Fatalf("Failed to play rush minion: %v", err)
		}

		// Check that the minion is on the field
		if len(player.Field) != 1 {
			t.Fatalf("Expected 1 minion on player's field, got %d", len(player.Field))
		}

		// Get the played rush minion
		playedRush := player.Field[0]

		// Attempt to attack the opponent's hero with the newly played rush minion
		err = engine.Attack(playedRush, opponent.Hero, false)

		// Assert attack fails because rush minions can't attack heroes on their first turn
		if err == nil {
			t.Error("Expected attack with rush minion against a hero to fail on the first turn, but it succeeded")
		}
	})

	t.Run("Minion with rush can attack a hero after one turn", func(t *testing.T) {
		// Setup
		g := test.CreateTestGame()
		engine := engine.NewEngine(g)
		engine.StartGame()
		player := g.Players[0]
		opponent := g.Players[1]

		// Create a rush minion entity for the field
		rushMinionEntity := test.CreateTestMinionEntity(g, player,
			test.WithName("Rush Minion"),
			test.WithAttack(3),
			test.WithHealth(2),
			test.WithCost(3),
			test.WithTag(game.TAG_RUSH, true))

		// Add rush minion to player's field directly (not from hand)
		player.Field = append(player.Field, rushMinionEntity)

		// Set NumTurnInPlay to 1 (not the first turn anymore)
		rushMinionEntity.NumTurnInPlay = 1

		// Reset exhausted state for the turn
		rushMinionEntity.Exhausted = false
		rushMinionEntity.NumAttackThisTurn = 0

		// Store hero's initial health
		initialHeroHealth := opponent.Hero.Health

		// Attempt to attack the opponent's hero with the rush minion
		err := engine.Attack(rushMinionEntity, opponent.Hero, false)

		// Assert attack succeeds
		if err != nil {
			t.Errorf("Expected attack with rush minion against a hero to succeed after first turn, got error: %v", err)
		}

		// Verify the attack dealt damage to the hero
		if opponent.Hero.Health != initialHeroHealth-rushMinionEntity.Attack {
			t.Errorf("Expected hero health to decrease by %d, but got %d (was %d)",
				rushMinionEntity.Attack, opponent.Hero.Health, initialHeroHealth)
		}
	})
}
