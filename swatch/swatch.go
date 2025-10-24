// Package swatch provides color swatch widgets for palette management.
package swatch

import (
	"image/color"
	"github.com/carlomunguia/pel/apptype"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
)

// Swatch display constants
const (
	SwatchSize         = 20  // Size of the swatch in pixels
	SwatchBorderWidth  = 2   // Border width when selected
	SwatchCornerRadius = 4   // Rounded corner radius
	MinSwatchSize      = 10  // Minimum swatch size
	MaxSwatchSize      = 100 // Maximum swatch size
)

// Default colors
var (
	SwatchBorderColor  = color.NRGBA{R: 255, G: 255, B: 255, A: 255} // White border for selected
	SwatchDefaultColor = color.NRGBA{R: 128, G: 128, B: 128, A: 255} // Gray default color
)

// Swatch represents a single color swatch in the palette
type Swatch struct {
	widget.BaseWidget
	Selected     bool            // Whether this swatch is currently selected
	Color        color.Color     // The color this swatch represents
	SwatchIndex  int             // Index in the palette
	clickHandler func(s *Swatch) // Handler called when swatch is clicked
}

// swatchRenderer handles the rendering of a color swatch widget
type swatchRenderer struct {
	square  *canvas.Rectangle   // The colored square
	border  *canvas.Rectangle   // Selection border
	objects []fyne.CanvasObject // All objects to render
	parent  *Swatch             // Reference to parent widget
}

// SetColor updates the swatch color and refreshes the display
func (s *Swatch) SetColor(c color.Color) {
	if c == nil {
		c = SwatchDefaultColor
	}
	s.Color = c
	s.Refresh()
}

// SetClickHandler updates the click handler function
func (s *Swatch) SetClickHandler(handler func(s *Swatch)) {
	s.clickHandler = handler
}

// GetIndex returns the swatch index in the palette
func (s *Swatch) GetIndex() int {
	return s.SwatchIndex
}

// NewSwatch creates a new color swatch widget
func NewSwatch(state *apptype.State, color color.Color, swatchIndex int, clickHandler func(s *Swatch)) *Swatch {
	if color == nil {
		color = SwatchDefaultColor
	}

	if swatchIndex < 0 {
		swatchIndex = 0
	}

	swatch := &Swatch{
		Selected:     false,
		Color:        color,
		clickHandler: clickHandler,
		SwatchIndex:  swatchIndex,
	}
	swatch.ExtendBaseWidget(swatch)
	return swatch
}

// CreateRenderer creates the renderer for the swatch widget
func (swatch *Swatch) CreateRenderer() fyne.WidgetRenderer {
	// Create the colored square
	square := canvas.NewRectangle(swatch.Color)

	// Create the selection border (initially invisible)
	border := canvas.NewRectangle(color.Transparent)
	border.StrokeColor = SwatchBorderColor
	border.StrokeWidth = SwatchBorderWidth

	// Border should be behind the square
	objects := []fyne.CanvasObject{border, square}

	return &swatchRenderer{
		square:  square,
		border:  border,
		objects: objects,
		parent:  swatch,
	}
}

// MinSize returns the minimum size for the swatch
func (renderer *swatchRenderer) MinSize() fyne.Size {
	return fyne.NewSize(SwatchSize, SwatchSize)
}

// Layout positions and sizes the swatch elements
func (renderer *swatchRenderer) Layout(size fyne.Size) {
	if renderer.square == nil || renderer.border == nil {
		return
	}

	// Square fills the entire area
	renderer.square.Resize(size)
	renderer.square.Move(fyne.NewPos(0, 0))

	// Border around the square
	renderer.border.Resize(size)
	renderer.border.Move(fyne.NewPos(0, 0))
}

// Refresh updates the swatch display based on current state
func (renderer *swatchRenderer) Refresh() {
	if renderer.parent == nil || renderer.square == nil || renderer.border == nil {
		return
	}

	// Update square color
	renderer.square.FillColor = renderer.parent.Color
	canvas.Refresh(renderer.square)

	// Update border visibility based on selection state
	if renderer.parent.Selected {
		renderer.border.FillColor = color.Transparent
		renderer.border.StrokeColor = SwatchBorderColor
		renderer.border.StrokeWidth = SwatchBorderWidth
	} else {
		renderer.border.FillColor = color.Transparent
		renderer.border.StrokeColor = color.Transparent
		renderer.border.StrokeWidth = 0
	}
	canvas.Refresh(renderer.border)
}

// Objects returns all canvas objects that need to be rendered
func (renderer *swatchRenderer) Objects() []fyne.CanvasObject {
	return renderer.objects
}

// Destroy cleans up any resources used by the renderer
func (renderer *swatchRenderer) Destroy() {
	renderer.square = nil
	renderer.border = nil
	renderer.objects = nil
	renderer.parent = nil
}

// DeselectAll deselects all swatches in the provided list
func DeselectAll(swatches []*Swatch) {
	for _, s := range swatches {
		if s != nil && s.Selected {
			s.Selected = false
			s.Refresh()
		}
	}
}

// SelectSwatch selects a swatch and deselects all others
func SelectSwatch(target *Swatch, allSwatches []*Swatch) {
	if target == nil {
		return
	}

	// Deselect all others
	for _, s := range allSwatches {
		if s != nil && s != target && s.Selected {
			s.Selected = false
			s.Refresh()
		}
	}

	// Select the target
	target.Selected = true
	target.Refresh()
}

// FindSwatchByIndex returns the swatch at the specified index, or nil if not found
func FindSwatchByIndex(swatches []*Swatch, index int) *Swatch {
	for _, s := range swatches {
		if s != nil && s.SwatchIndex == index {
			return s
		}
	}
	return nil
}

// GetSwatchCount returns the number of non-nil swatches
func GetSwatchCount(swatches []*Swatch) int {
	count := 0
	for _, s := range swatches {
		if s != nil {
			count++
		}
	}
	return count
}
