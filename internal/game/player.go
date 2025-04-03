package game

const (
	DefaultMaxMana      = 10
	DefaultStartingMana = 0
)

// Player represents a player in the game
type Player struct {
	Deck      []Card
	Hand      []Card
	Field     []Card
	Graveyard []Card
	Hero      Card
	HeroPower Card
	Weapon    Card

	Mana          int
	MaxMana       int
	TotalMana     int
	FatigueDamage int
	HandSize      int
	HasWeapon     bool
}

// NewPlayer creates a new player from a configuration
func NewPlayer() *Player {
	return &Player{
		Deck:      make([]Card, 0),
		Hand:      make([]Card, 0),
		Field:     make([]Card, 0),
		Graveyard: make([]Card, 0),
		HandSize:  10,
		MaxMana:   DefaultMaxMana,
		Mana:      DefaultStartingMana,
		TotalMana: DefaultStartingMana,
	}
}
