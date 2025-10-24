package main

import (
	"fmt"
	"image/color"
	"log"
	"github.com/carlomunguia/pel/apptype"
	"github.com/carlomunguia/pel/pelcanvas"
	"github.com/carlomunguia/pel/swatch"
	"github.com/carlomunguia/pel/ui"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

// Application metadata
const (
	AppName    = "Pel - Pixel Art Editor"
	AppVersion = "1.0.0"
	AppID      = "com.carlomunguia.pel"
)

// Configuration constants
const (
	DefaultCanvasWidth  = 600
	DefaultCanvasHeight = 600
	DefaultPxRows       = 10
	DefaultPxCols       = 10
	DefaultPxSize       = 30
	MaxSwatches         = 64
	MinWindowWidth      = 800
	MinWindowHeight     = 600
)

// Default color palette
var (
	DefaultBrushColor = color.NRGBA{R: 255, G: 255, B: 255, A: 255} // White
	DefaultBrushType  = apptype.BrushTypePencil                     // Pencil tool
)

func main() {
	// Initialize logger
	log.SetPrefix(fmt.Sprintf("[%s v%s] ", AppName, AppVersion))
	log.Println("Starting application...")

	// Create application with error recovery
	defer func() {
		if r := recover(); r != nil {
			log.Fatalf("Fatal error: %v", r)
		}
	}()

	// Create application
	pelApp := app.NewWithID(AppID)
	pelApp.SetIcon(nil) // TODO: Set app icon resource here

	// Create main window
	pelWindow := pelApp.NewWindow(AppName)

	// Configure window properties
	pelWindow.Resize(fyne.NewSize(MinWindowWidth, MinWindowHeight))
	pelWindow.SetFixedSize(false)
	pelWindow.SetMaster()
	pelWindow.CenterOnScreen()

	// Initialize application state with validated defaults
	state, err := initializeState()
	if err != nil {
		log.Fatalf("Failed to initialize application state: %v", err)
	}
	// Log state initialization with safe color extraction
	if nrgba, ok := state.BrushColor.(color.NRGBA); ok {
		log.Printf("Application state initialized: BrushType=%s, BrushColor=RGBA(%d,%d,%d,%d)",
			state.BrushType.String(), nrgba.R, nrgba.G, nrgba.B, nrgba.A)
	} else {
		log.Printf("Application state initialized: BrushType=%s", state.BrushType.String())
	}

	log.Printf("Application state initialized: BrushType=%s, BrushColor=RGBA(%d,%d,%d,%d)",
		state.BrushType.String(),
		state.BrushColor.(color.NRGBA).R,
		state.BrushColor.(color.NRGBA).G,
		state.BrushColor.(color.NRGBA).B,
		state.BrushColor.(color.NRGBA).A)

	// Configure canvas with constants
	pelCanvasConfig := apptype.PelCanvasConfig{
		DrawingArea:  fyne.NewSize(DefaultCanvasWidth, DefaultCanvasHeight),
		CanvasOffset: fyne.NewPos(0, 0),
		PxRows:       DefaultPxRows,
		PxCols:       DefaultPxCols,
		PxSize:       DefaultPxSize,
	}

	// Validate canvas configuration using built-in method
	if err := pelCanvasConfig.Validate(); err != nil {
		log.Fatalf("Invalid canvas configuration: %v", err)
	}

	// Create canvas
	pelCanvas := pelcanvas.NewPelCanvas(&state, pelCanvasConfig)
	log.Printf("Canvas initialized: %dx%d pixels (%dx%d grid), Total pixels: %d",
		DefaultCanvasWidth, DefaultCanvasHeight,
		DefaultPxRows, DefaultPxCols,
		pelCanvasConfig.TotalPixels())

	// Initialize application components
	appInit := ui.AppInit{
		PelCanvas: pelCanvas,
		PelWindow: pelWindow,
		State:     &state,
		Swatches:  make([]*swatch.Swatch, 0, MaxSwatches),
	}

	ui.Setup(&appInit)

	log.Println("UI setup complete")

	// Set cleanup handler
	pelWindow.SetOnClosed(func() {
		log.Println("Application closing...")
		// TODO: Add any cleanup logic here (save state, close resources, etc.)
		if state.HasFilePath() {
			log.Printf("Last opened file: %s", state.FilePath)
		}
	})

	// Start the application
	log.Println("Showing main window")
	appInit.PelWindow.ShowAndRun()
}

// initializeState creates and validates the initial application state
func initializeState() (apptype.State, error) {
	state := apptype.State{
		BrushColor:     DefaultBrushColor,
		BrushType:      DefaultBrushType,
		SwatchSelected: 0,
		FilePath:       "", // Empty for new project
	}

	// Validate state using built-in validation method
	if err := state.Validate(); err != nil {
		return state, fmt.Errorf("state validation failed: %w", err)
	}

	return state, nil
}
