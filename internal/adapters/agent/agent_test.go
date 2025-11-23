package agent

import (
	"context"
	"testing"

	"github.com/PundarikakshNTripathi/Kortex/internal/core/domain"
)

// MockBrowser implements ports.Browser for testing.
type MockBrowser struct {
	navigatedURL string
	clicked      string
	typed        string
	highlighted  string
}

func (m *MockBrowser) Navigate(url string) error {
	m.navigatedURL = url
	return nil
}

func (m *MockBrowser) GetSnapshot() (string, error) {
	return "<html><body><button id='submit'>Submit</button></body></html>", nil
}

func (m *MockBrowser) Highlight(selector, message string) error {
	m.highlighted = selector
	return nil
}

func (m *MockBrowser) Click(selector string) error {
	m.clicked = selector
	return nil
}

func (m *MockBrowser) Type(selector, text string) error {
	m.typed = selector + ":" + text
	return nil
}

// MockVectorStore implements ports.VectorStore for testing.
type MockVectorStore struct{}

func (m *MockVectorStore) Save(ctx context.Context, fragment *domain.MemoryFragment) error {
	return nil
}

func (m *MockVectorStore) Search(ctx context.Context, queryVector []float32, limit int) ([]domain.MemoryFragment, error) {
	return []domain.MemoryFragment{}, nil
}

func TestAgentInitialization(t *testing.T) {
	browser := &MockBrowser{}
	vectorStore := &MockVectorStore{}
	apiKey := "fake-api-key"

	agent := NewAgent(browser, vectorStore, apiKey)

	if agent == nil {
		t.Fatal("NewAgent returned nil")
	}
}

func TestTools(t *testing.T) {
	browser := &MockBrowser{}

	// Test NavigateTool
	navTool := &NavigateTool{Browser: browser}
	if navTool.IsLongRunning() {
		t.Error("NavigateTool should not be long running")
	}
	_, err := navTool.Run(context.Background(), struct{ URL string }{URL: "http://example.com"})
	if err != nil {
		t.Errorf("NavigateTool failed: %v", err)
	}
	if browser.navigatedURL != "http://example.com" {
		t.Errorf("Expected URL http://example.com, got %s", browser.navigatedURL)
	}

	// Test ClickTool
	clickTool := &ClickTool{Browser: browser}
	if clickTool.IsLongRunning() {
		t.Error("ClickTool should not be long running")
	}
	_, err = clickTool.Run(context.Background(), struct{ Selector string }{Selector: "#submit"})
	if err != nil {
		t.Errorf("ClickTool failed: %v", err)
	}
	if browser.clicked != "#submit" {
		t.Errorf("Expected click on #submit, got %s", browser.clicked)
	}
}
