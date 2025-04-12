package ui

import (
	"context"
	"fmt"
	"log"

	"github.com/getlantern/systray"

	"github.com/tarasowski/autospeech/pkg/clipboard"
	"github.com/tarasowski/autospeech/pkg/config"
)

// TrayMenu manages the system tray interface
type TrayMenu struct {
	ctx        context.Context
	cancel     context.CancelFunc
	state      *config.AppState
	clipMgr    *clipboard.Manager
	notifyMgr  *NotificationManager
	menuItems  map[string]*systray.MenuItem
	onStart    func()
	onStop     func()
	onQuit     func()
	onPaste    func(string)
}

// NewTrayMenu creates a new system tray interface
func NewTrayMenu(state *config.AppState) *TrayMenu {
	ctx, cancel := context.WithCancel(context.Background())
	return &TrayMenu{
		ctx:       ctx,
		cancel:    cancel,
		state:     state,
		clipMgr:   clipboard.NewManager(),
		notifyMgr: NewNotificationManager(),
		menuItems: make(map[string]*systray.MenuItem),
	}
}

// SetCallbacks sets the menu item callbacks
func (tm *TrayMenu) SetCallbacks(onStart, onStop, onQuit func(), onPaste func(string)) {
	tm.onStart = onStart
	tm.onStop = onStop
	tm.onQuit = onQuit
	tm.onPaste = onPaste
}

// Start initializes and shows the system tray
func (tm *TrayMenu) Start() {
	go systray.Run(
		func() { tm.setupTray() }, 
		func() { 
			fmt.Println("Cleaning up...")
			if tm.cancel != nil {
				tm.cancel()
			}
		},
	)
}

// setupTray initializes the system tray menu
func (tm *TrayMenu) setupTray() {
	systray.SetTitle("Speech-to-Text")
	systray.SetTooltip("Speech Recognition")

	mRecord := systray.AddMenuItem("Start Recording", "Start speech recognition")
	
	// Cache menu items for later use
	tm.menuItems["Start Recording"] = mRecord
	
	systray.AddSeparator()
	mQuit := systray.AddMenuItem("Quit", "Quit the app")

	go func() {
		for {
			select {
			case <-tm.ctx.Done():
				return
			case <-mRecord.ClickedCh:
				if !tm.state.IsRecording() {
					// Start recording
					tm.state.SetRecording(true)
					mRecord.SetTitle("Stop Recording")
					// Add this to the cache with the new title
					tm.menuItems["Stop Recording"] = mRecord
					
					// Clear previous transcribed text
					tm.state.SetTranscribedText("")
					
					if tm.onStart != nil {
						tm.onStart()
					}
				} else {
					// Stop recording and transcribe
					tm.state.SetRecording(false)
					mRecord.SetTitle("Processing...")
					
					if tm.onStop != nil {
						tm.onStop()
					}
				}
			case <-mQuit.ClickedCh:
				log.Println("Quit requested")
				systray.Quit()
				if tm.onQuit != nil {
					tm.onQuit()
				}
				return
			}
		}
	}()
}

// UpdateMenuTitle updates a menu item's title
func (tm *TrayMenu) UpdateMenuTitle(key, title string) {
	if item, ok := tm.menuItems[key]; ok {
		item.SetTitle(title)
	}
}

// UpdateMenuTooltip updates a menu item's tooltip
func (tm *TrayMenu) UpdateMenuTooltip(key, tooltip string) {
	if item, ok := tm.menuItems[key]; ok {
		item.SetTooltip(tooltip)
	}
}

// EnableMenuItem enables a menu item
func (tm *TrayMenu) EnableMenuItem(key string) {
	if item, ok := tm.menuItems[key]; ok {
		item.Enable()
	}
}

// DisableMenuItem disables a menu item
func (tm *TrayMenu) DisableMenuItem(key string) {
	if item, ok := tm.menuItems[key]; ok {
		item.Disable()
	}
}

// UpdateTrayTitle updates the system tray title
func (tm *TrayMenu) UpdateTrayTitle(title string) {
	systray.SetTitle(title)
}

// UpdateTrayTooltip updates the system tray tooltip
func (tm *TrayMenu) UpdateTrayTooltip(tooltip string) {
	systray.SetTooltip(tooltip)
}

// GetMenuItem gets a menu item by key
func (tm *TrayMenu) GetMenuItem(key string) (*systray.MenuItem, bool) {
	item, ok := tm.menuItems[key]
	return item, ok
}

// AddMenuItemToCache adds a menu item to the cache
func (tm *TrayMenu) AddMenuItemToCache(key string, item *systray.MenuItem) {
	tm.menuItems[key] = item
}

// RemoveMenuItemFromCache removes a menu item from the cache
func (tm *TrayMenu) RemoveMenuItemFromCache(key string) {
	delete(tm.menuItems, key)
}

// SetupForTranscriptionComplete updates the UI after transcription is complete
func (tm *TrayMenu) SetupForTranscriptionComplete(text string) {
	// Store the transcribed text but don't display it in the tray
	if text != "" {
		tm.state.SetTranscribedText(text)
		
		// Copy to clipboard right away
		tm.clipMgr.CopyToClipboard(text)
	}
	
	// Just reset the recording button regardless of text content
	if mRecord, ok := tm.menuItems["Stop Recording"]; ok {
		mRecord.SetTitle("Start Recording")
		// Update the cache
		tm.menuItems["Start Recording"] = mRecord
		// Remove the old title from the cache
		delete(tm.menuItems, "Stop Recording")
	}
	
	// Keep title simple
	systray.SetTitle("Speech-to-Text")
	systray.SetTooltip("Speech Recognition")
}