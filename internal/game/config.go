package game

// CardConfig represents the configuration for a card
type CardConfig struct {
	Name   string      `json:"name"`
	ZhName string      `json:"zh_name"`
	Cost   int         `json:"cost"`
	Attack int         `json:"attack"`
	Health int         `json:"health"`
	Type   CardType    `json:"type"`
	Tags   []TagConfig `json:"tags,omitempty"`
}

// TagConfig represents the configuration for a card tag
type TagConfig struct {
	Type  string      `json:"type"`
	Value interface{} `json:"value,omitempty"`
}
