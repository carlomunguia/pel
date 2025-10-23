// Package ui provides swatch panel construction for the Pel pixel art editor.
package ui

import (
	"image/color"
	"log"
	"pel/swatch"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// Swatch panel constants
const (
	MaxSwatches    = 64  // Maximum number of color swatches
	SwatchGridSize = 20  // Size of each swatch in the grid
	SwatchesPerRow = 8   // Number of swatches per row
	DefaultSwatchR = 255 // Default red value
	DefaultSwatchG = 255 // Default green value
	DefaultSwatchB = 255 // Default blue value
	DefaultSwatchA = 255 // Default alpha value
)

// Default swatch color
var DefaultSwatchColor = color.NRGBA{
	R: DefaultSwatchR,
	G: DefaultSwatchG,
	B: DefaultSwatchB,
	A: DefaultSwatchA,
}

// BuildSwatches creates and configures the color swatch panel
func BuildSwatches(app *AppInit) *fyne.Container {
	if app == nil {
		log.Println("Warning: Cannot build swatches - app is nil")
		return container.NewVBox()
	}

	if app.State == nil {
		log.Println("Warning: Cannot build swatches - state is nil")
		return container.NewVBox()
	}

	log.Println("Building swatches panel...")

	// Validate swatch capacity
	swatchCapacity := cap(app.Swatches)
	if swatchCapacity == 0 || swatchCapacity > MaxSwatches {
		log.Printf("Warning: Invalid swatch capacity %d, using default %d", swatchCapacity, MaxSwatches)
		swatchCapacity = MaxSwatches
	}

	// Pre-allocate canvas objects for all swatches
	canvasSwatches := make([]fyne.CanvasObject, 0, swatchCapacity)

	// Create swatches
	for i := 0; i < swatchCapacity; i++ {
		s := createSwatch(app, i)
		if s == nil {
			log.Printf("Warning: Failed to create swatch at index %d", i)
			continue
		}

		// Select first swatch by default
		if i == 0 {
			s.Selected = true
			app.State.SetSwatchSelected(0)
			app.State.SetBrushColor(s.Color)
		}

		app.Swatches = append(app.Swatches, s)
		canvasSwatches = append(canvasSwatches, s)
	}

	log.Printf("Created %d swatches", len(app.Swatches))

	// Create swatch grid
	swatchGrid := container.NewGridWrap(
		fyne.NewSize(SwatchGridSize, SwatchGridSize),
		canvasSwatches...,
	)

	// Create titled container
	title := widget.NewLabelWithStyle("Color Palette",
		fyne.TextAlignCenter,
		fyne.TextStyle{Bold: true})

	swatchPanel := container.NewBorder(
		container.NewVBox(title, widget.NewSeparator()),
		nil,
		nil,
		nil,
		swatchGrid,
	)

	return swatchPanel
}

// createSwatch creates a single swatch with the proper click handler
func createSwatch(app *AppInit, index int) *swatch.Swatch {
	if app == nil || app.State == nil {
		return nil
	}

	// Create swatch with click handler
	s := swatch.NewSwatch(
		app.State,
		DefaultSwatchColor,
		index,
		func(clickedSwatch *swatch.Swatch) {
			handleSwatchClick(app, clickedSwatch)
		},
	)

	return s
}

// handleSwatchClick handles when a swatch is clicked
func handleSwatchClick(app *AppInit, clickedSwatch *swatch.Swatch) {
	if app == nil || clickedSwatch == nil {
		log.Println("Warning: Invalid swatch click parameters")
		return
	}

	// Use helper function to select the swatch
	swatch.SelectSwatch(clickedSwatch, app.Swatches)

	// Update application state using setter methods
	app.State.SetSwatchSelected(clickedSwatch.SwatchIndex)
	app.State.SetBrushColor(clickedSwatch.Color)

	log.Printf("Selected swatch %d with color %v", clickedSwatch.SwatchIndex, clickedSwatch.Color)
}

// RefreshSwatches refreshes all swatches in the panel
func RefreshSwatches(app *AppInit) {
	if app == nil || len(app.Swatches) == 0 {
		return
	}

	for _, s := range app.Swatches {
		if s != nil {
			s.Refresh()
		}
	}
}

// UpdateSwatchColor updates a specific swatch's color
func UpdateSwatchColor(app *AppInit, index int, c color.Color) error {
	if app == nil || len(app.Swatches) == 0 {
		return nil
	}

	if index < 0 || index >= len(app.Swatches) {
		return nil
	}

	s := app.Swatches[index]
	if s != nil {
		s.SetColor(c)
	}

	return nil
}

// GetSelectedSwatch returns the currently selected swatch
func GetSelectedSwatch(app *AppInit) *swatch.Swatch {
	if app == nil || app.State == nil || len(app.Swatches) == 0 {
		return nil
	}

	index := app.State.SwatchSelected
	if index < 0 || index >= len(app.Swatches) {
		return nil
	}

	return app.Swatches[index]
}

// ClearSwatches resets all swatches to default color
func ClearSwatches(app *AppInit) {
	if app == nil || len(app.Swatches) == 0 {
		return
	}

	for _, s := range app.Swatches {
		if s != nil {
			s.SetColor(DefaultSwatchColor)
		}
	}

	log.Println("Cleared all swatches to default color")
}
