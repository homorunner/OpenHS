# OpenHS

OpenHS is a simulator core for the game Hearthstone, implemented in Go. This project aims to provide a robust and efficient implementation of Hearthstone's core game mechanics.

## Implemented Features

- **Game Engine Framework**:
  - Core game loop with phase management (BeginFirst, BeginShuffle, BeginDraw, etc.)
  - Turn-based gameplay flow with player switching

- **Card System**: 
  - Card template loading from JSON configurations
  - Basic card types (Minion, Spell, Weapon, Hero, HeroPower)
  - Card effect framework with triggers and conditions
  - Card tags implementation (Poisonous, Rush, Windfury, etc.)

- **Player Management**:
  - Player state tracking including hero, deck, hand, and board
  
- **Game Mechanics**:
  - Initial game setup with card drawing
  - Turn management with phase transitions
  - Card drawing mechanism with fatigue damage
  - Playing cards from hand to the field
  - Full combat system with minion/hero attacks
  - Mana crystal management
  - Death processing and graveyard management

## TODO List

- **Game Mechanics**:
  - Complete mulligan phase implementation
  - Hero powers implementation
  - More tags implementation (See `docs/tags.md`)
  - More advanced card effects and interactions
  - Secrets and auras implementation

- **Card Library**:
  - Implement more cards from the basic and classic sets
  - Add card rarity and classes

- **Game Features**:
  - Game history and replay system
  - Game state serialization/deserialization
  - Network play support
  - AI opponents

## Project Structure

```
openhs/
├── cmd/                  # Main applications
│   └── openhs/           # Main entry point
├── internal/             # Private application and library code
│   ├── game/             # Core game mechanics
│   ├── card/             # Card system implementation
│   ├── engine/           # Game rules engine
│   ├── types/            # Common type definitions
│   ├── config/           # Configuration handling
│   ├── logger/           # Logging utilities
│   ├── bootstrap/        # Application initialization
│   └── util/             # Utility functions
├── cards/                # Card definition files
├── config/               # Configuration files
├── games/                # Game scenario definitions
├── tests/                # Test files and utilities
└── third_party/          # Third-party dependencies
```

## Requirements

- Go 1.24 or later
- Git

## Getting Started

1. Clone the repository:
   ```bash
   git clone https://github.com/openhs/openhs.git
   cd openhs
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Run the project:
   ```bash
   go run cmd/openhs
   ```

## Development

This project follows standard Go project layout and best practices. To contribute:

1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create a new Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details. 