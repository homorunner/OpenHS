package card

import (
	"github.com/openhs/internal/game"
	"github.com/openhs/internal/logger"
)

// CardManager handles card creation and management
type CardManager struct {
	cardTemplates map[string]game.Card
}

// NewCardManager creates a new card manager
func NewCardManager() *CardManager {
	return &CardManager{
		cardTemplates: make(map[string]game.Card),
	}
}

// RegisterCard registers a new card template
func (cm *CardManager) RegisterCard(card game.Card) {
	logger.Debug("Registering card template", logger.String("name", card.Name))
	cm.cardTemplates[card.Name] = card
}

// CreateCard creates a new instance of a card from a template
func (cm *CardManager) CreateCard(name string) (*game.Card, error) {
	logger.Debug("Creating card instance", logger.String("name", name))
	template, exists := cm.cardTemplates[name]
	if !exists {
		logger.Error("Card template not found", logger.String("name", name))
		return nil, nil // TODO: Return proper error
	}
	
	card := template
	return &card, nil
}

// LoadCardDatabase loads all card templates
func (cm *CardManager) LoadCardDatabase() error {
	logger.Info("Loading card database")
	// TODO: Load cards from configuration files
	return nil
}

// GetCardTemplate returns a card template by name
func (cm *CardManager) GetCardTemplate(name string) (*game.Card, error) {
	logger.Debug("Retrieving card template", logger.String("name", name))
	template, exists := cm.cardTemplates[name]
	if !exists {
		logger.Error("Card template not found", logger.String("name", name))
		return nil, nil // TODO: Return proper error
	}
	return &template, nil
} 