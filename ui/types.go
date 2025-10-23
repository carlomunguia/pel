// Package ui provides user interface types and structures for the Pel pixel art editor.
package ui

import (
	"errors"
	"pel/apptype"
	"pel/pelcanvas"
	"pel/swatch"

	"fyne.io/fyne/v2"
)

// Common errors
var (
	ErrInvalidAppInit = errors.New("invalid app initialization")
	ErrNilComponent   = errors.New("required component is nil")
)

// AppInit holds all components needed for application initialization and UI setup.
// It serves as a central container for the main application components and their interactions.
type AppInit struct {
	PelCanvas *pelcanvas.PelCanvas // The main drawing canvas widget
	PelWindow fyne.Window          // The application's main window
	State     *apptype.State       // Global application state
	Swatches  []*swatch.Swatch     // Color palette swatches
}

// NewAppInit creates a new AppInit instance with the provided components.
// It validates that all required components are present.
func NewAppInit(canvas *pelcanvas.PelCanvas, window fyne.Window, state *apptype.State) (*AppInit, error) {
	if canvas == nil {
		return nil, errors.New("canvas cannot be nil")
	}
	if window == nil {
		return nil, errors.New("window cannot be nil")
	}
	if state == nil {
		return nil, errors.New("state cannot be nil")
	}

	return &AppInit{
		PelCanvas: canvas,
		PelWindow: window,
		State:     state,
		Swatches:  make([]*swatch.Swatch, 0),
	}, nil
}

// Validate checks if the AppInit structure has all required components properly initialized.
// Returns an error describing what is missing or invalid.
func (a *AppInit) Validate() error {
	if a == nil {
		return ErrInvalidAppInit
	}
	if a.PelWindow == nil {
		return errors.New("window is nil")
	}
	if a.PelCanvas == nil {
		return errors.New("canvas is nil")
	}
	if a.State == nil {
		return errors.New("state is nil")
	}
	return nil
}

// IsValid returns true if all required components are present.
func (a *AppInit) IsValid() bool {
	return a.Validate() == nil
}

// GetWindow returns the application window.
func (a *AppInit) GetWindow() fyne.Window {
	if a == nil {
		return nil
	}
	return a.PelWindow
}

// GetCanvas returns the drawing canvas.
func (a *AppInit) GetCanvas() *pelcanvas.PelCanvas {
	if a == nil {
		return nil
	}
	return a.PelCanvas
}

// GetState returns the application state.
func (a *AppInit) GetState() *apptype.State {
	if a == nil {
		return nil
	}
	return a.State
}

// GetSwatches returns all color swatches.
func (a *AppInit) GetSwatches() []*swatch.Swatch {
	if a == nil {
		return nil
	}
	return a.Swatches
}

// GetSwatchCount returns the number of swatches.
func (a *AppInit) GetSwatchCount() int {
	if a == nil || a.Swatches == nil {
		return 0
	}
	return len(a.Swatches)
}

// AddSwatch adds a new swatch to the collection.
func (a *AppInit) AddSwatch(s *swatch.Swatch) error {
	if a == nil {
		return ErrInvalidAppInit
	}
	if s == nil {
		return errors.New("cannot add nil swatch")
	}
	a.Swatches = append(a.Swatches, s)
	return nil
}

// GetSwatch returns the swatch at the specified index, or nil if out of bounds.
func (a *AppInit) GetSwatch(index int) *swatch.Swatch {
	if a == nil || index < 0 || index >= len(a.Swatches) {
		return nil
	}
	return a.Swatches[index]
}

// GetSelectedSwatch returns the currently selected swatch based on state.
func (a *AppInit) GetSelectedSwatch() *swatch.Swatch {
	if a == nil || a.State == nil {
		return nil
	}
	return a.GetSwatch(a.State.SwatchSelected)
}

// ClearSwatches removes all swatches from the collection.
func (a *AppInit) ClearSwatches() {
	if a == nil {
		return
	}
	a.Swatches = make([]*swatch.Swatch, 0)
}

// SetTitle updates the window title.
func (a *AppInit) SetTitle(title string) error {
	if a == nil || a.PelWindow == nil {
		return ErrNilComponent
	}
	a.PelWindow.SetTitle(title)
	return nil
}

// ShowWindow displays the application window.
func (a *AppInit) ShowWindow() error {
	if a == nil || a.PelWindow == nil {
		return ErrNilComponent
	}
	a.PelWindow.Show()
	return nil
}

// CloseWindow closes the application window.
func (a *AppInit) CloseWindow() error {
	if a == nil || a.PelWindow == nil {
		return ErrNilComponent
	}
	a.PelWindow.Close()
	return nil
}

// RefreshCanvas triggers a canvas refresh.
func (a *AppInit) RefreshCanvas() error {
	if a == nil || a.PelCanvas == nil {
		return ErrNilComponent
	}
	a.PelCanvas.Refresh()
	return nil
}

// Cleanup performs cleanup operations before shutdown.
// This can be extended to clean up resources, save state, etc.
func (a *AppInit) Cleanup() error {
	if a == nil {
		return nil
	}
	// Future: Add cleanup logic here
	// - Save unsaved work
	// - Close file handles
	// - Clear temporary data
	return nil
}
