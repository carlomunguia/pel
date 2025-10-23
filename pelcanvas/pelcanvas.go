// Package pelcanvas provides the pixel canvas widget for the Pel pixel art editor.
package pelcanvas

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"pel/apptype"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

// Canvas rendering constants
const (
	BorderStrokeWidth = 2
	DefaultGrayValue  = 128
)

// Default colors
var (
	BorderColor       = color.NRGBA{R: 100, G: 100, B: 100, A: 255}
	DefaultCanvasGray = color.NRGBA{R: DefaultGrayValue, G: DefaultGrayValue, B: DefaultGrayValue, A: 255}
)

// PelCanvasMouseState tracks the mouse state for pan/drag operations
type PelCanvasMouseState struct {
	previousCoord *fyne.PointEvent
}

// PelCanvas is the main canvas widget for drawing pixel art
type PelCanvas struct {
	widget.BaseWidget
	apptype.PelCanvasConfig
	renderer    *PelCanvasRenderer
	PixelData   image.Image
	mouseState  PelCanvasMouseState
	appState    *apptype.State
	reloadImage bool
}

// Bounds returns the current bounds of the canvas in screen coordinates
func (pelCanvas *PelCanvas) Bounds() image.Rectangle {
	x0 := int(pelCanvas.CanvasOffset.X)
	y0 := int(pelCanvas.CanvasOffset.Y)
	x1 := int(pelCanvas.PxCols*pelCanvas.PxSize + int(pelCanvas.CanvasOffset.X))
	y1 := int(pelCanvas.PxRows*pelCanvas.PxSize + int(pelCanvas.CanvasOffset.Y))
	return image.Rect(x0, y0, x1, y1)
}

// InBounds checks if a position is within the given bounds
func (pelCanvas *PelCanvas) InBounds(pos fyne.Position) bool {
	bounds := pelCanvas.Bounds()
	return pos.X >= float32(bounds.Min.X) &&
		pos.X < float32(bounds.Max.X) &&
		pos.Y >= float32(bounds.Min.Y) &&
		pos.Y < float32(bounds.Max.Y)
}

// NewBlankImage creates a new image filled with the specified color
func NewBlankImage(cols, rows int, c color.Color) (image.Image, error) {
	if cols <= 0 || rows <= 0 {
		return nil, fmt.Errorf("invalid image dimensions: %dx%d", cols, rows)
	}

	img := image.NewNRGBA(image.Rect(0, 0, cols, rows))
	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			img.Set(x, y, c)
		}
	}
	return img, nil
}

// NewPelCanvas creates a new pixel canvas with the given configuration
func NewPelCanvas(state *apptype.State, config apptype.PelCanvasConfig) *PelCanvas {
	if state == nil {
		log.Fatal("Cannot create PelCanvas: state is nil")
	}

	// Validate configuration
	if err := config.Validate(); err != nil {
		log.Fatalf("Invalid canvas configuration: %v", err)
	}

	pelCanvas := &PelCanvas{
		PelCanvasConfig: config,
		appState:        state,
	}

	// Create initial blank image
	img, err := NewBlankImage(pelCanvas.PxCols, pelCanvas.PxRows, DefaultCanvasGray)
	if err != nil {
		log.Fatalf("Failed to create blank image: %v", err)
	}
	pelCanvas.PixelData = img

	pelCanvas.ExtendBaseWidget(pelCanvas)
	log.Printf("Created PelCanvas: %dx%d grid", pelCanvas.PxCols, pelCanvas.PxRows)

	return pelCanvas
}

// CreateRenderer creates the renderer for the canvas widget
func (pelCanvas *PelCanvas) CreateRenderer() fyne.WidgetRenderer {
	canvasImage := canvas.NewImageFromImage(pelCanvas.PixelData)
	canvasImage.ScaleMode = canvas.ImageScalePixels
	canvasImage.FillMode = canvas.ImageFillContain

	// Create border lines
	canvasBorder := make([]canvas.Line, 4)
	for i := 0; i < len(canvasBorder); i++ {
		canvasBorder[i].StrokeColor = BorderColor
		canvasBorder[i].StrokeWidth = BorderStrokeWidth
	}

	renderer := &PelCanvasRenderer{
		pelCanvas:    pelCanvas,
		canvasImage:  canvasImage,
		canvasBorder: canvasBorder,
	}
	pelCanvas.renderer = renderer

	return renderer
}

// TryPan attempts to pan the canvas if the middle mouse button is pressed
func (pelCanvas *PelCanvas) TryPan(previousCoord *fyne.PointEvent, ev *desktop.MouseEvent) {
	if ev == nil {
		return
	}

	if previousCoord != nil && ev.Button == desktop.MouseButtonTertiary {
		pelCanvas.Pan(*previousCoord, ev.PointEvent)
	}
}

// SetColor sets the color of a pixel at the specified coordinates
// Returns an error if the coordinates are out of bounds or if the operation fails
func (pelCanvas *PelCanvas) SetColor(c color.Color, x, y int) error {
	if c == nil {
		return fmt.Errorf("color cannot be nil")
	}

	// Validate coordinates
	if x < 0 || x >= pelCanvas.PxCols || y < 0 || y >= pelCanvas.PxRows {
		return fmt.Errorf("coordinates out of bounds: (%d, %d)", x, y)
	}

	// Try to set color on the image
	success := false
	if nrgba, ok := pelCanvas.PixelData.(*image.NRGBA); ok {
		nrgba.Set(x, y, c)
		success = true
	} else if rgba, ok := pelCanvas.PixelData.(*image.RGBA); ok {
		rgba.Set(x, y, c)
		success = true
	}

	if !success {
		return fmt.Errorf("unsupported image type: %T", pelCanvas.PixelData)
	}

	pelCanvas.Refresh()
	return nil
}

// MouseToCanvasXY converts mouse event coordinates to canvas pixel coordinates
// Returns nil pointers if the coordinates are outside the canvas
func (pelCanvas *PelCanvas) MouseToCanvasXY(ev *desktop.MouseEvent) (*int, *int) {
	if ev == nil {
		return nil, nil
	}

	if !pelCanvas.InBounds(ev.Position) {
		return nil, nil
	}

	pxSize := float32(pelCanvas.PxSize)
	xOffset := pelCanvas.CanvasOffset.X
	yOffset := pelCanvas.CanvasOffset.Y

	x := int((ev.Position.X - xOffset) / pxSize)
	y := int((ev.Position.Y - yOffset) / pxSize)

	return &x, &y
}

// LoadImage loads an image into the canvas
// The canvas dimensions will be adjusted to match the image
func (pelCanvas *PelCanvas) LoadImage(img image.Image) error {
	if img == nil {
		return fmt.Errorf("cannot load nil image")
	}

	dimensions := img.Bounds()
	cols := dimensions.Dx()
	rows := dimensions.Dy()

	if cols <= 0 || rows <= 0 {
		return fmt.Errorf("invalid image dimensions: %dx%d", cols, rows)
	}

	pelCanvas.PelCanvasConfig.PxCols = cols
	pelCanvas.PelCanvasConfig.PxRows = rows
	pelCanvas.PixelData = img
	pelCanvas.reloadImage = true

	log.Printf("Loaded image: %dx%d pixels", cols, rows)
	pelCanvas.Refresh()

	return nil
}

// NewDrawing creates a new blank drawing with the specified dimensions
func (pelCanvas *PelCanvas) NewDrawing(cols, rows int) error {
	if cols <= 0 || rows <= 0 {
		return fmt.Errorf("invalid drawing dimensions: %dx%d", cols, rows)
	}

	// Clear file path to indicate unsaved new drawing
	pelCanvas.appState.SetFilePath("")

	// Update dimensions
	pelCanvas.PxCols = cols
	pelCanvas.PxRows = rows

	// Create blank image
	pixelData, err := NewBlankImage(cols, rows, DefaultCanvasGray)
	if err != nil {
		return fmt.Errorf("failed to create blank image: %w", err)
	}

	// Load the new image
	if err := pelCanvas.LoadImage(pixelData); err != nil {
		return fmt.Errorf("failed to load image: %w", err)
	}

	log.Printf("Created new drawing: %dx%d", cols, rows)
	return nil
}

// GetPixelColor returns the color at the specified pixel coordinates
func (pelCanvas *PelCanvas) GetPixelColor(x, y int) (color.Color, error) {
	if x < 0 || x >= pelCanvas.PxCols || y < 0 || y >= pelCanvas.PxRows {
		return nil, fmt.Errorf("coordinates out of bounds: (%d, %d)", x, y)
	}

	return pelCanvas.PixelData.At(x, y), nil
}

// Clear fills the entire canvas with the specified color
func (pelCanvas *PelCanvas) Clear(c color.Color) error {
	if c == nil {
		c = DefaultCanvasGray
	}

	img, err := NewBlankImage(pelCanvas.PxCols, pelCanvas.PxRows, c)
	if err != nil {
		return fmt.Errorf("failed to clear canvas: %w", err)
	}

	return pelCanvas.LoadImage(img)
}
