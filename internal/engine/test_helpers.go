package engine

import "github.com/openhs/internal/game"

func createTestPlayer() *game.Player {
	player := game.NewPlayer()
	player.Hero = createTestHeroEntity(player)

	deck := []*game.Entity{}
	for i := 0; i < 10; i++ {
		deck = append(deck, createTestMinionEntity(player, withName("Test Card")))
	}
	player.Deck = deck

	return player
}

// createTestGame creates a simple game with two players for testing
func createTestGame() *game.Game {
	g := game.NewGame()

	player1 := createTestPlayer()
	player2 := createTestPlayer()

	g.Players = append(g.Players, player1, player2)
	return g
}

// createTestMinionEntity creates a test minion with customizable properties
func createTestMinionEntity(player *game.Player, opts ...func(*game.Card)) *game.Entity {
	card := &game.Card{
		Name:      "Test Minion",
		Type:      game.Minion,
		Cost:      2,
		Attack:    2,
		Health:    3,
		MaxHealth: 3,
	}
	// Apply any option functions
	for _, opt := range opts {
		opt(card)
	}

	return game.NewEntity(card, player)
}

// createTestSpellEntity creates a test spell with customizable properties
func createTestSpellEntity(player *game.Player, opts ...func(*game.Card)) *game.Entity {
	card := &game.Card{
		Name: "Test Spell",
		Type: game.Spell,
		Cost: 1,
	}

	// Apply any option functions
	for _, opt := range opts {
		opt(card)
	}

	return game.NewEntity(card, player)
}

// createTestWeaponEntity creates a test weapon with customizable properties
func createTestWeaponEntity(player *game.Player, opts ...func(*game.Card)) *game.Entity {
	card := &game.Card{
		Name:      "Test Weapon",
		Type:      game.Weapon,
		Cost:      1,
		Attack:    1,
		Health:    4,
		MaxHealth: 4,
	}

	// Apply any option functions
	for _, opt := range opts {
		opt(card)
	}

	return game.NewEntity(card, player)
}

// createTestHeroEntity creates a test hero with customizable properties
func createTestHeroEntity(player *game.Player, opts ...func(*game.Card)) *game.Entity {
	card := &game.Card{
		Name:      "Test Hero",
		Type:      game.Hero,
		Attack:    0,
		Health:    30,
		MaxHealth: 30,
	}

	// Apply any option functions
	for _, opt := range opts {
		opt(card)
	}

	return game.NewEntity(card, player)
}

// Helper functions for common entity customizations
func withName(name string) func(*game.Card) {
	return func(c *game.Card) {
		c.Name = name
	}
}

func withCost(cost int) func(*game.Card) {
	return func(c *game.Card) {
		c.Cost = cost
	}
}

func withAttack(attack int) func(*game.Card) {
	return func(c *game.Card) {
		c.Attack = attack
	}
}

func withHealth(health int) func(*game.Card) {
	return func(c *game.Card) {
		c.Health = health
		c.MaxHealth = health
	}
}

func withTag(tagType game.TagType, value interface{}) func(*game.Card) {
	return func(c *game.Card) {
		c.Tags = append(c.Tags, game.NewTag(tagType, value))
	}
}

// addToHand adds an entity to player's hand
func addToHand(player *game.Player, entity *game.Entity) {
	player.Hand = append(player.Hand, entity)
}

// addToField adds a minion to player's field
func addToField(player *game.Player, entity *game.Entity) {
	player.Field = append(player.Field, entity)
}
