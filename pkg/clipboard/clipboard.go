package clipboard

import (
	"log"
	"runtime"
	"time"

	"github.com/go-vgo/robotgo"
)

// Manager handles clipboard operations and text pasting
type Manager struct{}

// NewManager creates a new clipboard manager
func NewManager() *Manager {
	return &Manager{}
}

// CopyToClipboard copies text to the system clipboard
func (m *Manager) CopyToClipboard(text string) {
	log.Println("Copying text to clipboard...")
	robotgo.WriteAll(text)
}

// Clear clears the clipboard
func (m *Manager) Clear() {
	robotgo.WriteAll("")
	log.Println("Clipboard cleared")
}

// PasteAtCursor pastes the text at the current cursor position
func (m *Manager) PasteAtCursor(text string) {
	log.Println("Manual paste button clicked")
	
	// Run paste operation in a separate goroutine to avoid blocking UI
	go func(clipText string) {
		// Using robotgo for manual paste
		log.Println("Starting paste sequence with robotgo...")
		
		// 1. Make sure text is in clipboard
		log.Println("Ensuring text is in clipboard...")
		robotgo.WriteAll(clipText)
		time.Sleep(300 * time.Millisecond)
		
		// 2. Get current mouse position
		log.Println("Getting current mouse position...")
		x, y := robotgo.GetMousePos()
		log.Printf("Current mouse position: x=%d, y=%d", x, y)
		
		// 3. Simulate a mouse click at current position to ensure focus
		log.Println("Simulating mouse click to ensure focus...")
		robotgo.MouseClick("left", false)
		time.Sleep(300 * time.Millisecond)
		
		// 4. Simulate paste keystroke (Ctrl+V or Command+V)
		if runtime.GOOS == "darwin" {
			log.Println("Using Command+V for paste on macOS...")
			robotgo.KeyTap("v", "command")
		} else {
			log.Println("Using Ctrl+V for paste...")
			robotgo.KeyTap("v", "ctrl")
			
			// If first attempt doesn't work, try the key sequence approach
			time.Sleep(300 * time.Millisecond)
			log.Println("Trying alternate key sequence approach...")
			robotgo.KeyDown("ctrl")
			time.Sleep(150 * time.Millisecond)
			robotgo.KeyTap("v")
			time.Sleep(150 * time.Millisecond)
			robotgo.KeyUp("ctrl")
		}
		
		log.Println("Manual paste completed")
	}(text)
}