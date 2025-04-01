package config

import (
	"github.com/openhs/internal/types"
)

// CardConfig represents the configuration for a card
type CardConfig struct {
	Name    string         `json:"name"`
	Cost    int            `json:"cost"`
	Attack  int            `json:"attack"`
	Health  int            `json:"health"`
	Type    types.CardType `json:"type"`
	Effects []EffectConfig `json:"effects,omitempty"`
}

// EffectConfig represents the configuration for a card effect
type EffectConfig struct {
	Trigger    types.Trigger `json:"trigger"`
	Conditions []string      `json:"conditions,omitempty"`
	Action     string        `json:"action"`
} 