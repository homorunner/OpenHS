package game

type Game struct {
	Players     []Player
	CurrentTurn int
	Phase       GamePhase
}

type Player struct {
	Mana       int
	MaxMana    int
	Deck       []Card
	Hand       []Card
	Board      []Card
	Hero       Card
	HeroPower  Card
}

type Card struct {
	Name        string
	Cost        int
	Attack      int
	Health      int
	Type        CardType
	Effects     []Effect
}

type CardType int

const (
	Minion CardType = iota
	Spell
	Weapon
	Hero
	HeroPower
)

type GamePhase int

const (
	StartGame GamePhase = iota
	StartTurn
	Draw
	Play
	Combat
	EndTurn
	EndGame
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
	Execute(*Game, *Card)
}

// Condition represents a requirement for an effect to trigger
type Condition interface {
	IsMet(*Game, *Card) bool
}

func NewGame() *Game {
	return &Game{
		Players:     make([]Player, 2),
		CurrentTurn: 1,
		Phase:       StartGame,
	}
} 