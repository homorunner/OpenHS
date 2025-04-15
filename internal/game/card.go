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
	Tags        []Tag                    // Card tags like Taunt, Divine Shield, etc.
	Load        func(g *Game, e *Entity) // Load functions register triggers to Game for Entity of this card
	Unload      func(g *Game, e *Entity) // Unload functions remove triggers from Game when Entity is removed/silenced/...
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
