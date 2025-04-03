package engine

import (
	"testing"

	"github.com/openhs/internal/game"
)

func createTestPlayer() *game.Player {
	player := game.NewPlayer()

	// Create a test hero card definition
	heroCard := &game.Card{
		Name:      "Jaina",
		Health:    30,
		MaxHealth: 30,
		Type:      game.Hero,
	}

	// Create hero entity
	hero := game.NewEntity(heroCard, player)
	player.Hero = hero

	// Create test deck
	for i := 0; i < 10; i++ {
		// Create card definition
		card := &game.Card{
			Name: "Test Card",
			Cost: 1,
			Type: game.Minion,
			Attack: 1,
			Health: 1,
			MaxHealth: 1,
		}
		
		// Create entity and add to deck
		entity := game.NewEntity(card, player)
		player.Deck = append(player.Deck, entity)
	}

	return player
}

// createTestGame creates a simple game with two players for testing
func createTestGame() *game.Game {
	g := game.NewGame()

	player1 := createTestPlayer()
	player2 := createTestPlayer()

	g.Players = append(g.Players, player1, player2)
	return g
}

// TestBeginDraw tests if beginDraw correctly draws cards for both players
func TestBeginDraw(t *testing.T) {
	g := createTestGame()
	e := NewEngine(g)

	// Before draw
	if len(g.Players[0].Hand) != 0 {
		t.Errorf("Expected player 1 hand size to be 0, got %d", len(g.Players[0].Hand))
	}
	if len(g.Players[1].Hand) != 0 {
		t.Errorf("Expected player 2 hand size to be 0, got %d", len(g.Players[1].Hand))
	}

	// Execute beginDraw
	err := e.beginDraw()
	if err != nil {
		t.Fatalf("beginDraw returned an error: %v", err)
	}

	// After draw - verify player 1 has 3 cards and player 2 has 4 cards
	if len(g.Players[0].Hand) != 3 {
		t.Errorf("Expected player 1 hand size to be 3, got %d", len(g.Players[0].Hand))
	}
	if len(g.Players[1].Hand) != 4 {
		t.Errorf("Expected player 2 hand size to be 4, got %d", len(g.Players[1].Hand))
	}

	// Verify decks have decreased by the correct amount
	if len(g.Players[0].Deck) != 7 {
		t.Errorf("Expected player 1 deck size to be 7, got %d", len(g.Players[0].Deck))
	}
	if len(g.Players[1].Deck) != 6 {
		t.Errorf("Expected player 2 deck size to be 6, got %d", len(g.Players[1].Deck))
	}

	// Verify next phase is set
	if e.nextPhase != game.MainBegin && e.nextPhase != game.BeginMulligan {
		t.Errorf("Expected next phase to be MainBegin or BeginMulligan, got %v", e.nextPhase)
	}
}

// TestMainDraw tests if mainDraw correctly draws a card for the current player
func TestMainDraw(t *testing.T) {
	g := createTestGame()
	e := NewEngine(g)

	// Setup current player
	g.CurrentPlayerIndex = 0
	g.CurrentPlayer = g.Players[0]

	// Initial state
	initialHandSize := len(g.CurrentPlayer.Hand)
	initialDeckSize := len(g.CurrentPlayer.Deck)

	// Execute mainDraw
	err := e.mainDraw()
	if err != nil {
		t.Fatalf("mainDraw returned an error: %v", err)
	}

	// Verify current player drew 1 card
	if len(g.CurrentPlayer.Hand) != initialHandSize+1 {
		t.Errorf("Expected hand size to increase by 1, got %d (was %d)",
			len(g.CurrentPlayer.Hand), initialHandSize)
	}

	// Verify deck decreased by 1
	if len(g.CurrentPlayer.Deck) != initialDeckSize-1 {
		t.Errorf("Expected deck size to decrease by 1, got %d (was %d)",
			len(g.CurrentPlayer.Deck), initialDeckSize)
	}

	// Verify next phase is set
	if e.nextPhase != game.MainStart {
		t.Errorf("Expected next phase to be MainStart, got %v", e.nextPhase)
	}
}

// TestPlayerSwitching tests if the player switching logic works correctly
func TestPlayerSwitching(t *testing.T) {
	g := createTestGame()
	e := NewEngine(g)

	// Setup initial state
	g.CurrentPlayerIndex = 0
	g.CurrentPlayer = g.Players[0]
	g.CurrentTurn = 1

	// Execute mainNext to switch players
	err := e.mainNext()
	if err != nil {
		t.Fatalf("mainNext returned an error: %v", err)
	}

	// Verify player index and player changed
	if g.CurrentPlayerIndex != 1 {
		t.Errorf("Expected current player index to be 1, got %d", g.CurrentPlayerIndex)
	}
	if g.CurrentPlayer != g.Players[1] {
		t.Errorf("Expected current player to be player 2")
	}

	// Verify turn counter increased
	if g.CurrentTurn != 2 {
		t.Errorf("Expected turn counter to be 2, got %d", g.CurrentTurn)
	}

	// Execute mainNext again to switch back to player 1
	err = e.mainNext()
	if err != nil {
		t.Fatalf("mainNext returned an error: %v", err)
	}

	// Verify player switched back to player 1
	if g.CurrentPlayerIndex != 0 {
		t.Errorf("Expected current player index to be 0, got %d", g.CurrentPlayerIndex)
	}
	if g.CurrentPlayer != g.Players[0] {
		t.Errorf("Expected current player to be player 1")
	}

	// Verify turn counter increased again
	if g.CurrentTurn != 3 {
		t.Errorf("Expected turn counter to be 3, got %d", g.CurrentTurn)
	}
}

// TestPhaseProgression tests that the game phases progress correctly
func TestPhaseProgression(t *testing.T) {
	g := createTestGame()
	e := NewEngine(g)

	// Set initial phase and disable auto-run
	e.nextPhase = game.BeginFirst
	e.SetAutoRun(false)

	// Process BeginFirst phase
	err := e.ProcessNextPhase()
	if err != nil {
		t.Fatalf("ProcessNextPhase returned an error: %v", err)
	}

	// Verify current phase is BeginFirst
	if g.Phase != game.BeginFirst {
		t.Errorf("Expected current phase to be BeginFirst, got %v", g.Phase)
	}

	// Verify next phase is BeginShuffle
	if e.nextPhase != game.BeginShuffle {
		t.Errorf("Expected next phase to be BeginShuffle, got %v", e.nextPhase)
	}

	// Process BeginShuffle phase
	err = e.ProcessNextPhase()
	if err != nil {
		t.Fatalf("ProcessNextPhase returned an error: %v", err)
	}

	// Verify current phase is BeginShuffle
	if g.Phase != game.BeginShuffle {
		t.Errorf("Expected current phase to be BeginShuffle, got %v", g.Phase)
	}

	// Verify next phase is BeginDraw
	if e.nextPhase != game.BeginDraw {
		t.Errorf("Expected next phase to be BeginDraw, got %v", e.nextPhase)
	}

	// Process BeginDraw phase
	err = e.ProcessNextPhase()
	if err != nil {
		t.Fatalf("ProcessNextPhase returned an error: %v", err)
	}

	// Verify current phase is BeginDraw
	if g.Phase != game.BeginDraw {
		t.Errorf("Expected current phase to be BeginDraw, got %v", g.Phase)
	}

	// Verify next phase is either MainBegin or BeginMulligan
	if e.nextPhase != game.MainBegin && e.nextPhase != game.BeginMulligan {
		t.Errorf("Expected next phase to be MainBegin or BeginMulligan, got %v", e.nextPhase)
	}
}

// TestAutoRunAndPhaseProcessing tests the auto-run functionality and ProcessUntil
func TestAutoRunAndPhaseProcessing(t *testing.T) {
	// Test with auto-run enabled (default)
	g := createTestGame()
	e := NewEngine(g)

	// Verify auto-run is enabled by default
	if !e.autoRun {
		t.Error("Expected auto-run to be enabled by default")
	}

	// Start the game - it should automatically progress
	err := e.StartGame()
	if err != nil {
		t.Fatalf("Failed to start game: %v", err)
	}

	// Verify the game has progressed past the initial phases
	if g.Phase == game.InvalidPhase || g.Phase == game.BeginFirst {
		t.Errorf("Expected auto-run to progress past initial phases, but got %v", g.Phase)
	}

	// Test with auto-run disabled
	g = createTestGame()
	e = NewEngine(g)

	// Disable auto-run
	e.SetAutoRun(false)

	// Start game - it should set nextPhase but not progress
	err = e.StartGame()
	if err != nil {
		t.Fatalf("Failed to start game: %v", err)
	}

	// Verify that game is still in InvalidPhase
	if g.Phase != game.InvalidPhase {
		t.Errorf("Expected game to stay at InvalidPhase with auto-run disabled, but got %v", g.Phase)
	}

	// Verify nextPhase is set correctly
	if e.nextPhase != game.BeginFirst {
		t.Errorf("Expected nextPhase to be BeginFirst, but got %v", e.nextPhase)
	}

	// Test manual progression
	err = e.ProcessNextPhase()
	if err != nil {
		t.Fatalf("Failed to process next phase: %v", err)
	}

	// Verify that game progressed to BeginFirst
	if g.Phase != game.BeginFirst {
		t.Errorf("Expected game to progress to BeginFirst, but got %v", g.Phase)
	}

	// Test ProcessUntil
	err = e.ProcessUntil(game.MainAction)
	if err != nil {
		t.Fatalf("Failed to process until MainAction: %v", err)
	}

	// Verify phase is now MainAction
	if g.Phase != game.MainAction {
		t.Errorf("Expected phase to be MainAction after ProcessUntil, but got %v", g.Phase)
	}
}

// TestEndPlayerTurn tests the end turn functionality
func TestEndPlayerTurn(t *testing.T) {
	g := createTestGame()
	e := NewEngine(g)

	// Process until first turn
	err := e.StartGame()
	if err != nil {
		t.Fatalf("Failed to start game: %v", err)
	}
	if g.Phase != game.MainAction {
		t.Errorf("Expected phase to be MainAction after StartGame, but got %v", g.Phase)
	}

	// End turn
	err = e.EndPlayerTurn()
	if err != nil {
		t.Fatalf("EndPlayerTurn returned an error: %v", err)
	}

	// Verify player has switched
	if g.CurrentPlayerIndex != 1 {
		t.Errorf("Expected current player index to be 1 after end turn, got %d", g.CurrentPlayerIndex)
	}

	if g.CurrentPlayer != g.Players[1] {
		t.Errorf("Expected current player to be player 2 after end turn")
	}

	// Verify turn counter increased
	if g.CurrentTurn != 2 {
		t.Errorf("Expected turn to be 2 after end turn, got %d", g.CurrentTurn)
	}

	// Test error case - ending turn when not in action phase
	g.Phase = game.MainDraw
	err = e.EndPlayerTurn()
	if err == nil {
		t.Error("Expected an error when ending turn outside of action phase, but got nil")
	}
}

// TestManaHandling tests the mana handling functionality
func TestManaHandling(t *testing.T) {
	g := createTestGame()
	e := NewEngine(g)

	// Setup current player
	g.CurrentPlayerIndex = 0
	g.CurrentPlayer = g.Players[0]

	// Verify initial mana values (from config.DefaultStartingMana)
	if g.CurrentPlayer.Mana != 0 {
		t.Errorf("Expected initial Mana to be 0, got %d", g.CurrentPlayer.Mana)
	}
	if g.CurrentPlayer.TotalMana != 0 {
		t.Errorf("Expected initial TotalMana to be 0, got %d", g.CurrentPlayer.TotalMana)
	}
	if g.CurrentPlayer.MaxMana != 10 {
		t.Errorf("Expected MaxMana to be 10, got %d", g.CurrentPlayer.MaxMana)
	}

	// Test mana increase in first turn
	err := e.mainResource()
	if err != nil {
		t.Fatalf("mainResource returned an error: %v", err)
	}

	// Verify mana increased after first turn
	if g.CurrentPlayer.TotalMana != 1 {
		t.Errorf("Expected TotalMana to be 1 after first turn, got %d", g.CurrentPlayer.TotalMana)
	}
	if g.CurrentPlayer.Mana != 1 {
		t.Errorf("Expected Mana to be 1 after first turn, got %d", g.CurrentPlayer.Mana)
	}

	// Simulate spending mana
	g.CurrentPlayer.Mana = 0

	// Test mana restoration in second turn
	err = e.mainResource()
	if err != nil {
		t.Fatalf("mainResource returned an error: %v", err)
	}

	// Verify TotalMana increased and Mana restored
	if g.CurrentPlayer.TotalMana != 2 {
		t.Errorf("Expected TotalMana to be 2 after second turn, got %d", g.CurrentPlayer.TotalMana)
	}
	if g.CurrentPlayer.Mana != 2 {
		t.Errorf("Expected Mana to be restored to 2, got %d", g.CurrentPlayer.Mana)
	}

	// Test mana increase up to max
	// Set TotalMana to 9 (one below max)
	g.CurrentPlayer.TotalMana = 9
	g.CurrentPlayer.Mana = 5 // Simulate some spent mana

	err = e.mainResource()
	if err != nil {
		t.Fatalf("mainResource returned an error: %v", err)
	}

	// Verify TotalMana increased to max and Mana restored
	if g.CurrentPlayer.TotalMana != 10 {
		t.Errorf("Expected TotalMana to reach 10, got %d", g.CurrentPlayer.TotalMana)
	}
	if g.CurrentPlayer.Mana != 10 {
		t.Errorf("Expected Mana to be restored to 10, got %d", g.CurrentPlayer.Mana)
	}

	// Test mana doesn't exceed max
	err = e.mainResource()
	if err != nil {
		t.Fatalf("mainResource returned an error: %v", err)
	}

	// Verify TotalMana doesn't exceed max
	if g.CurrentPlayer.TotalMana != 10 {
		t.Errorf("Expected TotalMana to remain at 10, got %d", g.CurrentPlayer.TotalMana)
	}
	if g.CurrentPlayer.Mana != 10 {
		t.Errorf("Expected Mana to remain at 10, got %d", g.CurrentPlayer.Mana)
	}
}

// TestManaHandlingAcrossPlayerTurns tests the mana system works correctly across player turns
func TestManaHandlingAcrossPlayerTurns(t *testing.T) {
	g := createTestGame()
	e := NewEngine(g)

	// Setup game state
	g.CurrentPlayerIndex = 0
	g.CurrentPlayer = g.Players[0]
	g.CurrentTurn = 1

	// First player's turn - increase mana to 1
	err := e.mainResource()
	if err != nil {
		t.Fatalf("mainResource returned an error: %v", err)
	}

	// Verify first player's mana
	if g.Players[0].TotalMana != 1 {
		t.Errorf("Expected player 1 TotalMana to be 1, got %d", g.Players[0].TotalMana)
	}
	if g.Players[0].Mana != 1 {
		t.Errorf("Expected player 1 Mana to be 1, got %d", g.Players[0].Mana)
	}

	// Verify second player's mana is still at default
	if g.Players[1].TotalMana != 0 {
		t.Errorf("Expected player 2 TotalMana to be 0, got %d", g.Players[1].TotalMana)
	}
	if g.Players[1].Mana != 0 {
		t.Errorf("Expected player 2 Mana to be 0, got %d", g.Players[1].Mana)
	}

	// Switch to second player
	err = e.mainNext()
	if err != nil {
		t.Fatalf("mainNext returned an error: %v", err)
	}

	// Verify player switched
	if g.CurrentPlayer != g.Players[1] {
		t.Errorf("Expected current player to be player 2")
	}

	// Second player's turn - increase mana to 1
	err = e.mainResource()
	if err != nil {
		t.Fatalf("mainResource returned an error: %v", err)
	}

	// Verify second player's mana
	if g.Players[1].TotalMana != 1 {
		t.Errorf("Expected player 2 TotalMana to be 1, got %d", g.Players[1].TotalMana)
	}
	if g.Players[1].Mana != 1 {
		t.Errorf("Expected player 2 Mana to be 1, got %d", g.Players[1].Mana)
	}

	// First player should still have their previous values
	if g.Players[0].TotalMana != 1 {
		t.Errorf("Expected player 1 TotalMana to remain 1, got %d", g.Players[0].TotalMana)
	}

	// Simulate player 1 spending mana
	g.Players[0].Mana = 0

	// Switch back to first player
	err = e.mainNext()
	if err != nil {
		t.Fatalf("mainNext returned an error: %v", err)
	}

	// Verify player switched back
	if g.CurrentPlayer != g.Players[0] {
		t.Errorf("Expected current player to be player 1")
	}

	// First player's second turn - increase mana to 2 and restore spent mana
	err = e.mainResource()
	if err != nil {
		t.Fatalf("mainResource returned an error: %v", err)
	}

	// Verify first player's mana increased and restored
	if g.Players[0].TotalMana != 2 {
		t.Errorf("Expected player 1 TotalMana to be 2, got %d", g.Players[0].TotalMana)
	}
	if g.Players[0].Mana != 2 {
		t.Errorf("Expected player 1 Mana to be restored to 2, got %d", g.Players[0].Mana)
	}

	// Second player's mana should be unchanged
	if g.Players[1].TotalMana != 1 {
		t.Errorf("Expected player 2 TotalMana to remain 1, got %d", g.Players[1].TotalMana)
	}
}
