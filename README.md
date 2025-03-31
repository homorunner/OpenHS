# OpenHS

OpenHS is a simulator core for the game Hearthstone, implemented in Go. This project aims to provide a robust and efficient implementation of Hearthstone's core game mechanics.

## Features (Planned)

- Card system implementation
- Game state management
- Combat mechanics
- Spell system
- Hero powers
- Card effects and interactions
- Game rules engine

## Project Structure

```
openhs/
├── cmd/                    # Main applications
│   └── openhs/            # Main entry point
├── internal/              # Private application and library code
│   ├── game/             # Core game mechanics
│   ├── card/             # Card system implementation
│   └── engine/           # Game rules engine
├── pkg/                   # Library code that's ok to use by external applications
└── test/                 # Additional external test applications and test data
```

## Requirements

- Go 1.21 or later
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