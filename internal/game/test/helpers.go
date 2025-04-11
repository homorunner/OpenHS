package test

import "github.com/openhs/internal/game"

func CreateTestPlayer(g *game.Game) *game.Player {
	player := game.NewPlayer()
	player.Hero = CreateTestHeroEntity(g, player)
	player.Hero.CurrentZone = game.ZONE_PLAY

	deck := []*game.Entity{}
	for i := 0; i < 10; i++ {
		entity := CreateTestMinionEntity(g, player, WithName("Test Card"))
		entity.CurrentZone = game.ZONE_DECK
		deck = append(deck, entity)
	}
	player.Deck = deck

	return player
}

// CreateTestGame creates a simple game with two players for testing
func CreateTestGame() *game.Game {
	g := game.NewGame()

	player1 := CreateTestPlayer(g)
	player2 := CreateTestPlayer(g)

	g.Players = append(g.Players, player1, player2)
	return g
}

// CreateTestMinionEntity creates a test minion with customizable properties
func CreateTestMinionEntity(g *game.Game, player *game.Player, opts ...func(*game.Card)) *game.Entity {
	card := &game.Card{
		Name:   "Test Minion",
		Type:   game.Minion,
		Cost:   2,
		Attack: 2,
		Health: 3,
	}
	// Apply any option functions
	for _, opt := range opts {
		opt(card)
	}

	return game.NewEntity(card, g, player)
}

// CreateTestSpellEntity creates a test spell with customizable properties
func CreateTestSpellEntity(g *game.Game, player *game.Player, opts ...func(*game.Card)) *game.Entity {
	card := &game.Card{
		Name: "Test Spell",
		Type: game.Spell,
		Cost: 1,
	}

	// Apply any option functions
	for _, opt := range opts {
		opt(card)
	}

	return game.NewEntity(card, g, player)
}

// CreateTestWeaponEntity creates a test weapon with customizable properties
func CreateTestWeaponEntity(g *game.Game, player *game.Player, opts ...func(*game.Card)) *game.Entity {
	card := &game.Card{
		Name:   "Test Weapon",
		Type:   game.Weapon,
		Cost:   1,
		Attack: 1,
		Health: 4,
	}

	// Apply any option functions
	for _, opt := range opts {
		opt(card)
	}

	return game.NewEntity(card, g, player)
}

// CreateTestHeroEntity creates a test hero with customizable properties
func CreateTestHeroEntity(g *game.Game, player *game.Player, opts ...func(*game.Card)) *game.Entity {
	card := &game.Card{
		Name:   "Test Hero",
		Type:   game.Hero,
		Attack: 0,
		Health: 30,
	}

	// Apply any option functions
	for _, opt := range opts {
		opt(card)
	}

	return game.NewEntity(card, g, player)
}

// Helper functions for common entity customizations
func WithName(name string) func(*game.Card) {
	return func(c *game.Card) {
		c.Name = name
	}
}

func WithCost(cost int) func(*game.Card) {
	return func(c *game.Card) {
		c.Cost = cost
	}
}

func WithAttack(attack int) func(*game.Card) {
	return func(c *game.Card) {
		c.Attack = attack
	}
}

func WithHealth(health int) func(*game.Card) {
	return func(c *game.Card) {
		c.Health = health
	}
}

func WithTag(tagType game.TagType, value interface{}) func(*game.Card) {
	return func(c *game.Card) {
		c.Tags = append(c.Tags, game.NewTag(tagType, value))
	}
}

// AddToHand adds an entity to player's hand
func AddToHand(player *game.Player, entity *game.Entity) {
	player.Hand = append(player.Hand, entity)
	entity.CurrentZone = game.ZONE_HAND
}

// AddToField adds a minion to player's field
func AddToField(player *game.Player, entity *game.Entity) {
	player.Field = append(player.Field, entity)
	entity.CurrentZone = game.ZONE_PLAY
}
