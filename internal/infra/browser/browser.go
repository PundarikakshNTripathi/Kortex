package browser

import (
	"encoding/json"
	"fmt"

	"github.com/playwright-community/playwright-go"
)

// PlaywrightBrowser implements the Browser interface using Microsoft Playwright.
// Playwright is a tool that lets code control a web browser (Chrome, Firefox, etc.).
type PlaywrightBrowser struct {
	pw      *playwright.Playwright // The main Playwright instance
	browser playwright.Browser     // The browser application (e.g., Chromium)
	page    playwright.Page        // The specific tab or window we are controlling
}

// NewPlaywrightBrowser creates a new instance of our browser adapter.
func NewPlaywrightBrowser() *PlaywrightBrowser {
	return &PlaywrightBrowser{}
}

// Init starts the browser.
// If headless is true, the browser runs in the background (invisible).
// If headless is false, you can see the browser window (useful for debugging).
func (pb *PlaywrightBrowser) Init(headless bool) error {
	pw, err := playwright.Run()
	if err != nil {
		return fmt.Errorf("could not start playwright: %v", err)
	}
	pb.pw = pw

	// Launch Chromium (open-source Chrome)
	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(headless),
	})
	if err != nil {
		return fmt.Errorf("could not launch browser: %v", err)
	}
	pb.browser = browser

	// Open a new tab
	page, err := browser.NewPage()
	if err != nil {
		return fmt.Errorf("could not create page: %v", err)
	}
	pb.page = page

	return nil
}

// Navigate tells the browser to go to a specific website.
func (pb *PlaywrightBrowser) Navigate(url string) error {
	if pb.page == nil {
		return fmt.Errorf("browser not initialized")
	}
	// Goto waits for the page to load before returning
	if _, err := pb.page.Goto(url); err != nil {
		return fmt.Errorf("could not navigate to %s: %v", url, err)
	}
	return nil
}

// Highlight injects JavaScript into the page to draw a colored box around an element.
// This helps the user see what the agent is focusing on.
func (pb *PlaywrightBrowser) Highlight(selector, message string) error {
	if pb.page == nil {
		return fmt.Errorf("browser not initialized")
	}

	// Wait up to 10 seconds for the element to appear
	_, err := pb.page.WaitForSelector(selector, playwright.PageWaitForSelectorOptions{
		Timeout: playwright.Float(10000),
	})
	if err != nil {
		return fmt.Errorf("element not found: %s", selector)
	}

	// This JavaScript runs INSIDE the browser page.
	// It removes old highlights and adds a new one with a tooltip.
	js := fmt.Sprintf(`
		// Clean old Kortex highlights
		document.querySelectorAll('.kortex-highlight').forEach(e => e.remove());
		document.querySelectorAll('.kortex-tooltip').forEach(e => e.remove());
		
		const el = document.querySelector('%s');
		if (el) {
			// Add visual styling (Cyan outline)
			el.style.outline = '4px solid #00E5FF'; 
			el.style.boxShadow = '0 0 20px rgba(0, 229, 255, 0.6)';
			el.scrollIntoView({behavior: 'smooth', block: 'center'});
			
			// Create and style the tooltip
			const tip = document.createElement('div');
			tip.className = 'kortex-tooltip';
			tip.innerText = '%s';
			tip.style.position = 'absolute';
			tip.style.background = '#333';
			tip.style.color = '#fff';
			tip.style.padding = '5px 10px';
			tip.style.borderRadius = '4px';
			tip.style.zIndex = '10000';
			tip.style.fontSize = '12px';
			tip.style.marginTop = '5px';
			
			// Position tooltip relative to element
			const rect = el.getBoundingClientRect();
			tip.style.left = (rect.left + window.scrollX) + 'px';
			tip.style.top = (rect.bottom + window.scrollY) + 'px';

			document.body.appendChild(tip);
		}
	`, selector, message)

	_, err = pb.page.Evaluate(js)
	if err != nil {
		return fmt.Errorf("failed to inject highlight script: %v", err)
	}

	return nil
}

// Click simulates a mouse click on an element.
func (pb *PlaywrightBrowser) Click(selector string) error {
	if pb.page == nil {
		return fmt.Errorf("browser not initialized")
	}

	err := pb.page.Click(selector, playwright.PageClickOptions{
		Timeout: playwright.Float(10000),
	})
	if err != nil {
		return fmt.Errorf("failed to click %s: %v", selector, err)
	}

	return nil
}

// Type simulates typing text into an input field.
func (pb *PlaywrightBrowser) Type(selector, text string) error {
	if pb.page == nil {
		return fmt.Errorf("browser not initialized")
	}

	err := pb.page.Fill(selector, text, playwright.PageFillOptions{
		Timeout: playwright.Float(10000),
	})
	if err != nil {
		return fmt.Errorf("failed to type into %s: %v", selector, err)
	}

	return nil
}

// AccessibilityNode is a simplified representation of a DOM element.
// We use this to give the AI a "clean" view of the page, removing noise like CSS classes.
type AccessibilityNode struct {
	Role     string              `json:"role"`               // e.g., "button", "link", "heading"
	Name     string              `json:"name"`               // The text or label
	Selector string              `json:"selector"`           // CSS selector to reach this element
	Children []AccessibilityNode `json:"children,omitempty"` // Nested elements
}

// GetSnapshot scans the page and returns a JSON tree of interactive elements.
// This is crucial because raw HTML is too messy for the AI to process efficiently.
func (pb *PlaywrightBrowser) GetSnapshot() (string, error) {
	if pb.page == nil {
		return "", fmt.Errorf("browser not initialized")
	}

	// Inject JS to traverse DOM and build simplified tree
	// This script generates a unique selector for each element and extracts role/name
	js := `
		(function() {
			// Helper to generate a unique CSS selector for an element
			function getSelector(el) {
				if (el.id) return '#' + el.id;
				if (el === document.body) return 'body';
				let idx = 1;
				let sib = el.previousElementSibling;
				while (sib) {
					if (sib.tagName === el.tagName) idx++;
					sib = sib.previousElementSibling;
				}
				return getSelector(el.parentElement) + ' > ' + el.tagName.toLowerCase() + ':nth-of-type(' + idx + ')';
			}

			function getRole(el) {
				return el.getAttribute('role') || el.tagName.toLowerCase();
			}

			function getName(el) {
				return el.getAttribute('aria-label') || el.innerText || '';
			}

			// Recursive function to walk the DOM tree
			function traverse(el) {
				// Skip hidden elements or scripts/styles
				const style = window.getComputedStyle(el);
				if (style.display === 'none' || style.visibility === 'hidden' || 
					['SCRIPT', 'STYLE', 'NOSCRIPT', 'META', 'HEAD', 'TITLE', 'LINK'].includes(el.tagName)) {
					return null;
				}

				const node = {
					role: getRole(el),
					name: getName(el).substring(0, 50).replace(/\s+/g, ' ').trim(), // Truncate and clean text
					selector: getSelector(el),
					children: []
				};

				for (let child of el.children) {
					const childNode = traverse(child);
					if (childNode) {
						node.children.push(childNode);
					}
				}
				
				return node;
			}

			return traverse(document.body);
		})()
	`

	result, err := pb.page.Evaluate(js)
	if err != nil {
		return fmt.Errorf("failed to evaluate snapshot script: %v", err)
	}

	// Convert the result to a pretty-printed JSON string
	bytes, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal snapshot: %v", err)
	}

	return string(bytes), nil
}
