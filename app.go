package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/PundarikakshNTripathi/Kortex/internal/adapters/agent"
	"github.com/PundarikakshNTripathi/Kortex/internal/infra/browser"
	"github.com/PundarikakshNTripathi/Kortex/internal/infra/sqlite"
	"github.com/joho/godotenv"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct holds the Kortex core components
type App struct {
	ctx         context.Context
	agent       *agent.AgentAdapter
	browser     *browser.PlaywrightBrowser
	vectorStore *sqlite.SQLiteVectorStore
	mu          sync.Mutex
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	apiKey := os.Getenv("GOOGLE_API_KEY")
	if apiKey == "" {
		a.emitLog("ERROR", "GOOGLE_API_KEY not set. Please create a .env file with your API key.")
		return
	}

	// Initialize Browser
	a.emitLog("INIT", "Initializing Playwright browser...")
	browserInstance := browser.NewPlaywrightBrowser()
	if err := browserInstance.Init(false); err != nil { // Desktop app uses visible browser
		a.emitLog("ERROR", fmt.Sprintf("Failed to initialize browser: %v", err))
		return
	}
	a.browser = browserInstance
	a.emitLog("INIT", "Browser initialized successfully ✓")

	// Initialize Vector Store
	a.emitLog("INIT", "Initializing vector store...")
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "./kortex.db"
	}
	vectorStore, err := sqlite.NewSQLiteVectorStore(dbPath)
	if err != nil {
		a.emitLog("ERROR", fmt.Sprintf("Failed to initialize vector store: %v", err))
		return
	}
	a.vectorStore = vectorStore
	a.emitLog("INIT", "Vector store initialized successfully ✓")

	// Initialize Agent
	a.emitLog("INIT", "Initializing Kortex agent...")
	a.agent = agent.NewAgent(a.browser, a.vectorStore, apiKey)
	a.emitLog("INIT", "🚀 Kortex agent ready! Awaiting your command...")
}

// shutdown is called when the app is closing
func (a *App) shutdown(ctx context.Context) {
	if a.browser != nil {
		a.emitLog("SHUTDOWN", "Closing browser...")
		// Browser cleanup would go here if needed
	}
}

// SendPrompt handles user input and executes the agent task
func (a *App) SendPrompt(prompt string) string {
	if a.agent == nil {
		return "Error: Agent not initialized. Please check your API key."
	}

	a.emitLog("USER", fmt.Sprintf("📝 %s", prompt))

	// Execute agent task in a goroutine to avoid blocking
	go func() {
		a.mu.Lock()
		defer a.mu.Unlock()

		a.emitLog("PLANNING", "🧠 Analyzing task and preparing execution plan...")

		// Create a custom context for the agent execution
		agentCtx := context.Background()

		// Execute the task
		err := a.agent.ExecuteTask(agentCtx, prompt)
		if err != nil {
			a.emitLog("ERROR", fmt.Sprintf("❌ Task execution failed: %v", err))
			return
		}

		a.emitLog("COMPLETE", "✅ Task completed successfully!")
	}()

	return "Task started. Watch the Mission Control for updates."
}

// emitLog sends a log event to the frontend
func (a *App) emitLog(level, message string) {
	if a.ctx != nil {
		runtime.EventsEmit(a.ctx, "kortex:log", map[string]string{
			"level":   level,
			"message": message,
		})
	}
	log.Printf("[%s] %s", level, message)
}

// GetStatus returns the current status of the agent
func (a *App) GetStatus() string {
	if a.agent == nil {
		return "Not initialized"
	}
	return "Ready"
}
