package browser

import (
	"strings"
	"testing"
)

func TestPlaywrightBrowser(t *testing.T) {
	// Skip if short mode is enabled, as this launches a real browser
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	browser := NewPlaywrightBrowser()
	err := browser.Init()
	if err != nil {
		t.Fatalf("Failed to init browser: %v", err)
	}
	// Ensure we close the browser/playwright at the end
	defer func() {
		if browser.browser != nil {
			browser.browser.Close()
		}
		if browser.pw != nil {
			browser.pw.Stop()
		}
	}()

	t.Run("Navigate and Highlight", func(t *testing.T) {
		// Use a data URL to avoid network dependency and ensure consistent content
		html := `
		<html>
			<body>
				<div id="target">Hello World</div>
				<button>Click me</button>
			</body>
		</html>
		`
		url := "data:text/html," + html

		err := browser.Navigate(url)
		if err != nil {
			t.Fatalf("Failed to navigate: %v", err)
		}

		err = browser.Highlight("#target", "Here is the target")
		if err != nil {
			t.Errorf("Failed to highlight: %v", err)
		}
	})

	t.Run("GetSnapshot", func(t *testing.T) {
		snapshot, err := browser.GetSnapshot()
		if err != nil {
			t.Fatalf("Failed to get snapshot: %v", err)
		}

		// Basic validation
		if !strings.Contains(snapshot, "Hello World") {
			t.Errorf("Snapshot should contain 'Hello World', got: %s", snapshot)
		}
		if !strings.Contains(snapshot, "#target") {
			t.Errorf("Snapshot should contain selector '#target', got: %s", snapshot)
		}
		if !strings.Contains(snapshot, "button") { // Role or tag name
			t.Errorf("Snapshot should contain button, got: %s", snapshot)
		}
	})
}
