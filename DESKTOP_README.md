# Kortex Desktop App - Setup Guide

## Prerequisites

Before running the Kortex desktop app, ensure you have:

1. **Go 1.23+** installed
2. **Node.js 18+** and npm installed
3. **Wails CLI v2** installed
4. **Playwright browsers** installed
5. **Google Gemini API Key**

## Quick Start

### 1. Set Up Environment Variables

Create a `.env` file in the root directory:

```bash
cp .env.example .env
```

Edit `.env` and add your Google API key:

```
GOOGLE_API_KEY=your-actual-api-key-here
```

**How to get your API key:**
1. Visit [Google AI Studio](https://aistudio.google.com/app/apikey)
2. Click "Get API Key" or "Create API Key"
3. Copy the key (starts with `AIza...`)

### 2. Install Dependencies

```powershell
# Install Go dependencies
go mod download

# Install frontend dependencies
cd frontend
npm install
cd ..

# Install Playwright browsers (if not already done)
go run github.com/playwright-community/playwright-go/cmd/playwright@latest install --with-deps
```

### 3. Run the App

```powershell
# Development mode (with hot reload)
wails dev

# Or build and run production version
wails build
.\build\bin\Kortex.exe
```

## Features

### 🎨 Cyberpunk UI
- **Dark Mode**: Sleek dark theme with cyan/purple accents
- **Dual-Pane Layout**: 
  - Left: Chat interface for commands
  - Right: Mission Control terminal for real-time logs

### 🤖 AI Agent
- Powered by Google Gemini 3 Pro
- Autonomous web navigation
- Visual element highlighting
- Real-time execution feedback

### 📊 Mission Control
- Live log streaming
- Color-coded log levels:
  - 🔵 NAVIGATE - Cyan
  - 🟢 CLICK - Green
  - 🟡 TYPE - Yellow
  - 🟣 HIGHLIGHT - Magenta
  - 🔷 GET_SNAPSHOT - Blue
  - 🟪 PLANNING - Purple
  - ✅ COMPLETE - Green
  - ❌ ERROR - Red

## Usage Examples

Try these commands in the chat interface:

```
Navigate to google.com
```

```
Search for "AI news"
```

```
Click the first result
```

## Troubleshooting

### Browser doesn't open
- Ensure Playwright browsers are installed:
  ```powershell
  go run github.com/playwright-community/playwright-go/cmd/playwright@latest install --with-deps
  ```

### API Key error
- Check that `GOOGLE_API_KEY` is set in your `.env` file
- Verify the key is valid at [Google AI Studio](https://aistudio.google.com/app/apikey)

### Build errors
- Run `go mod tidy` to clean up dependencies
- Delete `frontend/node_modules` and run `npm install` again
- Ensure you're using Go 1.23+ and Node.js 18+

## Development

### Project Structure

```
Kortex/
├── app.go              # Wails App struct (Go bridge)
├── main.go             # Wails entry point
├── frontend/
│   └── src/
│       ├── App.tsx     # Main React component
│       ├── App.css     # Main styles
│       ├── components/
│       │   ├── FlightRecorder.tsx  # Terminal component
│       │   └── FlightRecorder.css  # Terminal styles
│       └── style.css   # Global styles
├── internal/
│   ├── adapters/agent/ # AI agent logic
│   ├── infra/browser/  # Playwright integration
│   └── infra/sqlite/   # Vector store
└── .env                # Environment variables (create this!)
```

### Hot Reload

Run `wails dev` for automatic reloading during development. Changes to Go code or React components will trigger a rebuild.

## Building for Production

```powershell
# Build for Windows
wails build

# The executable will be in: build/bin/Kortex.exe
```

## Contributing

This project follows the Hexagonal Architecture pattern. When adding features:

1. Define interfaces in `internal/core/ports/`
2. Implement adapters in `internal/adapters/` or `internal/infra/`
3. Update the Wails bridge in `app.go` if needed
4. Add UI components in `frontend/src/components/`

---

**Built with ❤️ using Wails v2, React, and Google Gemini**
