package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/PundarikakshNTripathi/Kortex/internal/core/ports"
	"github.com/google/uuid"
	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/model/gemini"
	"google.golang.org/adk/runner"
	"google.golang.org/adk/session"
	"google.golang.org/adk/tool"
	"google.golang.org/genai"
)

// AgentAdapter implements the "Brain" of Kortex.
// It uses Google Gemini (via the ADK) to understand user goals and decide which tools to use.
type AgentAdapter struct {
	browser     ports.Browser     // The "Hands" and "Eyes"
	vectorStore ports.VectorStore // The "Memory"
	apiKey      string            // Google Gemini API Key
	modelName   string            // e.g., "gemini-3-pro-preview"
}

// NewAgent creates a new AgentAdapter.
// It connects the core logic to the browser and database.
func NewAgent(browser ports.Browser, vectorStore ports.VectorStore, apiKey string) *AgentAdapter {
	return &AgentAdapter{
		browser:     browser,
		vectorStore: vectorStore,
		apiKey:      apiKey,
		modelName:   "gemini-3-pro-preview", // Using Gemini 3 Pro for advanced reasoning
	}
}

// ExecuteTask is the main entry point for the agent.
// It takes a user's goal (e.g., "Find cheap flights to Tokyo") and runs the ReAct loop.
func (a *AgentAdapter) ExecuteTask(ctx context.Context, goal string) error {
	// 1. RAG: Search for context (Placeholder)
	// In the future, this will look up past conversations to understand preferences.
	ragContext := ""
	// TODO: Implement actual embedding and search

	// 2. Initialize Gemini Model
	// We create a client to talk to Google's AI servers.
	model, err := gemini.NewModel(ctx, a.modelName, &genai.ClientConfig{
		APIKey: a.apiKey,
	})
	if err != nil {
		return fmt.Errorf("failed to create model: %w", err)
	}

	// 3. Define Tools
	// These are the capabilities we give the AI. It can't do anything else.
	tools := []tool.Tool{
		&NavigateTool{Browser: a.browser},    // Go to a URL
		&ClickTool{Browser: a.browser},       // Click something
		&TypeTool{Browser: a.browser},        // Type text
		&HighlightTool{Browser: a.browser},   // Show the user what you're looking at
		&GetSnapshotTool{Browser: a.browser}, // Read the page
	}

	// 4. Create ADK Agent
	// We give the AI a persona and instructions.
	systemInstruction := `You are Kortex, the Autonomous Interface Layer. Your goal is to navigate websites for users. 
You MUST use the Highlight tool to show the user where you are looking before you click. Speak simply.`
	if ragContext != "" {
		systemInstruction += "\n\nContext from memory:\n" + ragContext
	}

	adkAgent, err := llmagent.New(llmagent.Config{
		Name:        "kortex_agent",
		Model:       model,
		Description: "An autonomous agent that navigates the web.",
		Instruction: systemInstruction,
		Tools:       tools,
	})
	if err != nil {
		return fmt.Errorf("failed to create agent: %w", err)
	}

	// 5. Run the Agent using Runner
	// The Runner manages the conversation loop:
	// User Goal -> AI Thinks -> AI Calls Tool -> Tool Runs -> AI Sees Result -> AI Thinks...
	sessionService := session.InMemoryService()
	runnerCfg := runner.Config{
		Agent:          adkAgent,
		SessionService: sessionService,
	}

	r, err := runner.New(runnerCfg)
	if err != nil {
		return fmt.Errorf("failed to create runner: %w", err)
	}

	sessionID := uuid.New().String()
	userContent := &genai.Content{
		Parts: []*genai.Part{genai.NewPartFromText(goal)},
	}

	// Run returns an iterator that streams events as they happen.
	// We loop through these events to see what the agent is doing.
	for event, err := range r.Run(ctx, "user", sessionID, userContent, agent.RunConfig{}) {
		if err != nil {
			return fmt.Errorf("runner execution failed: %w", err)
		}

		// If the agent speaks (not a tool call), print it.
		if event.IsFinalResponse() {
			if event.Content != nil {
				for _, part := range event.Content.Parts {
					if part.Text != "" {
						fmt.Printf("Agent Output: %s\n", part.Text)
					}
				}
			}
		}
	}

	return nil
}

// --- Tool Wrappers ---
// These structs wrap the core Browser interface methods so the ADK can understand them.
// Each tool has a Name, Description, and Run method.

type NavigateTool struct {
	Browser ports.Browser
}

func (t *NavigateTool) Name() string        { return "navigate" }
func (t *NavigateTool) Description() string { return "Navigates to a specified URL." }
func (t *NavigateTool) IsLongRunning() bool { return false }
func (t *NavigateTool) Run(ctx context.Context, args struct{ URL string }) (string, error) {
	logFlightRecorder("navigate", args) // Log for debugging
	err := t.Browser.Navigate(args.URL)
	if err != nil {
		return "", err
	}
	return "Navigated to " + args.URL, nil
}

type ClickTool struct {
	Browser ports.Browser
}

func (t *ClickTool) Name() string        { return "click" }
func (t *ClickTool) Description() string { return "Clicks on an element specified by the selector." }
func (t *ClickTool) IsLongRunning() bool { return false }
func (t *ClickTool) Run(ctx context.Context, args struct{ Selector string }) (string, error) {
	logFlightRecorder("click", args)
	err := t.Browser.Click(args.Selector)
	if err != nil {
		return "", err
	}
	return "Clicked " + args.Selector, nil
}

type TypeTool struct {
	Browser ports.Browser
}

func (t *TypeTool) Name() string { return "type" }
func (t *TypeTool) Description() string {
	return "Types text into an element specified by the selector."
}
func (t *TypeTool) IsLongRunning() bool { return false }
func (t *TypeTool) Run(ctx context.Context, args struct {
	Selector string
	Text     string
}) (string, error) {
	logFlightRecorder("type", args)
	err := t.Browser.Type(args.Selector, args.Text)
	if err != nil {
		return "", err
	}
	return "Typed into " + args.Selector, nil
}

type HighlightTool struct {
	Browser ports.Browser
}

func (t *HighlightTool) Name() string { return "highlight" }
func (t *HighlightTool) Description() string {
	return "Highlights an element to show where the agent is looking."
}
func (t *HighlightTool) IsLongRunning() bool { return false }
func (t *HighlightTool) Run(ctx context.Context, args struct {
	Selector string
	Message  string
}) (string, error) {
	logFlightRecorder("highlight", args)
	err := t.Browser.Highlight(args.Selector, args.Message)
	if err != nil {
		return "", err
	}
	return "Highlighted " + args.Selector, nil
}

type GetSnapshotTool struct {
	Browser ports.Browser
}

func (t *GetSnapshotTool) Name() string { return "get_snapshot" }
func (t *GetSnapshotTool) Description() string {
	return "Gets a text snapshot of the current page accessibility tree."
}
func (t *GetSnapshotTool) IsLongRunning() bool { return false }
func (t *GetSnapshotTool) Run(ctx context.Context, args struct{}) (string, error) {
	logFlightRecorder("get_snapshot", args)
	return t.Browser.GetSnapshot()
}

// --- Flight Recorder ---
// This is a simple logging system that saves every action to a file.
// It helps us debug what the agent did and why.

type FlightRecord struct {
	Tool string      `json:"tool"`
	Args interface{} `json:"args"`
}

func logFlightRecorder(toolName string, args interface{}) {
	record := FlightRecord{
		Tool: toolName,
		Args: args,
	}
	data, _ := json.Marshal(record)
	// Open file in append mode so we don't lose history
	f, err := os.OpenFile("kortex_flight_recorder.jsonl", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("Failed to write to flight recorder: %v", err)
		return
	}
	defer f.Close()
	f.Write(data)
	f.WriteString("\n")
}
