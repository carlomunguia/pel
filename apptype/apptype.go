// Package apptype defines the core types and interfaces for the Pel pixel art editor.
package apptype

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
)

// BrushType represents the type of brush tool being used
type BrushType int

// Brush type constants
const (
	BrushTypePencil BrushType = iota
	BrushTypeEraser
	BrushTypeFill
	BrushTypeLine
	BrushTypeRectangle
	BrushTypeCircle
)

// String returns a human-readable name for the brush type
func (bt BrushType) String() string {
	switch bt {
	case BrushTypePencil:
		return "Pencil"
	case BrushTypeEraser:
		return "Eraser"
	case BrushTypeFill:
		return "Fill"
	case BrushTypeLine:
		return "Line"
	case BrushTypeRectangle:
		return "Rectangle"
	case BrushTypeCircle:
		return "Circle"
	default:
		return "Unknown"
	}
}

// IsValid checks if the brush type is valid
func (bt BrushType) IsValid() bool {
	return bt >= BrushTypePencil && bt <= BrushTypeCircle
}

// Brushable defines the interface for objects that can be painted on
type Brushable interface {
	// SetColor sets the color at the specified canvas coordinates
	// Returns an error if the operation fails
	SetColor(c color.Color, x, y int) error

	// MouseToCanvasXY converts mouse event coordinates to canvas coordinates
	// Returns nil pointers if the coordinates are outside the canvas
	MouseToCanvasXY(ev *desktop.MouseEvent) (*int, *int)
}

// PelCanvasConfig holds the configuration for the pixel canvas
type PelCanvasConfig struct {
	DrawingArea  fyne.Size     // Size of the drawing area in pixels
	CanvasOffset fyne.Position // Offset of the canvas from the window origin
	PxRows       int           // Number of pixel rows in the grid
	PxCols       int           // Number of pixel columns in the grid
	PxSize       int           // Size of each pixel in screen pixels
}

// Validate checks if the canvas configuration is valid
func (c PelCanvasConfig) Validate() error {
	if c.PxRows <= 0 {
		return fmt.Errorf("pixel rows must be positive, got: %d", c.PxRows)
	}
	if c.PxCols <= 0 {
		return fmt.Errorf("pixel columns must be positive, got: %d", c.PxCols)
	}
	if c.PxSize <= 0 {
		return fmt.Errorf("pixel size must be positive, got: %d", c.PxSize)
	}
	if c.DrawingArea.Width <= 0 || c.DrawingArea.Height <= 0 {
		return fmt.Errorf("drawing area dimensions must be positive, got: %.2fx%.2f",
			c.DrawingArea.Width, c.DrawingArea.Height)
	}
	return nil
}

// TotalPixels returns the total number of pixels in the canvas
func (c PelCanvasConfig) TotalPixels() int {
	return c.PxRows * c.PxCols
}

// State represents the current state of the application
type State struct {
	BrushColor     color.Color // Current brush color
	BrushType      BrushType   // Current brush tool type
	SwatchSelected int         // Index of the currently selected color swatch
	FilePath       string      // Path to the currently open file (empty if new/unsaved)
}

// SetFilePath updates the file path for the current project
func (s *State) SetFilePath(path string) {
	s.FilePath = path
}

// SetBrushColor updates the current brush color
func (s *State) SetBrushColor(c color.Color) {
	s.BrushColor = c
}

// SetBrushType updates the current brush type
func (s *State) SetBrushType(bt BrushType) {
	if bt.IsValid() {
		s.BrushType = bt
	}
}

// SetSwatchSelected updates the selected swatch index
func (s *State) SetSwatchSelected(index int) {
	if index >= 0 {
		s.SwatchSelected = index
	}
}

// HasUnsavedChanges returns true if there's a file path (indicating the project has been saved)
func (s *State) HasFilePath() bool {
	return s.FilePath != ""
}

// Validate checks if the state is valid
func (s *State) Validate() error {
	if !s.BrushType.IsValid() {
		return fmt.Errorf("invalid brush type: %d", s.BrushType)
	}
	if s.SwatchSelected < 0 {
		return fmt.Errorf("swatch index cannot be negative: %d", s.SwatchSelected)
	}
	if s.BrushColor == nil {
		return fmt.Errorf("brush color cannot be nil")
	}
	return nil
}
