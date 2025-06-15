# Mexe-mexe 🃁 🂡 🃋 - Multiplayer Card Game 


A real-time multiplayer card game built in Go with WebSocket connections. Players compete to be the first to empty their hand by playing melds (sequences or books) and manipulating cards on the table through the unique "mexe-mexe" mechanic.

## Game Rules

### Objective
- **Win condition**: Be the first player to have 0 cards in your hand
- **Alternative win condition**: If the deck runs out, the player with fewer points/cards wins *(not yet implemented)*

### Gameplay
- **Turn-based**: Players take turns in sequence
- **Meld playing**: On your turn, you can play a meld (sequence 🃊 🃋 🃍 or book 🂱 🃑 🃁) from your hand to the table
- **Drawing**: If you cannot play a meld, you must draw a card from the deck
- **Turn ending**: If no meld is available and no mexe-mexe moves are possible, your turn ends

### Mexe-mexe Mechanic
The unique feature of this game! Once any meld is played on the table (by any player), "mexe-mexe" becomes available.

**Mexe-mexe allows you to:**
- Manipulate cards already on the table as if they were your own
- Rearrange existing melds to create new valid combinations
- Play cards from your hand by incorporating them with table cards
- **Requirement**: All resulting melds on the table must remain valid (proper sequences or books)

## Controls

| Key | Action |
|-----|--------|
| `<- or ->` | Navigate cards |
| `s` | Select/deselect cards |
| `p` | Play selected meld |
| `d` | Draw a card |
| `e` | End turn |
| `q` | Quit game |

## Architecture

### Server Layer
**Responsibilities:**
- WebSocket connection management
- User authentication
- Game room creation and matchmaking
- Player pairing (2 players per room)

### Game Engine
**Responsibilities:**
- Core game logic and rule validation
- Game state management (player hands, table cards, deck)
- Turn flow control
- Meld validation (sequences and books)
- Communication through input/output providers

### Client Layer
**Responsibilities:**
- WebSocket connection to server
- Terminal-based user interface
- User input handling
- Game state display

### Key Design Patterns

**Interface-Driven Design**: The game engine works with abstractions (`InputProvider`, `OutputProvider`) rather than concrete WebSocket implementations, making it easily testable and adaptable.

**Provider Pattern**: 
- `WebsocketInputProvider`: Handles player input via WebSocket
- `WebsocketOutputProvider`: Sends game updates to clients

## Building and Running

### Prerequisites
- Go 1.19 or higher
- Terminal environment (Linux/macOS recommended - Windows not tested)

### Server
```bash
go build cmd/server/main.go
./main
```

### Client
```bash
go build cmd/client/main.go
./main
```

**Note**: The client uses terminal-based UI and has not been tested on Windows.

## Architecture Flow

### 1. Connection & Lobby Phase
```
Client → WebSocket → Server.HandleConnections()
├── Read JoinServerMessage 
├── Authenticate user
├── Send WelcomeMessage
├── Read StartGameMessage  
├── Create/join GameRoom
└── When room full (2 players) → Start game
```

### 2. Game Handoff
```
Server creates:
├── engine.Game (with shuffled player order)
├── WebsocketInputProvider for each client  
├── WebsocketOutputProvider for each client
└── Calls room.StartGame() → game.Start() in goroutine
```

### 3. Game Phase
```
Game Engine Loop:
├── For each player's turn:
│   ├── Call inputProvider.GetPlay() → blocks on WebSocket
│   ├── Validate play with IsValid()
│   ├── Execute play with MakePlay()
│   └── Update game state
└── Continue until win/quit condition
```

## Project Structure

```
├── cmd/
│   ├── server/         # Server entry point
│   └── client/         # Client entry point
├── internal/
│   ├── engine/         # Game logic and rules
│   ├── server/         # WebSocket server implementation
│   └── client/         # Client implementation
└── └── service/        # Logger

```


## Highlights

- **Concurrent Design**: Each game runs in its own goroutine
- **Real-time Communication**: WebSocket connections provide instant updates
- **State Management**: Centralized game state with atomic operations
- **Input Validation**: Robust meld and game rule validation
- **Clean Interfaces**: Abstract providers enable flexible I/O handling
