package card

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/openhs/internal/config"
	"github.com/openhs/internal/logger"
	"github.com/openhs/internal/types"
)

var globalCardManager *CardManager

// GetCardManager returns the global card manager instance
func GetCardManager() *CardManager {
	if globalCardManager == nil {
		globalCardManager = NewCardManager()
	}
	return globalCardManager
}

// InitializeCardManager initializes the global card manager and loads card database
func InitializeCardManager(configDir string) error {
	cm := GetCardManager()
	return cm.LoadCardDatabase(configDir)
}

// CardManager handles card creation and management
type CardManager struct {
	cardTemplates map[string]types.Card
}

// NewCardManager creates a new card manager
func NewCardManager() *CardManager {
	return &CardManager{
		cardTemplates: make(map[string]types.Card),
	}
}

// RegisterCard registers a new card template
func (cm *CardManager) RegisterCard(card types.Card) {
	logger.Debug("Registering card template", logger.String("name", card.Name))
	cm.cardTemplates[card.Name] = card
}

// CreateCard creates a new instance of a card from a template
func (cm *CardManager) CreateCard(name string) (*types.Card, error) {
	logger.Debug("Creating card instance", logger.String("name", name))
	template, exists := cm.cardTemplates[name]
	if !exists {
		err := NewCardError(ErrCardNotFound, fmt.Sprintf("card template not found: %s", name))
		logger.Error("Failed to create card", logger.String("name", name), logger.Err(err))
		return nil, err
	}

	card := template
	return &card, nil
}

// GetCardTemplate returns a card template by name
func (cm *CardManager) GetCardTemplate(name string) (*types.Card, error) {
	logger.Debug("Retrieving card template", logger.String("name", name))
	template, exists := cm.cardTemplates[name]
	if !exists {
		err := NewCardError(ErrCardNotFound, fmt.Sprintf("card template not found: %s", name))
		logger.Error("Failed to get card template", logger.String("name", name), logger.Err(err))
		return nil, err
	}
	return &template, nil
}

// LoadCardDatabase loads all card templates from the specified directory
func (cm *CardManager) LoadCardDatabase(cardConfigDir string) error {
	logger.Info("Loading card database", logger.String("dir", cardConfigDir))

	// Ensure the card config directory exists
	if _, err := os.Stat(cardConfigDir); os.IsNotExist(err) {
		err := NewCardError(ErrCardNotFound, fmt.Sprintf("card config directory not found: %s", cardConfigDir))
		logger.Error("Failed to load card database", logger.String("dir", cardConfigDir), logger.Err(err))
		return err
	}

	// Read all JSON files in the cards directory
	files, err := os.ReadDir(cardConfigDir)
	if err != nil {
		logger.Error("Failed to read card config directory", logger.String("dir", cardConfigDir), logger.Err(err))
		return err
	}

	for _, file := range files {
		data, err := os.ReadFile(filepath.Join(cardConfigDir, file.Name()))
		if err != nil {
			logger.Error("Failed to read card config file",
				logger.String("file", file.Name()),
				logger.Err(err))
			continue
		}

		var cardConfig config.CardConfig
		if err := json.Unmarshal(data, &cardConfig); err != nil {
			logger.Error("Failed to parse card config",
				logger.String("file", file.Name()),
				logger.Err(err))
			continue
		}

		// Convert config to types.Card
		card := types.Card{
			Name:   cardConfig.Name,
			Cost:   cardConfig.Cost,
			Attack: cardConfig.Attack,
			Health: cardConfig.Health,
			Type:   cardConfig.Type,
		}

		// Register the card template
		cm.RegisterCard(card)
		logger.Info("Loaded card template",
			logger.String("name", card.Name),
			logger.String("file", file.Name()))
	}

	logger.Info("Card database loaded successfully")
	return nil
}
