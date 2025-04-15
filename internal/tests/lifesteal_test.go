package test

import (
	"testing"

	"github.com/openhs/internal/engine"
	"github.com/openhs/internal/game"
)

// TestLifesteal tests that entities with lifesteal heal their controller's hero
// when they deal damage
func TestLifesteal(t *testing.T) {
	g := game.CreateTestGame()
	e := engine.NewEngine(g)
	e.StartGame()

	player1 := g.Players[0]
	player2 := g.Players[1]

	// Test 1: Entity with lifesteal deals damage
	// Set up a minion with lifesteal
	source := game.CreateTestMinionEntity(g, player1,
		game.WithName("Lifesteal Minion"),
		game.WithAttack(3),
		game.WithHealth(3),
		game.WithTag(game.TAG_LIFESTEAL, true))

	// Set up target and damage it
	target := game.CreateTestMinionEntity(g, player2,
		game.WithName("Target Minion"),
		game.WithHealth(5))

	// Reduce hero health for testing healing
	player1.Hero.Health = 20
	player1.Hero.MaxHealth = 30

	// Deal damage with lifesteal minion
	e.DealDamage(source, target, 3)

	// Verify target took damage
	if target.Health != 2 {
		t.Errorf("Expected target health to be 2, got %d", target.Health)
	}

	// Verify hero was healed by the same amount
	if player1.Hero.Health != 23 {
		t.Errorf("Expected hero health to be 23 after lifesteal, got %d", player1.Hero.Health)
	}

	// Test 2: Entity without lifesteal doesn't heal
	normalSource := game.CreateTestMinionEntity(g, player1,
		game.WithName("Normal Minion"),
		game.WithAttack(2))

	target2 := game.CreateTestMinionEntity(g, player2,
		game.WithName("Target Minion 2"),
		game.WithHealth(4))

	// Set hero health again
	player1.Hero.Health = 20

	// Deal damage with normal minion
	e.DealDamage(normalSource, target2, 2)

	// Verify target took damage
	if target2.Health != 2 {
		t.Errorf("Expected target2 health to be 2, got %d", target2.Health)
	}

	// Verify hero was NOT healed
	if player1.Hero.Health != 20 {
		t.Errorf("Expected hero health to remain 20 without lifesteal, got %d", player1.Hero.Health)
	}

	// Test 3: Healing beyond max health is capped
	player1.Hero.Health = 28
	player1.Hero.MaxHealth = 30

	e.DealDamage(source, target, 3)

	// Verify hero health is capped at max health
	if player1.Hero.Health != 30 {
		t.Errorf("Expected hero health to be capped at 30 after lifesteal, got %d", player1.Hero.Health)
	}

	// Test 4: Test nil source (shouldn't panic)
	targetForNil := game.CreateTestMinionEntity(g, player2,
		game.WithHealth(5))

	// This should not panic and should just deal damage
	e.DealDamage(nil, targetForNil, 2)

	if targetForNil.Health != 3 {
		t.Errorf("Expected target health to be 3 after nil source damage, got %d", targetForNil.Health)
	}
}

// TestWeaponLifesteal tests that weapons with lifesteal heal their controller's hero
// when the hero attacks with the weapon
func TestWeaponLifesteal(t *testing.T) {
	g := game.CreateTestGame()
	e := engine.NewEngine(g)
	e.StartGame()

	player1 := g.Players[0]
	player2 := g.Players[1]

	// Set up a weapon with lifesteal
	weapon := game.CreateTestWeaponEntity(g, player1,
		game.WithName("Lifesteal Weapon"),
		game.WithAttack(4),
		game.WithHealth(2), // Durability
		game.WithTag(game.TAG_LIFESTEAL, true))

	// Equip the weapon
	player1.Weapon = weapon

	// Set the hero's attack to match the weapon
	player1.Hero.Attack = weapon.Attack

	// Set up target minion
	target := game.CreateTestMinionEntity(g, player2,
		game.WithName("Target Minion"),
		game.WithAttack(0),
		game.WithHealth(8))

	// Add target to the opponent's field
	e.AddEntityToField(player2, target, -1)

	// Reduce hero health for testing healing
	player1.Hero.Health = 15
	player1.Hero.MaxHealth = 30

	// Perform attack with the hero (using the weapon)
	err := e.Attack(player1.Hero, target, true)
	if err != nil {
		t.Errorf("Attack failed: %v", err)
	}

	// Verify target took damage
	if target.Health != 4 {
		t.Errorf("Expected target health to be 4, got %d", target.Health)
	}

	// Verify hero was healed by the weapon's lifesteal
	if player1.Hero.Health != 19 {
		t.Errorf("Expected hero health to be 19 after lifesteal, got %d", player1.Hero.Health)
	}

	// Verify weapon durability decreased
	if weapon.Health != 1 {
		t.Errorf("Expected weapon durability to be 1, got %d", weapon.Health)
	}

	// Test with a non-lifesteal weapon to verify no healing happens
	player1.Hero.Health = 15 // Reset hero health

	// Replace with a normal weapon without lifesteal
	normalWeapon := game.CreateTestWeaponEntity(g, player1,
		game.WithName("Normal Weapon"),
		game.WithAttack(3),
		game.WithHealth(2)) // Durability

	player1.Weapon = normalWeapon
	player1.Hero.Attack = normalWeapon.Attack

	// Create a new target
	target2 := game.CreateTestMinionEntity(g, player2,
		game.WithName("Target Minion 2"),
		game.WithAttack(0),
		game.WithHealth(6))

	e.AddEntityToField(player2, target2, -1)

	// Attack with the non-lifesteal weapon
	err = e.Attack(player1.Hero, target2, true)
	if err != nil {
		t.Errorf("Attack failed: %v", err)
	}

	// Verify target took damage
	if target2.Health != 3 {
		t.Errorf("Expected target2 health to be 3, got %d", target2.Health)
	}

	// Verify hero health did NOT change (no lifesteal)
	if player1.Hero.Health != 15 {
		t.Errorf("Expected hero health to remain 15 (no lifesteal), got %d", player1.Hero.Health)
	}
}

// TestLifestealDuringAttack tests that lifesteal works during combat
func TestLifestealDuringAttack(t *testing.T) {
	g := game.CreateTestGame()
	e := engine.NewEngine(g)
	e.StartGame()

	player1 := g.Players[0]
	player2 := g.Players[1]

	// Add minions to the board
	attacker := game.CreateTestMinionEntity(g, player1,
		game.WithName("Lifesteal Attacker"),
		game.WithAttack(4),
		game.WithHealth(5),
		game.WithTag(game.TAG_LIFESTEAL, true))

	defender := game.CreateTestMinionEntity(g, player2,
		game.WithName("Defender"),
		game.WithAttack(2),
		game.WithHealth(6))

	e.AddEntityToField(player1, attacker, -1)
	e.AddEntityToField(player2, defender, -1)

	// Set hero health
	player1.Hero.Health = 20
	player1.Hero.MaxHealth = 30

	// Perform attack
	err := e.Attack(attacker, defender, true)
	if err != nil {
		t.Errorf("Attack failed: %v", err)
	}

	// Verify damage was done to both minions
	if attacker.Health != 3 {
		t.Errorf("Expected attacker health to be 3, got %d", attacker.Health)
	}
	if defender.Health != 2 {
		t.Errorf("Expected defender health to be 2, got %d", defender.Health)
	}

	// Verify lifesteal healing occurred
	if player1.Hero.Health != 24 {
		t.Errorf("Expected player1 hero health to be 24 after lifesteal, got %d", player1.Hero.Health)
	}

	// Player2's minion doesn't have lifesteal, so hero health shouldn't change
	if player2.Hero.Health != 30 {
		t.Errorf("Expected player2 hero health to remain 30, got %d", player2.Hero.Health)
	}
}
