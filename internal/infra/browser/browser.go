package browser

import (
	"encoding/json"
	"fmt"

	"github.com/playwright-community/playwright-go"
)

type PlaywrightBrowser struct {
	pw      *playwright.Playwright
	browser playwright.Browser
	page    playwright.Page
}

func NewPlaywrightBrowser() *PlaywrightBrowser {
	return &PlaywrightBrowser{}
}

func (pb *PlaywrightBrowser) Init() error {
	pw, err := playwright.Run()
	if err != nil {
		return fmt.Errorf("could not start playwright: %v", err)
	}
	pb.pw = pw

	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(false),
	})
	if err != nil {
		return fmt.Errorf("could not launch browser: %v", err)
	}
	pb.browser = browser

	page, err := browser.NewPage()
	if err != nil {
		return fmt.Errorf("could not create page: %v", err)
	}
	pb.page = page

	return nil
}

func (pb *PlaywrightBrowser) Navigate(url string) error {
	if pb.page == nil {
		return fmt.Errorf("browser not initialized")
	}
	if _, err := pb.page.Goto(url); err != nil {
		return fmt.Errorf("could not navigate to %s: %v", url, err)
	}
	return nil
}

func (pb *PlaywrightBrowser) Highlight(selector, message string) error {
	if pb.page == nil {
		return fmt.Errorf("browser not initialized")
	}

	// Wait for selector with 10s timeout
	_, err := pb.page.WaitForSelector(selector, playwright.PageWaitForSelectorOptions{
		Timeout: playwright.Float(10000),
	})
	if err != nil {
		return fmt.Errorf("element not found: %s", selector)
	}

	js := fmt.Sprintf(`
		// Clean old Kortex highlights
		document.querySelectorAll('.kortex-highlight').forEach(e => e.remove());
		document.querySelectorAll('.kortex-tooltip').forEach(e => e.remove());
		
		const el = document.querySelector('%s');
		if (el) {
			el.style.outline = '4px solid #00E5FF'; // Cyan/Electric Blue
			el.style.boxShadow = '0 0 20px rgba(0, 229, 255, 0.6)';
			el.scrollIntoView({behavior: 'smooth', block: 'center'});
			
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

type AccessibilityNode struct {
	Role     string              `json:"role"`
	Name     string              `json:"name"`
	Selector string              `json:"selector"`
	Children []AccessibilityNode `json:"children,omitempty"`
}

func (pb *PlaywrightBrowser) GetSnapshot() (string, error) {
	if pb.page == nil {
		return "", fmt.Errorf("browser not initialized")
	}

	// Inject JS to traverse DOM and build simplified tree
	// This script generates a unique selector for each element and extracts role/name
	js := `
		(function() {
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

			function traverse(el) {
				// Skip hidden elements or scripts/styles
				const style = window.getComputedStyle(el);
				if (style.display === 'none' || style.visibility === 'hidden' || 
					['SCRIPT', 'STYLE', 'NOSCRIPT', 'META', 'HEAD', 'TITLE', 'LINK'].includes(el.tagName)) {
					return null;
				}

				const node = {
					role: getRole(el),
					name: getName(el).substring(0, 50).replace(/\s+/g, ' ').trim(), // Truncate and clean
					selector: getSelector(el),
					children: []
				};

				for (let child of el.children) {
					const childNode = traverse(child);
					if (childNode) {
						node.children.push(childNode);
					}
				}
				
				// If no children and no name, might be irrelevant container, but keeping structure for now
				return node;
			}

			return traverse(document.body);
		})()
	`

	result, err := pb.page.Evaluate(js)
	if err != nil {
		return "", fmt.Errorf("failed to evaluate snapshot script: %v", err)
	}

	// result is a map[string]interface{} representing the JSON object
	// We can marshal it directly
	bytes, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal snapshot: %v", err)
	}

	return string(bytes), nil
}
