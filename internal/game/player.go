package game

const (
	DefaultMaxMana      = 10
	DefaultStartingMana = 0
)

// Player represents a player in the game
type Player struct {
	Deck      []*Entity
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
