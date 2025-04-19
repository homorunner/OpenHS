# OpenHS Frontend

This is a web-based frontend for the OpenHS Hearthstone simulator. It provides a text-based visual interface to interact with the game engine without using copyrighted images or assets.

## Features

- Display game state including heroes, minions, weapons, and cards
- Play cards from hand to the board
- Attack with minions
- End turn functionality
- Game log to track actions
- CSS-based card designs with simple icons

## Running the Frontend

1. Make sure you have Go 1.24 or later installed

2. From the root directory of the OpenHS project, run:
   ```bash
   go run cmd/server/main.go
   ```

3. Open your web browser and navigate to:
   ```
   http://localhost:8080
   ```

## How to Play

1. **View Game State**: The game board shows your hand at the bottom and your opponent's hand (face down) at the top. The field areas in the middle display minions in play.

2. **Play a Card**: Click on a card in your hand to select it. For minion cards, you'll need to select a position on the board. For other cards, they will be played immediately.

3. **Attack**: Click on one of your minions that can attack (they will be highlighted), then click on an opponent's minion to attack it.

4. **End Turn**: Click the "End Turn" button to end your turn.

## Technical Details

The frontend consists of:

- A Go server (in cmd/frontend/main.go) that connects to the OpenHS game engine
- HTML/CSS for the game board layout
- JavaScript for handling game state and user interactions

The server provides two main API endpoints:
- `/api/game` - GET request to retrieve the current game state
- `/api/action` - POST request to perform game actions (play card, attack, end turn)

## Visual Design

The frontend uses a clean, CSS-based design for cards and game elements:
- Card borders and styles are implemented with proper CSS
- Simple Unicode icons represent different card types
- Color-coded elements for different types of cards and abilities

## Customization

You can customize the look and feel by modifying:
- `frontend/static/styles.css` - Styling for the game board and elements
- `frontend/static/script.js` - Client-side game logic and card icons 