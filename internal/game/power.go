package game

type Power struct {
	Type   PowerType
	Action func(g *Game, source, target *Entity)
}

type PowerType int

const (
	PowerTypeSpell PowerType = iota
	PowerTypeBattlecry
	PowerTypeDeathrattle
	PowerTypeHeroPower
)
