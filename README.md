# Kortex: The Autonomous Interface Layer

> **Copyright © 2025 Pundarikaksh N Tripathi**  
> This project is licensed under the CC BY-NC-ND 4.0 license. See the [LICENSE](LICENSE) file for details.

![Go Version](https://img.shields.io/badge/Go-1.23+-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![Wails](https://img.shields.io/badge/Wails-v2-00D9FF?style=for-the-badge&logo=go&logoColor=white)
![React](https://img.shields.io/badge/React-18+-61DAFB?style=for-the-badge&logo=react&logoColor=black)
![Playwright](https://img.shields.io/badge/Playwright-Go-45ba4b?style=for-the-badge&logo=playwright&logoColor=white)
![SQLite](https://img.shields.io/badge/SQLite-003B57?style=for-the-badge&logo=sqlite&logoColor=white)
![Gemini](https://img.shields.io/badge/Gemini-3_Pro-8E75B2?style=for-the-badge&logo=google-gemini&logoColor=white)
![License](https://img.shields.io/badge/License-CC_BY--NC--ND_4.0-lightgrey?style=for-the-badge)

## 📖 Introduction

**Kortex** is a local-first, cross-platform agentic protocol designed to decouple intent from action. Built on **Google Gemini 3 Pro** and the **Google Go Agent Development Kit (ADK)**, Kortex acts as an intelligent layer between users and the web, capable of autonomously navigating, understanding, and executing complex tasks.

Imagine a browser that doesn't just display pages, but *understands* them. Kortex uses advanced visual injection and accessibility tree analysis to interact with web content just like a human would—but at machine speed.

---

## 📑 Table of Contents

- [📖 Introduction](#-introduction)
- [💡 Solution & How It Works](#-solution--how-it-works)
- [🏗️ Architecture](#-architecture)
- [🛠️ Tech Stack & Engineering Decisions](#-tech-stack--engineering-decisions)
- [📂 Directory Structure](#-directory-structure)
- [🚀 Quick Start](#-quick-start)
- [🎨 Desktop App Features](#-desktop-app-features)
- [🐳 Docker Deployment (Web Version)](#-docker-deployment-web-version)
- [📚 Core Components](#-core-components)
- [❓ Troubleshooting](#-troubleshooting)
- [🤝 Contributing](#-contributing)
- [📄 License](#-license)

---

## 💡 Solution & How It Works

The modern web is complex. Automating it requires more than just scripts; it requires **agency**. Kortex provides a robust solution that bridges the gap between human intent and machine execution.

### How It Works

1.  **Intent Understanding**: You type a goal (e.g., "Find the cheapest flight to Tokyo"). Kortex uses **Gemini 3 Pro** to understand the semantics and break this down into logical steps.
2.  **Reasoning Loop (ReAct)**: Kortex enters a "Reason-Act" loop. It looks at the current state of the browser, decides what to do next (e.g., "Type 'Tokyo' into the search box"), and executes that action.
3.  **Visual Perception**: Unlike traditional scrapers that look at raw HTML, Kortex builds a simplified **Accessibility Tree**. This allows it to "see" the page structure (buttons, inputs, links) just like a screen reader or a human would, making it highly resilient to layout changes.
4.  **Action Execution**: Using **Playwright**, Kortex physically interacts with the page—clicking, typing, and scrolling. It even highlights elements visually so you can follow along.
5.  **Memory & Context**: Every interaction is stored in a local **SQLite Vector Database**. This allows Kortex to remember past searches and preferences, providing a personalized experience without sending data to the cloud.

---

## 🏗️ Architecture

Kortex follows the **Hexagonal Architecture (Ports & Adapters)** pattern. This ensures that our core logic (the "Brain") is isolated from external tools (the "Hands" and "Eyes"), making the system modular, testable, and easily extensible.

```mermaid
graph TD
    %% Styling
    classDef core fill:#ff9a9e,stroke:#333,stroke-width:2px,color:black;
    classDef infra fill:#a18cd1,stroke:#333,stroke-width:2px,color:black;
    classDef app fill:#fad0c4,stroke:#333,stroke-width:2px,color:black;
    classDef external fill:#84fab0,stroke:#333,stroke-width:2px,color:black;

    subgraph Application_Layer [Application Layer]
        direction TB
        Desktop[Desktop App\nWails v2 + React]:::app
        WebServer[Web Server\nFiber + WebSocket]:::app
    end

    subgraph Core_Domain [Core Domain (The Brain)]
        direction TB
        Logic[Business Logic\nReAct Loop]:::core
        Ports[Ports Interfaces\n(Abstractions)]:::core
    end

    subgraph Infrastructure_Layer [Infrastructure (Adapters)]
        direction TB
        AgentAdapter[Agent Adapter\nGoogle ADK + Gemini]:::infra
        BrowserAdapter[Browser Adapter\nPlaywright]:::infra
        VectorDB[Vector Store\nSQLite + Embeddings]:::infra
        FlightRecorder[Flight Recorder\nStructured Logger]:::infra
    end

    subgraph External_World [External World]
        User((User)):::external
        Web((Internet)):::external
        LLM((Gemini API)):::external
    end

    %% Connections
    User <-->|UI Events| Desktop
    User <-->|WebSocket| WebServer
    
    Desktop -->|Calls| AgentAdapter
    WebServer -->|Calls| AgentAdapter

    AgentAdapter -->|Implements| Ports
    BrowserAdapter -->|Implements| Ports
    VectorDB -->|Implements| Ports
    FlightRecorder -->|Implements| Ports

    Logic -->|Uses| Ports
    
    AgentAdapter <-->|API| LLM
    BrowserAdapter <-->|Navigates| Web
```

---

## 🛠️ Tech Stack & Engineering Decisions

We chose this stack to balance **performance**, **developer experience**, and **capabilities**.

| Component | Technology | Reasoning |
|-----------|------------|-----------|
| **Language** | **Go (Golang)** | Go offers exceptional concurrency (goroutines) which is crucial for handling multiple browser contexts and agent threads simultaneously. It's strictly typed, compiles to a single binary, and has a massive ecosystem. |
| **Desktop** | **Wails v2** | Unlike Electron which bundles a whole Chrome browser (heavy RAM usage), Wails uses the OS's native webview (WebView2 on Windows, WebKit on macOS). This results in a tiny binary (~15MB vs ~100MB) and low memory footprint. |
| **Frontend** | **React + Vite** | React provides a component-based UI that is easy to manage. Vite offers lightning-fast HMR (Hot Module Replacement) for a great dev experience. |
| **AI Model** | **Gemini 3 Pro** | Gemini's large context window and multimodal capabilities make it ideal for understanding complex web pages and reasoning about navigation steps. |
| **Agent** | **Google ADK** | The Agent Development Kit provides a standardized way to build agents, managing tool calls and state transitions reliability. |
| **Browser** | **Playwright** | Playwright is faster and more reliable than Selenium. It supports modern web features, handles dynamic content (SPAs) effortlessly, and allows for easy headless execution. |
| **Database** | **SQLite + Vec** | We use SQLite for a local-first approach. The `sqlite-vec` extension allows us to perform vector similarity search directly on the user's machine, keeping data private and fast. |

---

## 📂 Directory Structure

Here is a detailed breakdown of the project structure to help you navigate the codebase.

```text
Kortex/
├── app.go                  # Wails Bridge: Connects Go backend to React frontend.
├── main.go                 # Desktop Entry Point: Initializes Wails application.
├── cmd/
│   └── web/
│       └── main.go         # Web Server Entry Point: Runs Kortex as a Docker/Web service.
├── frontend/               # User Interface (React)
│   ├── src/
│   │   ├── App.tsx         # Main Layout: Dual-pane interface (Chat + Terminal).
│   │   ├── components/     # UI Components (FlightRecorder, ChatBox, etc.).
│   │   └── wailsjs/        # Generated Bindings: Go functions callable from JS.
├── internal/               # Private Application Code
│   ├── core/
│   │   ├── domain/         # Data Models: Defines Session, Message, Memory structs.
│   │   └── ports/          # Interfaces: Defines contracts for Adapters (Hexagonal Arch).
│   ├── adapters/
│   │   └── agent/          # The Brain: Implements the ReAct loop and Gemini integration.
│   └── infra/
│       ├── browser/        # The Hands: Playwright implementation for browser control.
│       ├── logger/         # The Black Box: Structured logging for the Flight Recorder.
│       └── sqlite/         # The Memory: Vector database implementation.
├── build/                  # Build Artifacts: Icons, manifests, and compiled binaries.
├── Dockerfile              # Container Config: For running Kortex in Docker.
├── .env.example            # Config Template: API keys and environment settings.
└── go.mod                  # Dependencies: Go module definitions.
```

---

## 🚀 Quick Start

### Prerequisites

1. **[Go 1.23+](https://go.dev/dl/)** installed
2. **[Node.js 18+](https://nodejs.org/)** and npm installed
3. **[Wails CLI v2](https://wails.io/docs/gettingstarted/installation)** installed:
   ```bash
   go install github.com/wailsapp/wails/v2/cmd/wails@latest
   ```
4. **Playwright browsers** installed:
   ```bash
   go run github.com/playwright-community/playwright-go/cmd/playwright@latest install --with-deps
   ```
5. **Google Gemini API Key** from [Google AI Studio](https://aistudio.google.com/app/apikey)

### Installation

1. **Clone the repository**:
   ```bash
   git clone https://github.com/PundarikakshNTripathi/Kortex.git
   cd Kortex
   ```

2. **Install Go dependencies**:
   ```bash
   go mod download
   ```

3. **Install frontend dependencies**:
   ```bash
   cd frontend
   npm install
   cd ..
   ```

4. **Set up environment variables**:
   ```bash
   cp .env.example .env
   ```
   
   Edit `.env` and add your API key:
   ```env
   GOOGLE_API_KEY=your-actual-api-key-here
   ```

### Running the Desktop App

**Development mode** (with hot reload):
```bash
wails dev
```

**Production build**:
```bash
wails build
.\build\bin\Kortex.exe  # Windows
./build/bin/Kortex      # macOS/Linux
```

### Running Tests

To verify the core functionality:
```bash
go test -v ./...
```

You should see a Chromium window pop up briefly as the tests run!

---

## 🎨 Desktop App Features

### Cyberpunk UI Theme

- **Dark Mode**: Deep space blue background (#0a0e27)
- **Neon Accents**: Electric cyan (#00d9ff) and purple (#a855f7)
- **Glassmorphism**: Semi-transparent panels with blur effects
- **Smooth Animations**: Fade-in, slide-in, and pulse effects
- **Matrix Terminal**: Scanning line effect in Mission Control

### Dual-Pane Layout

```
┌─────────────────────────────────────────────────────────┐
│  ⚡ KORTEX - Autonomous Interface Layer                 │
├──────────────────────────┬──────────────────────────────┤
│                          │  ⚙ MISSION CONTROL           │
│  Chat Interface (60%)    │                              │
│                          │  Terminal View (40%)         │
│  - User messages         │  - Real-time logs            │
│  - Agent responses       │  - Color-coded levels        │
│  - Input field           │  - Auto-scroll               │
│  - Loading states        │  - Timestamps                │
│                          │                              │
└──────────────────────────┴──────────────────────────────┘
```

### Mission Control Log Colors

| Log Level | Color | Icon | Description |
|-----------|-------|------|-------------|
| NAVIGATE | Cyan | 🔵 | Browser navigation |
| CLICK | Green | 🟢 | Element clicks |
| TYPE | Yellow | 🟡 | Text input |
| HIGHLIGHT | Magenta | 🟣 | Visual feedback |
| GET_SNAPSHOT | Blue | 🔷 | Page analysis |
| PLANNING | Purple | 🟪 | Agent reasoning |
| INIT | Cyan | ⚙️ | Initialization |
| USER | White | 📝 | User input |
| ERROR | Red | ❌ | Errors |
| COMPLETE | Green | ✅ | Task completion |

---

## 🐳 Docker Deployment (Web Version)

Kortex can also run as a **web server** in a Docker container, making it accessible via WebSocket without any desktop setup.

### Prerequisites

1. **[Docker](https://www.docker.com/get-started)** installed
2. **Google Gemini API Key** from [Google AI Studio](https://aistudio.google.com/app/apikey)

### Building the Docker Image

```bash
docker build -t kortex-web .
```

### Running the Container

**Basic usage:**
```bash
docker run -p 8080:8080 \
  -e GOOGLE_API_KEY=your-api-key-here \
  -e HEADLESS=true \
  kortex-web
```

**With persistent database:**
```bash
docker run -p 8080:8080 \
  -e GOOGLE_API_KEY=your-api-key-here \
  -e HEADLESS=true \
  -v $(pwd)/data:/app/data \
  kortex-web
```

### Connecting to the WebSocket

**Endpoint:** `ws://localhost:8080/ws/chat`

**Message format:**
```json
{
  "goal": "Navigate to google.com and search for AI news"
}
```

---

## 📚 Core Components

### The "Brain": Agent Adapter (`internal/adapters/agent`)

The Agent Adapter is the intelligence center of Kortex. It connects the reasoning capabilities of **Gemini** (via ADK) with the physical capabilities of the **Browser Adapter**.

*   **ReAct Loop**: Implements a Reason-Act loop where the agent observes the state, thinks about the next step, and executes a tool.
*   **Google ADK**: Uses `google.golang.org/adk` to manage the agent's lifecycle, session state, and tool execution.
*   **Tools**:
    *   `Navigate(url)`: Go to a website.
    *   `Click(selector)`: Interact with elements.
    *   `Type(selector, text)`: Input data.
    *   `Highlight(selector, message)`: Visually communicate intent to the user.
    *   `GetSnapshot()`: Read the page's accessibility tree.
*   **Flight Recorder**: Logs every tool execution for debugging and replay.

### The "Hands": Browser Adapter (`internal/infra/browser`)

The Browser Adapter is Kortex's primary way of interacting with the world. Implemented using **Playwright**, it allows Kortex to:

1.  **Navigate**: Open and traverse web pages.
2.  **See**: Generate a simplified **Accessibility Snapshot** of the page, converting raw HTML into a structured JSON tree of roles and names that the AI can understand.
3.  **Touch**: Highlight elements on the page to show the user what it's looking at or doing.

#### Key Features

*   **Visual Injection**: Kortex injects custom JavaScript to draw a "Cyan/Electric Blue" outline around elements and display floating tooltips.
    ```go
    // Example: Highlight the search bar with a message
    browser.Highlight("#search-input", "I am typing here...")
    ```
*   **Resilience**: Built-in timeouts and error handling ensure Kortex doesn't crash if a selector isn't found immediately.
*   **Visible Mode**: Runs in `Headless: false` mode by default, so you can watch Kortex work.

### The "Memory": Core & Vector Store (`internal/core`)

*   **Domain Models**: `Session`, `Message`, and `MemoryFragment` define how Kortex thinks and remembers.
*   **Vector Memory**: Uses cosine similarity search to retrieve relevant context from past interactions, giving Kortex long-term memory.

---

## ❓ Troubleshooting

**Q: The desktop app doesn't launch.**
*   **A**: Ensure you have Wails CLI installed: `go install github.com/wailsapp/wails/v2/cmd/wails@latest`
*   Run `wails doctor` to check your environment

**Q: The browser doesn't open.**
*   **A**: Ensure you've run the Playwright install command: `go run github.com/playwright-community/playwright-go/cmd/playwright@latest install --with-deps`

**Q: "GOOGLE_API_KEY not set" error.**
*   **A**: Create a `.env` file in the root directory with your API key:
    ```env
    GOOGLE_API_KEY=your-api-key-here
    ```
*   Get your key from [Google AI Studio](https://aistudio.google.com/app/apikey)

**Q: Frontend build errors.**
*   **A**: Delete `frontend/node_modules` and run `npm install` again
*   Ensure you're using Node.js 18+

**Q: "Vector search not supported" error.**
*   **A**: The current implementation uses a pure Go SQLite driver. For full vector search capabilities, ensure the `sqlite-vec` extension is properly loaded in your environment (or use the provided mock/fallback for basic testing).

---

## 🤝 Contributing

This project follows the Hexagonal Architecture pattern. When adding features:

1. Define interfaces in `internal/core/ports/`
2. Implement adapters in `internal/adapters/` or `internal/infra/`
3. Update the Wails bridge in `app.go` if needed
4. Add UI components in `frontend/src/components/`

---

## 📄 License

This project is distributed under the **Creative Commons Attribution-NonCommercial-NoDerivatives 4.0 International** license.

**What this means:**
- ✅ You can view and use this code for learning
- ✅ You can share this project with attribution
- ❌ You cannot use this commercially
- ❌ You cannot create modified versions

See [LICENSE](LICENSE) for the full legal text.

---

**Built with ❤️ using Wails v2, React, Go, and Google Gemini**
