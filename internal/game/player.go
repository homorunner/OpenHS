package game

// Entity represents a card instance in play with a reference to its definition and owner
type Entity struct {
	Card      *Card
	Owner     *Player
	Health    int
	MaxHealth int
	Attack    int
	// Track any modifications specific to this instance
	Buffs []Buff
}

// Buff represents a temporary modification to an entity
type Buff struct {
}

// NewEntity creates a new entity from a card
func NewEntity(card *Card, owner *Player) *Entity {
	return &Entity{
		Card:      card,
		Owner:     owner,
		Health:    card.Health,
		MaxHealth: card.MaxHealth,
		Attack:    card.Attack,
		Buffs:     make([]Buff, 0),
	}
}

const (
	DefaultMaxMana      = 10
	DefaultStartingMana = 0
)

// Player represents a player in the game
type Player struct {
	Deck      []*Entity // Using Entity for deck since cards can be modified even in the deck
	Hand      []*Entity
	Field     []*Entity
	Graveyard []*Entity
	Hero      *Entity
	HeroPower *Entity
	Weapon    *Entity

	Mana          int
	MaxMana       int
	TotalMana     int
	FatigueDamage int
	HandSize      int
	FieldSize     int
}

// NewPlayer creates a new player from a configuration
func NewPlayer() *Player {
	return &Player{
		Deck:      make([]*Entity, 0),
		Hand:      make([]*Entity, 0),
		Field:     make([]*Entity, 0),
		Graveyard: make([]*Entity, 0),
		HandSize:  10,
		FieldSize: 7,
		MaxMana:   DefaultMaxMana,
		Mana:      DefaultStartingMana,
		TotalMana: DefaultStartingMana,
	}
}
