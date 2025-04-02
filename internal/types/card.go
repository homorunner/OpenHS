package types

// Card represents a card in the game
// An empty card is a 0/0 minion with no cost and no effects
type Card struct {
	Name    string
	Cost    int
	Attack  int
	Health  int
	Type    CardType
	Effects []Effect
}

// CardType represents the type of a card
type CardType int

const (
	Minion CardType = iota
	Spell
	Weapon
	Hero
	HeroPower
)

// Effect represents a card effect or ability
type Effect struct {
	Trigger    Trigger
	Action     Action
	Conditions []Condition
}

// Trigger represents when an effect should activate
type Trigger int

const (
	OnPlay Trigger = iota
	OnDeath
	OnDamage
	OnHeal
	OnTurnStart
	OnTurnEnd
)

// Action represents what an effect does
type Action interface {
	Execute(*Card)
}

// Condition represents a requirement for an effect to trigger
type Condition interface {
	IsMet(*Card) bool
}
