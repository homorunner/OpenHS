package game

// Card represents a card in the game
// An empty card is a 0/0 minion with no cost and no effects
type Card struct {
	Name        string
	ZhName      string
	ID          string
	Description string
	Cost        int
	Attack      int
	Health      int
	Type        CardType
	Effects     []Effect
	Tags        []Tag // Card tags like Taunt, Divine Shield, etc.
	Load        func(g *Game, e *Entity)
	Unload      func(g *Game, e *Entity)
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

// String returns a string representation of the CardType
func (c CardType) String() string {
	switch c {
	case Minion:
		return "Minion"
	case Spell:
		return "Spell"
	case Weapon:
		return "Weapon"
	case Hero:
		return "Hero"
	case HeroPower:
		return "Hero Power"
	default:
		return "Unknown"
	}
}

func (c CardType) ZhString() string {
	switch c {
	case Minion:
		return "随从"
	case Spell:
		return "法术"
	case Weapon:
		return "武器"
	case Hero:
		return "英雄"
	case HeroPower:
		return "英雄技能"
	default:
		return "未知"
	}
}

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

// CardConfig represents the configuration for a card
type CardConfig struct {
	Name    string         `json:"name"`
	ZhName  string         `json:"zh_name"`
	Cost    int            `json:"cost"`
	Attack  int            `json:"attack"`
	Health  int            `json:"health"`
	Type    CardType       `json:"type"`
	Effects []EffectConfig `json:"effects,omitempty"`
	Tags    []TagConfig    `json:"tags,omitempty"`
}

// TagConfig represents the configuration for a card tag
type TagConfig struct {
	Type  string      `json:"type"`
	Value interface{} `json:"value,omitempty"`
}

// EffectConfig represents the configuration for a card effect
type EffectConfig struct {
	Trigger    Trigger  `json:"trigger"`
	Conditions []string `json:"conditions,omitempty"`
	Action     string   `json:"action"`
}
