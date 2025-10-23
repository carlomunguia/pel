// Package ui provides color picker functionality for the Pel pixel art editor.
package ui

import (
	"fmt"
	"image/color"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/lusingander/colorpicker"
)

// Color picker constants
const (
	PickerWidth    = 200 // Width of the color picker in pixels
	PreviewSize    = 40  // Size of the color preview square
	MinPickerWidth = 100
	MaxPickerWidth = 400
)

// SetupColorPicker creates and configures the color picker panel
func SetupColorPicker(app *AppInit) *fyne.Container {
	if app == nil {
		log.Println("Warning: Cannot setup color picker - app is nil")
		return container.NewVBox()
	}

	if app.State == nil {
		log.Println("Warning: Cannot setup color picker - state is nil")
		return container.NewVBox()
	}

	log.Println("Setting up color picker...")

	// Create the main color picker widget
	picker := colorpicker.New(PickerWidth, colorpicker.StyleHue)

	// Set initial color from state
	if app.State.BrushColor != nil {
		picker.SetColor(app.State.BrushColor)
	}

	// Create color preview
	preview := createColorPreview(app)

	// Create color info label
	colorInfo := widget.NewLabel(formatColorInfo(app.State.BrushColor))

	// Setup color change handler
	picker.SetOnChanged(func(c color.Color) {
		handleColorChange(app, c, preview, colorInfo)
	})

	// Create layout with title
	title := widget.NewLabelWithStyle("Color Picker",
		fyne.TextAlignCenter,
		fyne.TextStyle{Bold: true})

	pickerContainer := container.NewVBox(
		title,
		widget.NewSeparator(),
		picker,
		widget.NewSeparator(),
		container.NewVBox(
			widget.NewLabel("Current Color:"),
			preview,
			colorInfo,
		),
	)

	log.Println("Color picker setup complete")
	return pickerContainer
}

// createColorPreview creates a visual preview of the current color
func createColorPreview(app *AppInit) *canvas.Rectangle {
	if app == nil || app.State == nil || app.State.BrushColor == nil {
		return canvas.NewRectangle(color.Black)
	}

	preview := canvas.NewRectangle(app.State.BrushColor)
	preview.SetMinSize(fyne.NewSize(PreviewSize, PreviewSize))
	preview.StrokeColor = color.NRGBA{R: 128, G: 128, B: 128, A: 255}
	preview.StrokeWidth = 2

	return preview
}

// handleColorChange processes color picker changes and updates the application state
func handleColorChange(app *AppInit, c color.Color, preview *canvas.Rectangle, colorInfo *widget.Label) {
	if app == nil || app.State == nil || c == nil {
		log.Println("Warning: Invalid color change parameters")
		return
	}

	// Update application state using setter method
	app.State.SetBrushColor(c)

	// Update the current swatch if valid
	if err := updateCurrentSwatch(app, c); err != nil {
		log.Printf("Warning: Failed to update swatch: %v", err)
		// Continue anyway - this is not critical
	}

	// Update preview
	if preview != nil {
		preview.FillColor = c
		canvas.Refresh(preview)
	}

	// Update color info label
	if colorInfo != nil {
		colorInfo.SetText(formatColorInfo(c))
	}

	log.Printf("Color changed to: %s", formatColorInfo(c))
}

// updateCurrentSwatch updates the currently selected swatch with the new color
func updateCurrentSwatch(app *AppInit, c color.Color) error {
	if app == nil || len(app.Swatches) == 0 {
		return nil // No swatches to update
	}

	selectedIndex := app.State.SwatchSelected
	if selectedIndex < 0 || selectedIndex >= len(app.Swatches) {
		return fmt.Errorf("swatch index %d out of bounds (max: %d)",
			selectedIndex, len(app.Swatches)-1)
	}

	swatch := app.Swatches[selectedIndex]
	if swatch == nil {
		return fmt.Errorf("swatch at index %d is nil", selectedIndex)
	}

	swatch.SetColor(c)
	return nil
}

// formatColorInfo returns a human-readable string representation of a color
func formatColorInfo(c color.Color) string {
	if c == nil {
		return "No color"
	}

	r, g, b, a := c.RGBA()
	// RGBA() returns 16-bit values, convert to 8-bit
	r8 := uint8(r >> 8)
	g8 := uint8(g >> 8)
	b8 := uint8(b >> 8)
	a8 := uint8(a >> 8)

	// Show hex format
	if a8 == 255 {
		return fmt.Sprintf("#%02X%02X%02X", r8, g8, b8)
	}
	return fmt.Sprintf("#%02X%02X%02X%02X", r8, g8, b8, a8)
}
