package game

// CreateTestPlayer creates a test player with a hero and deck
func CreateTestPlayer(g *Game) *Player {
	player := NewPlayer()
	player.Hero = CreateTestHeroEntity(g, player)
	player.Hero.CurrentZone = ZONE_PLAY

	deck := []*Entity{}
	for i := 0; i < 10; i++ {
		entity := CreateTestMinionEntity(g, player, WithName("Test Card"))
		entity.CurrentZone = ZONE_DECK
		deck = append(deck, entity)
	}
	player.Deck = deck

	return player
}

// CreateTestGame creates a simple game with two players for testing
// and sets up basic game state (Phase, CurrentTurn, CurrentPlayer)
func CreateTestGame() *Game {
	g := NewGame()

	player1 := CreateTestPlayer(g)
	player2 := CreateTestPlayer(g)

	g.Players = append(g.Players, player1, player2)
	
	// Set up basic game state
	g.Phase = MainAction
	g.CurrentTurn = 1
	g.CurrentPlayerIndex = 0
	g.CurrentPlayer = g.Players[0]
	
	return g
}

// CreateTestMinionEntity creates a test minion with customizable properties
func CreateTestMinionEntity(g *Game, player *Player, opts ...func(*Card)) *Entity {
	card := &Card{
		Name:   "Test Minion",
		Type:   Minion,
		Cost:   2,
		Attack: 2,
		Health: 3,
	}
	// Apply any option functions
	for _, opt := range opts {
		opt(card)
	}

	return NewEntity(card, g, player)
}

// CreateTestSpellEntity creates a test spell with customizable properties
func CreateTestSpellEntity(g *Game, player *Player, opts ...func(*Card)) *Entity {
	card := &Card{
		Name: "Test Spell",
		Type: Spell,
		Cost: 1,
	}

	// Apply any option functions
	for _, opt := range opts {
		opt(card)
	}

	return NewEntity(card, g, player)
}

// CreateTestWeaponEntity creates a test weapon with customizable properties
func CreateTestWeaponEntity(g *Game, player *Player, opts ...func(*Card)) *Entity {
	card := &Card{
		Name:   "Test Weapon",
		Type:   Weapon,
		Cost:   1,
		Attack: 1,
		Health: 4,
	}

	// Apply any option functions
	for _, opt := range opts {
		opt(card)
	}

	return NewEntity(card, g, player)
}

// CreateTestHeroEntity creates a test hero with customizable properties
func CreateTestHeroEntity(g *Game, player *Player, opts ...func(*Card)) *Entity {
	card := &Card{
		Name:   "Test Hero",
		Type:   Hero,
		Attack: 0,
		Health: 30,
	}

	// Apply any option functions
	for _, opt := range opts {
		opt(card)
	}

	return NewEntity(card, g, player)
}

// Helper functions for common entity customizations
func WithName(name string) func(*Card) {
	return func(c *Card) {
		c.Name = name
	}
}

func WithCost(cost int) func(*Card) {
	return func(c *Card) {
		c.Cost = cost
	}
}

func WithAttack(attack int) func(*Card) {
	return func(c *Card) {
		c.Attack = attack
	}
}

func WithHealth(health int) func(*Card) {
	return func(c *Card) {
		c.Health = health
	}
}

func WithTag(tagType TagType, value interface{}) func(*Card) {
	return func(c *Card) {
		c.Tags = append(c.Tags, NewTag(tagType, value))
	}
} 