# Kortex

**The Autonomous Interface Layer.**

A local-first, cross-platform agentic protocol built on Google Gemini 3 Pro and Go-ADK. Kortex decouples intent from action, using MCP to autonomously navigate, explain, and execute complex web tasks for any user.

## Architecture

Kortex follows the **Hexagonal Architecture (Ports & Adapters)** pattern to ensure modularity and testability.

- **Core**: Contains the Domain logic and Ports (interfaces).
- **Infrastructure**: Contains the Adapters (implementations) for the Ports.

## Mission 1: The "Hexagonal" Core & Memory

We have successfully initialized the Kortex brain by building the Core Domain and Infrastructure layers.

### 1. Domain Layer (`internal/core/domain`)
Defined the core data structures:
- **Session**: Represents a user session (`ID`, `Context`, `CreatedAt`).
- **Message**: Represents a message in a session (`ID`, `SessionID`, `Role`, `Content`, `ToolCallID`, `ToolResult`).
- **MemoryFragment**: Represents a piece of memory with vector embeddings (`ID`, `Content`, `Embedding` [768 dim], `Tags`).

### 2. Ports Layer (`internal/core/ports`)
Defined the interfaces for external interactions:
- **AIProvider**: `GenerateStream(ctx, history, tools)`
- **VectorStore**: `Save(ctx, fragment)`, `Search(ctx, queryVector, limit)`
- **Browser**: `Maps(url)`, `GetSnapshot()`, `Highlight(selector)`

### 3. Infrastructure Layer (`internal/infra`)
Implemented the adapters:
- **SQLite VectorStore** (`internal/infra/sqlite`):
    - Uses `glebarez/sqlite` (pure Go driver) for compatibility.
    - Implements `VectorStore` interface.
    - **Note**: `Search` uses Raw SQL for `vec_distance_cosine` to support the `sqlite-vec` extension.
- **Flight Recorder** (`internal/infra/logger`):
    - Uses `log/slog` for structured logging.
    - Writes logs to `kortex_flight_recorder.jsonl`.

## Setup & Usage

### Prerequisites
- Go 1.23+

### Installation
1.  Clone the repository.
2.  Run setup:
    ```bash
    make setup
    ```

### Running Tests
```bash
make test
```

> [!NOTE]
> **sqlite-vec Support**: The `VectorStore` implementation is designed to work with the `sqlite-vec` extension. However, the current environment uses a pure Go SQLite driver which does not include this extension by default. Vector search queries will generate the correct SQL but may fail if the extension is not loaded in the SQLite environment.
