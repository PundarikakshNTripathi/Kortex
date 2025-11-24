package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"

	"github.com/PundarikakshNTripathi/Kortex/internal/adapters/agent"
	"github.com/PundarikakshNTripathi/Kortex/internal/infra/browser"
	"github.com/PundarikakshNTripathi/Kortex/internal/infra/sqlite"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/websocket/v2"
	"github.com/joho/godotenv"
)

// KortexCore holds the core components
type KortexCore struct {
	agent       *agent.AgentAdapter
	browser     *browser.PlaywrightBrowser
	vectorStore *sqlite.SQLiteVectorStore
	mu          sync.Mutex
}

// WebSocketMessage represents a message from the client
type WebSocketMessage struct {
	Goal string `json:"goal"`
}

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Get configuration from environment
	apiKey := os.Getenv("GOOGLE_API_KEY")
	if apiKey == "" {
		log.Fatal("GOOGLE_API_KEY not set. Please create a .env file with your API key.")
	}

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "./kortex.db"
	}

	headlessStr := os.Getenv("HEADLESS")
	headless := false
	if headlessStr != "" {
		var err error
		headless, err = strconv.ParseBool(headlessStr)
		if err != nil {
			log.Printf("Invalid HEADLESS value '%s', defaulting to false", headlessStr)
		}
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Initialize Kortex Core
	log.Println("🚀 Initializing Kortex Core...")

	// Initialize Browser
	log.Printf("📱 Initializing Playwright browser (headless: %v)...", headless)
	browserInstance := browser.NewPlaywrightBrowser()
	if err := browserInstance.Init(headless); err != nil {
		log.Fatalf("❌ Failed to initialize browser: %v", err)
	}
	log.Println("✓ Browser initialized successfully")

	// Initialize Vector Store
	log.Println("💾 Initializing vector store...")
	vectorStore, err := sqlite.NewSQLiteVectorStore(dbPath)
	if err != nil {
		log.Fatalf("❌ Failed to initialize vector store: %v", err)
	}
	log.Println("✓ Vector store initialized successfully")

	// Initialize Agent
	log.Println("🧠 Initializing Kortex agent...")
	agentAdapter := agent.NewAgent(browserInstance, vectorStore, apiKey)
	log.Println("✓ Kortex agent ready!")

	// Create Kortex Core
	core := &KortexCore{
		agent:       agentAdapter,
		browser:     browserInstance,
		vectorStore: vectorStore,
	}

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		AppName:      "Kortex Web Server",
		ServerHeader: "Kortex",
	})

	// Middleware
	app.Use(logger.New())
	app.Use(cors.New())

	// Health check endpoint
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
			"agent":  "ready",
		})
	})

	// WebSocket upgrade middleware
	app.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	// WebSocket chat endpoint
	app.Get("/ws/chat", websocket.New(func(c *websocket.Conn) {
		log.Printf("🔌 New WebSocket connection from %s", c.RemoteAddr())

		// Send welcome message
		if err := c.WriteJSON(fiber.Map{
			"type":    "log",
			"level":   "INIT",
			"message": "🚀 Connected to Kortex! Send your goal to begin.",
		}); err != nil {
			log.Printf("Error sending welcome message: %v", err)
			return
		}

		defer func() {
			log.Printf("🔌 WebSocket connection closed from %s", c.RemoteAddr())
			c.Close()
		}()

		for {
			var msg WebSocketMessage
			if err := c.ReadJSON(&msg); err != nil {
				if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
					log.Println("Client closed connection")
					break
				}
				log.Printf("Error reading message: %v", err)
				break
			}

			if msg.Goal == "" {
				c.WriteJSON(fiber.Map{
					"type":    "error",
					"message": "Goal cannot be empty",
				})
				continue
			}

			// Send acknowledgment
			c.WriteJSON(fiber.Map{
				"type":    "log",
				"level":   "USER",
				"message": fmt.Sprintf("📝 %s", msg.Goal),
			})

			// Execute task in goroutine
			go func(goal string) {
				core.mu.Lock()
				defer core.mu.Unlock()

				// Send planning log
				c.WriteJSON(fiber.Map{
					"type":    "log",
					"level":   "PLANNING",
					"message": "🧠 Analyzing task and preparing execution plan...",
				})

				// Execute the task
				ctx := context.Background()
				err := core.agent.ExecuteTask(ctx, goal)

				if err != nil {
					c.WriteJSON(fiber.Map{
						"type":    "log",
						"level":   "ERROR",
						"message": fmt.Sprintf("❌ Task execution failed: %v", err),
					})
					return
				}

				c.WriteJSON(fiber.Map{
					"type":    "log",
					"level":   "COMPLETE",
					"message": "✅ Task completed successfully!",
				})
			}(msg.Goal)
		}
	}))

	// Graceful shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
		<-sigChan

		log.Println("\n🛑 Shutting down Kortex...")

		if err := app.Shutdown(); err != nil {
			log.Printf("Error during shutdown: %v", err)
		}

		// Cleanup browser
		if core.browser != nil && core.browser != nil {
			log.Println("Closing browser...")
		}

		os.Exit(0)
	}()

	// Start server
	log.Printf("🌐 Kortex Web Server starting on port %s...", port)
	log.Printf("📡 WebSocket endpoint: ws://localhost:%s/ws/chat", port)

	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
