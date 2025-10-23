// Package ui provides user interface layout and initialization for the Pel pixel art editor.
package ui

import (
	"fmt"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

// Layout configuration constants
const (
	ToolbarHeight    = 40  // Height of the toolbar in pixels
	StatusBarHeight  = 25  // Height of the status bar in pixels
	SwatchPanelWidth = 200 // Width of the swatch panel in pixels
)

// Setup initializes the complete user interface layout
// This is the main entry point for UI initialization
func Setup(app *AppInit) {
	if app == nil {
		log.Fatal("Cannot setup UI: AppInit is nil")
		return
	}

	// Validate required components
	if err := validateAppInit(app); err != nil {
		log.Fatalf("UI setup validation failed: %v", err)
		return
	}

	log.Println("Setting up UI components...")

	// Setup menus (File, Edit, View, Help, etc.)
	SetupMenus(app)
	log.Println("Menus initialized")

	// Build color swatch panel
	swatchesContainer := BuildSwatches(app)
	if swatchesContainer == nil {
		log.Println("Warning: Swatches container is nil, using empty container")
		swatchesContainer = container.NewVBox()
	}

	// Setup color picker
	colorPicker := SetupColorPicker(app)
	if colorPicker == nil {
		log.Println("Warning: Color picker is nil, using empty container")
		colorPicker = container.NewVBox()
	}

	// Build toolbar (optional - for future implementation)
	toolbar := buildToolbar(app)

	// Build status bar (optional - for future implementation)
	statusBar := buildStatusBar(app)

	// Create main layout:
	// - Top: toolbar (if available)
	// - Bottom: swatches and status bar
	// - Left: (reserved for future tools panel)
	// - Right: color picker
	// - Center: canvas

	// Combine status bar with swatches if status bar exists
	bottomContainer := swatchesContainer
	if statusBar != nil {
		bottomContainer = container.NewBorder(statusBar, nil, nil, nil, swatchesContainer)
	}

	// Create the main application layout
	appLayout := container.NewBorder(
		toolbar,         // top
		bottomContainer, // bottom
		nil,             // left (reserved for future tools panel)
		colorPicker,     // right
		app.PelCanvas,   // center
	)

	// Set the window content
	app.PelWindow.SetContent(appLayout)

	log.Println("UI setup complete")
}

// validateAppInit checks if the AppInit structure has all required fields
func validateAppInit(app *AppInit) error {
	if app.PelWindow == nil {
		return fmt.Errorf("window is nil")
	}
	if app.PelCanvas == nil {
		return fmt.Errorf("canvas is nil")
	}
	if app.State == nil {
		return fmt.Errorf("state is nil")
	}
	return nil
}

// buildToolbar creates the toolbar with common tool buttons
// TODO: Implement toolbar with brush type selection buttons
func buildToolbar(app *AppInit) fyne.CanvasObject {
	// Return nil for now - toolbar not yet implemented
	// Future: Add buttons for Pencil, Eraser, Fill, Line, Rectangle, Circle
	return nil
}

// buildStatusBar creates the status bar with info display
// TODO: Implement status bar showing:
// - Current tool name
// - Canvas size
// - Mouse position
// - Zoom level
func buildStatusBar(app *AppInit) fyne.CanvasObject {
	// Return nil for now - status bar not yet implemented
	return nil
}

// RebuildLayout rebuilds the entire UI layout
// Useful after configuration changes or major state updates
func RebuildLayout(app *AppInit) {
	log.Println("Rebuilding UI layout...")
	Setup(app)
}
