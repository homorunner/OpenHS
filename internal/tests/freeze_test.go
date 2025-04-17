package test

import (
	"testing"

	"github.com/openhs/internal/engine"
	"github.com/openhs/internal/game"
)

func TestFreeze(t *testing.T) {
	t.Run("Frozen entity cannot attack", func(t *testing.T) {
		// Setup game
		g := game.CreateTestGame()
		e := engine.NewEngine(g)
		e.StartGame()

		// Setup players
		player1 := g.Players[0]
		player2 := g.Players[1]

		// Create attacker and defender
		attacker := game.CreateTestMinionEntity(g, player1, 
			game.WithName("Test Attacker"),
			game.WithAttack(2),
			game.WithHealth(2))
		defender := game.CreateTestMinionEntity(g, player2,
			game.WithName("Test Defender"),
			game.WithAttack(2),
			game.WithHealth(2))

		// Add minions to field
		g.AddEntityToField(player1, attacker, 0)
		g.AddEntityToField(player2, defender, 0)

		e.EndPlayerTurn()
		e.EndPlayerTurn()

		// Make sure attacker can attack normally
		err := g.Attack(attacker, defender, false)
		if err != nil {
			t.Errorf("Expected attack to succeed, but got error: %v", err)
		}

		// Reset health
		attacker.Health = 2
		defender.Health = 2

		// Freeze the attacker
		g.Freeze(attacker)

		// Check that the attacker is now frozen
		if !game.HasTag(attacker.Tags, game.TAG_FROZEN) {
			t.Errorf("Expected attacker to be frozen")
		}

		// Try to attack again - should fail
		err = g.Attack(attacker, defender, false)
		if err == nil {
			t.Errorf("Expected attack to fail due to freeze, but it succeeded")
		}
	})

	t.Run("Frozen entity remains frozen if exhausted at end of turn", func(t *testing.T) {
		// Setup game
		g := game.CreateTestGame()
		e := engine.NewEngine(g)
		e.StartGame()

		// Set current player to player 1
		player1 := g.Players[0]
		g.CurrentPlayer = player1
		g.CurrentPlayerIndex = 0

		// Create a minion for player 1
		minion := game.CreateTestMinionEntity(g, player1,
			game.WithName("Test Minion"),
			game.WithAttack(2),
			game.WithHealth(2))
		g.AddEntityToField(player1, minion, 0)

		// Make the minion exhausted (as if it has already attacked)
		minion.Exhausted = true
		
		// Freeze the minion
		g.Freeze(minion)

		// Verify minion is frozen
		if !game.HasTag(minion.Tags, game.TAG_FROZEN) {
			t.Errorf("Expected minion to be frozen")
		}

		// End player 1's turn (this should not unfreeze the exhausted minion)
		e.EndPlayerTurn()

		// Verify minion is still frozen after player 1's turn ends
		if !game.HasTag(minion.Tags, game.TAG_FROZEN) {
			t.Errorf("Expected exhausted minion to still be frozen after turn ends")
		}
	})
	
	t.Run("Frozen entity unfreezes if not exhausted at end of turn", func(t *testing.T) {
		// Setup game
		g := game.CreateTestGame()
		e := engine.NewEngine(g)
		e.StartGame()

		// Set current player to player 1
		player1 := g.Players[0]
		g.CurrentPlayer = player1
		g.CurrentPlayerIndex = 0

		// Create a minion for player 1
		minion := game.CreateTestMinionEntity(g, player1,
			game.WithName("Test Minion"),
			game.WithAttack(2),
			game.WithHealth(2))
		g.AddEntityToField(player1, minion, 0)

		// Make sure the minion is not exhausted
		minion.Exhausted = false
		
		// Freeze the minion
		g.Freeze(minion)

		// Verify minion is frozen
		if !game.HasTag(minion.Tags, game.TAG_FROZEN) {
			t.Errorf("Expected minion to be frozen")
		}

		// End player 1's turn (this should unfreeze the non-exhausted minion)
		e.EndPlayerTurn()

		// Verify minion is unfrozen after player 1's turn ends
		if game.HasTag(minion.Tags, game.TAG_FROZEN) {
			t.Errorf("Expected non-exhausted minion to be unfrozen after turn ends")
		}
	})
} 