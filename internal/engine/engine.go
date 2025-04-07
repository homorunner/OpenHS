package engine

import (
	"errors"

	"github.com/openhs/internal/game"
	"github.com/openhs/internal/logger"
)

// Engine handles the game rules and mechanics
type Engine struct {
	game      *game.Game
	nextPhase game.GamePhase
	autoRun   bool
}

// NewEngine creates a new game engine
func NewEngine(g *game.Game) *Engine {
	return &Engine{
		game:      g,
		nextPhase: game.InvalidPhase,
		autoRun:   true,
	}
}

// SetAutoRun sets whether the engine should automatically progress to the next phase
func (e *Engine) SetAutoRun(autoRun bool) {
	e.autoRun = autoRun
}

// StartGame begins the game flow
func (e *Engine) StartGame() error {
	// Trigger start of game effects if needed

	// Set the first phase
	e.nextPhase = game.BeginFirst

	// Process the first phase if autoRun is enabled
	if e.autoRun {
		return e.ProcessNextPhase()
	}

	return nil
}

// ProcessNextPhase processes the next game phase
func (e *Engine) ProcessNextPhase() error {
	if e.nextPhase == game.InvalidPhase {
		return errors.New("invalid next phase")
	}

	// Update current phase
	e.game.Phase = e.nextPhase

	// Process phase based on type
	var err error
	switch e.nextPhase {
	case game.BeginFirst:
		err = e.beginFirst()
	case game.BeginShuffle:
		err = e.beginShuffle()
	case game.BeginDraw:
		err = e.beginDraw()
	case game.BeginMulligan:
		err = e.beginMulligan()
	case game.MainBegin:
		err = e.mainBegin()
	case game.MainReady:
		err = e.mainReady()
	case game.MainStartTriggers:
		err = e.mainStartTriggers()
	case game.MainResource:
		err = e.mainResource()
	case game.MainDraw:
		err = e.mainDraw()
	case game.MainStart:
		err = e.mainStart()
	case game.MainAction:
		err = e.mainAction()
	case game.MainEnd:
		err = e.mainEnd()
	case game.MainCleanup:
		err = e.mainCleanup()
	case game.MainNext:
		err = e.mainNext()
	case game.FinalWrapup:
		err = e.finalWrapup()
	case game.FinalGameover:
		err = e.finalGameover()
	default:
		err = errors.New("unknown game phase")
	}

	if err != nil {
		return err
	}

	// Process the next phase if autoRun is enabled
	if e.autoRun && e.game.Phase != game.MainAction && e.game.Phase != game.FinalGameover {
		return e.ProcessNextPhase()
	}

	return nil
}

// ProcessUntil processes game phases until the specified phase
func (e *Engine) ProcessUntil(phase game.GamePhase) error {
	autoRunBackup := e.autoRun
	e.autoRun = false

	for e.game.Phase != phase {
		if err := e.ProcessNextPhase(); err != nil {
			e.autoRun = autoRunBackup
			return err
		}
	}

	e.autoRun = autoRunBackup
	return nil
}

// Phase transition implementations

func (e *Engine) beginFirst() error {
	// Initial game setup
	logger.Debug("Phase: Begin First")

	// Set next phase
	e.nextPhase = game.BeginShuffle
	return nil
}

func (e *Engine) beginShuffle() error {
	logger.Debug("Phase: Begin Shuffle")

	// Shuffle player decks
	// This would involve shuffling the Deck slice for each player
	// TODO: Implement deck shuffling logic

	// Set next phase
	e.nextPhase = game.BeginDraw
	return nil
}

func (e *Engine) beginDraw() error {
	logger.Debug("Phase: Begin Draw")

	// Draw initial hands for players
	// First player (index 0) draws 3 cards, second player (index 1) draws 4 cards
	if len(e.game.Players) >= 2 {
		// Draw 3 cards for the first player
		for i := 0; i < 3; i++ {
			e.DrawCard(e.game.Players[0])
		}

		// Draw 4 cards for the second player
		for i := 0; i < 4; i++ {
			e.DrawCard(e.game.Players[1])
		}
	}

	// Set next phase based on whether to skip mulligan or not
	// For now we'll always skip mulligan
	// TODO: Add skipMulligan configuration to game config
	skipMulligan := true
	if skipMulligan {
		e.nextPhase = game.MainBegin
	} else {
		e.nextPhase = game.BeginMulligan
	}

	return nil
}

func (e *Engine) beginMulligan() error {
	logger.Debug("Phase: Begin Mulligan")

	// Handle mulligan phase
	// Players select cards to redraw
	// TODO: Implement mulligan logic

	// Set next phase
	e.nextPhase = game.MainBegin
	return nil
}

func (e *Engine) mainBegin() error {
	logger.Debug("Phase: Main Begin")

	// Main phase begins

	// Set the current turn and player
	e.game.CurrentTurn = 1
	e.game.CurrentPlayerIndex = 0
	e.game.CurrentPlayer = e.game.Players[0]

	// Give "The Coin" to second player
	// TODO: Implement coin logic

	// Set next phase
	e.nextPhase = game.MainReady
	return nil
}

func (e *Engine) mainReady() error {
	logger.Debug("Phase: Main Ready")

	// Reset player states for new turn
	player := e.game.CurrentPlayer

	// Reset attack counters and exhaustion for all entities controlled by the player
	// Hero
	if player.Hero != nil {
		player.Hero.NumAttackThisTurn = 0
		player.Hero.Exhausted = false
	}

	// Field minions
	for _, minion := range player.Field {
		minion.NumAttackThisTurn = 0
		minion.Exhausted = false
	}

	// Set next phase
	e.nextPhase = game.MainStartTriggers
	return nil
}

func (e *Engine) mainStartTriggers() error {
	logger.Debug("Phase: Main Start Triggers")

	// Process start-of-turn triggers
	// TODO: Implement trigger system

	// Set next phase
	e.nextPhase = game.MainResource
	return nil
}

func (e *Engine) mainResource() error {
	logger.Debug("Phase: Main Resource")

	// Give mana crystal to current player
	player := e.game.CurrentPlayer

	// Increase total mana by 1, but do not exceed MaxMana
	if player.TotalMana < player.MaxMana {
		player.TotalMana++
	}

	// Restore mana to current total
	player.Mana = player.TotalMana

	// Set next phase
	e.nextPhase = game.MainDraw
	return nil
}

func (e *Engine) mainDraw() error {
	logger.Debug("Phase: Main Draw")

	// Draw a card for the current player
	e.DrawCard(e.game.CurrentPlayer)

	// Set next phase
	e.nextPhase = game.MainStart
	return nil
}

func (e *Engine) mainStart() error {
	logger.Debug("Phase: Main Start")

	// Process any destroyed entities and update auras
	// TODO: Implement destruction and aura logic

	// Set next phase
	e.nextPhase = game.MainAction
	return nil
}

func (e *Engine) mainAction() error {
	logger.Debug("Phase: Main Action")

	// Player action phase - no automatic transition
	// The game waits for player input

	// This phase doesn't automatically transition
	return nil
}

func (e *Engine) mainEnd() error {
	logger.Debug("Phase: Main End")

	// Process end-of-turn triggers
	// TODO: Implement end-of-turn logic

	// Set next phase
	e.nextPhase = game.MainCleanup
	return nil
}

func (e *Engine) mainCleanup() error {
	logger.Debug("Phase: Main Cleanup")

	// Clean up one-turn effects
	// TODO: Implement cleanup logic

	// Set next phase
	e.nextPhase = game.MainNext
	return nil
}

func (e *Engine) mainNext() error {
	logger.Debug("Phase: Main Next")

	// Switch to next player
	e.game.CurrentPlayerIndex = (e.game.CurrentPlayerIndex + 1) % len(e.game.Players)
	e.game.CurrentPlayer = e.game.Players[e.game.CurrentPlayerIndex]

	// Increment turn counter
	e.game.CurrentTurn++

	// Set next phase
	e.nextPhase = game.MainReady
	return nil
}

func (e *Engine) finalWrapup() error {
	logger.Debug("Phase: Final Wrapup")

	// Determine game result
	// TODO: Implement game result logic

	// Set next phase
	e.nextPhase = game.FinalGameover
	return nil
}

func (e *Engine) finalGameover() error {
	logger.Debug("Phase: Final Gameover")

	// Game is over
	return nil
}

// EndPlayerTurn transitions from the action phase to the end phase
func (e *Engine) EndPlayerTurn() error {
	if e.game.Phase != game.MainAction {
		return errors.New("can only end turn during action phase")
	}

	e.nextPhase = game.MainEnd
	return e.ProcessNextPhase()
}

// PerformPlayerAction processes a player's action during the action phase
// This would handle playing cards, attacking, using hero power, etc.
func (e *Engine) PerformPlayerAction() error {
	if e.game.Phase != game.MainAction {
		return errors.New("can only perform actions during action phase")
	}

	// Process the action
	// TODO: Implement action processing

	return nil
}

// CheckGameOver checks if the game is over and transitions to the appropriate phase
func (e *Engine) CheckGameOver() bool {
	// Check if any player has lost
	// TODO: Implement game over conditions

	// If game is over, transition to final wrap up
	// e.nextPhase = game.FinalWrapup
	// e.ProcessNextPhase()
	// return true

	return false
}
