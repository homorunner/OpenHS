package tests

import (
	"testing"

	edr "github.com/openhs/cards/edr"
	"github.com/openhs/internal/engine"
	"github.com/openhs/internal/game"
	"github.com/openhs/internal/game/test"
)

var card *game.Card

func init() {
	(&edr.ScorchingObserver{}).Register(game.GetCardManager())
	card, _ = game.GetCardManager().CreateCardInstance("Scorching Observer")
}

// TestScorchingObserverProperties tests that Scorching Observer has the correct properties
func TestScorchingObserverProperties(t *testing.T) {
	// Setup
	g := test.CreateTestGame()
	engine := engine.NewEngine(g)
	engine.StartGame()
	player := g.Players[0]

	// Create a Scorching Observer entity
	entity := game.NewEntity(card, g, player)
	engine.AddEntityToField(player, entity, -1)

	// Verify the properties
	if entity.Card.Cost != 9 {
		t.Errorf("Expected Scorching Observer cost to be 9, got %d", entity.Card.Cost)
	}
	if entity.Attack != 7 {
		t.Errorf("Expected Scorching Observer attack to be 7, got %d", entity.Attack)
	}
	if entity.Health != 9 {
		t.Errorf("Expected Scorching Observer health to be 9, got %d", entity.Health)
	}

	// Verify the tags
	if !game.HasTag(entity.Tags, game.TAG_LIFESTEAL) {
		t.Error("Expected Scorching Observer to have Lifesteal")
	}
	if !game.HasTag(entity.Tags, game.TAG_RUSH) {
		t.Error("Expected Scorching Observer to have Rush")
	}
}

// TestScorchingObserverRush tests the Rush ability of Scorching Observer
func TestScorchingObserverRush(t *testing.T) {
	t.Run("Scorching Observer can attack a minion the same turn it is played", func(t *testing.T) {
		// Setup
		g := test.CreateTestGame()
		engine := engine.NewEngine(g)
		engine.StartGame()
		player := g.Players[0]
		opponent := g.Players[1]

		// Create a Scorching Observer entity for the hand
		entity := game.NewEntity(card, g, player)

		// Add to player's hand
		player.Hand = []*game.Entity{entity}
		player.Mana = 10 // Ensure enough mana

		// Create a target minion for the opponent
		targetMinion := test.CreateTestMinionEntity(g, opponent,
			test.WithName("Target Minion"),
			test.WithAttack(5),
			test.WithHealth(10))

		// Add target minion to opponent's field
		engine.AddEntityToField(opponent, targetMinion, -1)

		// Play the Scorching Observer
		err := engine.PlayCard(player, 0, nil, -1, 0)
		if err != nil {
			t.Fatalf("Failed to play Scorching Observer: %v", err)
		}

		// Check that the minion is on the field
		if len(player.Field) != 1 {
			t.Fatalf("Expected 1 minion on player's field, got %d", len(player.Field))
		}

		// Get the played Scorching Observer
		playedObserver := player.Field[0]

		// Check that it's not exhausted due to rush
		if playedObserver.Exhausted {
			t.Error("Expected Scorching Observer to not be exhausted when played")
		}

		// Attempt to attack a minion with the newly played Scorching Observer
		err = engine.Attack(playedObserver, targetMinion, false)

		// Assert attack is successful
		if err != nil {
			t.Errorf("Expected attack with Scorching Observer to succeed against a minion, got error: %v", err)
		}

		// Verify the attack had an effect
		if targetMinion.Health != 3 { // 10 - 7 = 3
			t.Errorf("Expected target health to be 3 after attack, got %d", targetMinion.Health)
		}

		// Check that the Scorching Observer is exhausted after attacking
		if !playedObserver.Exhausted {
			t.Error("Expected Scorching Observer to be exhausted after attacking")
		}
	})

	t.Run("Scorching Observer cannot attack a hero the same turn it is played", func(t *testing.T) {
		// Setup
		g := test.CreateTestGame()
		engine := engine.NewEngine(g)
		engine.StartGame()
		player := g.Players[0]
		opponent := g.Players[1]

		// Create a Scorching Observer entity for the hand
		entity := game.NewEntity(card, g, player)

		// Add to player's hand
		player.Hand = []*game.Entity{entity}
		player.Mana = 10 // Ensure enough mana

		// Play the Scorching Observer
		err := engine.PlayCard(player, 0, nil, -1, 0)
		if err != nil {
			t.Fatalf("Failed to play Scorching Observer: %v", err)
		}

		// Get the played Scorching Observer
		playedObserver := player.Field[0]

		// Attempt to attack the opponent's hero
		err = engine.Attack(playedObserver, opponent.Hero, false)

		// Assert attack fails because rush minions can't attack heroes on their first turn
		if err == nil {
			t.Error("Expected attack with Scorching Observer against a hero to fail on the first turn, but it succeeded")
		}
	})
}

// TestScorchingObserverLifesteal tests the Lifesteal ability of Scorching Observer
func TestScorchingObserverLifesteal(t *testing.T) {
	// Setup
	g := test.CreateTestGame()
	engine := engine.NewEngine(g)
	engine.StartGame()
	player := g.Players[0]
	opponent := g.Players[1]

	// Create a Scorching Observer entity
	entity := game.NewEntity(card, g, player)

	// Add to player's field directly
	engine.AddEntityToField(player, entity, -1)

	// Create a target minion for the opponent
	targetMinion := test.CreateTestMinionEntity(g, opponent,
		test.WithName("Target Minion"),
		test.WithAttack(4),
		test.WithHealth(8))

	// Add target minion to opponent's field
	engine.AddEntityToField(opponent, targetMinion, -1)

	// Set hero health for testing healing
	player.Hero.Health = 20
	player.Hero.MaxHealth = 30

	// Perform attack
	err := engine.Attack(entity, targetMinion, false)
	if err != nil {
		t.Errorf("Attack failed: %v", err)
	}

	// Verify damage was done to both minions
	if entity.Health != 5 { // 9 - 4 = 5
		t.Errorf("Expected Scorching Observer health to be 5, got %d", entity.Health)
	}
	if targetMinion.Health != 1 { // 8 - 7 = 1
		t.Errorf("Expected target minion health to be 1, got %d", targetMinion.Health)
	}

	// Verify lifesteal healing occurred
	if player.Hero.Health != 27 { // 20 + 7 = 27
		t.Errorf("Expected player hero health to be 27 after lifesteal, got %d", player.Hero.Health)
	}
}

// TestScorchingObserverCombined tests both Rush and Lifesteal together
func TestScorchingObserverCombined(t *testing.T) {
	// Setup
	g := test.CreateTestGame()
	engine := engine.NewEngine(g)
	engine.StartGame()
	player := g.Players[0]
	opponent := g.Players[1]

	// Create a Scorching Observer entity for the hand
	entity := game.NewEntity(card, g, player)

	// Add to player's hand
	player.Hand = []*game.Entity{entity}
	player.Mana = 10 // Ensure enough mana

	// Create a target minion for the opponent
	targetMinion := test.CreateTestMinionEntity(g, opponent,
		test.WithName("Target Minion"),
		test.WithAttack(3),
		test.WithHealth(8))

	// Add target minion to opponent's field
	engine.AddEntityToField(opponent, targetMinion, -1)

	// Set hero health for testing healing
	player.Hero.Health = 20
	player.Hero.MaxHealth = 30

	// Play the Scorching Observer
	err := engine.PlayCard(player, 0, nil, -1, 0)
	if err != nil {
		t.Fatalf("Failed to play Scorching Observer: %v", err)
	}

	// Get the played Scorching Observer
	playedObserver := player.Field[0]

	// Attempt to attack a minion with the newly played Scorching Observer
	err = engine.Attack(playedObserver, targetMinion, false)

	// Assert attack is successful
	if err != nil {
		t.Errorf("Expected attack with Scorching Observer to succeed against a minion, got error: %v", err)
	}

	// Verify damage was done to both minions
	if playedObserver.Health != 6 { // 9 - 3 = 6
		t.Errorf("Expected Scorching Observer health to be 6, got %d", playedObserver.Health)
	}
	if targetMinion.Health != 1 { // 8 - 7 = 1
		t.Errorf("Expected target minion health to be 1, got %d", targetMinion.Health)
	}

	// Verify lifesteal healing occurred
	if player.Hero.Health != 27 { // 20 + 7 = 27
		t.Errorf("Expected player hero health to be 27 after lifesteal, got %d", player.Hero.Health)
	}

	// Check that the Scorching Observer is exhausted after attacking
	if !playedObserver.Exhausted {
		t.Error("Expected Scorching Observer to be exhausted after attacking")
	}
}
